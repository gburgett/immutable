package queue

import "fmt"

// An immutable copy-on-write queue
type Queue struct {
	// two immutable linked stacks that make up the queue.
	popper, pusher *node

	// incremented every time the popper stack is remade.  This assists in peek iteration.
	generation uint64
	// the most recently generated sequential ID
	lastSeqId uint64
}

type PeekIterator struct {
	// The value at the location in the queue.  Nil if iterated past the end of queue, or if the value at this location is nil.
	Value interface{}
	// True if the value exists, false iff iterated past the end of queue.
	HasValue bool

	// the current node that the iterator is pointing to
	current *node
	// the generation of the queue that this iterator's node belongs in
	generation uint64
}

// A node in one of the two stacks making up the queue.
type node struct {
	value interface{}

	next *node
	// the sequential id of the node.  A newly pushed node always has the highest sequential ID in the queue.
	seqId uint64
}

func NewQueue() *Queue {
	return &Queue{}
}

func NewQueueFrom(items ...interface{}) *Queue {

	var popper *node = nil
	for i := len(items) - 1; i >= 0; i-- {
		popper = &node{
			value: items[i],
			next:  popper,
			seqId: uint64(i + 1),
		}
	}

	return &Queue{
		popper:     popper,
		pusher:     nil,
		generation: 0,
		lastSeqId:  uint64(len(items)),
	}
}

// Returns a new queue with the given value pushed to the end
func (q *Queue) Push(value interface{}) *Queue {

	popper := q.popper
	pusher := q.pusher
	if popper == nil {
		// the popper should only be nil for an empty queue
		popper = &node{
			value: value,
			next:  nil,
			seqId: q.lastSeqId + 1,
		}
	} else {
		//prepend it onto the pusher list
		pusher = &node{
			value: value,
			next:  q.pusher,
			seqId: q.lastSeqId + 1,
		}
	}

	return &Queue{
		popper:     popper,
		pusher:     pusher,
		generation: q.generation,
		lastSeqId:  q.lastSeqId + 1,
	}
}

// Returns a new queue with the given value popped off the front.
// If the queue is empty returns nil & the same queue pointer
func (q *Queue) Pop() (interface{}, *Queue) {
	if q.popper == nil {
		return nil, q
	}

	popper := q.popper
	pusher := q.pusher
	gen := q.generation

	ret := popper.value
	popper = popper.next
	if popper == nil {
		//rebuild pusher onto popper & increment generation
		popper = pusher.reverse()
		gen++
	}

	return ret, &Queue{
		popper:     popper,
		pusher:     pusher,
		generation: gen,
		lastSeqId:  q.lastSeqId,
	}
}

func (q *Queue) Count() uint64 {
	if q.popper == nil {
		return 0
	}

	return q.lastSeqId - q.popper.seqId + 1
}

func (q *Queue) Peek() PeekIterator {
	if q.popper == nil {
		return PeekIterator{}
	}

	return PeekIterator{
		Value:      q.popper.value,
		HasValue:   true,
		current:    q.popper,
		generation: q.generation,
	}
}

func (q *Queue) PeekNext(current PeekIterator) (PeekIterator, *Queue) {
	if current.generation > q.generation {
		panic(fmt.Sprintf("Unexpected queue generation - encountered %d expected %d", q.generation, current.generation))
	}
	//at this point we know for sure the current node has either been popped off, or is in the popper stack somewhere (though the stack could have been rebuilt)

	c := current.current

	if q.popper == nil || !current.HasValue {
		//queue is empty or iterator is at end of queue
		return PeekIterator{}, q
	}

	if c.seqId < q.popper.seqId {
		//the node we were pointing to has been popped off, continue from the beginning
		return PeekIterator{
			Value:      q.popper.value,
			HasValue:   true,
			current:    q.popper,
			generation: q.generation,
		}, q
	}

	if current.generation < q.generation {
		//the node we were pointing to is still on queue, but the stack has been rebuilt.
		//got to find it again.
		for n := q.popper; n != nil; n = n.next {
			if n.seqId == c.seqId {
				//found the node
				c = n
				break
			}
		}
	}

	//advance to the next node
	if c.next != nil {
		//the simple case - we were already pointing to the correct node & next exists
		return PeekIterator{
			Value:      c.next.value,
			HasValue:   true,
			current:    c.next,
			generation: q.generation,
		}, q

	}

	//need to rebuild from the pusher
	if q.pusher == nil {
		//end of queue
		return PeekIterator{}, q
	}
	bottom := q.pusher.reverse()
	popper := rebuildOnto(q.popper, bottom)
	q = &Queue{
		popper:     popper,
		pusher:     nil,
		generation: q.generation + 1, //always increment generation when we rebuild the popper
		lastSeqId:  q.lastSeqId,
	}
	return PeekIterator{
		Value:      bottom.value,
		HasValue:   true,
		current:    bottom,
		generation: q.generation,
	}, q

}

func (n *node) reverse() *node {

	var ret *node = nil
	for t := n; t != nil; t = t.next {
		ret = &node{
			value: t.value,
			next:  ret,
			seqId: t.seqId,
		}
	}

	return ret
}

func rebuildOnto(top *node, bottom *node) *node {
	if top == nil {
		return bottom
	}

	bottom = rebuildOnto(top.next, bottom)
	return &node{
		value: top.value,
		next:  bottom,
		seqId: top.seqId,
	}
}
