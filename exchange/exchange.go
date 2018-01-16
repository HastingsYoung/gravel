package exchange

import (
	"errors"
	. "github.com/gravel/models"
)

type Exchange struct {
	// brokers
	pool map[string]*Broker
	// stocks
	stocks map[string]*Stock
	// orderbooks
	books map[string]*OrderBook
	exit  chan bool
}

func NewExchange() *Exchange {
	return &Exchange{
		pool:   map[string]*Broker{},
		stocks: map[string]*Stock{},
		books:  map[string]*OrderBook{},
		exit:   make(chan bool),
	}
}

func (ex *Exchange) Register(b *Broker) {
	ex.pool[b.BrokerId] = b
}

func (ex *Exchange) DeRegister(id string) {
	if b, ok := ex.pool[id]; ok {
		b.Stop()
		delete(ex.pool, id)
	}
}

func (ex *Exchange) Broadcast() *Message {

	var (
		payload []*Summary
	)

	for _, book := range ex.books {
		payload = append(payload, book.Sum())
	}

	return &Message{
		Command:   MESSAGE_COMMAND_SUMMARY,
		Summaries: payload,
	}
}

func (ex *Exchange) Buy(
	code, market string,
	price, amount float64,
) error {
	if book, ok := ex.books[code]; ok {
		if queue := book.GetQueue(market); queue != nil {
			queue.Add(
				NewOrder(
					market,
					ORDER_TYPE_BID,
					code,
					price,
					amount,
				),
			)
			return nil
		}
		return errors.New("Market not exist")
	}

	return errors.New("Stock code not exist")
}

func (ex *Exchange) Sell(
	code, market string,
	price, amount float64,
) error {
	if book, ok := ex.books[code]; ok {
		if queue := book.GetQueue(market); queue != nil {
			queue.Add(
				NewOrder(
					market,
					ORDER_TYPE_ASK,
					code,
					price,
					amount,
				),
			)
			return nil
		}
		return errors.New("Market not exist")
	}

	return errors.New("Stock code not exist")
}

func (ex *Exchange) Issue(s *Stock, num ...int) error {

	var (
		count = -1
		max   = 0
	)

	if len(num) > 0 {
		max = num[0]
	}

	ex.stocks[s.Code] = s
	ex.books[s.Code] = NewBook()
	// todo: add flexibility to  book type
	// e.g. BTC_ASK, BTC_BID, ETH_ASK, ETH_BID

	ask := NewQueueAsk()
	bid := NewQueueBid()

	ex.books[s.Code].SetQueue("ASK", ask)
	ex.books[s.Code].SetQueue("BID", bid)

	for _, b := range ex.pool {
		if b.IsIdle() {
			b.Watch(ask, bid, ex.books[s.Code].Deals)
			b.Start()

			count++

			if count == max {
				break
			}
		}
	}

	if count < 0 {
		return errors.New("No broker available at the moment, please re-try after a while")
	}

	return nil
}

func (ex *Exchange) Stop() {
	ex.exit <- true
}

func (ex *Exchange) Start() {
	for {
		select {
		case <-ex.exit:
			for _, b := range ex.pool {
				b.Stop()
			}
			return
		default:
			for _, b := range ex.books {
				b.Update()
			}
		}
	}
}
