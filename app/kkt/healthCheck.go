package kkt

import (
	"errors"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/re-star-ru/shtrih-m-driver/app/commands"
)

// TODO: FN STATE

var ErrUnknownState = errors.New("unknown state")

func UpdateState(kkt *KKT) error {
	defer func() {
		log.Printf(
			"health check kkt %v_%v, state: %v, substate: %v",
			kkt.Organization, kkt.Place, kkt.State.Current(), kkt.Substate.Current(),
		)
	}()

	resp, err := kkt.m.SendMessage(commands.CreateShortStatus())
	if err != nil {
		return err
	}

	if err = kkt.parseCmd(resp); err != nil {
		return err
	}

	err = prepareState(kkt)
	if err != nil {
		err = fmt.Errorf("error update shift %w", err)
		return err
	}

	// TODO after all: Substate
	// 0, 1, 2, 3, 4, 5
	// 5,4 - wait 5sec, retry
	// 3 - send continue print cmd
	// 2 - event "NOT INSERT PAPER" then state 3
	// 1 - instert paper, just for fun
	// 0 - ok!
	// switch kkt.Substate.Current() {
	// case 0:
	//	kkt.Substate.SetState("paperLoaded")
	// default:
	//	kkt.Substate.SetState("wrongSubstate")
	//}

	return nil
}

func prepareState(kkt *KKT) error {

	switch kkt.State.Current() {
	case "shiftOpen", "shiftClosed", "shiftExpired":
		return updateShift(kkt, kkt.State.Current())

	case "checkOpen":
		log.Print("check open in check state, cancel check")
		return cancelCheck(kkt)

	default:
		return ErrUnknownState
	}
}

func updateShift(kkt *KKT, state string) error {
	switch state {
	case "shiftOpen":
		t := time.Now()
		if t.Hour() <= 6 || t.Hour() >= 23 {
			log.Print("day ended, closing shift")
			return closeSession(kkt)
		}

		return nil

	case "shiftClosed":
		t := time.Now()
		if t.Hour() >= 7 && t.Hour() <= 22 {
			log.Print("day goes, open shift")
			return openSession(kkt)
		}

		return nil

	case "shiftExpired":
		log.Print("shift expired, closing")
		return closeSession(kkt)

	default:
		return nil
	}
}
