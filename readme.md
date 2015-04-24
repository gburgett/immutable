This project contains a set of immutable data structures written in go.

### flist
package flist contains an immutable functional linked list.  All write operations return a new list, making read operations inherently
thread-safe.

### trie
package trie contains an immutable patricia trie (aka radix trie) implementation, where keys are byte slices and values are `interface{}`.
It is copy on write, so read operations are inherently thread-safe.

The radix of the trie is fixed at 256, to make the comparison logic simpler.