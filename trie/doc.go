/*
package trie contains an immutable copy-on-write radix trie implementation.
The keys of the trie are byte slices, and the values are `interface{}`.
To keep the logic simple, the radix of the trie is fixed at 256 (i.e. 1 byte).

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

	BenchmarkSet_32bit_Add_EmptyTree	   1000000	      1031 ns/op
	BenchmarkSet_32bit_Add_100Items	  		100000	     15300 ns/op
	BenchmarkSet_32bit_Add_10kItems	   		 30000	     38297 ns/op
	BenchmarkSet_32bit_Add_100kItems	     30000	     59470 ns/op
	BenchmarkSet_64bit_Add_EmptyTree	   1000000	      1230 ns/op
	BenchmarkSet_64bit_Add_100Items	  		100000	     16799 ns/op
	BenchmarkSet_64bit_Add_10kItems	   		 50000	     35907 ns/op
	BenchmarkSet_64bit_Add_100kItems	     30000	     60793 ns/op
	BenchmarkSet_128bit_Add_EmptyTree	   1000000	      1232 ns/op
	BenchmarkSet_128bit_Add_100Items	    100000	     17203 ns/op
	BenchmarkSet_128bit_Add_10kItems	     50000	     36022 ns/op
	BenchmarkSet_128bit_Add_100kItems	     30000	     58697 ns/op
	BenchmarkSet_1kbyte_Add_EmptyTree	   1000000	      1239 ns/op
	BenchmarkSet_1kbyte_Add_100Items	    100000	     16985 ns/op
	BenchmarkSet_1kbyte_Add_10kItems	     50000	     31270 ns/op
	BenchmarkGet_32bit_SingleItem		  30000000	        42.7 ns/op
	BenchmarkGet_32bit_100Items			  20000000	        64.9 ns/op
	BenchmarkGet_32bit_10kItems			  10000000	       141 ns/op
	BenchmarkGet_32bit_100kItems	 	   5000000	       371 ns/op
	BenchmarkGet_64bit_SingleItem		  30000000	        40.2 ns/op
	BenchmarkGet_64bit_100Items			  20000000	        59.4 ns/op
	BenchmarkGet_64bit_10kItems			  10000000	       137 ns/op
	BenchmarkGet_64bit_100kItems	 	   5000000	       339 ns/op
	BenchmarkGet_128bit_SingleItem		  30000000	        42.0 ns/op
	BenchmarkGet_128bit_100Items		  20000000	        62.6 ns/op
	BenchmarkGet_128bit_10kItems		  10000000	       136 ns/op
	BenchmarkGet_128bit_100kItems	 	   5000000	       400 ns/op
	BenchmarkGet_1kbyte_SingleItem		  30000000	        42.2 ns/op
	BenchmarkGet_1kbyte_100Items		  30000000	        60.5 ns/op
	BenchmarkGet_1kbyte_10kItems	 	   5000000	       273 ns/op
	BenchmarkDelete_32bit_SingleItem	   3000000	       516 ns/op
	BenchmarkDelete_32bit_100Items	        200000	      9155 ns/op
	BenchmarkDelete_32bit_10kItems	   	     50000	     30767 ns/op
	BenchmarkDelete_32bit_100kItems	   		 30000	     56899 ns/op
	BenchmarkDelete_64bit_SingleItem	   3000000	       504 ns/op
	BenchmarkDelete_64bit_100Items	  		200000	     10441 ns/op
	BenchmarkDelete_64bit_10kItems	   	     50000	     33566 ns/op
	BenchmarkDelete_64bit_100kItems	   		 30000	     52457 ns/op
	BenchmarkDelete_128bit_SingleItem	   3000000	       530 ns/op
	BenchmarkDelete_128bit_100Items	  		200000	      9670 ns/op
	BenchmarkDelete_128bit_10kItems	   		 50000	     30439 ns/op
	BenchmarkDelete_128bit_100kItems	     30000	     51552 ns/op
	BenchmarkDelete_1kbyte_SingleItem	   3000000	       530 ns/op
	BenchmarkDelete_1kbyte_100Items	  		200000	     10135 ns/op
	BenchmarkDelete_1kbyte_10kItems	   	 	 50000	     28029 ns/op
*/
package trie
