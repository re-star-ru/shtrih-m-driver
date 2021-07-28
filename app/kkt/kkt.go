package kkt

import "fmt"

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


func NewKKTFSM() *

// -------------------------------------------------------------------------------

type KKT struct {
	state int
}

func New() *KKT {
	return &KKT{}
}

func (kkt *KKT) Exec(cmd interface{ Exec() }) {
	var s PrintCheckCommand

	switch kkt.state {
	case ShiftOpened:

	}

	s = func() {

	}

}

type PrintCheckCommand func()

//
//func (c PrintCheckCommand) Exec() {
//
//}

func (kkt *KKT) PrintCheck() {

}
