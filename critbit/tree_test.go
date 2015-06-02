package critbit

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/rand"
)

func TestSet_Insert_Root(t *testing.T) {
	instance := NilTrie()

	//act
	result, old := instance.Set([]byte{0x01, 0x02, 0x03}, 123)

	//assert
	assert.Nil(t, old, "nothing should exist in tree")
	assert.NotEqual(t, result, instance, "result should be new list")
	assert.Equal(t, 1, result.Len(), "result should have length 1")
	assert.NotNil(t, result.root, "result should have a root node")
	assert.Equal(t, []byte{0x01, 0x02, 0x03}, result.root.key, "root node should be the leaf")
	assert.Equal(t, 123, result.root.value, "root node should be the leaf")

	assert.Equal(t, 0, instance.Len(), "original tree should remain immutable")
	assert.Nil(t, instance.root, "original tree should remain immutable")
}

func TestSet_Insert_SecondNodeAtZero(t *testing.T) {
	instance, _ := NilTrie().Set([]byte{0x01, 0x02, 0x03}, 123)

	//act
	result, old := instance.Set([]byte{0x01, 0x02, 0x02}, 122)

	//assert
	assert.Nil(t, old, "nothing should be overwritten")
	assert.NotEqual(t, result, instance, "result should be new list")
	assert.Equal(t, 2, result.Len(), "result should have length 1")
	assert.Nil(t, result.root.key, "root node should be a critbit node")
	assert.Equal(t, 2, result.root.critbyte, "root critbyte")
	assert.Equal(t, 0xFE, result.root.critbit, "root critbit")
	assert.Equal(t, []byte{0x01, 0x02, 0x03}, result.root.children[1].key, "first child key")
	assert.Equal(t, []byte{0x01, 0x02, 0x02}, result.root.children[0].key, "zero child key")

	assert.Equal(t, 1, instance.Len(), "original tree should remain immutable")
	assert.Equal(t, []byte{0x01, 0x02, 0x03}, instance.root.key, "original tree should remain immutable")
}

func TestSet_Insert_SecondNodeAtOne(t *testing.T) {
	instance, _ := NilTrie().Set([]byte{0x01, 0x02, 0x03}, 123)

	//act
	result, old := instance.Set([]byte{0x01, 0x04, 0x02}, 142)

	//assert
	assert.Nil(t, old, "nothing should be overwritten")
	assert.NotEqual(t, result, instance, "result should be new list")
	assert.Equal(t, 2, result.Len(), "result should have length 1")
	assert.Nil(t, result.root.key, "root node should be a critbit node")
	assert.Equal(t, 1, result.root.critbyte, "root critbyte")
	assert.Equal(t, 0xFB, result.root.critbit, "root critbit")
	assert.Equal(t, []byte{0x01, 0x04, 0x02}, result.root.children[1].key, "first child key")
	assert.Equal(t, []byte{0x01, 0x02, 0x03}, result.root.children[0].key, "zero child key")

	assert.Equal(t, 1, instance.Len(), "original tree should remain immutable")
	assert.Equal(t, []byte{0x01, 0x02, 0x03}, instance.root.key, "original tree should remain immutable")
}

func TestSet_Insert_ThirdNodeBeforeNode(t *testing.T) {
	instance, _ := NilTrie().Set([]byte{0x01, 0x02, 0x03}, 123)
	instance, _ = instance.Set([]byte{0x01, 0x02, 0x04}, 123)

	//act
	result, old := instance.Set([]byte{0x01, 0x80, 0x02}, 182)

	//assert
	assert.Nil(t, old, "nothing should be overwritten")
	assert.NotEqual(t, result, instance, "result should be new list")
	assert.Equal(t, 3, result.Len(), "result should have length 1")
	assert.Nil(t, result.root.key, "root node should be a critbit node")
	assert.Equal(t, 1, result.root.critbyte, "root critbyte")
	assert.Equal(t, 0x7F, result.root.critbit, "root critbit")
	assert.Equal(t, []byte{0x01, 0x80, 0x02}, result.root.children[1].key, "first child key")
	assert.Nil(t, result.root.children[0].key, "zero child should be node")

}

func TestSet_Insert_ThirdNodeAfterNode(t *testing.T) {
	instance, _ := NilTrie().Set([]byte{0x01, 0x02, 0x03}, 123)
	instance, _ = instance.Set([]byte{0x01, 0x01, 0x04}, 123)

	//act
	result, old := instance.Set([]byte{0x01, 0x02, 0x02}, 182)

	//assert
	assert.Nil(t, old, "nothing should be overwritten")
	assert.NotEqual(t, result, instance, "result should be new list")
	assert.Equal(t, 3, result.Len(), "result should have length 1")
	assert.Nil(t, result.root.key, "root node should be a critbit node")
	assert.Equal(t, 1, result.root.critbyte, "root critbyte")
	assert.Equal(t, 0xFD, result.root.critbit, "root critbit")
	assert.Equal(t, []byte{0x01, 0x01, 0x04}, result.root.children[0].key, "0 child key")
	assert.Nil(t, result.root.children[1].key, "1 child should be node")
	assert.Equal(t, 2, result.root.children[1].critbyte, "child critbyte")
	assert.Equal(t, 0xFE, result.root.children[1].critbit, "child critbit")

}

func TestSet_Insert_BestLeafCritbitAfterSpecialCaseNode(t *testing.T) {
	instance, _ := NilTrie().Set([]byte{0x65, 0x7b, 0x1b}, 1)
	instance, _ = instance.Set([]byte{0x65}, 2)
	instance, _ = instance.Set([]byte{0x65, 0x03, 0xec, 0x04}, 3)

	//act
	got, ok := instance.Get([]byte{0x65})

	//assert
	require.True(t, ok)
	assert.Equal(t, 2, got.(int))
}

func TestSet_OverwriteRoot_ReturnsOriginal(t *testing.T) {
	instance, _ := NilTrie().Set([]byte{0x01, 0x02, 0x03}, 123)

	//act
	result, old := instance.Set([]byte{0x01, 0x02, 0x03}, 124)

	//assert
	assert.Equal(t, 123, old.(int), "old")
	assert.NotEqual(t, instance, result, "should return new tree")
	assert.Equal(t, 124, result.root.value.(int))
	assert.Equal(t, 1, result.Len(), "len")
}

func TestSet_OverwriteDeep_ReturnsOriginal(t *testing.T) {
	instance, _ := NilTrie().Set([]byte{0x01, 0x02, 0x03}, 123)
	instance, _ = instance.Set([]byte{0x01, 0x01, 0x04}, 123)
	instance, _ = instance.Set([]byte{0x01, 0x02, 0x02}, 123)

	//act
	result, old := instance.Set([]byte{0x01, 0x02, 0x03}, 124)

	//assert
	assert.Equal(t, 123, old.(int), "old")
	assert.NotEqual(t, instance, result, "should return new tree")
	assert.Equal(t, 124, result.root.children[1].children[1].value.(int))
	assert.Equal(t, 3, result.Len(), "len")
}

func TestSet_Prefix_CreatesSpecialCaseNode(t *testing.T) {
	instance, _ := NilTrie().Set([]byte{0x01, 0x02, 0x03}, 123)

	//act
	result, old := instance.Set([]byte{0x01, 0x02}, 12)

	//assert
	assert.Nil(t, old, "old")
	assert.NotEqual(t, instance, result, "should return new tree")
	assert.Equal(t, 12, result.root.children[0].value.(int))
	assert.Equal(t, 123, result.root.children[1].value.(int))
	assert.Equal(t, 2, result.Len(), "len")
}

func TestSet_Suffix_CreatesSpecialCaseNode(t *testing.T) {
	instance, _ := NilTrie().Set([]byte{0x01, 0x02, 0x03}, 123)

	//act
	result, old := instance.Set([]byte{0x01, 0x02, 0x03, 0x04}, 1234)

	//assert
	assert.Nil(t, old, "old")
	assert.NotEqual(t, instance, result, "should return new tree")
	assert.Equal(t, 123, result.root.children[0].value.(int))
	assert.Equal(t, 1234, result.root.children[1].value.(int))
	assert.Equal(t, 2, result.Len(), "len")
}

func TestSet_LengthLessThanCritbit_InsertsAheadOfNode(t *testing.T) {
	instance, _ := NilTrie().Set([]byte("ffffff"), "f")
	instance, _ = instance.Set([]byte("fffffg"), "g")

	//act
	fmt.Println("inserting...\n")
	instance, old := instance.Set([]byte("aaa"), "aaa")

	//assert
	assert.Nil(t, old, "nothing should be overwritten")
	got, ok := instance.Get([]byte("aaa"))
	require.True(t, ok)
	assert.NotNil(t, got)
	assert.Equal(t, "aaa", got)
}

func TestSet_NilValue_Panics(t *testing.T) {
	instance := NilTrie()

	defer func() {
		err := recover()
		if err == nil {
			assert.Fail(t, "should have panicked")
		}
	}()

	//act
	_, _ = instance.Set([]byte{0x01}, nil)

	//assert
	assert.Fail(t, "should have panicked")
}

func TestGet_NilTrie_ReturnsNothing(t *testing.T) {

	instance := NilTrie()

	//act
	result, ok := instance.Get([]byte{0x01, 0x02, 0x03})

	//assert
	assert.False(t, ok)
	assert.Nil(t, result)
}

func TestGet_SingleNode_Gets(t *testing.T) {

	instance, _ := NilTrie().Set([]byte{0x01, 0x02, 0x03}, 123)

	//act
	result, ok := instance.Get([]byte{0x01, 0x02, 0x03})

	//assert
	assert.True(t, ok)
	assert.Equal(t, 123, result)
}

func TestGet_SingleNode_Misses(t *testing.T) {

	instance, _ := NilTrie().Set([]byte{0x01, 0x02, 0x04}, 124)

	//act
	result, ok := instance.Get([]byte{0x01, 0x02, 0x03})

	//assert
	assert.False(t, ok)
	assert.Nil(t, result)
}

func TestGet_Deep_Gets(t *testing.T) {

	instance, _ := NilTrie().Set([]byte{0x01, 0x02, 0x03}, 123)
	instance, _ = instance.Set([]byte{0x01, 0x01, 0x04}, 114)
	instance, _ = instance.Set([]byte{0x01, 0x02, 0x02}, 122)

	//act
	result, ok := instance.Get([]byte{0x01, 0x02, 0x02})

	//assert
	assert.True(t, ok)
	assert.Equal(t, 122, result)
}

func TestGet_Deep_Misses(t *testing.T) {

	instance, _ := NilTrie().Set([]byte{0x01, 0x02, 0x03}, 123)
	instance, _ = instance.Set([]byte{0x01, 0x01, 0x04}, 114)
	instance, _ = instance.Set([]byte{0x01, 0x02, 0x02}, 122)

	//act
	result, ok := instance.Get([]byte{0x01, 0x01, 0x02})

	//assert
	assert.False(t, ok)
	assert.Nil(t, result)
}

func TestGet_Prefix_Gets(t *testing.T) {
	instance, _ := NilTrie().Set([]byte{0x01, 0x02, 0x03}, 123)
	instance, _ = instance.Set([]byte{0x01, 0x02}, 12)

	//act
	result, ok := instance.Get([]byte{0x01, 0x02})

	//assert
	assert.True(t, ok)
	assert.Equal(t, 12, result.(int))
}

func TestGet_Suffix_Gets(t *testing.T) {
	instance, _ := NilTrie().Set([]byte{0x01, 0x02, 0x03}, 123)
	instance, _ = instance.Set([]byte{0x01, 0x02}, 12)

	//act
	result, ok := instance.Get([]byte{0x01, 0x02, 0x03})

	//assert
	assert.True(t, ok)
	assert.Equal(t, 123, result.(int))
}

func TestDelete_NilTrie_Fails(t *testing.T) {
	instance := NilTrie()

	//act
	result, was := instance.Delete([]byte{0x01, 0x02})

	//assert
	assert.Equal(t, result, instance, "should return same tree")
	assert.Nil(t, was, "should have deleted nothing")
}

func TestDelete_Root_Success(t *testing.T) {
	instance, _ := NilTrie().Set([]byte{0x01, 0x02, 0x03}, 123)

	//act
	result, was := instance.Delete([]byte{0x01, 0x02, 0x03})

	//assert
	assert.NotEqual(t, result, instance, "should make new tree")
	assert.Equal(t, 123, was.(int), "should return old value")

	_, ok := result.Get([]byte{0x01, 0x02, 0x03})
	assert.False(t, ok, "tree should no longer contain value")
	assert.Equal(t, 0, result.Len(), "len")
	_, ok = instance.Get([]byte{0x01, 0x02, 0x03})
	assert.True(t, ok, "expect immutability")
}

func TestDelete_SingleNode_Fails(t *testing.T) {
	instance, _ := NilTrie().Set([]byte{0x01, 0x02, 0x03}, 123)

	//act
	result, was := instance.Delete([]byte{0x01, 0x02, 0x04})

	//assert
	assert.Equal(t, result, instance, "should return same tree")
	assert.Nil(t, was, "should return no value")
	assert.Equal(t, 1, result.Len(), "len")

	_, ok := instance.Get([]byte{0x01, 0x02, 0x03})
	assert.True(t, ok, "item should remain")
}

func TestDelete_Deep_Success(t *testing.T) {
	instance, _ := NilTrie().Set([]byte{0x01, 0x02, 0x03}, 123)
	instance, _ = instance.Set([]byte{0x01, 0x02, 0x04}, 124)
	instance, _ = instance.Set([]byte{0x01, 0x03, 0x03}, 133)

	//act
	result, was := instance.Delete([]byte{0x01, 0x02, 0x04})

	//assert
	assert.NotEqual(t, result, instance, "should make new tree")
	assert.Equal(t, 124, was.(int), "should return old value")

	_, ok := result.Get([]byte{0x01, 0x02, 0x04})
	assert.False(t, ok, "tree should no longer contain value")
	assert.Equal(t, 2, result.Len(), "len")
	_, ok = instance.Get([]byte{0x01, 0x02, 0x04})
	assert.True(t, ok, "expect immutability")

	//node structure
	assert.Equal(t, 123, result.root.children[0].value.(int), "123")
	assert.Equal(t, 133, result.root.children[1].value.(int), "133")
}

func TestDelete_Deep_Failure(t *testing.T) {
	instance, _ := NilTrie().Set([]byte{0x01, 0x02, 0x03}, 123)
	instance, _ = instance.Set([]byte{0x01, 0x02, 0x04}, 124)
	instance, _ = instance.Set([]byte{0x01, 0x03, 0x03}, 133)

	//act
	result, was := instance.Delete([]byte{0x01, 0x02, 0x01})

	//assert
	assert.Equal(t, result, instance, "should return same tree")
	assert.Nil(t, was, "should return no value")
}

func TestDelete_Prefix_Success(t *testing.T) {
	instance, _ := NilTrie().Set([]byte{0x01, 0x02, 0x03}, 123)
	instance, _ = instance.Set([]byte{0x01, 0x02}, 12)

	//act
	result, was := instance.Delete([]byte{0x01, 0x02})

	//assert
	assert.NotEqual(t, result, instance, "should make new tree")
	assert.Equal(t, 12, was.(int), "should return old value")

	_, ok := result.Get([]byte{0x01, 0x02})
	assert.False(t, ok, "tree should no longer contain value")
	assert.Equal(t, 1, result.Len(), "len")
	_, ok = instance.Get([]byte{0x01, 0x02})
	assert.True(t, ok, "expect immutability")

	//node structure
	assert.Equal(t, 123, result.root.value.(int), "123")
}

func TestDelete_Suffix_Success(t *testing.T) {
	instance, _ := NilTrie().Set([]byte{0x01, 0x02, 0x03}, 123)
	instance, _ = instance.Set([]byte{0x01, 0x02}, 12)

	//act
	result, was := instance.Delete([]byte{0x01, 0x02, 0x03})

	//assert
	assert.NotEqual(t, result, instance, "should make new tree")
	assert.Equal(t, 123, was.(int), "should return old value")

	_, ok := result.Get([]byte{0x01, 0x02, 0x03})
	assert.False(t, ok, "tree should no longer contain value")
	assert.Equal(t, 1, result.Len(), "len")
	_, ok = instance.Get([]byte{0x01, 0x02, 0x03})
	assert.True(t, ok, "expect immutability")

	//node structure
	assert.Equal(t, 12, result.root.value.(int), "123")
}

//-- Internal functions --//

func TestFindCritbit_LowestBit(t *testing.T) {
	u := []byte{0x00, 0x00, 0x00, 0x00}
	p := []byte{0x01, 0x00, 0x00, 0x00}

	//act
	critbyte, critbit := findCritbit(u, p)

	//assert
	assert.Equal(t, 0, critbyte, "critbyte should be initial byte")
	assert.Equal(t, 0xFE, critbit, "1111 1110")
}

func TestFindCritbit_HighestBit(t *testing.T) {
	u := []byte{0x27, 0x00, 0x00, 0x00} // 0000 0011
	p := []byte{0x27, 0x80, 0x00, 0x00} // 0000 0000

	//act
	critbyte, critbit := findCritbit(u, p)

	//assert
	assert.Equal(t, 1, critbyte, "critbyte should be second byte")
	assert.Equal(t, 0x7F, critbit, "1101 1111")
}

func TestFindCritbit_AllBits(t *testing.T) {
	for i := uint8(0); i < 8; i++ {
		for j := 0x01 << i; j < 0x01<<(i+1); j++ {
			//test every possible bitmask combination having critbit at position i
			u := []byte{0x01 << i}
			p := []byte{0x00}

			//act
			_, critbit := findCritbit(u, p)
			_, critbit2 := findCritbit(p, u)

			//assert
			assert.Equal(t, ^(uint8(0x01) << i), critbit, fmt.Sprintf("u: %x p: %x index: %d", u[0], p[0], i))
			assert.Equal(t, ^(uint8(0x01) << i), critbit2, fmt.Sprintf("u: %x p: %x index: %d", u[0], p[0], i))
		}
	}
}

func TestFindCritbit_PIsLonger_CritbitIs255(t *testing.T) {
	u := []byte{0x27, 0x00, 0x00, 0x94}
	p := []byte{0x27, 0x00, 0x00}

	//act
	critbyte, critbit := findCritbit(u, p)

	//assert
	assert.Equal(t, 2, critbyte, "critbyte should be index of last byte in p")
	assert.Equal(t, 0xFF, critbit, "255")
}

func TestFindCritbit_UIsLonger_CritbitIs255(t *testing.T) {
	u := []byte{0x27, 0x00, 0x00, 0x94}
	p := []byte{0x27, 0x00, 0x00, 0x94, 0x00}

	//act
	critbyte, critbit := findCritbit(u, p)

	//assert
	assert.Equal(t, 3, critbyte, "critbyte should be index of last byte in u")
	assert.Equal(t, 0xFF, critbit, "255")
}

func TestFindDirection_InitialByte_ComparesCorrectBit(t *testing.T) {
	u := []byte{0x80, 0x00}
	p := []byte{0x00, 0x00}

	critbyte := 0
	critbit := uint8(0x7F)

	//act
	directionU := findDirection(u, critbyte, critbit)
	directionP := findDirection(p, critbyte, critbit)

	//assert
	assert.Equal(t, 1, directionU, "0x80")
	assert.Equal(t, 0, directionP, "0x00")
}

func TestFindDirection_LastByte_ComparesCorrectBit(t *testing.T) {
	u := []byte{0x00, 0x01}
	p := []byte{0x00, 0x80}

	critbyte := 1
	critbit := uint8(0x01)

	//act
	directionU := findDirection(u, critbyte, critbit)
	directionP := findDirection(p, critbyte, critbit)

	//assert
	assert.Equal(t, 1, directionU, "0x01")
	assert.Equal(t, 0, directionP, "0x80")
}

func TestFindDirection_AllBits_ComparesCorrectBit(t *testing.T) {
	u := []byte{0xFF}
	p := []byte{0x00}

	for i := uint8(0); i < 8; i++ {
		critbyte := 0
		critbit := ^(uint8(0x01) << i)

		//act
		directionU := findDirection(u, critbyte, critbit)
		directionP := findDirection(p, critbyte, critbit)

		//assert
		assert.Equal(t, 1, directionU, "0x01")
		assert.Equal(t, 0, directionP, "0x80")

	}
}

func TestInsert_RandomBytes_DoesNotFail(t *testing.T) {
	var number = 1000
	keys := make([][]byte, number)
	for i := 0; i < number-1; i++ {
		keys[i] = randBytes()
	}
	keys[number-1] = []byte{0x8d}

	tree := NilTrie()

	var current []byte
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Error with %x! %v\n%s", current, err, tree.DumpTrie())
			panic(err)
		}
	}()
	for i := 0; i < number; i++ {
		current = keys[i]
		tree, _ = tree.Set(current, i)
	}
	for i := number - 1; i >= 0; i-- {
		current = keys[i]
		_, ok := tree.Get(current)
		if !ok {
			assert.Fail(t, fmt.Sprintf("Could not get %x.  Tree: \n%s", current, tree.DumpTrie()))
		}
	}

	snapshot := tree

	seen := make(map[string]bool)
	for i := 0; i < number; i++ {
		current = keys[i]
		if seen[string(current)] {
			continue
		}
		seen[string(current)] = true
		var was interface{}
		tree, was = tree.Delete(current)
		if was == nil {
			assert.Fail(t, fmt.Sprintf("deleted empty value for key [%x].  Tree: \n %s", current, tree.DumpTrie()))
		}
	}
	assert.Equal(t, 0, tree.Len())

	expect := len(seen)
	count := 0
	snapshot.VisitAscend(nil, func(key []byte, val interface{}) bool {
		count++
		if !seen[string(key)] {
			assert.Fail(t, fmt.Sprintf("visited key that was never seen: [%x].  Tree: \n%s", current, tree.DumpTrie()))
		}
		delete(seen, string(key))
		return true
	})
	assert.Equal(t, expect, count)
}

func (t *Trie) DumpTrie() string {
	if t.root == nil {
		return "empty"
	}
	return t.root.DumpNode("")
}

func (n *node) DumpNode(prefix string) string {
	head := fmt.Sprintf("[%x](%d %x)\n", n.key, n.critbyte, n.critbit)

	p2 := prefix + "  "
	if n.children[0] != nil {
		ch0 := fmt.Sprintf("%s[0]: %s", p2, n.children[0].DumpNode(p2))
		ch1 := fmt.Sprintf("%s[1]: %s", p2, n.children[1].DumpNode(p2))
		return head + ch0 + ch1
	}
	return head
}

func randBytes() []byte {
	bytes := make([]byte, rand.Intn(32))
	for i := 0; i < len(bytes); i++ {
		bytes[i] = byte(rand.Intn(256))
	}
	return bytes
}
