package flist

// A singly-linked list of queue nodes which is immutable.  This provides built-in
// thread safety for reads.
type List struct {
	count int

	head interface{}
	tail *List
}

var nilList *List = &List{
	count: 0,
	head:  nil,
}

//-- Construction --//

// The empty list.
func NilList() *List {
	return nilList
}

// Constructs a list from the given queue node items.  The resulting list
// has the same iteration order as the items.
func ConsFromSlice(items []interface{}) *List {
	var l = nilList
	for i := len(items) - 1; i >= 0; i-- {
		l = Cons(items[i], l)
	}
	return l
}

// Creates a new list by prepending the given node to the head.  This is the only way to add an item to an immutable list.
func Cons(item interface{}, list *List) *List {
	return &List{
		count: list.count + 1,
		head:  item,
		tail:  list,
	}
}

// Creates a new list by prepending the given list of items onto the given list.  This is equivalent to
// items.Count() cons operations
func Prepend(items *List, onto *List) *List {
	if items.IsNil() {
		return onto
	}

	return Cons(items.head, Prepend(items.tail, onto))
}

//-- Properties --//

// Gets the count of the list.  This is an O(1) operation.
func (l *List) Count() int {
	return l.count
}

// Returns true if this is the end of the list, aka the nil list.
func (l *List) IsNil() bool {
	return l.count == 0
}

// Gets the interface{} at this index in the list.
func (l *List) Head() interface{} {
	return l.head
}

// Gets the next item in the list.  Can be used in a for loop to iterate the items in the list, ex:
//   for l := mylist; !l.IsNil(); l = l.Next() {
//		n := l.Tail(); //do something
//	 }
//
func (l *List) Tail() *List {
	return l.tail
}

//-- Functions --//

// Applies a transformation function over the list, mapping the values in the list to a new value.
// Results in a new list containing the transformed values for each index.
func (l *List) Map(f func(interface{}) interface{}) *List {
	if l.IsNil() {
		return l
	}

	head := f(l.head)
	return Cons(head, l.tail.Map(f))
}

// Applies a filtering function to the list, returning a new list containing only the items
// for which the filter returns true.
func (l *List) Filter(f func(interface{}) bool) *List {
	if l.IsNil() {
		return l
	}

	if f(l.head) {
		return Cons(l.head, l.tail.Filter(f))
	} else {
		return l.tail.Filter(f)
	}
}

// Returns a new list which is in the reverse order of this list.
func (l *List) Reverse() *List {

	ret := nilList
	for c := l; !c.IsNil(); c = c.tail {
		ret = Cons(c.head, ret)
	}

	return ret
}

// Performs an aggregation over the list.  The aggregation function's first parameter is the
// running result, whose initial value is the 'init' parameter.  The aggregation function's second
// parameter is the current head, and the return value is the new running total.
// The return value of Aggregate is the final result of the aggregation function.
func (l *List) Aggregate(init interface{}, agg func(a, b interface{}) interface{}) interface{} {
	ret := init
	for c := l; !c.IsNil(); c = c.tail {
		ret = agg(ret, c.head)
	}
	return ret
}

// Returns a slice which contains the items of this list in order.
func (l *List) ToSlice() []interface{} {

	ret := make([]interface{}, l.count)
	i := 0

	for c := l; !c.IsNil(); c = c.tail {
		ret[i] = c.head
		i++
	}

	return ret
}

// Returns a channel which will receive all the items in this list.  This is
// useful for larger lists that you don't want to turn into a slice.
func (l *List) ToChan() <-chan interface{} {
	ch := make(chan interface{})
	go func() {
		for c := l; !c.IsNil(); c = c.tail {
			ch <- c.head
		}
		close(ch)
	}()
	return ch
}
