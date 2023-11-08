package server

import (
	"time"
)

func StartTimer(matchId string) {
	go func() {
		lastMoveTime, getLastMoveTimeErr := GetMatchManager().GetLastMoveOccurredTime(matchId)
		if getLastMoveTimeErr != nil {
			return
		}
		for {
			minTimeNanos, getMinTimeSecondsErr := GetMatchManager().GetMatchMinTimeout(matchId)
			if getMinTimeSecondsErr != nil {
				return
			}
			GetLogManager().Log("timer", "sleeping")
			time.Sleep(*minTimeNanos)
			GetLogManager().Log("timer", "checking match time")
			newLastMoveTime, getNewLastMoveTimeErr := GetMatchManager().CheckMatchTime(matchId, *lastMoveTime)
			if getNewLastMoveTimeErr != nil {
				GetLogManager().LogRed("timer", getNewLastMoveTimeErr.Error())
				return
			}
			lastMoveTime = newLastMoveTime
		}
	}()
}
