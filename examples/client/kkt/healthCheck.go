package kkt

import (
	"github.com/fess932/shtrih-m-driver/examples/client/commands"
	"log"
	"time"
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
			log.Println("day ended, closing shift")
			return closeSession(kkt)
		}

		return

	case "shiftClosed":
		t := time.Now()
		if t.Hour() >= 7 && t.Hour() <= 22 {
			log.Println("day goes, open shift")
			return openSession(kkt)
		}

		return

	case "shiftExpired":
		log.Println("shift expired, closing")
		return closeSession(kkt)

	case "checkOpen":
		log.Println("check open in check state, cancel check")
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
