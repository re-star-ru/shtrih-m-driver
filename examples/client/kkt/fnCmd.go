package kkt

import (
	"fmt"

	"github.com/fess932/shtrih-m-driver/examples/client/commands"
)

func parseFNcmd(fncmd []byte, kkt *KKT) error {
	if fncmd[1] != 0x00 {
		return errCheck(fncmd[1])
	}
	if len(fncmd) <= 2 { // если длинна команды 2 то это пустая команда не требующая обработки имеющая лишь код ошибки
		return nil
	}

	f, ok := fnRoutes[fncmd[0]]
	if !ok {
		return fmt.Errorf("not found cmd handler for: %v", fncmd[0])
	}

	f(fncmd[2:], kkt)

	return nil
}

var fnRoutes = map[byte]func(cmd []byte, kkt *KKT){
	commands.FnCloseCheckV2: closeCheckV2,
}

func closeCheckV2(fncmd []byte, kkt *KKT) {

}
