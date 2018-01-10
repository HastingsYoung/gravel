package models

// Acknowledgement: https://github.com/jupp0r/go-priority-queue/blob/master/priorty_queue.go

import (
	"container/heap"
	"errors"
)

const (
	ORDER_BOOK_TYPE_ASK = "ASK"
	ORDER_BOOK_TYPE_BID = "BID"
)

func NewBook() *OrderBook {
	return &OrderBook{
		queues: map[string]OrderQueue{},
	}
}

type OrderBook struct {
	queues map[string]OrderQueue
}

func (ob *OrderBook) Sum() *Summary {
	return &Summary{
		Queues: ob.queues,
	}
}

func (ob *OrderBook) List(key string, queue OrderQueue) {
	ob.queues[key] = queue
}

type Summary struct {
	Queues map[string]OrderQueue `json:"queues"`
}

type OrderQueue interface {
	Init()
	Add(o *Order)
	Update(id string, n *Order) error
	Peek(i int) *Order
	Len() int
	Next() *Order
	IsEmpty() bool
}

func NewQueueAsk() *OrderQueueAsk {
	return &OrderQueueAsk{
		items:  []*Order{},
		lookup: map[interface{}]*Order{},
	}
}

type OrderQueueAsk struct {
	items  []*Order
	lookup map[interface{}]*Order // two-way mappings
}

// Interface method for heap
func (ask *OrderQueueAsk) Len() int {
	return len(ask.items)
}

// Interface method for heap
func (ask *OrderQueueAsk) Less(i, j int) bool {
	return ask.items[i].Price < ask.items[j].Price
}

// Interface method for heap
func (ask *OrderQueueAsk) Swap(i, j int) {
	ask.items[i], ask.items[j] = ask.items[j], ask.items[i]
	ask.items[i].Index = i
	ask.items[j].Index = j
}

// Interface method for heap
func (ask *OrderQueueAsk) Push(x interface{}) {
	order := x.(*Order)
	order.Index = len(ask.items)
	ask.items = append(ask.items, order)
}

// Interface method for heap
func (ask *OrderQueueAsk) Pop() interface{} {
	n := len(ask.items)
	order := ask.items[n-1]
	order.Index = -1
	ask.items = ask.items[0 : n-1]
	return order
}

// Initialise the order queue
func (ask *OrderQueueAsk) Init() {
	heap.Init(ask)
}

// Add a new order in the orderbook
func (ask *OrderQueueAsk) Add(o *Order) {
	if _, ok := ask.lookup[o.OrderId]; ok {
		return
	}

	heap.Push(ask, o)
	ask.lookup[o.OrderId] = o
}

func (ask *OrderQueueAsk) Next() *Order {
	if ask.IsEmpty() {
		return nil
	}

	order := heap.Pop(ask).(*Order)
	delete(ask.lookup, order.OrderId)
	return order
}

func (ask *OrderQueueAsk) Update(id string, n *Order) error {
	if _, ok := ask.lookup[id]; !ok {
		return errors.New("order does not exist")
	}

	index := ask.lookup[id].Index

	*(ask.lookup[id]) = *n
	heap.Fix(ask, index)
	return nil
}

func (ask *OrderQueueAsk) Peek(i int) *Order {
	if ok := ask.items[i]; ok == nil {
		return ok
	}

	return ask.items[i]
}

func (ask *OrderQueueAsk) IsEmpty() bool {
	return len(ask.items) == 0
}

func NewQueueBid() *OrderQueueBid {
	return &OrderQueueBid{
		items:  []*Order{},
		lookup: map[interface{}]*Order{},
	}
}

type OrderQueueBid struct {
	items  []*Order
	lookup map[interface{}]*Order // two-way mappings
}

// Interface method for heap
func (bid *OrderQueueBid) Len() int {
	return len(bid.items)
}

// Interface method for heap
func (bid *OrderQueueBid) Less(i, j int) bool {
	return bid.items[i].Price > bid.items[j].Price
}

// Interface method for heap
func (bid *OrderQueueBid) Swap(i, j int) {
	bid.items[i], bid.items[j] = bid.items[j], bid.items[i]
	bid.items[i].Index = i
	bid.items[j].Index = j
}

// Interface method for heap
func (bid *OrderQueueBid) Push(x interface{}) {
	order := x.(*Order)
	order.Index = len(bid.items)
	bid.items = append(bid.items, order)
}

// Interface method for heap
func (bid *OrderQueueBid) Pop() interface{} {
	n := len(bid.items)
	order := bid.items[n-1]
	order.Index = -1
	bid.items = bid.items[0 : n-1]
	return order
}

// Initialise the order queue
func (bid *OrderQueueBid) Init() {
	heap.Init(bid)
}

// Add a new order in the orderbook
func (bid *OrderQueueBid) Add(o *Order) {
	if _, ok := bid.lookup[o.OrderId]; ok {
		return
	}

	heap.Push(bid, o)
	bid.lookup[o.OrderId] = o
}

func (bid *OrderQueueBid) Next() *Order {
	if bid.IsEmpty() {
		return nil
	}

	order := heap.Pop(bid).(*Order)
	delete(bid.lookup, order.OrderId)
	return order
}

func (bid *OrderQueueBid) Update(id string, n *Order) error {
	if _, ok := bid.lookup[id]; !ok {
		return errors.New("order does not exist")
	}

	index := bid.lookup[id].Index

	*(bid.lookup[id]) = *n
	heap.Fix(bid, index)
	return nil
}

func (bid *OrderQueueBid) Peek(i int) *Order {
	if ok := bid.items[i]; ok == nil {
		return ok
	}

	return bid.items[i]
}

func (bid *OrderQueueBid) IsEmpty() bool {
	return len(bid.items) == 0
}
