package kkt

import (
	"fmt"

	"github.com/looplab/fsm"
)

//// states ---------------------------------------------------------------------
//
//type StateType int
//
//const (
//	ShiftOpened StateType = iota + 1
//	ShiftClosed
//)
//
//// events -----------------------------------------------------------------------
//
//type EventType int
//
//const NoOp EventType = -1
//const (
//	ShiftOpen EventType = iota + 1
//	ShiftClose
//)
//
//// ------------------------------------------------------------------------------
//
//type ShiftOpenAction struct{}
//
//func (a *ShiftOpenAction) Execute() EventType {
//	fmt.Println("Shift has been opened")
//	return NoOp
//}
//
//type ShiftCloseAction struct{}
//
//func (a *ShiftCloseAction) Execute() EventType {
//	fmt.Println("Shift has been closed")
//	return NoOp
//}
//
//// -------------------------------------------------------------------------------

type KKT struct {
	Addr string
	fsm  *fsm.FSM
}

func New(addr string) *KKT {
	k := &KKT{Addr: addr}

	// init fsm
	k.fsm = fsm.NewFSM(
		"closed",
		fsm.Events{
			{Name: "open", Src: []string{"closed"}, Dst: "open"},
			{Name: "close", Src: []string{"open"}, Dst: "closed"},
			{Name: "printCheck", Src: []string{"open"}},
		},
		fsm.Callbacks{
			"enter_state":  func(e *fsm.Event) { k.enterState(e) },
			"before_close": func(e *fsm.Event) { k.beforeCloseShift(e) },
			"before_open":  func(e *fsm.Event) { k.beforeOpenShift(e) },
		},
	)

	return k
}

func (kkt KKT) enterState(e *fsm.Event) {
	fmt.Printf("the kkt with addr %s is %s\n", kkt.Addr, e.Dst)
}

func (kkt KKT) OpenShift() error {
	return kkt.fsm.Event("open")
}

func (kkt KKT) CloseShift() error {
	return kkt.fsm.Event("close")
}

func (kkt KKT) PrintCheck() error {
	if !kkt.fsm.Can("printCheck") {
		return fmt.Errorf("cannot print check this state: %s", kkt.fsm.Current())
	}

	println("print check")

	return nil
}

// callbacks

func (kkt KKT) beforeCloseShift(e *fsm.Event) {
	fmt.Printf("before close callback shift %s is %s\n", kkt.Addr, e.Dst)
	fmt.Printf("todo: closing shift")

	//err := fmt.Errorf("error befdore closing")
	//if err != nil {
	//	e.Cancel(err)
	//}
}

func (kkt KKT) beforeOpenShift(e *fsm.Event) {
	fmt.Printf("before open callback shift %s is %s\n", kkt.Addr, e.Dst)
	fmt.Printf("todo: opening shift")

	//err := fmt.Errorf("error befdore closing")
	//if err != nil {
	//	e.Cancel(err)
	//}
}
