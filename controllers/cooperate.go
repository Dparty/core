package controllers

import (
	"net/http"

	model "github.com/Dparty/model/restaurant"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

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
	go func() {
		for {

		}
	}()
	go func() {
		for {
			conn.WriteMessage(websocket.TextMessage, []byte("Hello, WebSocket!"))
		}
	}()
}
