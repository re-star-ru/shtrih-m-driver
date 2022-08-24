package configs

import "time"

const DefaultTimeout = time.Second * 120

type ConfKKT map[string]Ck

type Ck struct {
	Addr string // ip address for kkt
	Inn  string // default inn for fiscal operations
}
