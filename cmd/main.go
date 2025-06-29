package main

import (
	"io"
	"strings"
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"gop2p-vault/p2p"
	"gop2p-vault/store"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go [port] (e.g. :9000 or :9001)")
	}
	selfAddr := os.Args[1]

	transport := p2p.NewTCPTransport(selfAddr)

	// 设置收到消息时的回调函数
	// transport.OnMessage(func(peerID string, data []byte) {
	// 	fmt.Printf("[recv from %s]: %s\n", peerID, string(data))
	// })

	// 设置收到gob时的回调函数
	// transport.OnMessage(func(peerID string, msg *p2p.Message) {
	// 	fmt.Printf("[recv from %s] Type: %s | Data: %s\n", peerID, msg.Type, string(msg.Data))
	// })


	transport.OnMessage(func(peerID string, msg *p2p.Message) {
		fmt.Printf("[recv from %s] Type: %s | Len: %d bytes\n", peerID, msg.Type, len(msg.Data))
	
		switch msg.Type {
		case "text":
			fmt.Println("[text msg]:", string(msg.Data))
	
		case "upload":
			fmt.Println("[debug] Received upload message")

			// Compute file hash as key
			key, err := store.HashKey(strings.NewReader(string(msg.Data)))
			if err != nil {
				fmt.Println("Failed to hash content:", err)
				return
			}
			s := store.New("./data")
			err = s.Write(key, strings.NewReader(string(msg.Data)))
			if err != nil {
				fmt.Println("Failed to write file:", err)
				return
			}
			fmt.Println("[upload]: File stored with key", key)
		}
	})
	
	

	// 启动监听
	go func() {
		if err := transport.ListenAndAccept(); err != nil {
			log.Fatalf("listen error: %v", err)
		}
	}()

	// 等待监听启动完成
	time.Sleep(time.Second)

	// 如果是发送端节点，我们连接另一个 peer 并发送消息
	if len(os.Args) == 3 {
		targetAddr := os.Args[2] // e.g. :9000
		reader := bufio.NewReader(os.Stdin)

		// for {
		// 	fmt.Print("Enter message: ")
		// 	msg, _ := reader.ReadString('\n')
		// 	err := transport.Send(targetAddr, []byte(msg))
		// 	if err != nil {
		// 		fmt.Printf("Send error: %v\n", err)
		// 	}
		// }

		for {
			fmt.Print("Enter command (text <msg> | upload <file>): ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)
		
			if strings.HasPrefix(input, "text ") {
				text := strings.TrimPrefix(input, "text ")
				msg := &p2p.Message{Type: "text", Data: []byte(text)}
				// payload, _ := p2p.Encode(msg)
				transport.Send(targetAddr, msg)
		
			} else if strings.HasPrefix(input, "upload ") {
				filename := strings.TrimPrefix(input, "upload ")
				f, err := os.Open(filename)
				if err != nil {
					fmt.Println("Failed to open file:", err)
					continue
				}
				defer f.Close()
		
				content, _ := io.ReadAll(f)
				msg := &p2p.Message{Type: "upload", Data: content}
				// payload, _ := p2p.Encode(msg)
				transport.Send(targetAddr, msg)
				
			}
		}

	}

	select {} // 阻止主线程退出（用于监听模式）
}
