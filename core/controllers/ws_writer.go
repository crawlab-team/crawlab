package controllers

import (
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab/trace"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"io"
	http2 "net/http"
)

type WsWriter struct {
	io.Writer
	io.Closer
	conn *websocket.Conn
}

func (w *WsWriter) Write(data []byte) (n int, err error) {
	log.Infof("websocket write: %s", string(data))
	err = w.conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		return 0, err
	}
	return len(data), nil
}

func (w *WsWriter) Close() (err error) {
	return w.conn.Close()
}

func (w *WsWriter) CloseWithText(text string) {
	_ = w.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, text))
}

func (w *WsWriter) CloseWithError(err error) {
	_ = w.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, err.Error()))
}

func NewWsWriter(c *gin.Context) (writer *WsWriter, err error) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http2.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Errorf("websocket open connection error: %v", err)
		trace.PrintError(err)
	}

	return &WsWriter{
		conn: conn,
	}, nil
}
