package kkt

import (
	"fmt"
	"log"
)

func errCheck(e byte) error {
	err := fmt.Errorf("error in receive message: %X", e)
	log.Println(err)

	return err
}
