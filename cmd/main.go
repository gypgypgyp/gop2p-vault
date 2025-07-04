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
	"gop2p-vault/server"
	"gop2p-vault/protocol"
)

import "strconv"
import "gop2p-vault/store"

func parseInt64(s string) int64 {
	n, _ := strconv.ParseInt(s, 10, 64)
	return n
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go [port] (e.g. :9000 or :9001)")
	}
	selfAddr := os.Args[1]

	transport := p2p.NewTCPTransport(selfAddr)


	transport.OnMessage(func(peerID string, msg *p2p.Message) {
		fmt.Printf("[recv from %s] Type: %s | Len: %d bytes\n", peerID, msg.Type, len(msg.Data))
	
		switch msg.Type {
		case "text":
			fmt.Println("[text msg]:", string(msg.Data))
	
		case "upload":
			fmt.Println("[debug] Received upload message")

			// Compute file hash as key
			key, err := server.HandleUpload(msg.Data)
			if err != nil {
				fmt.Println("Failed to hash content:", err)
				return
			}
			fmt.Println("[upload]: File stored with key", key)

		case "download":
			fileKey := string(msg.Data)
			resp, err := server.HandleDownload(fileKey)
			if err != nil {
				fmt.Println("[download]: Requested file not found:", err)
				return
			}

			err = transport.Send(peerID, resp)
			if err != nil {
				fmt.Println("Failed to send download_result:", err)
			} else {
				fmt.Println("[download]: Sent file", fileKey, "to", peerID)
			}

		case "download_result":
			path, err := server.HandleDownloadResult(msg.Data)
			if err != nil {
				fmt.Println("Failed to save downloaded file:", err)
			} else {
				fmt.Println("[download_result]: File saved to", path)
			}

		case "delete":
			fileKey := string(msg.Data)
			err := server.HandleDelete(fileKey)
			if err != nil {
				fmt.Println("[delete]: Failed to delete file:", err)
			} else {
				fmt.Println("[delete]: File", fileKey, "deleted successfully")
			}

		case "metadata":
			md, err := protocol.DecodeMetadata(msg.Data)
			if err != nil {
				fmt.Println("Failed to decode metadata:", err)
				return
			}
			fmt.Printf("[metadata] Received metadata: Name=%s, Size=%d, Hash=%s, Mime=%s\n",
				md.Name, md.Size, md.Hash, md.MimeType)

		case "have":
			fileKey := string(msg.Data)
			filePath := store.HashPath("./data", fileKey)
    		_, err := os.Stat(filePath)
			reply := "no"
			if err == nil {
				reply = "yes"
			}
			resp := &p2p.Message{
				Type: "have_result",
				Data: []byte(reply),
			}
			transport.Send(peerID, resp)
		
		case "have_result":
			fmt.Printf("[have_result] Peer %s has file: %s\n", peerID, string(msg.Data))
			

		default:
			fmt.Println("Unknown message type:", msg.Type)
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

		for {
			fmt.Print("Enter command (text <msg> | upload <file> | download <key> | metadata): ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)
		
			if strings.HasPrefix(input, "text ") {
				text := strings.TrimPrefix(input, "text ")
				msg := &p2p.Message{Type: "text", Data: []byte(text)}
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
				transport.Send(targetAddr, msg)
				
			} else if strings.HasPrefix(input, "download ") {
				key := strings.TrimPrefix(input, "download ")
				msg := &p2p.Message{Type: "download", Data: []byte(key)}
				err := transport.Send(targetAddr, msg)
				if err != nil {
					fmt.Printf("Send error: %v\n", err)
				} else {
					fmt.Println("[debug] download request sent for key:", key)
				}
			} else if strings.HasPrefix(input, "delete ") {
				key := strings.TrimPrefix(input, "delete ")
				msg := &p2p.Message{Type: "delete", Data: []byte(key)}
				err := transport.Send(targetAddr, msg)
				if err != nil {
					fmt.Printf("Send error: %v\n", err)
				} else {
					fmt.Println("[debug] delete request sent for key:", key)
				}
			} else if strings.HasPrefix(input, "meta ") {
				fields := strings.Fields(strings.TrimPrefix(input, "meta "))
				if len(fields) != 4 {
					fmt.Println("Usage: meta <Name> <Size> <Hash> <Mime>")
					continue
				}
				md := &protocol.Metadata{
					Name:     fields[0],
					Size:     parseInt64(fields[1]),
					Hash:     fields[2],
					MimeType: fields[3],
				}
				data, err := protocol.EncodeMetadata(md)
				if err != nil {
					fmt.Println("Encoding error:", err)
					continue
				}
				msg := &p2p.Message{Type: "metadata", Data: data}
				transport.Send(targetAddr, msg)
			}else if strings.HasPrefix(input, "have ") {
				key := strings.TrimPrefix(input, "have ")
				msg := &p2p.Message{Type: "have", Data: []byte(key)}
				transport.Send(targetAddr, msg)
			}
		}
	}
	select {} // 阻止主线程退出（用于监听模式）
}