package oracle

const (
	CURRENCY_EXCHANGE = iota
	FLIGHT
)

// fetch data from given url, and send rate at time t to end client.

type CurrencyExchange struct {
	Time string  `json:"time"`
	Rate float64 `json:"rate"`
}

type FlightDetails struct {
	Base    string `json:"base"`
	Price   string `json:"price"`
	Time    string `json:"time"`
	Details string `json:"details"`
}

// clients can only send external world transactions related to the type of access reponses defined here
// because size of data is large. We consider only important details of access response to have small size data

/*
type Response struct {

}


type  struct {
	Time string `json:"time"`
	Rate string `json:"rate"`
}


var URL map[string]
*/
