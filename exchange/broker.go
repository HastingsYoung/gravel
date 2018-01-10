package exchange

import (
	. "github.com/gravel/models"
	"github.com/satori/go.uuid"
	"math"
)

// A broker will match the orders listed in an exchange
type Broker struct {
	BrokerId string
	Ask      OrderQueue
	Bid      OrderQueue
	Deals    chan *Deal
	idle     bool
	exit     chan bool
}

func NewBroker() *Broker {
	return &Broker{
		BrokerId: uuid.NewV4().String(),
		idle:     true,
	}
}

func (b *Broker) Watch(ask OrderQueue, bid OrderQueue, deals chan *Deal) {
	b.Ask = ask
	b.Bid = bid
	b.Deals = deals
}

func (b *Broker) Start() {

	// no watched queue exist or broker already start
	if b.Ask == nil || b.Bid == nil || !b.IsIdle() {
		return
	}

	b.idle = false

	// start looping
	go func() {
		for {

			select {
			case <-b.exit:
				b.idle = true
				return
			default:

				if oask, obid := b.Ask.Peek(0), b.Bid.Peek(0); oask != nil && obid != nil {
					if oask.Price < obid.Price {
						if oask.Amount == obid.Amount {
							b.Deals <- Match(b.Ask.Next(), b.Bid.Next())
						} else if oask.Amount < obid.Amount {
							b.Deals <- Match(b.Ask.Next(), obid)
						} else {
							b.Deals <- Match(oask, b.Bid.Next())
						}
					}
				}
			}
		}
	}()
}

func (b *Broker) Stop() {
	b.Ask = nil
	b.Bid = nil
	b.exit <- true
}

func (b *Broker) IsIdle() bool {
	return b.idle
}

func Match(ask *Order, bid *Order) *Deal {
	var (
		amount = math.Min(ask.Amount, bid.Amount)
	)

	ask.Amount -= amount
	bid.Amount -= amount

	return NewDeal(amount, ask.Price)
}
