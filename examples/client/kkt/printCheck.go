package kkt

import (
	"fmt"
	"github.com/fess932/shtrih-m-driver/examples/client/commands"
	"github.com/fess932/shtrih-m-driver/pkg/driver/models"
	"log"
)

// send printCmd[]
// if specified send dontPrintOneCheck
// send writeCashierInn
// send closeCheck

func PrintCheckHandler(check models.CheckPackage) func(kkt *KKT) error {
	return func(kkt *KKT) (err error) {
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
			if err = notPrintOneCheck(kkt); err != nil {
				log.Println(err)
				return
			}
		}

		// add operationV2 to check
		for _, v := range check.Operations {
			if err = sendOperationsV2(kkt, v); err != nil {
				log.Println(err)
				return err
			}
		}

		if err = writeCashierINN(kkt, check.CashierINN); err != nil {
			log.Println(err)
			return err
		}

		return sendCloseCheckV2(kkt, check)
	}
}

func sendOperationsV2(kkt *KKT, o models.Operation) error {
	data, err := commands.CreateFNOperationV2(o)
	if err != nil {
		return err
	}

	resp, err := kkt.m.SendMessage(data)
	if err != nil {
		return err
	}

	return kkt.parseCmd(resp)
}

func sendCloseCheckV2(kkt *KKT, check models.CheckPackage) error {
	data, err := commands.CreateFNCloseCheck(check)
	if err != nil {
		return err
	}

	resp, err := kkt.m.SendMessage(data)
	if err != nil {
		return err
	}

	return kkt.parseCmd(resp)
}

func writeCashierINN(kkt *KKT, inn string) error {
	data, err := commands.CreateWriteCashierINN(inn)
	if err != nil {
		return err
	}

	resp, err := kkt.m.SendMessage(data)
	if err != nil {
		return err
	}

	return kkt.parseCmd(resp)
}

func notPrintOneCheck(kkt *KKT) (err error) {
	resp, err := kkt.m.SendMessage(commands.CreateNotPrintOneCheck())
	if err != nil {
		return err
	}

	return kkt.parseCmd(resp)
}
