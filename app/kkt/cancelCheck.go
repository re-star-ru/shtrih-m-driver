package kkt

import (
	"github.com/re-star-ru/shtrih-m-driver/app/commands"
)

func cancelCheck(kkt *KKT) error {
	data := commands.CreateCancelCheck()

	resp, err := kkt.m.SendMessage(data)
	if err != nil {
		return err
	}

	return kkt.parseCmd(resp)
}
