package kkt

import (
	"fmt"
	"log"

	"github.com/fess932/shtrih-m-driver/examples/client/commands"
)

func (kkt *KKT) parseCmd(cmd []byte) error {
	if cmd[0] == 0xFF {
		return parseFNcmd(cmd[1:])
	}
	if cmd[1] != 0x00 {
		return errCheck(cmd[1])
	}
	if len(cmd) <= 2 { // если длинна команды 2 то это пустая команда не требующая обработки имеющая лишь код ошибки
		return nil
	}

	f, ok := routes[cmd[0]]
	if !ok {
		return fmt.Errorf("not found cmd handler for: %v", cmd[0])
	}

	f(cmd[2:], kkt)

	return nil
}

var routes = map[byte]func(cmd []byte, kkt *KKT){
	commands.ShortStatus: updateState,
}

func updateState(cmd []byte, kkt *KKT) {
	st := status(cmd)
	log.Println("cmd:", st, kkt.Addr)

	// main State
	switch st.state {
	case 2:
		kkt.State.SetState("shiftOpen")
	case 3:
		kkt.State.SetState("shiftExpired")
	case 4:
		kkt.State.SetState("shiftClosed")
	case 8:
		kkt.State.SetState("checkOpen")

	default:
		kkt.State.SetState("wrongState")
	}

	// Substate
	switch st.substate {
	case 0:
		kkt.Substate.SetState("paperLoaded")
	default:
		kkt.Substate.SetState("wrongSubstate")
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
