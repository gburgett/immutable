package flist

import (
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCons_SingleItemToNil_CountIs1(t *testing.T) {

	itemId := makeRandomItemID(t, 8)

	//act
	l := Cons(itemId, NilList())

	//assert
	assert.Equal(t, 1, l.Count(), "Count")
	assert.Equal(t, itemId, l.Head().([]byte), "ItemID")
	assert.False(t, l.IsNil(), "IsNil")

	l = l.Tail()
	assert.True(t, l.IsNil(), "IsNil")
}

func TestCons_MultipleItems_ItemsExistInOrder(t *testing.T) {

	itemIds := make([]interface{}, 8)
	for i := 0; i < len(itemIds); i++ {
		itemIds[i] = makeRandomItemID(t, 8)
	}

	//act
	list := ConsFromSlice(itemIds)

	//assert
	assert.Equal(t, 8, list.Count(), "Count")
	i := 0
	for l := list; !l.IsNil(); l = l.Tail() {
		assert.Equal(t, itemIds[i].([]byte), l.Head().([]byte), fmt.Sprintf("ItemID %d", i))
		i++
	}
	assert.Equal(t, 8, i, "Iterated all")
}

func TestPrepend_NilList_ReturnsOriginalList(t *testing.T) {
	itemIds := make([]interface{}, 8)
	for i := 0; i < len(itemIds); i++ {
		itemIds[i] = makeRandomItemID(t, 8)
	}
	list := ConsFromSlice(itemIds)

	//act
	result := Prepend(NilList(), list)

	//assert
	assert.Equal(t, itemIds, result.ToSlice(), "should be original list")
}

func TestPrepend_HasItems_ReturnsNewlList(t *testing.T) {
	itemIds := make([]interface{}, 8)
	for i := 0; i < len(itemIds); i++ {
		itemIds[i] = makeRandomItemID(t, 8)
	}
	list := ConsFromSlice(itemIds)

	list2 := ConsFromSlice(itemIds[0:2])

	//act
	result := Prepend(list2, list)

	//assert
	expect := append(itemIds[0:2], itemIds...)
	assert.Equal(t, expect, result.ToSlice(), "should be new list with items in order")
	assert.NotEqual(t, list2, result, "Should not be same pointer as original list")
}

func TestFilter_AllFalse_ReturnsNil(t *testing.T) {
	list := NilList()
	for i := 0; i < 17; i++ {
		list = Cons(makeRandomItemID(t, 8), list)
	}

	//act
	list = list.Filter(func(n interface{}) bool { return false })

	//assert
	assert.True(t, list.IsNil(), "Expect nil")
	assert.Nil(t, list.Tail(), "No next")
}

func TestFilter_AllTrue_ReturnsAll(t *testing.T) {
	itemIds := make([]interface{}, 8)
	for i := 0; i < len(itemIds); i++ {
		itemIds[i] = makeRandomItemID(t, 8)
	}

	list := ConsFromSlice(itemIds)

	//act
	list = list.Filter(func(n interface{}) bool { return true })

	//assert
	assert.Equal(t, 8, list.Count(), "Count")
	i := 0
	for l := list; !l.IsNil(); l = l.Tail() {
		assert.Equal(t, itemIds[i].([]byte), l.Head().([]byte), fmt.Sprintf("ItemID %d", i))
		i++
	}
	assert.Equal(t, 8, i, "Iterated all")
}

func TestFilter_FirstHalf_ReturnsFrontHalf(t *testing.T) {
	itemIds := make([]interface{}, 8)
	for i := 0; i < len(itemIds); i++ {
		itemIds[i] = makeRandomItemID(t, 8)
	}

	list := ConsFromSlice(itemIds)

	//act
	i := 0
	list = list.Filter(func(n interface{}) bool { i++; return i <= 4 })

	//assert
	assert.Equal(t, 4, list.Count(), "Count")
	i = 0
	for l := list; !l.IsNil(); l = l.Tail() {
		assert.Equal(t, itemIds[i].([]byte), l.Head().([]byte), fmt.Sprintf("ItemID %d", i))
		i++
	}
	assert.Equal(t, 4, i, "Iterated all")
}

func TestReverse_ReturnsInReverseOrder(t *testing.T) {
	itemIds := make([]interface{}, 27)
	for i := 0; i < len(itemIds); i++ {
		itemIds[i] = makeRandomItemID(t, 8)
	}

	list := ConsFromSlice(itemIds)

	//act
	list = list.Reverse()

	//assert
	assert.Equal(t, 27, list.Count(), "Count")
	i := len(itemIds) - 1
	for l := list; !l.IsNil(); l = l.Tail() {
		assert.Equal(t, itemIds[i].([]byte), l.Head().([]byte), fmt.Sprintf("ItemID %d", i))
		i--
	}
	assert.Equal(t, -1, i, "Iterated all")
}

func TestMap_ZeroingFunction_ReturnsZeroedItems(t *testing.T) {
	itemIds := make([]interface{}, 15)
	for i := 0; i < len(itemIds); i++ {
		itemIds[i] = makeRandomItemID(t, 8)
	}

	list := ConsFromSlice(itemIds)

	//act
	result := list.Map(func(item interface{}) interface{} {
		bytes := make([]byte, 8)
		copy(bytes, item.([]byte))
		for i, _ := range bytes {
			bytes[i] = 0
		}
		return bytes
	})

	//assert
	assert.Equal(t, 15, result.Count(), "count")
	i := 0
	for l := result; !l.IsNil(); l = l.Tail() {
		assert.Equal(t, make([]byte, 8), l.Head(), "value at "+string(i))
		i++
	}
	assert.Equal(t, 15, i, "iterated whole list")
}

func TestAggregate_Sum_ReturnsSum(t *testing.T) {
	items := make([]interface{}, 12)
	for i := 0; i < len(items); i++ {
		items[i] = len(items) - i
	}
	// 12, 11, 10...
	list := ConsFromSlice(items)

	//act
	sum := list.Aggregate(0, func(a, b interface{}) interface{} {
		return a.(int) + b.(int)
	})

	//assert
	assert.Equal(t, 78, sum, "sum")
}

func TestToSlice_ReturnsSliceInOrder(t *testing.T) {
	itemIds := make([]interface{}, 14)
	for i := 0; i < len(itemIds); i++ {
		itemIds[i] = makeRandomItemID(t, 8)
	}

	list := ConsFromSlice(itemIds)

	//act
	slice := list.ToSlice()

	//assert
	assert.Equal(t, 14, len(slice), "Count")
	i := 0
	for ; i < len(itemIds); i++ {
		assert.Equal(t, itemIds[i], slice[i], fmt.Sprintf("ItemID %d", i))
	}

	assert.Equal(t, 14, i, "Iterated all")
}

func TestToChan_ReturnsChanInOrder(t *testing.T) {
	itemIds := make([]interface{}, 14)
	for i := 0; i < len(itemIds); i++ {
		itemIds[i] = makeRandomItemID(t, 8)
	}

	list := ConsFromSlice(itemIds)

	//act
	ch := list.ToChan()

	//assert
	i := 0
	for node := range ch {
		assert.Equal(t, itemIds[i], node, fmt.Sprintf("ItemID %d", i))
		i++
	}

	assert.Equal(t, 14, i, "Iterated all")
}

func makeRandomItemID(t *testing.T, idLen int) []byte {
	b := make([]byte, idLen)
	if _, err := rand.Read(b); err != nil {
		t.Error(err)
	}
	return b
}
