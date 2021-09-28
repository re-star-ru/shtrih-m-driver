package kkt

import (
	"context"
	"fmt"

	"github.com/re-star-ru/shtrih-m-driver/app/commands"
	"github.com/re-star-ru/shtrih-m-driver/app/models"
	"github.com/rs/zerolog/log"
)

// send printCmd[]
// if specified send dontPrintOneCheck
// send writeCashierInn
// send closeCheck

func PrintCheckHandler(check models.CheckPackage) func(context context.Context, kkt *KKT) error {
	return func(context context.Context, kkt *KKT) (err error) {
		// check state
		if !kkt.canPrintCheck() { // check State
			if kkt.Substate.Is("paperEmpty") {
				return fmt.Errorf("закончилась бумага, заправьте принтер, состояние: %v", kkt.Substate)
			}

			return fmt.Errorf("невозможно напечатать чек, неправильное состояние кассы %v", kkt.State.Current())
		}

		if err = commands.ValidateINN(check.CashierINN); err != nil {
			return
		}

		// set state not print one check if specified
		if check.NotPrint {
			if err = notPrintOneCheck(context, kkt); err != nil {
				log.Err(err).Send()
				return
			}
		}

		// add operationV2 to check
		for _, v := range check.Operations {
			if err = sendOperationsV2(context, kkt, v); err != nil {
				log.Err(err).Send()
				return err
			}
		}

		if err = writeCashierINN(context, kkt, check.CashierINN); err != nil {
			log.Err(err).Send()
			return err
		}

		return sendCloseCheckV2(context, kkt, check)
	}
}

func sendOperationsV2(context context.Context, kkt *KKT, o models.Operation) error {
	data, err := commands.CreateFNOperationV2(o)
	if err != nil {
		return err
	}

	resp, err := kkt.m.SendMessage(context, data)
	if err != nil {
		return err
	}

	return kkt.parseCmd(resp)
}

func sendCloseCheckV2(context context.Context, kkt *KKT, check models.CheckPackage) error {
	data, err := commands.CreateFNCloseCheck(check)
	if err != nil {
		return err
	}

	resp, err := kkt.m.SendMessage(context, data)
	if err != nil {
		return err
	}

	return kkt.parseCmd(resp)
}

func writeCashierINN(context context.Context, kkt *KKT, inn string) error {
	data, err := commands.CreateWriteCashierINN(inn)
	if err != nil {
		return err
	}

	resp, err := kkt.m.SendMessage(context, data)
	if err != nil {
		return err
	}

	return kkt.parseCmd(resp)
}

func notPrintOneCheck(context context.Context, kkt *KKT) (err error) {
	resp, err := kkt.m.SendMessage(context, commands.CreateNotPrintOneCheck())
	if err != nil {
		return err
	}

	return kkt.parseCmd(resp)
}
