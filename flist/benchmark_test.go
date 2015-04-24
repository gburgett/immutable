package flist

import "testing"

func benchmarkCons(b *testing.B, numItems int) {

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		list := NilList()
		for j := 0; j < numItems; j++ {
			list = Cons(j, list)
		}
	}
}

func BenchmarkCons_SingleItem(b *testing.B) {
	benchmarkCons(b, 1)
}

func BenchmarkCons_100Items(b *testing.B) {
	benchmarkCons(b, 100)
}

func BenchmarkCons_10kItems(b *testing.B) {
	benchmarkCons(b, 10*1000)
}

func BenchmarkCons_1MItems(b *testing.B) {
	benchmarkCons(b, 1000*1000)
}

func benchmarkConsFromSlice(b *testing.B, numItems int) {
	slice := make([]interface{}, numItems)
	for i := 0; i < numItems; i++ {
		slice[i] = numItems - i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ConsFromSlice(slice)
	}
}

func BenchmarkConsFromSlice_SingleItem(b *testing.B) {
	benchmarkConsFromSlice(b, 1)
}

func BenchmarkConsFromSlice_100Items(b *testing.B) {
	benchmarkConsFromSlice(b, 100)
}

func BenchmarkConsFromSlice_10kItems(b *testing.B) {
	benchmarkConsFromSlice(b, 10*1000)
}

func BenchmarkConsFromSlice_1MItems(b *testing.B) {
	benchmarkConsFromSlice(b, 1000*1000)
}

func benchmarkPrepend(b *testing.B, numItems int) {
	list := NilList()
	for i := 0; i < numItems; i++ {
		list = Cons(i, list)
	}
	onto := Cons(1, NilList()) //doesn't matter how big this is

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Prepend(list, onto)
	}
}

func BenchmarkPrepend_SingleItem(b *testing.B) {
	benchmarkPrepend(b, 1)
}

func BenchmarkPrepend_100Items(b *testing.B) {
	benchmarkPrepend(b, 100)
}

func BenchmarkPrepend_10kItems(b *testing.B) {
	benchmarkPrepend(b, 10*1000)
}

func BenchmarkPrepend_1MItems(b *testing.B) {
	benchmarkPrepend(b, 1000*1000)
}

func benchmarkMap_Invert(b *testing.B, numItems int) {
	list := NilList()
	for i := 0; i < numItems; i++ {
		list = Cons(i, list)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = list.Map(func(i interface{}) interface{} {
			return numItems - i.(int)
		})
	}
}

func BenchmarkMap_SingleItem(b *testing.B) {
	benchmarkMap_Invert(b, 1)
}

func BenchmarkMap_100Items(b *testing.B) {
	benchmarkMap_Invert(b, 100)
}

func BenchmarkMap_10kItems(b *testing.B) {
	benchmarkMap_Invert(b, 10*1000)
}

func BenchmarkMap_1MItems(b *testing.B) {
	benchmarkMap_Invert(b, 1000*1000)
}

func benchmarkAggregate_Sum(b *testing.B, numItems int) {
	list := NilList()
	for i := 0; i < numItems; i++ {
		list = Cons(i, list)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = list.Aggregate(0, func(a, b interface{}) interface{} {
			return a.(int) + b.(int)
		})
	}
}

func BenchmarkAggregate_SingleItem(b *testing.B) {
	benchmarkAggregate_Sum(b, 1)
}

func BenchmarkAggregate_100Items(b *testing.B) {
	benchmarkAggregate_Sum(b, 100)
}

func BenchmarkAggregate_10kItems(b *testing.B) {
	benchmarkAggregate_Sum(b, 10*1000)
}

func BenchmarkAggregate_1MItems(b *testing.B) {
	benchmarkAggregate_Sum(b, 1000*1000)
}

func benchmarkReverse(b *testing.B, numItems int) {
	list := NilList()
	for i := 0; i < numItems; i++ {
		list = Cons(i, list)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = list.Reverse()
	}
}

func BenchmarkReverse_SingleItem(b *testing.B) {
	benchmarkReverse(b, 1)
}

func BenchmarkReverse_100Items(b *testing.B) {
	benchmarkReverse(b, 100)
}

func BenchmarkReverse_10kItems(b *testing.B) {
	benchmarkReverse(b, 10*1000)
}

func BenchmarkReverse_1MItems(b *testing.B) {
	benchmarkReverse(b, 1000*1000)
}
