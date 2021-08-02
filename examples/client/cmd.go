package main

import (
	"fmt"
)

func (kkt *KKT) parseCmd(cmd []byte) error {
	if cmd[0] == 0xFF {
		return parseFNcmd(cmd[1:])
	}
	if cmd[1] != 0x00 {
		return errCheck()
	}

	f, ok := routes[cmd[0]]
	if !ok {
		return fmt.Errorf("not found cmd handler for: %v", cmd[0])
	}

	f(cmd[2:], kkt)

	return nil
}

var routes = map[byte]func(cmd []byte, kkt *KKT){
	0x10: updateState,
}

func updateState(cmd []byte, kkt *KKT) {
	st := status(cmd)

	// main state
	switch st.state {
	case 2:
		kkt.state.SetState("shiftOpen")
	default:
		kkt.state.SetState("wrong state")
	}

	// substate
	switch st.substate {
	case 0:
		kkt.substate.SetState("paperLoaded")
	default:
		kkt.substate.SetState("wrong substate")
	}
}

type state struct {
	state    byte
	substate byte
}

func status(cmd []byte) state {
	kktMode := cmd[3]
	kktSubMode := cmd[4]

	return state{
		state:    kktMode,
		substate: kktSubMode,
	}
}
