/*

Package queue implements a copy on write queue based on cons-lists.  This queue has average O(1) pushing & popping, and O(n) iteration.  I say average because
every (n - i)th pop requires an O(i) reversal of the stack of items that have accumulated on the push side.  Over the lifetime of the queue every pushed item will need to
be reversed once.

The copy on write queue allows iterating over a snapshot of the queue, and even resuming iteration inside a future snapshot of the same queue.  This means you can keep
your iterator around while you pop items off the queue, and push new items onto the queue, then continue iterating from where you left off.

Examples:

creating a queue and pushing some items -
    queue := queue.NewQueue().Push(10).Push(11).Push(12)

creating a new queue from a set of items -
    items := []interface{}{1, 5, 2, 7, 4, 6}
    queue := queue.NewQueueFrom(items...)

popping items off a queue -
	q := queue
	for v, q := queue.Pop(); v != nil; v, q = q.Pop() {
		// do something with v
	}
	q.Count() == 0 //true
	queue.Count() > 0 //true

iterating a queue (simple) -
	oldqueue := queue
	for i := queue.Peek(); i.HasValue; i, queue = queue.PeekNext(i) {
		// do something with i.Value
	}
	queue != oldqueue // possibly true depending on the internal balance of the queue.

complex iteration -
	It is possible to hold onto an old iterator and resume iteration in a later snapshot of the queue.  To understand when this is possible,
	it is important to understand that the copy-on-write nature of the queue creates a parent-child relationship between queue snapshots.
	When a Push, Pop, or PeekNext operation is performed, one of the return values is potentially a new instance of the queue.  This new
	instance is a 'child' of the queue that the operation was performed on, and the original queue instance is the 'parent'.  In the majority
	of cases, the parent should be discarded & garbage collected.  However, it is possible to split the geneaology of a queue by performing
	two different operations on the same parent queue.  This has consequences for resuming iteration.

	Example:
						grandparent := queue.NewQueueFrom(items...)
						it2 := queue.Peek()
									|
						parent := grandparent.Push(17)
						/ 						\
	v, left := parent.Pop()				 		 \
	itLeft, left := left.PeekNext()			  	  \
					 /							v, right := parent.Push(23)
					/							itRight := right.Peek()
				   /								 \
		// This is OK								// so is this
	it2, left := left.PeekNext(it1)				it3, right := right.PeekNext(it1)
				 |									    |
	 	// This may or may not panic					|
	it4, left := parent.PeekNext(it2)					|
													//this is totally undefined (no clue what could happen)
												it5, right := right.PeekNext(itLeft)

	So, when you have an iterator, please ensure you only use it to resume iteration from a child of the queue snapshot that
	generated the iterator.  PeekNext with any queue snapshot that is not the one that generated the iterator or one of its children
	is undefined, and may panic.


Some benchmarks on my Macbook Pro 2.5GHz Intel Core i7, 16GB 1600MHz DDR3
	BenchmarkPush_EmptyQueue				10000000	       225 ns/op
	BenchmarkPush_SingleItem				10000000	       229 ns/op
	BenchmarkPush_BalancedQueue				10000000	       228 ns/op

	BenchmarkPop_EmptyQueue			  	  1000000000	         2.58 ns/op
	BenchmarkPop_SingleItem			    	20000000	       100 ns/op
	BenchmarkPop_MustRebuild10Items	 	 	 1000000	      1081 ns/op
	BenchmarkPop_MustRebuild1kItems	   	   	   20000	     88423 ns/op
	BenchmarkPop_MustRebuild100kItems	    	 200	   8537657 ns/op

	BenchmarkPeekNext_BestCase		   	   100000000	        11.8 ns/op
	BenchmarkPeekNext_MustRebuild10Items	 1000000	      1491 ns/op
	BenchmarkPeekNext_MustRebuild1kItems	   20000	     93085 ns/op
	BenchmarkPeekNext_MustRebuild100kItems	     200	   8485226 ns/op

	BenchmarkIterateWholeQueue_10Items	 	 1000000	      1552 ns/op		155.2 ns/peek amoritized
	BenchmarkIterateWholeQueue_1kItems	   	   10000	    110540 ns/op		111   ns/peek amoritized
	BenchmarkIterateWholeQueue_100kItems	     100	  10373414 ns/op		104   ns/peek amoritized
	ok  	github.com/gburgett/immutable/queue	30.840s

Some insights:
	* performance really suffers when the queue must be rebuilt.  Don't throw this effort away.  Any operation which has the potential
	  to rebuild the queue will return the newly rebuilt queue.
	* rebuilding cost is amoritized over the lifetime of the queue.  If your usage pattern is a well-distributed mix of push and pop this
	  works great.  If your usage pattern is heavy write followed by heavy read you may have to work around a multi-millisecond rebuild.

*/
package queue
