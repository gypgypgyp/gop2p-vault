package p2p

import (
	// "fmt"
	"io"
	"log"
	"net"
	"sync"
)

type TCPTransport struct {
	address     string
	listener    net.Listener
	connections map[string]net.Conn
	lock        sync.RWMutex
	handler     HandlerFunc
}

func NewTCPTransport(addr string) *TCPTransport {
	return &TCPTransport{
		address:     addr, // ✅ 加上这个！
		connections: make(map[string]net.Conn),
	}
}

func (t *TCPTransport) ListenAndAccept() error {
	l, err := net.Listen("tcp", t.address)
	if err != nil {
		return err
	}
	t.listener = l
	log.Printf("[P2P] Listening on %s\n", t.address)

	for {
		conn, err := t.listener.Accept()
		if err != nil {
			return err
		}
		peerID := conn.RemoteAddr().String()
		t.lock.Lock()
		t.connections[peerID] = conn
		t.lock.Unlock()

		go t.handleConnection(conn, peerID)
	}
}

func (t *TCPTransport) handleConnection(conn net.Conn, peerID string) {
	defer conn.Close()

	buf := make([]byte, 4096)
	for {
		n, err := conn.Read(buf)
		if err == io.EOF {
			log.Printf("[P2P] Connection closed by %s\n", peerID)
			return
		} else if err != nil {
			log.Printf("[P2P] Read error from %s: %v\n", peerID, err)
			return
		}

		if t.handler != nil {
			// t.handler(peerID, buf[:n])
			msg, err := Decode[Message](buf[:n])
			if err != nil {
				log.Printf("[P2P] Failed to decode message from %s: %v\n", peerID, err)
				return
			}
			log.Printf("[debug] Decoded msg: type=%s, len=%d\n", msg.Type, len(msg.Data))
	
			t.handler(peerID, msg)
		}
	}
}

func (t *TCPTransport) Send(peerID string, msg *Message) error {
	payload, err := Encode(msg)
	if err != nil {
		return err
	}

	log.Printf("[debug] sending message to %s, size=%d\n", peerID, len(payload))

	// Try to reuse existing connection
	t.lock.RLock()
	conn, ok := t.connections[peerID]
	t.lock.RUnlock()

	// If no existing connection, dial it
	if !ok {
		conn, err = net.Dial("tcp", peerID)
		if err != nil {
			return err
		}
		t.lock.Lock()
		t.connections[peerID] = conn
		t.lock.Unlock()

		// Start reading from the new connection
		go t.handleConnection(conn, peerID)
	}

	// Send the encoded payload
	_, err = conn.Write(payload)
	return err

	// t.lock.RLock()
	// conn, ok := t.connections[peerID]
	// t.lock.RUnlock()
	// if !ok {
	// 	// 尝试建立连接
	// 	conn, err := net.Dial("tcp", peerID)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	t.lock.Lock()
	// 	t.connections[peerID] = conn
	// 	t.lock.Unlock()

	// 	// 启动监听协程
	// 	go t.handleConnection(conn, peerID)
	// }

	// // 再次获取连接（可能是新建的）
	// t.lock.RLock()
	// conn = t.connections[peerID]
	// t.lock.RUnlock()

	// _, err := conn.Write(data)
	// return err
}

func (t *TCPTransport) OnMessage(h HandlerFunc) {
	t.handler = h
}

func (t *TCPTransport) Close() error {
	t.lock.Lock()
	defer t.lock.Unlock()
	for _, conn := range t.connections {
		conn.Close()
	}
	if t.listener != nil {
		return t.listener.Close()
	}
	return nil
}
