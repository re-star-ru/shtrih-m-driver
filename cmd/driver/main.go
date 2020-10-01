package main

import (
	"log"
	"shtrih-drv/internal/fiscalprinter"
)

func main() {
	log.Println("golang shtrih")

	printer := fiscalprinter.NewPrinterProtocol()
	err := printer.Connect()
	if err != nil {
		log.Fatal(err)
	}
}
