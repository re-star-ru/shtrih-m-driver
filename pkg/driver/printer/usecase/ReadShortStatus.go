package printerUsecase

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/fess932/shtrih-m-driver/pkg/driver/models"

	"github.com/fess932/shtrih-m-driver/pkg/consts"
)

func (p *printerUsecase) ReadShortStatus() byte {
	p.logger.Debug("Send command ReadShortStatus")

	buf, cmdLen := p.createCommandBuffer(consts.ReadShortStatus, p.password)

	rFrame, err := p.send(buf.Bytes(), cmdLen)

	if err != nil {
		p.logger.Error(err)
		return 0
	}

	if err := models.CheckOnPrinterError(rFrame.ERR); err != nil {
		p.logger.Error(err)
		return 0
	}

	p.logger.Debug("frame in: \n", hex.Dump(rFrame.Bytes()))

	in := bufio.NewReader(bytes.NewReader(rFrame.DATA))

	operatorNumber, _ := in.ReadByte()

	flags := make([]byte, 2)
	in.Read(flags)

	mode, _ := in.ReadByte() // & 15; ? wft
	subMode, _ := in.ReadByte()
	receiptOperationsLo, _ := in.ReadByte()
	batteryState, _ := in.ReadByte()
	//double batteryVoltage = (double)batteryState / 255.0D * 100.0D * 5.0D / 100.0D;
	//double powerVoltage = (double)powerState * 24.0D / 216.0D * 100.0D / 100.0D;
	powerState, _ := in.ReadByte()
	receiptOperationsHi, _ := in.ReadByte()
	reserved1, _ := in.ReadByte()
	reserved2, _ := in.ReadByte()
	reserved3, _ := in.ReadByte()
	lastResult, _ := in.ReadByte()

	//int receiptOperations = receiptOperationsLo + (receiptOperationsHi << 8);
	str := fmt.Sprintf("\noperator number: %v\nflags: %v\nmode: %v\nsubMode: %v\n"+
		"receiptOperationsLo: %v\nbatteryState: %v\npowerState: %v\n"+
		"receiptOperationsHi: %v\nreserved1: %v\nreserved2: %v\nreserved3: %v\nlastResult: %v", operatorNumber, flags, mode, subMode,
		receiptOperationsLo, batteryState, powerState, receiptOperationsHi, reserved1, reserved2, reserved3, lastResult)

	p.logger.Debug(str)

	return mode
}
