package immutable

import "github.com/gburgett/immutable/trie"

// A marker interface for collection entries.
type Entry interface {
}

// Represents an immutable map.  All write methods return a copy of the map,
// the original is not modified.
type Map interface {
	// Returns a new immutable Map containing the given entry at the given key.
	// Also returns the previous entry at that key, or nil if no entry existed there.
	Set(key []byte, value Entry) (Map, Entry)

	// Returns the entry at the given key.
	// Also returns an OK indicator which is true if the key existed.
	Get(key []byte) (Entry, bool)

	// Returns a new immutable Map without the entry specified by the given key.
	// Also returns the entry that was deleted, or nil if no entry was there.
	Delete(key []byte) (Map, Entry)
}

// TrieMap is a wrapper arround trie.Trie which explicitly implements Map.
type TrieMap struct {
	*trie.Trie
}

func (m TrieMap) Set(key []byte, value Entry) (Map, Entry) {
	t, e := m.Trie.Set(key, value)
	return TrieMap{
		Trie: t,
	}, e
}

func (m TrieMap) Get(key []byte) (Entry, bool) {
	e, ok := m.Trie.Get(key)
	return e, ok
}

func (m TrieMap) Delete(key []byte) (Map, Entry) {
	t, e := m.Trie.Delete(key)
	return TrieMap{
		Trie: t,
	}, e
}
