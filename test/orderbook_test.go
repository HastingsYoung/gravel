package test

import (
	. "github.com/gravel/models"
	"math/rand"
	"testing"
)

func TestOrderList(t *testing.T) {
	stock := NewStock(
		"Test_Stock_Name",
		"Test_Code",
		"Test_Description",
		100000,
		90000,
		"/Test_Link",
	)
	r := rand.New(rand.NewSource(99))

	book := NewBook()
	ask := NewQueueAsk()
	bid := NewQueueBid()
	book.List("Ask", ask)
	book.List("Bid", bid)

	for i := 0; i < 10; i++ {
		ask.Add(
			NewOrder(
				"Test_Market",
				ORDER_TYPE_ASK,
				stock,
				r.Float64()*100,
				r.Float64()*100,
			),
		)
		bid.Add(
			NewOrder(
				"Test_Market",
				ORDER_TYPE_BID,
				stock,
				r.Float64()*100,
				r.Float64()*100,
			),
		)
	}

	v1 := book.Sum().Queues["Ask"]
	v2 := book.Sum().Queues["Bid"]

	var max float64 = -1.00
	for next := v1.Next(); next != nil; next = v1.Next() {
		if max > next.Price {
			t.Error(
				"Queue not sorted",
				"Expected",
				max,
				"greater than",
				next.Price,
			)
		}
		max = next.Price
	}

	var min float64 = 100.00
	for next := v2.Next(); next != nil; next = v2.Next() {
		if min < next.Price {
			t.Error(
				"Queue not sorted",
				"Expected",
				min,
				"less than",
				next.Price,
			)
		}
		min = next.Price
	}
}
