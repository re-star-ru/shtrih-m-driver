package main

import (
	"fmt"
	"log"
)

func parseCmd(cmd []byte) error {
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

	f(cmd[2:])

	return nil
}

var routes = map[byte]func(cmd []byte){
	0x10: status,
}

func status(cmd []byte) {
	log.Println("status:", cmd)

	kktFlags := cmd[1:3]
	kktMode := cmd[4]
	kktSubMode := cmd[5]
	positionsInCheck := cmd[6] + cmd[9]
	voltageBattarey := cmd[7]
	voltage := cmd[8]
	lastPrintStatus := cmd[13]

	str := fmt.Sprintf("kktMode: %v", kktMode)

}
