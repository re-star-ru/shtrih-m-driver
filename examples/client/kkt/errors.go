package kkt

import (
	"errors"
	"fmt"
	"log"
)

func errCheck(e byte) error {
	err := fmt.Errorf("error in receive message: %v", e)
	log.Println(err)

	return err
}

var errChecksum = errors.New("checksum does not match")
