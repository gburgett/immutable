// The Trie package contains an immutable patricia (or radix) trie implementation.  This radix trie
// has a fixed radix of 256 to make comparisons easier (i.e. we can compare bytewise instead of bitwise).
// It accepts byte slices as keys and `interface{}` as values.
package trie

import (
	"bytes"
	"fmt"
)

type node struct {
	//the leaf at this key - if not nil, the key for this item = keySlice.
	value interface{}
	//the key prefix of this node.  If this node is a leaf this is the full key.
	keySlice []byte
	//the child nodes, indexed by the next byte in their keys.
	children map[byte]*node
}

func (n *node) isLeaf() bool {
	return n.value != nil
}

/// ---- Public Members ---- //

// A Trie is a key-value store which holds items in a sorted tree structure based on the value of the key.
// This implementation is immutable; i.e. all modification methods return a new Trie which contains
// the result.  The old Trie object is left unchanged.
//
// Worst case reads/sets are O(len(key))
type Trie struct {
	//the root node.  This always has a keySlice of len 0.
	root  *node
	count uint32
}

// The nil trie.
func NilTrie() *Trie {
	return nilTrie
}

var nilTrie *Trie = &Trie{
	root: &node{
		keySlice: make([]byte, 0),
		children: newMap(0),
	},
	count: 0,
}

// Returns a new Trie with the given key set to the given value.
// Also returns the old value that was at that key.
func Set(t *Trie, key []byte, value interface{}) (*Trie, interface{}) {
	if value == nil {
		panic("value cannot be nil")
	}
	//log.Printf("setting key [% x]", key)

	if len(key) == 0 {
		//the root node is the leaf for this key
		count := t.count
		if !t.root.isLeaf() {
			count++
		}
		return &Trie{
			root: &node{
				value:    value,
				keySlice: make([]byte, 0),
				children: copyMap(t.root.children, len(t.root.children)),
			},
			count: count,
		}, t.root.value
	}

	//clone the byte key to ensure immutability
	cloned := make([]byte, len(key))
	copy(cloned, key)
	newNode, old := setNode(t.root, cloned, value)
	count := t.count
	if old == nil {
		count++
	}
	return &Trie{
		root:  newNode,
		count: count,
	}, old
}

// Gets the value for the given key.  If the key doesn't exist,
// the second return value will be false and the value will be nil.
func (t *Trie) Get(key []byte) (interface{}, bool) {
	if len(key) == 0 {
		return t.root.value, t.root.isLeaf()
	}

	n := t.root.getNode(key)
	if n == nil {
		return nil, false
	}
	return n.value, n.isLeaf()
}

// Gets the number of items in the trie.
func (t *Trie) Len() uint32 {
	return t.count
}

// Deletes an item out of the trie.  Returns a new trie with the item removed, and the
// value of the previous item.
func Delete(t *Trie, key []byte) (*Trie, interface{}) {
	if len(key) == 0 {
		//the root node is the leaf for this key
		if t.root.isLeaf() {
			return &Trie{
				root: &node{
					value:    nil,
					keySlice: make([]byte, 0),
					children: copyMap(t.root.children, len(t.root.children)),
				},
				count: t.count - 1,
			}, t.root.value
		}
		return t, nil
	}
	newRoot, val := deleteNode(t.root, key)

	if newRoot != t.root {
		return &Trie{
			root:  newRoot,
			count: t.count - 1,
		}, val
	}
	return t, nil
}

/// --- Set functions --- //
func setNode(t *node, key []byte, value interface{}) (*node, interface{}) {

	childKey := key[len(t.keySlice)] //childKey is the next byte in the array

	//log.Printf("looking in node [% x], childKey %x", t.keySlice, childKey)

	child, ok := t.children[childKey]
	childMap := copyMap(t.children, len(t.children))
	var old interface{} = nil
	if !ok {
		//log.Printf("did not find child, making new one")
		//make a new node for that key as the child
		childMap = copyMap(t.children, len(t.children)+1)
		n := &node{
			value:    value,
			keySlice: key,
			children: newMap(0),
		}
		childMap[childKey] = n
	} else {
		idx := firstDiffIndex(child.keySlice[len(t.keySlice):], key[len(t.keySlice):]) + len(t.keySlice)

		if idx == len(key) {
			if idx == len(child.keySlice) {
				//log.Printf("formatting node [% x]", child.keySlice)
				//found the node - set it
				n := &node{
					value:    value,
					keySlice: key,
					children: copyMap(child.children, len(child.children)),
				}
				childMap[childKey] = n
				old = child.value
			} else {
				//log.Printf("adding leaf above child [% x]", child.keySlice)
				//need a new leaf with the child as its child
				n := &node{
					value:    value,
					keySlice: key,
					children: newMap(1),
				}
				n.children[child.keySlice[idx]] = child
				childMap[childKey] = n
			}
		} else if idx == len(child.keySlice) {
			//the new leaf goes below the child - recurse into it
			//log.Printf("recursing into child [% x]", child.keySlice)
			//debugging
			if len(key) <= len(child.keySlice) {
				fmt.Println(t.printDbg(""))
				panic(fmt.Sprintf("key [% x] smaller than or equal to node slice [% x]", key, child.keySlice))
			}
			if !startsWith(key, child.keySlice) {
				panic(fmt.Sprintf("key [% x] doesn't start with node slice [% x]", key, child.keySlice))
			}

			//the child is a node with multiple children itself - recurse into it
			var n *node
			n, old = setNode(child, key, value)

			childMap[childKey] = n
		} else {
			//need to make a new node, containing the child and the new leaf as children
			idx := firstDiffIndex(child.keySlice, key)
			newChild := &node{
				keySlice: key[:idx],
				children: newMap(2),
			}
			newChild.children[child.keySlice[idx]] = child
			newChild.children[key[idx]] = &node{
				keySlice: key,
				value:    value,
				children: newMap(0),
			}
			childMap[childKey] = newChild
		}
	}

	//pointer to same value, with same key, but new children
	ret := &node{
		value:    t.value,
		keySlice: t.keySlice,
		children: childMap, //update the children
	}
	//log.Printf("making node: %s", ret.printDbg(""))
	return ret, old
}

func newMap(capacity int) map[byte]*node {
	return make(map[byte]*node, capacity)
}

func copyMap(m map[byte]*node, capacity int) map[byte]*node {
	ret := newMap(capacity)
	for k, v := range m {
		ret[k] = v
	}
	return ret
}

func startsWith(longer []byte, prefix []byte) bool {
	for i, b := range prefix {
		if longer[i] != b {
			return false
		}
	}
	return true
}

func firstDiffIndex(left []byte, right []byte) int {
	var longer, shorter []byte
	if len(left) > len(right) {
		longer = left
		shorter = right
	} else {
		longer = right
		shorter = left
	}

	for i, b := range shorter {
		if longer[i] != b {
			return i
		}
	}
	return len(shorter)
}

func deleteNode(t *node, key []byte) (*node, interface{}) {

	if len(key) < len(t.keySlice) {
		return t, nil
	}

	childKey := key[len(t.keySlice)] //childKey is the next byte in the array
	child, ok := t.children[childKey]
	if !ok {
		//no node with that key exists
		return t, nil
	}

	if bytes.Equal(child.keySlice[len(t.keySlice):], key[len(t.keySlice):]) {
		//this child is the node to delete.
		if !child.isLeaf() {
			//there's no value here, return no change
			return t, nil
		}

		var newChild *node
		length := len(child.children)
		if length > 1 {
			//just set this child's value to nil to make it no longer a leaf
			newChild = &node{
				value:    nil,
				keySlice: child.keySlice,
				children: child.children,
			}
			//replace the child in our node's children
			children := copyMap(t.children, len(t.children))
			children[childKey] = newChild
			return &node{
				value:    t.value,
				keySlice: t.keySlice,
				children: children,
			}, child.value
		} else if length == 1 {
			//replace it in our list with its child
			children := make(map[byte]*node, 1)
			for _, v := range child.children {
				children[childKey] = v
			}
			return &node{
				value:    t.value,
				keySlice: t.keySlice,
				children: children,
			}, child.value
		} else {
			//delete it from our list
			children := newMap(len(t.children) - 1)
			for k, v := range t.children {
				if k == childKey {
					continue
				}
				children[k] = v
			}
			return &node{
				value:    t.value,
				keySlice: t.keySlice,
				children: children,
			}, child.value
		}
	} else {
		newChild, old := deleteNode(child, key)
		if newChild == child {
			//didn't find it
			return t, nil
		}
		//replace the child in our node's children
		children := copyMap(t.children, len(t.children))
		children[childKey] = newChild
		return &node{
			value:    nil,
			keySlice: make([]byte, 0),
			children: children,
		}, old
	}
}

/// --- Get functions ---//
func (t *node) getNode(key []byte) *node {

	if len(key) < len(t.keySlice) {
		return nil
	}

	childKey := key[len(t.keySlice)] //childKey is the next byte in the array
	child, ok := t.children[childKey]
	if !ok {
		return nil
	}

	if bytes.Equal(child.keySlice, key) {
		return child
	} else {
		return child.getNode(key)
	}
}

func (t *node) printDbg(offset string) string {
	var buffer bytes.Buffer

	offset2 := offset + "  "
	offset3 := offset + "    "

	buffer.WriteString("{\n")
	buffer.WriteString(offset2)
	if len(t.keySlice) == 0 {
		buffer.WriteString("key: root\n")
	} else {
		buffer.WriteString(fmt.Sprintf("key: [% x],\n", t.keySlice))
	}
	if t.value != nil {
		buffer.WriteString(offset2)
		buffer.WriteString(fmt.Sprintf("value: %v,\n", t.value))
	}
	buffer.WriteString(offset2)
	buffer.WriteString("children: {\n")
	for k, v := range t.children {
		buffer.WriteString(offset3)
		buffer.WriteString(fmt.Sprintf("%x: ", k))
		buffer.WriteString(v.printDbg(offset3 + "  "))
		buffer.WriteString(",\n")
	}
	buffer.WriteString(offset2)
	buffer.WriteString("},\n")
	buffer.WriteString(offset)
	buffer.WriteString("}")

	return buffer.String()
}
