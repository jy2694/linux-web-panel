package handlers

import (
	"net/http"
	"os"
	"os/exec"
	"sync"

	"github.com/creack/pty"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 개발 환경에서는 모든 origin 허용
	},
}

type ShellHandler struct {
	connections sync.Map
}

func NewShellHandler() *ShellHandler {
	return &ShellHandler{}
}

func (h *ShellHandler) HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	// pty 생성
	cmd := exec.Command("/bin/bash")
	ptmx, err := pty.Start(cmd)
	if err != nil {
		return
	}
	defer ptmx.Close()

	// 연결 ID 생성 및 저장
	connID := c.GetString("user_id")
	h.connections.Store(connID, ptmx)
	defer h.connections.Delete(connID)

	// pty -> WebSocket
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := ptmx.Read(buf)
			if err != nil {
				return
			}
			if err := conn.WriteMessage(websocket.TextMessage, buf[:n]); err != nil {
				return
			}
		}
	}()

	// WebSocket -> pty
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			return
		}
		if _, err := ptmx.Write(message); err != nil {
			return
		}
	}
}

func (h *ShellHandler) CloseConnection(userID string) {
	if ptmx, ok := h.connections.Load(userID); ok {
		if f, ok := ptmx.(*os.File); ok {
			f.Close()
		}
		h.connections.Delete(userID)
	}
}
