package p2p

type Message struct {
	Type string // "ping", "upload", "download", "ack", etc.
	Data []byte // Payload (can be text, file content, or metadata)
}
