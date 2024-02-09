package matchmaking

import (
	"fmt"
	"github.com/CameronHonis/chess-arbitrator/builders"
	"github.com/CameronHonis/chess-arbitrator/matcher"
	"github.com/CameronHonis/chess-arbitrator/models"
	. "github.com/CameronHonis/log"
	. "github.com/CameronHonis/marker"
	. "github.com/CameronHonis/service"
	"math"
	"time"
)

type MatchmakingServiceI interface {
	ServiceI
	AddClient(client *models.ClientProfile) error
	RemoveClient(client *models.ClientProfile) error
}

type MatchmakingService struct {
	Service

	__dependencies__ Marker
	LogService       LoggerServiceI
	MatchService     matcher.MatcherServiceI

	__state__ Marker
	pool      *MatchmakingPool
}

func NewMatchmakingService(config *MatchmakingConfig) *MatchmakingService {
	matchmakingService := &MatchmakingService{
		pool: NewMatchmakingPool(),
	}
	matchmakingService.Service = *NewService(matchmakingService, config)

	return matchmakingService
}

func (mm *MatchmakingService) OnStart() {
	go mm.loopMatchmaking()
}

func (mm *MatchmakingService) AddClient(client *models.ClientProfile) error {
	mm.LogService.Log(models.ENV_MATCHMAKING, fmt.Sprintf("adding client %s to matchmaking pool", client.ClientKey))
	addErr := mm.pool.AddClient(client)
	if addErr != nil {
		return addErr
	}
	mm.LogService.Log(models.ENV_MATCHMAKING, fmt.Sprintf("%d clients in pool", len(mm.pool.nodeByClientKey)))
	return nil
}

func (mm *MatchmakingService) RemoveClient(client *models.ClientProfile) error {
	return mm.pool.RemoveClient(client.ClientKey)
}

func (mm *MatchmakingService) loopMatchmaking() {
	for {
		time.Sleep(time.Second)
		if mm.pool.Head() == mm.pool.Tail() {
			continue
		}
		currPoolNode := mm.pool.Head()
		for currPoolNode != nil && currPoolNode.next != nil {
			waitTime := time.Now().Unix() - currPoolNode.timeJoined
			bestMatchPoolNode := currPoolNode.next
			bestMatchWeight := float64(10000000000)
			nextPoolNode := currPoolNode.next
			for nextPoolNode != nil {
				nextPoolNodeMatchWeight := weightProfileDiff(currPoolNode.clientProfile, nextPoolNode.clientProfile, waitTime)
				if nextPoolNodeMatchWeight < bestMatchWeight {
					bestMatchPoolNode = nextPoolNode
					bestMatchWeight = nextPoolNodeMatchWeight
				}
				nextPoolNode = nextPoolNode.next
			}
			clientA := currPoolNode.clientProfile
			clientB := bestMatchPoolNode.clientProfile
			if IsMatchable(clientA, clientB, waitTime) {
				matchErr := mm.matchClients(clientA, clientB)
				if matchErr != nil {
					mm.LogService.LogRed(models.ENV_MATCHMAKING, fmt.Sprintf("error matching clients %s and %s: %s\n", clientA.ClientKey, clientB.ClientKey, matchErr))
				} else {
					mm.LogService.LogGreen(models.ENV_MATCHMAKING, fmt.Sprintf("matched clients %s and %s\n", clientA.ClientKey, clientB.ClientKey))
				}
			}
			currPoolNode = currPoolNode.next
		}
	}
}

func (mm *MatchmakingService) matchClients(clientA *models.ClientProfile, clientB *models.ClientProfile) error {
	removeErr := mm.pool.RemoveClient(clientA.ClientKey)
	if removeErr != nil {
		return fmt.Errorf("error removing client %s from matchmaking pool: %s", clientA.ClientKey, removeErr)
	}
	removeErr = mm.pool.RemoveClient(clientB.ClientKey)
	if removeErr != nil {
		return fmt.Errorf("error removing client %s from matchmaking pool: %s", clientB.ClientKey, removeErr)
	}
	match := builders.NewMatch(clientA.ClientKey, clientB.ClientKey, builders.NewBulletTimeControl(), models.MATCH_RESULT_IN_PROGRESS)
	addMatchErr := mm.MatchService.AddMatch(match)
	if addMatchErr != nil {
		return fmt.Errorf("error adding match %s: %s", match.Uuid, addMatchErr)
	}
	return nil
}

func IsMatchable(clientA *models.ClientProfile, clientB *models.ClientProfile, longestWaitSeconds int64) bool {
	return weightProfileDiff(clientA, clientB, longestWaitSeconds) <= 100
}

func weightProfileDiff(p1 *models.ClientProfile, p2 *models.ClientProfile, longestWaitSeconds int64) float64 {
	eloDiff := math.Abs(float64(p1.Elo - p2.Elo))
	eloCoeff := 100 / (float64(longestWaitSeconds) + 50) // asymptotic curve approaches 0, 2 @ t=0, 1 @ t=50, 0.5 @ t=100
	return eloDiff * eloCoeff
}
