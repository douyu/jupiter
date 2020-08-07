package xgin

import (
	"io"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

//WebSocketConn websocket conn, see websocket.Conn
type WebSocketConn interface {
	Subprotocol() string
	Close() error
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	WriteControl(messageType int, data []byte, deadline time.Time) error
	NextWriter(messageType int) (io.WriteCloser, error)
	WritePreparedMessage(pm *websocket.PreparedMessage) error
	WriteMessage(messageType int, data []byte) error
	SetWriteDeadline(t time.Time) error
	NextReader() (messageType int, r io.Reader, err error)
	ReadMessage() (messageType int, p []byte, err error)
	SetReadDeadline(t time.Time) error
	SetReadLimit(limit int64)
	CloseHandler() func(code int, text string) error
	SetCloseHandler(h func(code int, text string) error)
	PingHandler() func(appData string) error
	SetPingHandler(h func(appData string) error)
	PongHandler() func(appData string) error
	SetPongHandler(h func(appData string) error)
	UnderlyingConn() net.Conn
	EnableWriteCompression(enable bool)
	SetCompressionLevel(level int) error
}

//WebSocketFunc ..
type WebSocketFunc func(WebSocketConn, error)

//WebSocket ..
type WebSocket struct {
	Pattern string
	Handler WebSocketFunc
	*websocket.Upgrader
	Header http.Header
}

//Upgrade get upgrage request
func (ws *WebSocket) Upgrade(w http.ResponseWriter, r *http.Request) {
	conn, err := ws.Upgrader.Upgrade(w, r, ws.Header)
	if err == nil {
		defer conn.Close()
	}
	ws.Handler(conn, err)
}

//WebSocketOption ..
type WebSocketOption func(*WebSocket)

//WebSocketOptions ..
func WebSocketOptions(pattern string, handler WebSocketFunc, opts ...WebSocketOption) *WebSocket {
	ws := &WebSocket{
		Pattern:  pattern,
		Handler:  handler,
		Upgrader: &websocket.Upgrader{},
	}
	for _, opt := range opts {
		opt(ws)
	}
	return ws
}
