package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	. "github.com/gravel/app"
	. "github.com/gravel/exchange"
	. "github.com/gravel/models"
	"net/http"
	"os"
	"time"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	exchange = NewExchange()
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

func main() {
	const (
		NAME        = "Stock"
		CODE        = "STK"
		DESCRIPTION = "This is an example of using gravel to trade stock"
		TOTAL       = 100000
		CIRCULATING = 90000
		REF         = "/stocks/stk"
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
		hub = newHub()
	)

	for i := 0; i < 5; i++ {
		exchange.Register(NewBroker())
	}

	go exchange.Start()

	exchange.Issue(stock)

	go hub.run()

	http.HandleFunc("/websocket", func(w http.ResponseWriter, r *http.Request) {
		serve(hub, w, r)
	})

	http.HandleFunc("/", Index)
	http.Handle(
		"/static/",
		http.StripPrefix("/static/",
			http.FileServer(
				http.Dir(os.Getenv("GOPATH")+"/src/github.com/gravel/static/"),
			),
		),
	)

	panic(http.ListenAndServe("localhost:8080", nil))
}

type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan *Message

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan *Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {

	go func() {
		for {
			<-time.After(5 * time.Second)
			h.broadcast <- exchange.Broadcast()
		}
	}()

	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			fmt.Println("Hub", "register")
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			fmt.Println("Hub", "unregister")
		case message := <-h.broadcast:
			fmt.Println("Hub", "broadcast")
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan *Message
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		var (
			message Message
		)

		message.Order = &Order{}

		fmt.Println("Hub", "waiting for incoming connection")

		if err := c.conn.ReadJSON(&message); err != nil {
			fmt.Println("Hub", err)
			break
		}

		switch message.GetCommand() {
		case MESSAGE_COMMAND_OPEN:
			continue
		case MESSAGE_COMMAND_CLOSE:
			return
		case MESSAGE_COMMAND_BUY:
			if err := exchange.Buy(
				message.Order.StockCode,
				message.Order.Market,
				message.Order.Price,
				message.Order.Amount,
			); err != nil {
				panic(err)
			}
		case MESSAGE_COMMAND_SELL:
			if err := exchange.Sell(
				message.Order.StockCode,
				message.Order.Market,
				message.Order.Price,
				message.Order.Amount,
			); err != nil {
				panic(err)
			}
		case MESSAGE_COMMAND_NEW_STOCK:
			continue
		default:
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteJSON(NewErrorMessage("Connection disrupted"))
				return
			}

			n := len(c.send)

			c.conn.WriteJSON(message)

			for i := 0; i < n; i++ {
				c.conn.WriteJSON(<-c.send)
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func serve(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan *Message)}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
