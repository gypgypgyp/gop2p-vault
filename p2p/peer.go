package p2p

import "sync"

type PeerInfo struct {
	ID      string // e.g. "127.0.0.1:9000"
	Version string // optional
}

type Peer struct {
	ID      string
	Address string
}

type PeerStore struct {
	peers map[string]*Peer
	lock  sync.RWMutex
}

func NewPeerStore() *PeerStore {
	return &PeerStore{
		peers: make(map[string]*Peer),
	}
}

func (ps *PeerStore) Add(peer *Peer) {
	ps.lock.Lock()
	defer ps.lock.Unlock()
	ps.peers[peer.ID] = peer
}

func (ps *PeerStore) Get(id string) (*Peer, bool) {
	ps.lock.RLock()
	defer ps.lock.RUnlock()
	peer, ok := ps.peers[id]
	return peer, ok
}

func (ps *PeerStore) All() []*Peer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()
	var all []*Peer
	for _, p := range ps.peers {
		all = append(all, p)
	}
	return all
}

func (ps *PeerStore) Delete(id string) {
	ps.lock.Lock()
	defer ps.lock.Unlock()
	delete(ps.peers, id)
}
