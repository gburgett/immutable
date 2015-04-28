/*
package critbit contains a critbit trie implementation.  A critbit trie is essentially a binary search tree,
where unnecessary nodes are skipped.  Each node contains a bit index into the key n, and two children 0 and 1
such that the bits [0:n]0 and [0:n]1 are each prefixes of a leaf in the tree.  When navigating a critbit trie,
the bit index indicates the "critical bit", and that bit decides whether to go down the 0 child or the 1 child.

In the worst case a critbit tree is still O(log n), but the advantage is that the key comparison is much faster
than competing implementations.

Examples:

Getting an empty trie -
	tree := trie.NilTrie()

Adding an item to the trie -
	tree, _ = tree.Set([]byte{0x01}, 1)

Getting an item from the tree
	got, ok := tree.Get([]byte{0x01})
	//got.(int) == 1, ok == true

Updating an item in the tree
	var prev interface{}
	tree, prev = tree.Set([]byte{0x01}, 2)
	//prev.(int) == 1 (the old value)

Removing an item from the tree
	tree, prev = tree.Delete([]byte{0x01})
	//prev.(int) == 2
	got, ok = tree.Get([]byte{0x01})
	//got == nil, ok == false

Some benchmarks on my Macbook Pro 2.5GHz Intel Core i7, 16GB 1600MHz DDR3

	BenchmarkSet_32bit_Add_EmptyTree	 5000000	       269 ns/op
	BenchmarkSet_32bit_Add_100Items	 	 1000000	      1200 ns/op
	BenchmarkSet_32bit_Add_10kItems	  	  500000	      2191 ns/op
	BenchmarkSet_32bit_Add_100kItems	  300000	      3999 ns/op
	BenchmarkSet_64bit_Add_EmptyTree	 5000000	       360 ns/op
	BenchmarkSet_64bit_Add_100Items	 	 1000000	      1596 ns/op
	BenchmarkSet_64bit_Add_10kItems	  	  500000	      2232 ns/op
	BenchmarkSet_64bit_Add_100kItems	  500000	      3787 ns/op
	BenchmarkSet_128bit_Add_EmptyTree	 5000000	       329 ns/op
	BenchmarkSet_128bit_Add_100Items	 1000000	      1418 ns/op
	BenchmarkSet_128bit_Add_10kItems	 1000000	      2230 ns/op
	BenchmarkSet_128bit_Add_100kItems	  300000	      3688 ns/op
	BenchmarkSet_1kbyte_Add_EmptyTree	 5000000	       337 ns/op
	BenchmarkSet_1kbyte_Add_100Items	 1000000	      1606 ns/op
	BenchmarkSet_1kbyte_Add_10kItems	 1000000	      1519 ns/op
	BenchmarkGet_32bit_SingleItem		50000000	        24.3 ns/op
	BenchmarkGet_32bit_100Items			20000000	        82.1 ns/op
	BenchmarkGet_32bit_10kItems			10000000	       223 ns/op
	BenchmarkGet_32bit_100kItems	 	 3000000	       462 ns/op
	BenchmarkGet_64bit_SingleItem	   100000000	        23.7 ns/op
	BenchmarkGet_64bit_100Items			20000000	        73.0 ns/op
	BenchmarkGet_64bit_10kItems			10000000	       226 ns/op
	BenchmarkGet_64bit_100kItems	 	 3000000	       494 ns/op
	BenchmarkGet_128bit_SingleItem		50000000	        25.1 ns/op
	BenchmarkGet_128bit_100Items		20000000	        75.0 ns/op
	BenchmarkGet_128bit_10kItems		10000000	       230 ns/op
	BenchmarkGet_128bit_100kItems	 	 3000000	       506 ns/op
	BenchmarkGet_1kbyte_SingleItem		50000000	        25.5 ns/op
	BenchmarkGet_1kbyte_100Items		20000000	        79.0 ns/op
	BenchmarkGet_1kbyte_10kItems	 	 5000000	       315 ns/op
	BenchmarkDelete_32bit_SingleItem	20000000	        88.7 ns/op
	BenchmarkDelete_32bit_100Items	 	 5000000	       534 ns/op
	BenchmarkDelete_32bit_10kItems	 	 2000000	       904 ns/op
	BenchmarkDelete_32bit_100kItems	 	 1000000	      1653 ns/op
	BenchmarkDelete_64bit_SingleItem	20000000	        89.7 ns/op
	BenchmarkDelete_64bit_100Items	 	 3000000	       630 ns/op
	BenchmarkDelete_64bit_10kItems	 	 2000000	       901 ns/op
	BenchmarkDelete_64bit_100kItems	 	 1000000	      1583 ns/op
	BenchmarkDelete_128bit_SingleItem	20000000	        90.0 ns/op
	BenchmarkDelete_128bit_100Items	 	 5000000	       364 ns/op
	BenchmarkDelete_128bit_10kItems	 	 2000000	       896 ns/op
	BenchmarkDelete_128bit_100kItems	 1000000	      1544 ns/op
	BenchmarkDelete_1kbyte_SingleItem	20000000	        91.2 ns/op
	BenchmarkDelete_1kbyte_100Items	 	 5000000	       419 ns/op
	BenchmarkDelete_1kbyte_10kItems	 	 2000000	       667 ns/op

*/

package critbit
