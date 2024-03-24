package apiServer

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log/slog"
	"net/http"
)

func socketHandler(c *gin.Context) {
	upGrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		slog.Error("Failed to upgrade to websocket", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": ""})
		return
	}
	defer ws.Close()

	for {
		msgType, msg, err := ws.ReadMessage()
		if err != nil {
			slog.Error("Failed to read message from websocket", "err", err)
			return
		}
		go handleWSMessage(ws, msgType, msg)
	}
}

func handleWSMessage(ws *websocket.Conn, msgType int, msg []byte) {
	// TODO: Design the message format
}

func createRlpaServer() {

}
