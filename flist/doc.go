/*

Package flist implements an immutable cons-list with some functional methods.
The cons-list is a singly-linked list with immutable values.  A cons-list
can only be prepended to.  In order to append the entire list must be rebuilt.

Examples:

creating a list with a single element -
    list := flist.Cons(1, flist.NilList())
	//newList -> 1 -> NilList()

prepending an item onto a list -
	list = flist.Cons(2, list)
	//newList -> 2 -> 1 -> NilList()

prepending another list onto a list -
	onto := flist.ConsFromSlice([]{4})
	list := flist.ConsFromSlice([]{1, 2, 3})
	newList := flist.Prepend(list, onto)
	//newList -> 1 -> 2 -> 3 -> 4 -> NilList()

creating a list from a slice -
	list := flist.ConsFromSlice([]{1, 2, 3, 4})
	//newList -> 1 -> 2 -> 3 -> 4 -> NilList()

iterating over a list -
	for l := list; !l.IsNil(); l = l.Tail(){
		val := l.Head().(int)
	}

Some benchmarks on my Macbook Pro 2.5GHz Intel Core i7, 16GB 1600MHz DDR3

	BenchmarkCons_SingleItem			20000000	       114 ns/op
	BenchmarkCons_100Items	  			  200000	     11749 ns/op
	BenchmarkCons_10kItems	    			2000	   1111490 ns/op
	BenchmarkCons_1MItems	      			  10	 105085395 ns/op

	BenchmarkConsFromSlice_SingleItem	10000000	       125 ns/op
	BenchmarkConsFromSlice_100Items	  	  200000	     11652 ns/op
	BenchmarkConsFromSlice_10kItems	    	2000	    735558 ns/op
	BenchmarkConsFromSlice_1MItems	      	  20	  70324017 ns/op

	BenchmarkPrepend_SingleItem			10000000	       149 ns/op
	BenchmarkPrepend_100Items	  		  100000	     15145 ns/op
	BenchmarkPrepend_10kItems	    		1000	   1853248 ns/op
	BenchmarkPrepend_1MItems	       		   3	 349360372 ns/op

	BenchmarkMap_SingleItem				10000000	       224 ns/op
	BenchmarkMap_100Items	  			  100000	     22479 ns/op
	BenchmarkMap_10kItems	     			 500	   2459189 ns/op
	BenchmarkMap_1MItems	       			   3	 428336081 ns/op

	BenchmarkAggregate_SingleItem		10000000	       134 ns/op
	BenchmarkAggregate_100Items	  		  200000	      7417 ns/op
	BenchmarkAggregate_10kItems	    		2000	    570809 ns/op
	BenchmarkAggregate_1MItems	      		  30	  53183841 ns/op

	BenchmarkReverse_SingleItem			10000000	       154 ns/op
	BenchmarkReverse_100Items	  		  100000	     14649 ns/op
	BenchmarkReverse_10kItems	    		2000	    847436 ns/op
	BenchmarkReverse_1MItems	      		  20	  77569776 ns/op

*/
package flist
