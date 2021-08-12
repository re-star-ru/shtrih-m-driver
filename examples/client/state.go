package main

import "github.com/looplab/fsm"

// Events
const (
	printCheck  = "printCheck"
	shiftOpen   = "shiftOpen"
	shiftClose  = "shiftClose"
	shiftReopen = "shiftReopen"
)

func newState() *fsm.FSM {
	return fsm.NewFSM(
		"",
		fsm.Events{
			{Name: printCheck, Src: []string{"shiftOpen"}, Dst: "shiftOpen"},
			{Name: shiftOpen, Src: []string{"shiftClosed"}, Dst: "shiftOpen"},
			{Name: shiftClose, Src: []string{"shiftOpen"}, Dst: "shiftClosed"},
			{Name: shiftReopen, Src: []string{"shiftExpired"}, Dst: "shiftClosed"},
		},
		fsm.Callbacks{},
	)
}

func newSubstate() *fsm.FSM {
	return fsm.NewFSM(
		"",
		fsm.Events{{
			Name: printCheck, Src: []string{"paperLoaded"}, Dst: "paperLoaded",
		}},
		fsm.Callbacks{},
	)
}
