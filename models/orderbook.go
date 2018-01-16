package models

// Acknowledgement: https://github.com/jupp0r/go-priority-queue/blob/master/priorty_queue.go

import (
	"container/heap"
	"errors"
	// "fmt"
	"github.com/gravel/math"
	"sync"
)

// Important: unlike the operations under OrderQueue,
// OrderBook struct is thread unsafe, please use Exchange
// to handle higher-level concurrencies
func NewBook() *OrderBook {
	return &OrderBook{
		queues: map[string]OrderQueue{},
		Deals:  make(chan *Deal),
	}
}

// An orderbook lists the trading options for a specified stock
type OrderBook struct {
	queues    map[string]OrderQueue
	histories []*Deal
	Deals     chan *Deal
	sync.Mutex
}

func (ob *OrderBook) Sum() *Summary {
	length := len(ob.histories)

	if length == 0 {
		return &Summary{
			Queues:    ob.queues,
			Histories: []*Deal{},
		}
	}

	return &Summary{
		Queues:    ob.queues,
		Histories: ob.histories[math.MaxInt(length-100, 0):],
	}
}

func (ob *OrderBook) Update() {
	select {
	case deal := <-ob.Deals:
		// fmt.Println("Price:", deal.Price, "Amount:", deal.Amount, "Timestamp:", deal.Timestamp, "Total:", deal.Total)
		ob.Lock()
		ob.histories = append(ob.histories, deal)
		ob.Unlock()
	default:
	}
}

func (ob *OrderBook) SetQueue(key string, queue OrderQueue) {
	ob.queues[key] = queue
}

func (ob *OrderBook) GetQueue(key string) OrderQueue {
	return ob.queues[key]
}

func (ob *OrderBook) Queues() *map[string]OrderQueue {
	return &ob.queues
}

type Summary struct {
	Queues    map[string]OrderQueue `json:"queues"`
	Histories []*Deal               `json:"histories"`
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
		Items:  []*Order{},
		lookup: map[interface{}]*Order{},
	}
}

type OrderQueueAsk struct {
	Items  []*Order               `json:"items"`
	lookup map[interface{}]*Order // two-way mappings
	sync.RWMutex
}

// Interface method for heap
func (ask *OrderQueueAsk) Len() int {
	return len(ask.Items)
}

// Interface method for heap
func (ask *OrderQueueAsk) Less(i, j int) bool {
	return ask.Items[i].Price < ask.Items[j].Price
}

// Interface method for heap
func (ask *OrderQueueAsk) Swap(i, j int) {
	ask.Items[i], ask.Items[j] = ask.Items[j], ask.Items[i]
	ask.Items[i].Index = i
	ask.Items[j].Index = j
}

// Interface method for heap
func (ask *OrderQueueAsk) Push(x interface{}) {
	order := x.(*Order)
	order.Index = ask.Len()
	ask.Items = append(ask.Items, order)
}

// Interface method for heap
func (ask *OrderQueueAsk) Pop() interface{} {
	n := ask.Len()
	order := ask.Items[n-1]
	order.Index = -1
	ask.Items = ask.Items[0 : n-1]
	return order
}

// Initialise the order queue
func (ask *OrderQueueAsk) Init() {
	ask.Lock()
	defer ask.Unlock()
	heap.Init(ask)
}

// Add a new order in the orderbook
func (ask *OrderQueueAsk) Add(o *Order) {
	ask.Lock()
	defer ask.Unlock()

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

	ask.Lock()
	defer ask.Unlock()

	order := heap.Pop(ask).(*Order)
	delete(ask.lookup, order.OrderId)
	return order
}

func (ask *OrderQueueAsk) Update(id string, n *Order) error {
	if _, ok := ask.lookup[id]; !ok {
		return errors.New("Order does not exist")
	}

	ask.Lock()
	defer ask.Unlock()

	index := ask.lookup[id].Index

	*(ask.lookup[id]) = *n
	heap.Fix(ask, index)
	return nil
}

// debug: concurrent access
func (ask *OrderQueueAsk) Peek(i int) *Order {

	var (
		l = ask.Len()
	)

	if i >= l {
		return nil
	}

	return ask.Items[i]
}

func (ask *OrderQueueAsk) IsEmpty() bool {
	return ask.Len() == 0
}

func NewQueueBid() *OrderQueueBid {
	return &OrderQueueBid{
		Items:  []*Order{},
		lookup: map[interface{}]*Order{},
	}
}

type OrderQueueBid struct {
	Items  []*Order               `json:"items"`
	lookup map[interface{}]*Order // two-way mappings
	sync.RWMutex
}

// Interface method for heap
func (bid *OrderQueueBid) Len() int {
	return len(bid.Items)
}

// Interface method for heap
func (bid *OrderQueueBid) Less(i, j int) bool {
	return bid.Items[i].Price > bid.Items[j].Price
}

// Interface method for heap
func (bid *OrderQueueBid) Swap(i, j int) {
	bid.Items[i], bid.Items[j] = bid.Items[j], bid.Items[i]
	bid.Items[i].Index = i
	bid.Items[j].Index = j
}

// Interface method for heap
func (bid *OrderQueueBid) Push(x interface{}) {
	order := x.(*Order)
	order.Index = bid.Len()
	bid.Items = append(bid.Items, order)
}

// Interface method for heap
func (bid *OrderQueueBid) Pop() interface{} {
	n := bid.Len()
	order := bid.Items[n-1]
	order.Index = -1
	bid.Items = bid.Items[0 : n-1]
	return order
}

// Initialise the order queue
func (bid *OrderQueueBid) Init() {
	bid.Lock()
	defer bid.Unlock()
	heap.Init(bid)
}

// Add a new order in the orderbook
func (bid *OrderQueueBid) Add(o *Order) {
	bid.Lock()
	defer bid.Unlock()

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

	bid.Lock()
	defer bid.Unlock()

	order := heap.Pop(bid).(*Order)
	delete(bid.lookup, order.OrderId)
	return order
}

func (bid *OrderQueueBid) Update(id string, n *Order) error {
	if _, ok := bid.lookup[id]; !ok {
		return errors.New("Order does not exist")
	}

	bid.RLock()
	defer bid.RUnlock()

	index := bid.lookup[id].Index

	*(bid.lookup[id]) = *n
	heap.Fix(bid, index)
	return nil
}

func (bid *OrderQueueBid) Peek(i int) *Order {

	if i >= bid.Len() {
		return nil
	}

	return bid.Items[i]
}

func (bid *OrderQueueBid) IsEmpty() bool {
	return bid.Len() == 0
}
