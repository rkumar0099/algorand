package tests

import (
	"log"
)

func TestTopup() {
	status, msg := a1.TopUp(100)
	log.Println(status, msg)
}
