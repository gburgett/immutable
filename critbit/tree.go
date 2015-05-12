package critbit

import "bytes"

// A copy-on-write critbit Trie.  It stores key-value pairs where the key is a byte slice.
// The internal implementation is based on https://github.com/agl/critbit/blob/master/critbit.pdf
type Trie struct {
	root  *node
	count uint32
}

type node struct {
	// the index into the byte array of the byte containing the critical bit.
	critbyte int
	// a bitmask where all bits except the critical bit are true.  This allows more efficient
	// identification of the critical bit.
	critbit uint8

	children [2]*node

	//the key at this leaf.  All external leafs have a non-nil key.
	key []byte
	//The value at this leaf.  All external leafs have a non-nil value, and do not point to other nodes.
	value interface{}
}

var nilTrie *Trie = &Trie{}

// Gets the singleton nil trie.  It is a singleton because it is immutable.
func NilTrie() *Trie {
	return nilTrie
}

//-- read operations --//

// Gets an item out of the tree by its key.  Returns the item and a boolean which is
// true if the item existed.
func (t *Trie) Get(key []byte) (interface{}, bool) {
	if t.root == nil {
		return nil, false
	}

	n := t.root.findBestLeaf(key)
	if bytes.Equal(key, n.key) {
		return n.value, true
	}
	return nil, false
}

// Gets the number of items in the tree.
func (t *Trie) Len() uint32 {
	return t.count
}

//-- write operations --//

// Returns a new Trie with the given key set to the given value.
func (t *Trie) Set(key []byte, value interface{}) (*Trie, interface{}) {
	if value == nil {
		panic("value cannot be nil")
	}

	if t.root == nil {
		return &Trie{
			root: &node{
				key:   key,
				value: value,
			},
			count: 1,
		}, nil
	}

	n := t.root.findBestLeaf(key)
	if bytes.Equal(key, n.key) {
		return &Trie{
			root:  t.root.setLeaf(key, value),
			count: t.count,
		}, n.value
	}

	//insert node
	critbyte, critbit := findCritbit(key, n.key)
	return &Trie{
		root:  t.root.insertLeaf(key, value, critbyte, critbit),
		count: t.count + 1,
	}, nil
}

func (t *Trie) Delete(key []byte) (*Trie, interface{}) {
	if t.root == nil {
		return t, nil
	}

	n := t.root.findBestLeaf(key)
	if bytes.Equal(key, n.key) {
		return &Trie{
			root:  t.root.deleteLeaf(key),
			count: t.count - 1,
		}, n.value
	}

	return t, nil
}

//-- internal functions --//

func (n *node) findBestLeaf(key []byte) *node {
	if n.key != nil {
		//it's a leaf - return it
		return n
	}

	direction := findDirection(key, n.critbyte, n.critbit)
	return n.children[direction].findBestLeaf(key)
}

func (n *node) setLeaf(key []byte, value interface{}) *node {
	if n.key != nil {
		//it's the leaf - set it
		return &node{
			key:   key,
			value: value,
		}
	}

	//walk the tree, and create a new node to return pointing to our new deep child.
	direction := findDirection(key, n.critbyte, n.critbit)
	ret := &node{
		critbit:  n.critbit,
		critbyte: n.critbyte,
	}
	ret.children[1-direction] = n.children[1-direction]
	ret.children[direction] = n.children[direction].setLeaf(key, value)
	return ret
}

func (n *node) insertLeaf(key []byte, value interface{}, critbyte int, critbit uint8) *node {

	if n.key != nil ||
		n.critbyte > critbyte || (n.critbyte == critbyte && n.critbit < critbit) {
		//this is the leaf we calculated the critbit from OR
		//this node's critbit is bigger than the one we're trying to add, add a node before it
		dir := findDirection(key, critbyte, critbit)
		ret := &node{
			critbyte: critbyte,
			critbit:  critbit,
		}
		ret.children[dir] = &node{
			key:   key,
			value: value,
		}
		ret.children[1-dir] = n
		return ret
	}

	//this node's critbit is smaller than the one we're trying to add, insert after it
	dir := findDirection(key, n.critbyte, n.critbit)
	ret := &node{
		critbyte: n.critbyte,
		critbit:  n.critbit,
	}
	ret.children[dir] = n.children[dir].insertLeaf(key, value, critbyte, critbit)
	ret.children[1-dir] = n.children[1-dir]
	return ret
}

func (n *node) deleteLeaf(key []byte) *node {

	if n.key != nil {
		//this is the expected leaf delete it by returning nil
		return nil
	}

	dir := findDirection(key, n.critbyte, n.critbit)
	result := n.children[dir].deleteLeaf(key)
	if result == nil {
		//the child was deleted - this node is no longer necessary
		return n.children[1-dir]
	}

	//update the child in this node
	ret := &node{
		critbyte: n.critbyte,
		critbit:  n.critbit,
	}
	ret.children[dir] = result
	ret.children[1-dir] = n.children[1-dir]
	return ret

}

func findDirection(key []byte, critbyte int, critbit uint8) int {
	if critbit == 255 {
		//special case - length comparison
		if critbyte == len(key) {
			return 0
		}
		return 1
	}
	//identify correct child
	c := key[critbyte]
	r := (1 + (critbit | c)) >> 7
	return 1 - int(r)
}

func findCritbit(u []byte, p []byte) (int, uint8) {
	//find the critical byte

	found := false
	var newbyte int
	var newcritbit uint8

	//search through our new key u to find the first differing byte
	for newbyte = 0; newbyte < len(u); newbyte++ {

		if newbyte >= len(p) {
			return len(p), 255 //special case - u is longer
		}
		if p[newbyte] != u[newbyte] {
			newcritbit = p[newbyte] ^ u[newbyte]
			found = true
			break
		}
	}

	if !found {
		return len(u), 255 //special case - p is longer
	}

	//find the critical bit
	newcritbit |= newcritbit >> 1
	newcritbit |= newcritbit >> 2
	newcritbit |= newcritbit >> 4

	//at this point all bits including and below the highest bit are set
	//xor to unset all but the highest bit
	newcritbit = newcritbit ^ (newcritbit >> 1)
	//invert to create our critbit mask
	newcritbit = ^newcritbit

	//critbits (lowest to highest): fe, fd, fb, f7, ef, df, bf, 7f

	return newbyte, newcritbit
}
