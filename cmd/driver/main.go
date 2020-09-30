package main

import (
	"log"
	"shtrih-drv/internal/fiscalprinter"
)

func main() {
	log.Println("golang shtrih")

	//_, err := net.Dial("tcp", "10.51.0.71:7778")
	//if err != nil {
	//	log.Println(err.Error())
	//}

	fp := fiscalprinter.NewFiscalPrinter()

	log.Println(fp.GetSerialNumber())
}
