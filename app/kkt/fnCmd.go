package kkt

import (
	"encoding/hex"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/re-star-ru/shtrih-m-driver/app/commands"
)

func parseFNcmd(fncmd []byte, kkt *KKT) error {
	log.Print("parce fncmd")
	log.Print(hex.Dump(fncmd))

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
	commands.FnOperationV2:  fnOperationV2,
	commands.FNStatus:       fnStatus,
}

func closeCheckV2(fncmd []byte, kkt *KKT) {
	log.Printf("close check handler %X, %s\n", fncmd, kkt.Addr)
}

func fnOperationV2(fncmd []byte, kkt *KKT) {
	log.Printf("fn operation handler %X, %s\n", fncmd, kkt.Addr)
}

func fnStatus(fncmd []byte, kkt *KKT) {
	log.Printf("fn status handler %X, %s\n", fncmd, kkt.Addr)
}
