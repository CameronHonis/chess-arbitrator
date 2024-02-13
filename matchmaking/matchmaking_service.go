package matchmaking

import (
	"fmt"
	"github.com/CameronHonis/chess-arbitrator/builders"
	"github.com/CameronHonis/chess-arbitrator/matcher"
	"github.com/CameronHonis/chess-arbitrator/models"
	"github.com/CameronHonis/log"
	"github.com/CameronHonis/marker"
	"github.com/CameronHonis/service"
	"sync"
	"time"
)

type MatchmakingServiceI interface {
	service.ServiceI
	AddClient(client *models.ClientProfile, timeControl *models.TimeControl) error
	RemoveClient(clientKey models.Key) error
	GetClientCountByTimeControl(timeControl *models.TimeControl) int
}

type MatchmakingService struct {
	service.Service

	__dependencies__ marker.Marker
	LogService       log.LoggerServiceI
	MatchService     matcher.MatcherServiceI

	__state__             marker.Marker
	poolByTimeControlHash map[string]*MatchmakingPool
	poolByClientKey       map[models.Key]*MatchmakingPool
	mu                    sync.Mutex
}

func NewMatchmakingService(config *MatchmakingConfig) *MatchmakingService {
	matchmakingService := &MatchmakingService{
		poolByTimeControlHash: make(map[string]*MatchmakingPool),
		poolByClientKey:       make(map[models.Key]*MatchmakingPool),
		mu:                    sync.Mutex{},
	}
	matchmakingService.Service = *service.NewService(matchmakingService, config)

	return matchmakingService
}

func (mm *MatchmakingService) OnStart() {
	go mm.loopMatchmaking()
}

func (mm *MatchmakingService) AddClient(client *models.ClientProfile, timeControl *models.TimeControl) error {
	mm.LogService.Log(models.ENV_MATCHMAKING, fmt.Sprintf("adding client %s to matchmaking pool", client.ClientKey))

	timeControlHash := timeControl.Hash()

	mm.mu.Lock()
	defer mm.mu.Unlock()
	if _, ok := mm.poolByClientKey[client.ClientKey]; ok {
		return fmt.Errorf("client with key %s already in a different timeControl pool", client.ClientKey)
	}
	pool := mm.poolByTimeControlHash[timeControlHash]
	if pool == nil {
		pool = NewMatchmakingPool()
		mm.poolByTimeControlHash[timeControlHash] = pool
	}

	mm.poolByClientKey[client.ClientKey] = pool
	addErr := pool.AddClient(client, timeControl)
	if addErr != nil {
		return addErr
	}
	mm.LogService.Log(models.ENV_MATCHMAKING, fmt.Sprintf("%d clients in pool", len(pool.nodeByClientKey)))
	return nil
}

func (mm *MatchmakingService) RemoveClient(clientKey models.Key) error {
	mm.mu.Lock()
	defer mm.mu.Unlock()

	mm.LogService.Log(models.ENV_MATCHMAKING, fmt.Sprintf("removing client %s from matchmaking pool", clientKey))
	pool := mm.poolByClientKey[clientKey]
	if pool == nil {
		return fmt.Errorf("no pool found for client %s", clientKey)
	}

	delete(mm.poolByClientKey, clientKey)
	removeErr := pool.RemoveClient(clientKey)
	if removeErr != nil {
		return removeErr
	}

	mm.LogService.Log(models.ENV_MATCHMAKING, fmt.Sprintf("%d clients in pool", len(pool.nodeByClientKey)))
	return nil
}

func (mm *MatchmakingService) GetClientCountByTimeControl(timeControl *models.TimeControl) int {
	mm.mu.Lock()
	defer mm.mu.Unlock()

	pool := mm.poolByTimeControlHash[timeControl.Hash()]
	if pool == nil {
		return 0
	}
	return len(pool.nodeByClientKey)
}

func (mm *MatchmakingService) loopMatchmaking() {
	for {
		time.Sleep(time.Second)
		for _, pool := range mm.poolByTimeControlHash {
			mm.mu.Lock()
			head := pool.Head()
			tail := pool.Tail()
			mm.mu.Unlock()

			if head == tail {
				continue
			}
			currPoolNode := head
			for currPoolNode != nil && currPoolNode.next != nil {
				waitTime := time.Now().Unix() - currPoolNode.timeJoined

				clientA := currPoolNode.clientProfile
				clientB, _ := pool.GetBestMatch(currPoolNode, waitTime)
				if clientB == nil {
					continue
				}

				matchErr := mm.MatchClient(clientA, clientB)
				if matchErr != nil {
					mm.LogService.LogRed(models.ENV_MATCHMAKING, fmt.Sprintf("error matching clients %s and %s: %s\n", clientA.ClientKey, clientB.ClientKey, matchErr))
				} else {
					mm.LogService.LogGreen(models.ENV_MATCHMAKING, fmt.Sprintf("matched clients %s and %s\n", clientA.ClientKey, clientB.ClientKey))
				}

				currPoolNode = currPoolNode.next
			}
		}
	}
}

func (mm *MatchmakingService) MatchClient(clientA *models.ClientProfile, clientB *models.ClientProfile) error {
	removeErr := mm.RemoveClient(clientA.ClientKey)
	if removeErr != nil {
		return fmt.Errorf("error removing client %s from matchmaking pool: %s", clientA.ClientKey, removeErr)
	}
	removeErr = mm.RemoveClient(clientB.ClientKey)
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
