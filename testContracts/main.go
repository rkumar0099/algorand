package main

import (
	"github.com/rkumar0099/algorand/tests"
)

func main() {
	/*
		username := "Rabu"
		password := "abc"
		pk, sk := tests.TestCreate(username, password)
		log.Println(pk, sk)
		tests.TestLogin(username, password, pk)
		tests.TestTopup()
		tests.TestTransfer()
	*/
	tests.TestPriceFeed()
}
