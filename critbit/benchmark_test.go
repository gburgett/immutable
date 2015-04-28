package critbit

import (
	"crypto/rand"
	"testing"
)

func benchmarkSet_Add(b *testing.B, numItems int, keyLen int) {
	tree := NilTrie()

	for i := 0; i < numItems; i++ {
		key := makeRandomKey(b, keyLen)
		tree, _ = tree.Set(key, i)
	}

	keys := make([][]byte, 1000)
	for i := 0; i < len(keys); i++ {
		keys[i] = makeRandomKey(b, keyLen)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = tree.Set(keys[i%len(keys)], i)
	}
}

func BenchmarkSet_32bit_Add_EmptyTree(b *testing.B) {
	benchmarkSet_Add(b, 0, 32/8)
}

func BenchmarkSet_32bit_Add_100Items(b *testing.B) {
	benchmarkSet_Add(b, 100, 32/8)
}

func BenchmarkSet_32bit_Add_10kItems(b *testing.B) {
	benchmarkSet_Add(b, 10*1000, 32/8)
}

func BenchmarkSet_32bit_Add_100kItems(b *testing.B) {
	benchmarkSet_Add(b, 100*1000, 32/8)
}

func BenchmarkSet_64bit_Add_EmptyTree(b *testing.B) {
	benchmarkSet_Add(b, 0, 64/8)
}

func BenchmarkSet_64bit_Add_100Items(b *testing.B) {
	benchmarkSet_Add(b, 100, 64/8)
}

func BenchmarkSet_64bit_Add_10kItems(b *testing.B) {
	benchmarkSet_Add(b, 10*1000, 64/8)
}

func BenchmarkSet_64bit_Add_100kItems(b *testing.B) {
	benchmarkSet_Add(b, 100*1000, 64/8)
}

func BenchmarkSet_128bit_Add_EmptyTree(b *testing.B) {
	benchmarkSet_Add(b, 0, 128/8)
}

func BenchmarkSet_128bit_Add_100Items(b *testing.B) {
	benchmarkSet_Add(b, 100, 128/8)
}

func BenchmarkSet_128bit_Add_10kItems(b *testing.B) {
	benchmarkSet_Add(b, 10*1000, 128/8)
}

func BenchmarkSet_128bit_Add_100kItems(b *testing.B) {
	benchmarkSet_Add(b, 100*1000, 128/8)
}

func BenchmarkSet_1kbyte_Add_EmptyTree(b *testing.B) {
	benchmarkSet_Add(b, 0, 128/8)
}

func BenchmarkSet_1kbyte_Add_100Items(b *testing.B) {
	benchmarkSet_Add(b, 100, 128/8)
}

func BenchmarkSet_1kbyte_Add_10kItems(b *testing.B) {
	benchmarkSet_Add(b, 10*1000, 1024)
}

func benchmarkGet(b *testing.B, numItems int, keyLen int) {
	tree := NilTrie()

	keys := make([][]byte, numItems)
	for i := 0; i < len(keys); i++ {
		key := makeRandomKey(b, keyLen)
		keys[i] = key
		tree, _ = tree.Set(key, i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = tree.Get(keys[i%len(keys)])
	}
}

func BenchmarkGet_32bit_SingleItem(b *testing.B) {
	benchmarkGet(b, 1, 32/8)
}

func BenchmarkGet_32bit_100Items(b *testing.B) {
	benchmarkGet(b, 100, 32/8)
}

func BenchmarkGet_32bit_10kItems(b *testing.B) {
	benchmarkGet(b, 10*1000, 32/8)
}

func BenchmarkGet_32bit_100kItems(b *testing.B) {
	benchmarkGet(b, 100*1000, 32/8)
}

func BenchmarkGet_64bit_SingleItem(b *testing.B) {
	benchmarkGet(b, 1, 64/8)
}

func BenchmarkGet_64bit_100Items(b *testing.B) {
	benchmarkGet(b, 100, 64/8)
}

func BenchmarkGet_64bit_10kItems(b *testing.B) {
	benchmarkGet(b, 10*1000, 64/8)
}

func BenchmarkGet_64bit_100kItems(b *testing.B) {
	benchmarkGet(b, 100*1000, 64/8)
}

func BenchmarkGet_128bit_SingleItem(b *testing.B) {
	benchmarkGet(b, 1, 128/8)
}

func BenchmarkGet_128bit_100Items(b *testing.B) {
	benchmarkGet(b, 100, 128/8)
}

func BenchmarkGet_128bit_10kItems(b *testing.B) {
	benchmarkGet(b, 10*1000, 128/8)
}

func BenchmarkGet_128bit_100kItems(b *testing.B) {
	benchmarkGet(b, 100*1000, 128/8)
}

func BenchmarkGet_1kbyte_SingleItem(b *testing.B) {
	benchmarkGet(b, 1, 128/8)
}

func BenchmarkGet_1kbyte_100Items(b *testing.B) {
	benchmarkGet(b, 100, 128/8)
}

func BenchmarkGet_1kbyte_10kItems(b *testing.B) {
	benchmarkGet(b, 10*1000, 1024)
}

func benchmarkDelete(b *testing.B, numItems int, keyLen int) {
	tree := NilTrie()

	keys := make([][]byte, numItems)
	for i := 0; i < len(keys); i++ {
		key := makeRandomKey(b, keyLen)
		keys[i] = key
		tree, _ = tree.Set(key, i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = tree.Delete(keys[i%len(keys)])
	}
}

func BenchmarkDelete_32bit_SingleItem(b *testing.B) {
	benchmarkDelete(b, 1, 32/8)
}

func BenchmarkDelete_32bit_100Items(b *testing.B) {
	benchmarkDelete(b, 100, 32/8)
}

func BenchmarkDelete_32bit_10kItems(b *testing.B) {
	benchmarkDelete(b, 10*1000, 32/8)
}

func BenchmarkDelete_32bit_100kItems(b *testing.B) {
	benchmarkDelete(b, 100*1000, 32/8)
}

func BenchmarkDelete_64bit_SingleItem(b *testing.B) {
	benchmarkDelete(b, 1, 64/8)
}

func BenchmarkDelete_64bit_100Items(b *testing.B) {
	benchmarkDelete(b, 100, 64/8)
}

func BenchmarkDelete_64bit_10kItems(b *testing.B) {
	benchmarkDelete(b, 10*1000, 64/8)
}

func BenchmarkDelete_64bit_100kItems(b *testing.B) {
	benchmarkDelete(b, 100*1000, 64/8)
}

func BenchmarkDelete_128bit_SingleItem(b *testing.B) {
	benchmarkDelete(b, 1, 128/8)
}

func BenchmarkDelete_128bit_100Items(b *testing.B) {
	benchmarkDelete(b, 100, 128/8)
}

func BenchmarkDelete_128bit_10kItems(b *testing.B) {
	benchmarkDelete(b, 10*1000, 128/8)
}

func BenchmarkDelete_128bit_100kItems(b *testing.B) {
	benchmarkDelete(b, 100*1000, 128/8)
}

func BenchmarkDelete_1kbyte_SingleItem(b *testing.B) {
	benchmarkDelete(b, 1, 128/8)
}

func BenchmarkDelete_1kbyte_100Items(b *testing.B) {
	benchmarkDelete(b, 100, 128/8)
}

func BenchmarkDelete_1kbyte_10kItems(b *testing.B) {
	benchmarkDelete(b, 10*1000, 1024)
}

func makeRandomKey(t *testing.B, keyLen int) []byte {
	b := make([]byte, keyLen)
	if _, err := rand.Read(b); err != nil {
		t.Error(err)
		panic(err)
	}
	return b
}
