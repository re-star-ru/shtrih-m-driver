package kkt

import (
	"fmt"

	"github.com/looplab/fsm"
)

// states ---------------------------------------------------------------------

type StateType int

const (
	ShiftOpened StateType = iota + 1
	ShiftClosed
)

// events -----------------------------------------------------------------------

type EventType int

const NoOp EventType = -1
const (
	ShiftOpen EventType = iota + 1
	ShiftClose
)

// ------------------------------------------------------------------------------

type ShiftOpenAction struct{}

func (a *ShiftOpenAction) Execute() EventType {
	fmt.Println("Shift has been opened")
	return NoOp
}

type ShiftCloseAction struct{}

func (a *ShiftCloseAction) Execute() EventType {
	fmt.Println("Shift has been closed")
	return NoOp
}

// -------------------------------------------------------------------------------

type KKT struct {
	Addr string
	FSM  *fsm.FSM
}

func New(addr string) *KKT {
	k := &KKT{Addr: addr}

	k.FSM = fsm.NewFSM(
		"closed",
		fsm.Events{
			{Name: "open", Src: []string{"closed"}, Dst: "open"},
			{Name: "close", Src: []string{"open"}, Dst: "closed"},
		},
		fsm.Callbacks{
			"enter_state": func(e *fsm.Event) { k.enterState(e) },
		},
	)

	return k
}

func (kkt KKT) enterState(e *fsm.Event) {
	fmt.Printf("the kkt with addr %s is %s\n", kkt.Addr, e.Dst)
}

//func (kkt *KKT) Exec(cmd interface{ Exec() }) {
//	var s PrintCheckCommand
//
//	switch kkt.state {
//	case ShiftOpened:
//
//	}
//
//	s = func() {
//
//	}
//
//}

type PrintCheckCommand func()

//
//func (c PrintCheckCommand) Exec() {
//
//}

func (kkt *KKT) PrintCheck() {

}
