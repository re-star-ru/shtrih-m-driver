package kkt

import (
	"github.com/fess932/shtrih-m-driver/examples/client/commands"
)

func cancelCheck(kkt *KKT) error {
	data := commands.CreateCancelCheck()

	resp, err := kkt.m.SendMessage(data)
	if err != nil {
		return err
	}

	return kkt.parseCmd(resp)
}
