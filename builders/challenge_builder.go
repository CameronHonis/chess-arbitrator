package builders

import (
	"github.com/CameronHonis/chess-arbitrator/models"
	"github.com/google/uuid"
	"time"
)

func NewChallenge(challengerKey models.Key, challengedKey models.Key, isChallengerWhite bool,
	isChallengerBlack bool, timeControl *models.TimeControl, botName string, isActive bool) *models.Challenge {

	challengeId := uuid.New().String()
	now := time.Now()

	return &models.Challenge{
		Uuid:              challengeId,
		ChallengerKey:     challengerKey,
		ChallengedKey:     challengedKey,
		IsChallengerWhite: isChallengerWhite,
		IsChallengerBlack: isChallengerBlack,
		TimeControl:       timeControl,
		BotName:           botName,
		TimeCreated:       &now,
		IsActive:          isActive,
	}
}

type ChallengeBuilder struct {
	challenge *models.Challenge
}

func NewChallengeBuilder() *ChallengeBuilder {
	return &ChallengeBuilder{
		challenge: NewChallenge("", "", false, false, nil, "", true),
	}
}

func (b *ChallengeBuilder) WithUuid(uuid string) *ChallengeBuilder {
	b.challenge.Uuid = uuid
	return b
}

func (b *ChallengeBuilder) WithRandomUuid() *ChallengeBuilder {
	b.challenge.Uuid = uuid.New().String()
	return b
}

func (b *ChallengeBuilder) WithChallengerKey(key models.Key) *ChallengeBuilder {
	b.challenge.ChallengerKey = key
	return b
}

func (b *ChallengeBuilder) WithChallengedKey(key models.Key) *ChallengeBuilder {
	b.challenge.ChallengedKey = key
	return b
}

func (b *ChallengeBuilder) WithIsChallengerWhite(isWhite bool) *ChallengeBuilder {
	b.challenge.IsChallengerWhite = isWhite
	return b
}

func (b *ChallengeBuilder) WithIsChallengerBlack(isBlack bool) *ChallengeBuilder {
	b.challenge.IsChallengerBlack = isBlack
	return b
}

func (b *ChallengeBuilder) WithTimeControl(timeControl *models.TimeControl) *ChallengeBuilder {
	b.challenge.TimeControl = timeControl
	return b
}

func (b *ChallengeBuilder) WithBotName(botName string) *ChallengeBuilder {
	b.challenge.BotName = botName
	return b
}

func (b *ChallengeBuilder) WithTimeCreated(timeCreated *time.Time) *ChallengeBuilder {
	b.challenge.TimeCreated = timeCreated
	return b
}

func (b *ChallengeBuilder) WithIsActive(isActive bool) *ChallengeBuilder {
	b.challenge.IsActive = isActive
	return b
}

func (b *ChallengeBuilder) FromChallenge(challenge *models.Challenge) *ChallengeBuilder {
	challengeCopy := *challenge
	b.challenge = &challengeCopy
	return b
}

func (b *ChallengeBuilder) Build() *models.Challenge {
	return b.challenge
}
