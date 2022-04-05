package main

import (
	"log"
	"time"

	"github.com/rkumar0099/algorand/api"
)

func main() {

	//contract.Testing_Price_feed()
	//contract.Testing_FIN_Feed()
	//contract.Testing_Flight_Feed()
	TestCreate()
}

func TestCreate() {
	username := "Rabu"
	password := "abc"
	log.Println("[Debug] [contracts] Account created")
	//for i := 0; i < 20; i++ {
	//	go api.New().CreateAccount(username, password)
	//}
	//time.Sleep(1 * time.Minute)
	a := api.New()
	for j := 0; j < 5; j++ {
		time.Sleep(5 * time.Second)

		for i := 0; i < 10; i++ {
			go a.CreateAccount(username, password)
		}
	}
	time.Sleep(1 * time.Minute)
	//a.CreateAccount(username, password)
	//api.New().CreateAccount(username, password)
	//time.Sleep(20 * time.Second)

}
