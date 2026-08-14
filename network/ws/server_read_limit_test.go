/**
 * @Desc: 单条消息读取上限（SetReadLimit）回归测试
 */

package ws_test

import (
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/dobyte/due/network/ws/v2"
	"github.com/dobyte/due/v2/network"
	"github.com/gorilla/websocket"
)

func freeTCPPort(t *testing.T) int {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer ln.Close()
	return ln.Addr().(*net.TCPAddr).Port
}

func TestServer_MaxMessageBytes(t *testing.T) {
	port := freeTCPPort(t)

var received int

server := ws.NewServer(
		ws.WithServerAddr("127.0.0.1:"+strconv.Itoa(port)),
		ws.WithServerPath("/read-limit-test"),
		ws.WithServerMaxMessageBytes(1024),
	)
	server.OnReceive(func(conn network.Conn, data []byte) {
		received++
	})

	if err := server.Start(); err != nil {
		t.Fatalf("server start: %v", err)
	}
	defer server.Stop()

	// 等待端口就绪
	deadline := time.Now().Add(3 * time.Second)
	var conn *websocket.Conn
	for time.Now().Before(deadline) {
		var err error
		conn, _, err = (&websocket.Dialer{}).Dial("ws://127.0.0.1:"+strconv.Itoa(port)+"/read-limit-test", nil)
		if err == nil {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
	if conn == nil {
		t.Fatal("dial failed")
	}
	defer conn.Close()

	// 发送超过 1024 字节的消息，服务端必须以 1009 关闭连接且不投递
	_ = conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	if err := conn.WriteMessage(websocket.BinaryMessage, make([]byte, 64*1024)); err != nil {
		t.Fatalf("write oversized message: %v", err)
	}

	_, _, err := conn.ReadMessage()
	if err == nil {
		t.Fatal("expected read error after oversized message, got nil")
	}
	closeErr, ok := err.(*websocket.CloseError)
	if !ok || closeErr.Code != websocket.CloseMessageTooBig {
		t.Fatalf("expected close 1009 (message too big), got %v", err)
	}

	if received != 0 {
		t.Fatalf("oversized message should never be delivered, delivered=%d", received)
	}
}
