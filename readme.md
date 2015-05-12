This project contains a set of immutable copy-on-write data structures written in go.

### flist
package flist contains an immutable linked list, also known as a cons-list.  All write operations return a new list, 
making read operations inherently thread-safe.

### queue
package queue contains an immutable queue implementation based on the cons-list concept.  All operations which may modify the
structure of the queue return a queue pointer which may or may not be a new instance.

The queue implementation allows iteration through a queue snapshot which can be resumed on a later snapshot.  This behavior is only defined
for narrow cases, see the doc for details.  This particular functionality was needed for a job queueing system I'm working on.

### critbit
package critbit contains an immutable copy-on-write critbit tree.  The critbit tree is an unbalanced binary search tree which attempts to minimize the amount of time spent in navigating each individual node.  It ends up being much faster than a balanced AVL tree except in the worst-case, most unbalanced scenarios.  This is due to the fact that an AVL tree has to do a comparison of the whole key at each node, while the critbit tree compares only 1 bit per node.

The critbit tree can be a good replacement for map when the copy-on-write immutability is desired.