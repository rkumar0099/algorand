package tests

import (
	"log"
)

func TestTransfer() {
	username := "Rabi"
	password := "abcd"
	//a2 := api.New("127.0.0.1:9021")
	pk2, _, msg := a2.CreateAccount(username, password)
	log.Println(pk2, msg)

	status, msg := a2.LogIn(username, password, pk2)
	log.Println(status, msg)

	status, msg = a1.Transfer(pk2, 150)
	log.Println(status, msg)
}
