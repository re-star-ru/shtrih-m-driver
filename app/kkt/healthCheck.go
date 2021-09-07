package kkt

import (
	"time"

	"github.com/rs/zerolog/log"

	"github.com/re-star-ru/shtrih-m-driver/app/commands"
)

// TODO: FN STATE

func healhCheck(kkt *KKT) (err error) {
	resp, err := kkt.m.SendMessage(commands.CreateShortStatus())
	if err != nil {
		return err
	}

	if err = kkt.parseCmd(resp); err != nil {
		return err
	}

	switch kkt.State.Current() {
	case "shiftOpen":
		t := time.Now()
		if t.Hour() <= 6 || t.Hour() >= 23 {
			log.Print("day ended, closing shift")
			return closeSession(kkt)
		}

		return

	case "shiftClosed":
		t := time.Now()
		if t.Hour() >= 7 && t.Hour() <= 22 {
			log.Print("day goes, open shift")
			return openSession(kkt)
		}

		return

	case "shiftExpired":
		log.Print("shift expired, closing")
		return closeSession(kkt)

	case "checkOpen":
		log.Print("check open in check state, cancel check")
		return cancelCheck(kkt)
	}

	// TODO after all: Substate
	//switch kkt.Substate.Current() {
	//case 0:
	//	kkt.Substate.SetState("paperLoaded")
	//default:
	//	kkt.Substate.SetState("wrongSubstate")
	//}

	return
}
