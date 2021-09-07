package kkt

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/re-star-ru/shtrih-m-driver/app/commands"
)

func closeSession(kkt *KKT) error {
	data := commands.CreateCloseSession()

	resp, err := kkt.m.SendMessage(data)
	if err != nil {
		return err
	}

	return kkt.parseCmd(resp)
}

func openSession(kkt *KKT) error {
	if err := commands.ValidateINN(kkt.CashierInn); err != nil {
		return err
	}

	log.Print("OPEN SESSION validate inn ok")

	if err := beginOpenSession(kkt); err != nil {
		return fmt.Errorf("err begin open session: %w", err)
	}

	log.Print("OPEN SESSION begin open session ok")

	if err := writeCashierINN(kkt, kkt.CashierInn); err != nil {
		log.Err(err).Send()
		return err
	}

	log.Print("OPEN SESSION write cashier inn ok")

	if err := endOpenSession(kkt); err != nil {
		log.Err(err).Send()
		return err
	}

	log.Print("OPEN SESSION end open session ok")

	return nil
}

func beginOpenSession(kkt *KKT) error {
	data := commands.CreateFNBeginOpenSession()

	resp, err := kkt.m.SendMessage(data)
	if err != nil {
		return err
	}

	return kkt.parseCmd(resp)
}

func endOpenSession(kkt *KKT) error {
	data := commands.CreateOpenSession()

	resp, err := kkt.m.SendMessage(data)
	if err != nil {
		return err
	}

	return kkt.parseCmd(resp)
}
