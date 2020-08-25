package gocache

// PeerGetter is the interface that must be implemented
// by a peer
type PeerGetter interface {
	Get(group, key string) ([]byte, error)
}

// PeerPicker is the interface that must be implemented
// to locate the peer that specified by the key
type PeerPicker interface {
	PickPeer(key string) (PeerGetter, bool)
}
