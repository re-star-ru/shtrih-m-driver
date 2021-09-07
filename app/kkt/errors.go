package kkt

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

func errCheck(e byte) error {
	err := fmt.Errorf("error in receive message: %X", e)
	log.Err(err).Send()

	return err
}
