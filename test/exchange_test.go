package test

import (
	. "github.com/gravel/exchange"
	. "github.com/gravel/models"
	"math/rand"
	"testing"
	"time"
)

func TestExchange(t *testing.T) {
	const (
		NAME        = "Test_Stock_Name"
		CODE        = "Test_Code"
		DESCRIPTION = "Test_Description"
		TOTAL       = 100000
		CIRCULATING = 90000
		REF         = "/Test_Link"
		MARKET_ASK  = "ASK"
		MARKET_BID  = "BID"
	)

	var (
		stock = NewStock(
			NAME,
			CODE,
			DESCRIPTION,
			TOTAL,
			CIRCULATING,
			REF,
		)
		exchange = NewExchange()
		r        = rand.New(rand.NewSource(99))
	)

	for i := 0; i < 10; i++ {
		exchange.Register(NewBroker())
	}

	go exchange.Start()

	exchange.Issue(stock)

	go func() {
		for {
			err := exchange.Buy(
				CODE,
				MARKET_BID,
				15,
				15,
			)
			if err != nil {
				panic(err)
			}
		}
	}()

	go func() {
		for {
			err := exchange.Sell(
				CODE,
				MARKET_ASK,
				10+r.Float64()*10,
				10+r.Float64()*10,
			)
			if err != nil {
				panic(err)
			}
		}
	}()

	<-time.After(3 * time.Second)

	// log histories
	for _, his := range exchange.Broadcast().Summaries[0].Histories {
		t.Log(*his)
	}

	exchange.Stop()
}
