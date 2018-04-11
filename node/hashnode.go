package node

type HashNode struct {
	key   []byte
	trie  *Trie
	dirty bool
}

func NewHash(key []byte, trie *Trie) *HashNode {
	return &HashNode{key, trie, false}
}

func (self *HashNode) RlpData() interface{} {
	return self.key
}

func (self *HashNode) Hash() interface{} {
	return self.key
}

func (self *HashNode) setDirty(dirty bool) {
	self.dirty = dirty
}

// These methods will never be called but we have to satisfy Node interface
func (self *HashNode) Value() Node { return nil }
func (self *HashNode) Dirty() bool { return true }
