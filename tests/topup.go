package tests

import (
	"log"
)

func TestTopup() {
	status, msg := a.TopUp(100)
	log.Println(status, msg)
}
