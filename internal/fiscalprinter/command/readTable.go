package command

import (
	"errors"
	"fmt"
)

type ReadTable struct {
	*PrinterCommand
	tableNumber int
	rowNumber   int
	fieldNumber int
}

func NewReadTable() *ReadTable {
	rt := &ReadTable{
		NewPrinterCommand(),
		0, 0, 0,
	}

	rt.text = "Get table field value"
	rt.commandCode = 31
	return rt
}

func (rt *ReadTable) ReadTable(tableNumber, rowNumber, fieldNumber int) error {
	if err := rt.checkRange(tableNumber, 0, 255, "table number"); err != nil {
		return err
	}
	if err := rt.checkRange(rowNumber, 0, 255, "table number"); err != nil {
		return err
	}
	if err := rt.checkRange(fieldNumber, 0, 255, "table number"); err != nil {
		return err
	}

	rt.tableNumber = tableNumber
	rt.rowNumber = rowNumber
	rt.fieldNumber = fieldNumber

	return nil
}

func (rt *ReadTable) checkRange(value, min, max int, name string) error {
	if value < min {
		return errors.New(fmt.Sprintln(name, ": invalid parameter value (", value, " < ", min, ")"))
	}
	if value > max {
		return errors.New(fmt.Sprintln(name, ": invalid parameter value (", value, " > ", max, ")"))
	}

	return nil
}
