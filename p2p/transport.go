package p2p

// HandlerFunc is the callback for received messages.
// peerID is the remote peer's address (ip:port), msg is the decoded GOB Message.
type HandlerFunc func(peerID string, msg *Message)

// Transport is the interface for network communication between peers.
type Transport interface {
	// ListenAndAccept starts the transport and begins accepting incoming connections.
	ListenAndAccept() error

	// Send transmits a message to the specified peer.
	Send(peerID string, data []byte) error

	// Close shuts down all connections and the listener.
	Close() error

	// OnMessage registers a callback to handle incoming messages.
	OnMessage(handler HandlerFunc)
}
