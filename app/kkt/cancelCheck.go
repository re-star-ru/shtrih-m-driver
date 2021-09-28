package kkt

import (
	"context"

	"github.com/re-star-ru/shtrih-m-driver/app/commands"
)

func cancelCheck(context context.Context, kkt *KKT) error {
	data := commands.CreateCancelCheck()

	resp, err := kkt.m.SendMessage(context, data)
	if err != nil {
		return err
	}

	return kkt.parseCmd(resp)
}
