package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	api "github.com/Dparty/core-api"
	model "github.com/Dparty/model/restaurant"
	"github.com/cskr/pubsub"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var mu sync.Mutex

var ps *pubsub.PubSub = pubsub.New(0)

var tables map[uint][]api.Specification = make(map[uint][]api.Specification)

type Message struct {
	Action        string             `json:"action"`
	Specification *api.Specification `json:"specification"`
}

func GetOrders(id uint) []api.Specification {
	mu.Lock()
	orders, ok := tables[id]
	if !ok {
		orders = make([]api.Specification, 0)
		tables[id] = orders
	}
	mu.Unlock()
	return orders
}

func PushOrders(id uint, order api.Specification) []api.Specification {
	mu.Lock()
	orders, ok := tables[id]
	if !ok {
		orders = make([]api.Specification, 0)
		tables[id] = orders
	}
	orders = append(orders, order)
	tables[id] = orders
	mu.Unlock()
	return orders
}

func cooperate(c *gin.Context) {
	tableId := c.Param("id")
	var table model.Table
	ctx := db.Find(&table, tableId)
	if ctx.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, "")
		return
	}
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	submit := make(chan struct{})
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				fmt.Println(err)
				conn.Close()
				break
			}
			var m Message
			json.Unmarshal(message, &m)
			var orders []api.Specification
			switch m.Action {
			case "ADD":
				orders = PushOrders(table.ID, *m.Specification)
			case "SUBMIT":
				submit <- struct{}{}
				conn.Close()
				return
			}
			ps.Pub(orders, tableId)
		}
	}()
	go func() {
		j, _ := json.Marshal(GetOrders(table.ID))
		conn.WriteMessage(websocket.TextMessage, j)
		sub := ps.Sub(tableId)
		for {
			select {
			case orders := (<-sub):
				j, _ := json.Marshal(orders.([]api.Specification))
				conn.WriteMessage(websocket.TextMessage, j)
			case <-submit:
				return
			}
		}
	}()
}
