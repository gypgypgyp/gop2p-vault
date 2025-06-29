package p2p

import (
	"encoding/gob"
	"net"
	"time"
	"log"
)

const HandshakeTimeout = 5 * time.Second

// PerformHandshake acts as both client & server depending on initiator flag
func PerformHandshake(conn net.Conn, selfID string, isInitiator bool) (*PeerInfo, error) {
	conn.SetDeadline(time.Now().Add(HandshakeTimeout))
	defer conn.SetDeadline(time.Time{}) // clear timeout

	encoder := gob.NewEncoder(conn)
	decoder := gob.NewDecoder(conn)

	if isInitiator {
		// 1. Send our info
		if err := encoder.Encode(PeerInfo{ID: selfID}); err != nil {
			return nil, err
		}
		// 2. Wait for response
		var remote PeerInfo
		if err := decoder.Decode(&remote); err != nil {
			return nil, err
		}
		log.Printf("[P2P] Handshake complete with %s (initiator)", remote.ID)
		return &remote, nil
	} else {
		// 1. Wait for initiator info
		var remote PeerInfo
		if err := decoder.Decode(&remote); err != nil {
			return nil, err
		}
		// 2. Respond with our info
		if err := encoder.Encode(PeerInfo{ID: selfID}); err != nil {
			return nil, err
		}
		log.Printf("[P2P] Handshake complete with %s (responder)", remote.ID)
		
		return &remote, nil
	}
}
