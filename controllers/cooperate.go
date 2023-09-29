package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	model "github.com/Dparty/model/restaurant"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type cooperateMessage struct {
	Orders []model.Order
}

func (c cooperateMessage) ToJson() []byte {
	s, _ := json.Marshal(c)
	return s
}

var tables = make(map[string][]model.Order)

func cooperate(c *gin.Context) {
	tableId := c.Param("id")
	var table model.Table
	ctx := db.Find(&table, tableId)
	if ctx.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, "")
		return
	}
	fmt.Println(table)
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	go func() {
		time.Sleep(time.Second * 10)
		err := conn.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()
	go func() {
		for {
			mt, message, err := conn.ReadMessage()
			if err != nil {
				fmt.Println(err)
				conn.Close()
				break
			}
			fmt.Println(mt, string(message), err)
		}
	}()
	go func() {
		for {
			if conn == nil {
				break
			}
			conn.WriteMessage(websocket.TextMessage, []byte("Hello, WebSocket!"))
			time.Sleep(time.Second)
		}
	}()
}
