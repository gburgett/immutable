package critbit

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVisitAscend_NilTrie(t *testing.T) {

	//act
	keys := visitToSlice(NilTrie(), nil)

	//assert
	require.Equal(t, 0, len(keys), "len")
}

func TestVisitAscend_Root(t *testing.T) {
	instance, _ := NilTrie().Set([]byte{0x01, 0x02, 0x03}, 123)

	//act
	keys := visitToSlice(instance, nil)

	//assert
	require.Equal(t, 1, len(keys), "len")
	assert.Equal(t, []byte{0x01, 0x02, 0x03}, keys[0])
}

func TestVisitAscend_Root_SkipsLessThanFrom(t *testing.T) {
	instance, _ := NilTrie().Set([]byte{0x01, 0x02, 0x03}, 123)

	//act
	keys := visitToSlice(instance, []byte{0x02})

	//assert
	require.Equal(t, 0, len(keys), "len")
}

func TestVisitAscend_Root_DoesntSkipGreaterThanFrom(t *testing.T) {
	instance, _ := NilTrie().Set([]byte{0x01, 0x02, 0x03}, 123)

	//act
	keys := visitToSlice(instance, []byte{0x01})

	//assert
	require.Equal(t, 1, len(keys), "len")
}

func TestVisitAscend_SecondNodeAtZero(t *testing.T) {
	instance, _ := NilTrie().Set([]byte{0x01, 0x02, 0x03}, 123)
	instance, _ = instance.Set([]byte{0x01, 0x02, 0x02}, 122)

	//act
	keys := visitToSlice(instance, nil)

	//assert
	require.Equal(t, 2, len(keys), "len")
	assert.Equal(t, []byte{0x01, 0x02, 0x02}, keys[0])
	assert.Equal(t, []byte{0x01, 0x02, 0x03}, keys[1])
}

func TestVisitAscend_SecondNodeAtZero_SkipsLessThanFrom(t *testing.T) {
	instance, _ := NilTrie().Set([]byte{0x01, 0x02, 0x03}, 123)
	instance, _ = instance.Set([]byte{0x01, 0x02, 0x02}, 122)

	//act
	keys := visitToSlice(instance, []byte{0x01, 0x02, 0x03})

	//assert
	require.Equal(t, 1, len(keys), "len")
	assert.Equal(t, []byte{0x01, 0x02, 0x03}, keys[0])
}

func TestVisitAscend_SecondNodeAtOne(t *testing.T) {
	instance, _ := NilTrie().Set([]byte{0x01, 0x02, 0x03}, 123)
	instance, _ = instance.Set([]byte{0x01, 0x04, 0x02}, 142)

	//act
	keys := visitToSlice(instance, nil)

	//assert
	require.Equal(t, 2, len(keys), "len")
	assert.Equal(t, []byte{0x01, 0x02, 0x03}, keys[0])
	assert.Equal(t, []byte{0x01, 0x04, 0x02}, keys[1])
}

func TestVisitAscend_ThirdNodeBeforeNode(t *testing.T) {
	instance, _ := NilTrie().Set([]byte{0x01, 0x02, 0x03}, 123)
	instance, _ = instance.Set([]byte{0x01, 0x02, 0x04}, 123)
	instance, _ = instance.Set([]byte{0x01, 0x80, 0x02}, 182)

	//act
	keys := visitToSlice(instance, nil)

	//assert
	require.Equal(t, 3, len(keys), "len")
	assert.Equal(t, []byte{0x01, 0x02, 0x03}, keys[0])
	assert.Equal(t, []byte{0x01, 0x02, 0x04}, keys[1])
	assert.Equal(t, []byte{0x01, 0x80, 0x02}, keys[2])
}

func TestVisitAscend_ThirdNodeAfterNode(t *testing.T) {
	instance, _ := NilTrie().Set([]byte{0x01, 0x02, 0x03}, 123)
	instance, _ = instance.Set([]byte{0x01, 0x01, 0x04}, 123)
	instance, _ = instance.Set([]byte{0x01, 0x02, 0x02}, 182)

	//act
	keys := visitToSlice(instance, nil)

	//assert
	require.Equal(t, 3, len(keys), "len")
	assert.Equal(t, []byte{0x01, 0x01, 0x04}, keys[0])
	assert.Equal(t, []byte{0x01, 0x02, 0x02}, keys[1])
	assert.Equal(t, []byte{0x01, 0x02, 0x03}, keys[2])
}

func TestVisitAscend_PrefixSpecialCaseNode(t *testing.T) {
	instance, _ := NilTrie().Set([]byte{0x01, 0x02, 0x03}, 123)
	instance, _ = instance.Set([]byte{0x01, 0x02}, 12)

	//act
	keys := visitToSlice(instance, nil)

	//assert
	require.Equal(t, 2, len(keys), "len")
	assert.Equal(t, []byte{0x01, 0x02}, keys[0])
	assert.Equal(t, []byte{0x01, 0x02, 0x03}, keys[1])
}

func TestVisitAscend_PrefixSpecialCaseNode_SkipsLessThanFrom(t *testing.T) {
	instance, _ := NilTrie().Set([]byte{0x01, 0x02, 0x03}, 123)
	instance, _ = instance.Set([]byte{0x01, 0x02}, 12)

	//act
	keys := visitToSlice(instance, []byte{0x01, 0x02, 0x00})

	//assert
	require.Equal(t, 1, len(keys), "len")
	assert.Equal(t, []byte{0x01, 0x02, 0x03}, keys[0])
}

func TestVisitAscend_SuffixSpecialCaseNode(t *testing.T) {
	instance, _ := NilTrie().Set([]byte{0x01, 0x02, 0x03}, 123)
	instance, _ = instance.Set([]byte{0x01, 0x02, 0x03, 0x04}, 1234)

	//act
	keys := visitToSlice(instance, nil)

	//assert
	require.Equal(t, 2, len(keys), "len")
	assert.Equal(t, []byte{0x01, 0x02, 0x03}, keys[0])
	assert.Equal(t, []byte{0x01, 0x02, 0x03, 0x04}, keys[1])
}

func TestVisitAscend_LengthLessThanCritbit(t *testing.T) {
	instance, _ := NilTrie().Set([]byte("ffffff"), "f")
	instance, _ = instance.Set([]byte("fffffg"), "g")
	instance, _ = instance.Set([]byte("aaa"), "aaa")

	//act
	keys := visitToSlice(instance, nil)

	//assert
	require.Equal(t, 3, len(keys), "len")
	assert.Equal(t, "aaa", string(keys[0]))
	assert.Equal(t, "ffffff", string(keys[1]))
	assert.Equal(t, "fffffg", string(keys[2]))
}

func TestVisitAscend_StopBeforeEnd(t *testing.T) {
	instance, _ := NilTrie().Set([]byte{0x01, 0x02, 0x03, 0x04}, 1234)
	instance, _ = instance.Set([]byte{0x01, 0x02, 0x03, 0x05}, 1235)
	instance, _ = instance.Set([]byte{0x01, 0x02, 0x03, 0x06}, 1236)

	//		    _[ 0x06 ]
	//		 _ =_
	// root =_   [ 0x05 ]
	//		  [ 0x04 ]

	//act
	vals := make([]interface{}, 0, 2)
	instance.VisitAscend(nil, func(key []byte, val interface{}) bool {
		vals = append(vals, val)
		if bytes.Equal(key, []byte{0x01, 0x02, 0x03, 0x05}) {
			return false
		}
		return true
	})

	//assert
	require.Equal(t, 2, len(vals), "len")
	assert.Equal(t, 1234, vals[0])
	assert.Equal(t, 1235, vals[1])
}

func TestVisitAscend_DeepNodeStillNeedsCompare(t *testing.T) {
	instance, _ := NilTrie().Set([]byte("1"), 1)
	instance, _ = instance.Set([]byte("2"), 2)
	instance, _ = instance.Set([]byte("3"), 3)

	fmt.Println(instance.DumpTrie())

	//act
	vals := make([]interface{}, 0, 1)
	instance.VisitAscend([]byte("3"), func(key []byte, val interface{}) bool {
		vals = append(vals, val)
		return true
	})

	//assert
	require.Equal(t, 1, len(vals), "len")
	assert.Equal(t, vals[0], 3)
}

func TestVisitAscend_FromPrefixOfPrefix(t *testing.T){
	instance, _ := NilTrie().Set([]byte("abcd/abcd1"), "2")
	instance, _ = instance.Set([]byte("abcd/abcd2"), "3")
	instance, _ = instance.Set([]byte("abcdefgh/abcd3"), "4")
	instance, _ = instance.Set([]byte("abc/ghi"), "-")
	instance, _ = instance.Set([]byte("abcd/a"), "1")

	//act
	fmt.Println(instance.DumpTrie())
	keys := visitToSlice(instance, []byte("abcd/"))

	//assert
	assert.Equal(t, "abcd/a", string(keys[0]))
	assert.Equal(t, "abcd/abcd1", string(keys[1]))
	assert.Equal(t, "abcd/abcd2", string(keys[2]))
	require.Equal(t, 4, len(keys), "len")
	assert.Equal(t, "abcdefgh/abcd3", string(keys[3]))
}

func visitToSlice(t *Trie, from []byte) [][]byte {
	ret := make([][]byte, 0, t.Len())
	t.VisitAscend(from, func(key []byte, val interface{}) bool {
		ret = append(ret, key)
		return true
	})
	return ret
}
