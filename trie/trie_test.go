package trie

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testInterface interface{}

type testItem struct {
	value string
}

func Test_Set_EmptyKey_AddedAtRoot(t *testing.T) {
	key := make([]byte, 0)
	value := &testItem{
		value: "hi",
	}

	//act
	result, old := Set(makeTrieSimple([]byte{0x08, 0x08, 0x01}), key, value)

	//assert
	assert.NotNil(t, result.root, "root")
	assert.Equal(t, value, result.root.value, "value")
	assert.Nil(t, old, "old")
	assert.Equal(t, 2, result.Len(), "len")

	assert.NotNil(t, result.root.children[0x08], "child")
	assert.Equal(t, "[08 08 01]", result.root.children[0x08].value.(*testItem).value, "child.value")
}

func Test_Set_EmptyKey_AddedAtRoot_Immutable(t *testing.T) {
	key := make([]byte, 0)
	value := &testItem{
		value: "hi",
	}
	tree := NilTrie()

	//act
	result, old := Set(tree, key, value)

	//assert
	assert.NotNil(t, result)
	assert.Nil(t, old, "old")

	assert.NotNil(t, tree.root, "root")
	assert.Nil(t, tree.root.value, "value")
	assert.Equal(t, 0, tree.Len(), "len")
}

func Test_Set_DeepKey_AddedAsLeaf(t *testing.T) {
	key := []byte{128, 128}
	value := &testItem{
		value: "hi 128",
	}
	tree := NilTrie()

	//act
	result, old := Set(tree, key, value)

	//assert
	assert.NotNil(t, result, "result")
	assert.Nil(t, old, "old")
	assert.Equal(t, 1, result.Len(), "len")

	//root
	assert.NotNil(t, result.root, "root")
	assert.Nil(t, result.root.value, "root.value")

	//leaf
	leaf := result.root.children[128]
	assert.NotNil(t, leaf, "leaf")
	assert.Equal(t, value, leaf.value, "leaf.value")
}

func Test_Set_DeepKey_AddedAsLeaf_Immutable(t *testing.T) {
	key := []byte{128, 128}
	value := &testItem{
		value: "hi 128",
	}
	tree := NilTrie()

	//act
	result, old := Set(tree, key, value)

	//assert
	assert.NotNil(t, result)
	assert.Nil(t, old, "old")

	assert.NotNil(t, tree.root, "root")
	assert.Nil(t, tree.root.value, "root.value")
	assert.Equal(t, 0, len(tree.root.children), "root.children")
	assert.Equal(t, 0, tree.Len(), "len")
}

func Test_Set_SecondDeepKey_AddedAsLeaf(t *testing.T) {
	tree := makeTrieSimple([]byte{128, 128, 128})

	key := []byte{128, 128, 127}
	value := &testItem{
		value: "hi 129",
	}

	//act
	result, old := Set(tree, key, value)

	//assert
	assert.NotNil(t, result, "result")
	assert.Nil(t, old, "old")
	assert.Equal(t, 2, result.Len(), "len")

	//root (key length 0)
	assert.NotNil(t, result.root, "root")
	log.Printf(result.root.printDbg(""))
	assert.Nil(t, result.root.value, "root.value")

	//node (key length 2)
	node := result.root.children[128]
	assert.NotNil(t, node, "node")
	assert.Nil(t, node.value, "node.value")
	assert.Equal(t, node.keySlice, []byte{128, 128}, node.keySlice)

	//original leaf
	leaf := node.children[128]
	assert.NotNil(t, leaf, "original")
	assert.Equal(t, "[80 80 80]", leaf.value.(*testItem).value, "original.value")

	//new leaf
	leaf = node.children[127]
	assert.NotNil(t, leaf, "new")
	assert.Equal(t, "hi 129", leaf.value.(*testItem).value, "new.value")
}

func Test_Set_SecondDeepKey_AddedAsLeaf_Immutable(t *testing.T) {
	tree := makeTrieSimple([]byte{128, 128, 128})

	key := []byte{128, 128, 127}
	value := &testItem{
		value: "hi 129",
	}

	//act
	result, old := Set(tree, key, value)

	//assert
	assert.NotNil(t, result)
	assert.Nil(t, old, "old")
	assert.Equal(t, 1, tree.Len(), "len")

	//root (key length 0)
	assert.NotNil(t, tree.root, "root")
	log.Printf(tree.root.printDbg(""))
	assert.Nil(t, tree.root.value, "root.value")
	assert.Equal(t, 1, len(tree.root.children), "root.children")
	assert.Equal(t, 1, len(tree.root.children), "root.children")

	//original leaf
	leaf := tree.root.children[128]
	assert.NotNil(t, leaf, "original")
	assert.Equal(t, "[80 80 80]", leaf.value.(*testItem).value, "original.value")
	assert.Equal(t, 0, len(leaf.children), "leaf.children")
	assert.Equal(t, 0, len(leaf.children), "leaf.children")

}

func Test_Set_LongerKey_LeafOfLeaf(t *testing.T) {
	tree := makeTrieSimple([]byte{32, 128, 128})

	key := []byte{32, 128, 128, 127}
	value := &testItem{
		value: "hi 127",
	}

	//act
	result, old := Set(tree, key, value)

	//assert
	assert.NotNil(t, result, "result")
	assert.Nil(t, old, "old")
	assert.Equal(t, 2, result.Len(), "len")

	//root (key length 0)
	assert.NotNil(t, result.root, "root")
	log.Printf(result.root.printDbg(""))
	assert.Nil(t, result.root.value, "root.value")

	//original leaf, which is also a node
	leaf := result.root.children[32]
	node := leaf
	assert.NotNil(t, leaf, "original")
	assert.Equal(t, "[20 80 80]", leaf.value.(*testItem).value, "original.value")

	//new leaf
	leaf = node.children[127]
	assert.NotNil(t, leaf, "new")
	assert.Equal(t, "hi 127", leaf.value.(*testItem).value, "new.value")
}

func Test_Set_LongerKey_LeafOfLeaf_Immutable(t *testing.T) {
	tree := makeTrieSimple([]byte{32, 128, 128})

	key := []byte{32, 128, 128, 127}
	value := &testItem{
		value: "hi 127",
	}

	//act
	result, _ := Set(tree, key, value)

	//assert
	assert.NotNil(t, result, "result")
	assert.Equal(t, 1, tree.Len(), "len")

	//root (key length 0)
	assert.NotNil(t, tree.root, "root")
	log.Printf(tree.root.printDbg(""))
	assert.Nil(t, tree.root.value, "root.value")
	assert.Equal(t, 1, len(tree.root.children), "root.children")
	assert.Equal(t, 1, len(tree.root.children), "root.children")

	//original leaf, which is also a node
	leaf := tree.root.children[32]
	assert.NotNil(t, leaf, "original")
	assert.Equal(t, "[20 80 80]", leaf.value.(*testItem).value, "original.value")
	assert.Equal(t, 0, len(leaf.children), "root.children")
	assert.Equal(t, 0, len(leaf.children), "root.children")
}

func Test_Set_ShorterKey_LeafOfLeaf(t *testing.T) {
	tree := makeTrieSimple([]byte{128, 128, 128, 127})

	key := []byte{128, 128, 128}
	value := &testItem{
		value: "hi 127",
	}

	//act
	result, old := Set(tree, key, value)

	//assert
	assert.NotNil(t, result, "result")
	assert.Nil(t, old, "old")
	assert.Equal(t, 2, result.Len(), "len")

	//root (key length 0)
	assert.NotNil(t, result.root, "root")
	log.Printf(result.root.printDbg(""))
	assert.Nil(t, result.root.value, "root.value")

	//new leaf, which is also a node
	leaf := result.root.children[128]
	node := leaf
	assert.NotNil(t, leaf, "original")
	assert.Equal(t, "hi 127", leaf.value.(*testItem).value, "original.value")

	//original leaf
	leaf = node.children[127]
	assert.NotNil(t, leaf, "new")
	assert.Equal(t, "[80 80 80 7f]", leaf.value.(*testItem).value, "new.value")
}

func Test_Set_ShorterKey_LeafOfLeaf_Immutable(t *testing.T) {
	tree := makeTrieSimple([]byte{128 + 64, 128, 128, 127})

	key := []byte{128 + 64, 128, 128}
	value := &testItem{
		value: "hi 127",
	}

	//act
	result, _ := Set(tree, key, value)

	//assert
	assert.NotNil(t, result, "result")
	assert.Equal(t, 1, tree.Len(), "len")

	//root (key length 0)
	assert.NotNil(t, tree.root, "root")
	log.Printf(tree.root.printDbg(""))
	assert.Nil(t, tree.root.value, "root.value")
	assert.Equal(t, 1, len(tree.root.children), "root.children")
	assert.Equal(t, 1, len(tree.root.children), "root.children")

	//original leaf
	leaf := tree.root.children[0xc0]
	assert.NotNil(t, leaf, "original")
	assert.Equal(t, "[c0 80 80 7f]", leaf.value.(*testItem).value, "original.value")
	assert.Equal(t, 0, len(leaf.children), "leaf.children")
	assert.Equal(t, 0, len(leaf.children), "leaf.children")
}

func Test_Set_ShorterKey_LeafOfRoot(t *testing.T) {
	tree := makeTrieSimple([]byte{128, 129, 128, 127})

	key := []byte{128}
	value := &testItem{
		value: "hi 128",
	}

	//act
	result, old := Set(tree, key, value)

	//assert
	assert.NotNil(t, result, "result")
	assert.Nil(t, old, "old")
	assert.Equal(t, 2, result.Len(), "len")

	//root (key length 0)
	assert.NotNil(t, result.root, "root")
	log.Printf(result.root.printDbg(""))
	assert.Nil(t, result.root.value, "root.value")

	//new leaf, which is also a node
	leaf := result.root.children[128]
	node := leaf
	assert.NotNil(t, leaf, "original")
	assert.Equal(t, "hi 128", leaf.value.(*testItem).value, "original.value")

	//original leaf
	leaf = node.children[129]
	assert.NotNil(t, leaf, "new")
	assert.Equal(t, "[80 81 80 7f]", leaf.value.(*testItem).value, "new.value")
}

func Test_Set_MultinodeChild_RecurseInto(t *testing.T) {
	/*
		{
			key: root,
			children: {
				80: {
					key: [80],
					children: {
						80: {
							key: [80 80]
							value: '[80 80]'
							children: {}
						},
						<--- new node goes here
						40: {
							key: [80 40]
							value: '[80 40]'
							children: {}
						}
					}
				}
			}
		}*/
	tree := makeTrieSimple([]byte{0x80, 0x80}, []byte{0x80, 0x40})

	key := []byte{0x80, 0x60}
	value := &testItem{
		value: "hi 0x60",
	}

	//act
	result, old := Set(tree, key, value)

	//assert
	assert.NotNil(t, result, "result")
	assert.Nil(t, old, "old")
	assert.Equal(t, 3, result.Len(), "len")

	log.Printf(result.root.printDbg(""))

	//node
	node := result.root.children[0x80]
	assert.NotNil(t, node, "node")

	//original children
	assert.NotNil(t, node.children[0x40], "node.0x40")
	assert.Equal(t, "[80 40]", node.children[0x40].value.(*testItem).value, "node.0x40.value")
	assert.NotNil(t, node.children[0x80], "node.0x80")
	assert.Equal(t, "[80 80]", node.children[0x80].value.(*testItem).value, "node.0x80.value")

	//new child
	assert.NotNil(t, node.children[0x60], "new leaf")
	assert.Equal(t, "hi 0x60", node.children[0x60].value.(*testItem).value, "new leaf.value")
}

func Test_Set_RecurseInto_SplitNode(t *testing.T) {
	/*10: {
	  key: [10 47],
	  children: {
	    9: {
	        key: [10 47 09 1c],
	        value: 0,
	        children: {
	        },
	      },
	    e0: {
	        key: [10 47 e0 54],
	        value: 0,
	        children: {
	        },
	      },
	  },
	},*/
	tree := makeTrieSimple([]byte{0x10, 0x47, 0x09, 0x1c}, []byte{0x10, 0x47, 0xe0, 0x54}, []byte{0x32, 0x00, 0x00, 0x00})
	key := []byte{0x10, 0x46, 0xdf, 0xf0}

	//act
	result, _ := Set(tree, key, &testItem{
		value: "new item",
	})

	//assert
	got, ok := result.Get(key)
	assert.True(t, ok, "exists")
	assert.Equal(t, "new item", got.(*testItem).value, "value")
}

func Test_Set_Existing_Replaced(t *testing.T) {
	tree := makeTrieSimple([]byte{0x80, 0x80, 0x80, 0x92}, []byte{0x80, 0x80, 0x80, 0x91}, []byte{0x40, 0x80, 0x32})

	key := []byte{0x80, 0x80, 0x80, 0x92}
	value := &testItem{
		value: "hi replaced",
	}

	//act
	result, old := Set(tree, key, value)

	//assert
	assert.NotNil(t, result, "result")
	assert.Equal(t, "[80 80 80 92]", old.(*testItem).value, "old")
	log.Printf(result.root.printDbg(""))
	assert.Equal(t, 3, result.Len(), "len")

	//leaf
	leaf := result.root.children[0x80].children[0x92]
	assert.NotNil(t, leaf, "leaf")
	assert.Equal(t, "hi replaced", leaf.value.(*testItem).value, "should have replaced value")

	//existing
	leaf = result.root.children[0x40]
	assert.NotNil(t, leaf, "existing 0x40")
	assert.Equal(t, "[40 80 32]", leaf.value.(*testItem).value, "existing 0x40")

	leaf = result.root.children[0x80].children[0x91]
	assert.NotNil(t, leaf, "existing 0x91")
	assert.Equal(t, "[80 80 80 91]", leaf.value.(*testItem).value, "existing 91")
}

func Test_Set_Existing_Replaced_Immutable(t *testing.T) {
	tree := makeTrieSimple([]byte{0x80, 0x80, 0x80, 0x92}, []byte{0x80, 0x80, 0x80, 0x91}, []byte{0x40, 0x80, 0x32})

	key := []byte{0x80, 0x80, 0x80, 0x92}
	value := &testItem{
		value: "hi replaced",
	}

	//act
	result, _ := Set(tree, key, value)

	//assert
	assert.NotNil(t, result, "result")
	log.Printf(tree.root.printDbg(""))
	assert.Equal(t, 3, tree.Len(), "len")

	//leaf
	leaf := tree.root.children[0x80].children[0x92]
	assert.NotNil(t, leaf, "leaf")
	assert.Equal(t, "[80 80 80 92]", leaf.value.(*testItem).value, "should have replaced value")

	//existing
	leaf = tree.root.children[0x40]
	assert.NotNil(t, leaf, "existing 0x40")
	assert.Equal(t, "[40 80 32]", leaf.value.(*testItem).value, "existing 0x40")

	leaf = tree.root.children[0x80].children[0x91]
	assert.NotNil(t, leaf, "existing 0x91")
	assert.Equal(t, "[80 80 80 91]", leaf.value.(*testItem).value, "existing 91")
}

func Test_Get_EmptyKey_Root(t *testing.T) {
	tree := makeTrieSimple([]byte{})

	//act
	result, ok := tree.Get([]byte{})

	//assert
	assert.True(t, ok, "ok")
	assert.Equal(t, "[]", result.(*testItem).value, "value")
}

func Test_Get_SingleItem(t *testing.T) {
	tree := makeTrieSimple([]byte{0x92})

	//act
	result, ok := tree.Get([]byte{0x92})

	//assert
	assert.True(t, ok, "ok")
	assert.Equal(t, "[92]", result.(*testItem).value, "value")
}

func Test_Delete_NotExist(t *testing.T) {
	tree := makeTrieSimple([]byte{0x93, 0x94}, []byte{0x93, 0x95})

	//act
	result, old := Delete(tree, []byte{0x94})

	//assert
	assert.Nil(t, old, "old")
	assert.Equal(t, 2, result.Len(), "len")
	assert.Equal(t, tree.root, result.root, "no change, should return same root nodes")

}

func Test_DeleteDeep_NotExist(t *testing.T) {
	tree := makeTrieSimple([]byte{0x93, 0x94}, []byte{0x93, 0x95})

	//act
	result, old := Delete(tree, []byte{0x93, 0x93})

	//assert
	assert.Nil(t, old, "old")
	assert.Equal(t, 2, result.Len(), "len")
	assert.Equal(t, tree.root, result.root, "no change, should return same root nodes")
}

func Test_Delete_Exist(t *testing.T) {
	tree := makeTrieSimple([]byte{0x93, 0x94}, []byte{0x93, 0x95})

	//act
	result, old := Delete(tree, []byte{0x93, 0x94})

	//assert
	assert.NotNil(t, old, "old")
	assert.Equal(t, "[93 94]", old.(*testItem).value, "old.value")
	assert.Equal(t, 1, result.Len(), "len")
	_, ok := result.Get([]byte{0x93, 0x94})
	assert.False(t, ok, "should not contain deleted item")
}

func Test_DeleteRoot_NotExist(t *testing.T) {
	tree := makeTrieSimple([]byte{0x93}, []byte{0x93, 0x95})

	//act
	result, old := Delete(tree, []byte{})

	//assert
	assert.Nil(t, old, "old")
	assert.Equal(t, 2, result.Len(), "len")
}

func Test_DeleteRoot_Exist(t *testing.T) {
	tree := makeTrieSimple([]byte{}, []byte{0x93, 0x95})

	//act
	result, old := Delete(tree, []byte{})

	//assert
	assert.NotNil(t, old, "old")
	assert.Equal(t, 1, result.Len(), "len")
	_, ok := result.Get([]byte{})
	assert.False(t, ok, "should not contain deleted item")
}

func Test_DeleteNode_NoChange(t *testing.T) {
	tree := makeTrieSimple([]byte{0x93, 0x94}, []byte{0x93, 0x95})

	//act
	result, old := Delete(tree, []byte{0x93})

	//assert
	assert.Nil(t, old, "old")
	assert.Equal(t, 2, result.Len(), "len")
	assert.Equal(t, tree, result, "no change")
}

func Test_DeleteLeafWithChildren_ValueUnset(t *testing.T) {
	tree := makeTrieSimple([]byte{0x93}, []byte{0x93, 0x94}, []byte{0x93, 0x95})

	//act
	result, old := Delete(tree, []byte{0x93})

	//assert
	assert.NotNil(t, old, "old")
	assert.Equal(t, 2, result.Len(), "len")
	_, ok := result.Get([]byte{0x93})
	assert.False(t, ok, "should not contain deleted item")
}

func Test_DeleteSingleChildLeaf_ChildMergedUp(t *testing.T) {
	tree := makeTrieSimple([]byte{0x93}, []byte{0x93, 0x94})

	//act
	result, old := Delete(tree, []byte{0x93})

	//assert
	assert.NotNil(t, old, "old")
	assert.Equal(t, 1, result.Len(), "len")
	_, ok := result.Get([]byte{0x93})
	assert.False(t, ok, "should not contain deleted item")

	fmt.Println(result.root.printDbg(""))
	child := result.root.children[0x93]
	assert.Equal(t, []byte{0x93, 0x94}, child.keySlice, "deep child should be merged up")

	val, ok := result.Get([]byte{0x93, 0x94})
	assert.True(t, ok, "merged child should exist")
	assert.Equal(t, "[93 94]", val.(*testItem).value, "merged child should be correct value")
}

func Test_DeletePureLeaf_RemovedFromParent(t *testing.T) {
	tree := makeTrieSimple([]byte{0x93}, []byte{0x93, 0x94})

	//act
	result, old := Delete(tree, []byte{0x93, 0x94})

	//assert
	assert.NotNil(t, old, "old")
	assert.Equal(t, 1, result.Len(), "len")
	_, ok := result.Get([]byte{0x93, 0x94})
	assert.False(t, ok, "should not contain deleted item")

	fmt.Println(result.root.printDbg(""))
	child := result.root.children[0x93]
	assert.Equal(t, []byte{0x93}, child.keySlice, "parent should be maintained")
	assert.Equal(t, 0, len(child.children), "deep child should be removed")
}

func disabled_Test_Get_FullTree(t *testing.T) {

	tree := makeGiantTrie(2)

	//act
	for i := 0; i <= 1; i++ {
		log.Printf("checking %d", i)
		for j := 0; j < 256; j++ {
			for k := 0; k < 256; k++ {
				bytes := []byte{byte(i), byte(j), byte(k)}
				result, ok := tree.Get(bytes)

				//assert
				assert.True(t, ok, "ok")
				assert.Equal(t, bytes, result, "value")
				if !ok {
					log.Printf("uh oh! %s", tree.root.printDbg(""))
					panic(fmt.Sprintf("!ok key: [% x]", bytes))
				}
			}
			bytes := []byte{byte(i), byte(j)}
			result, ok := tree.Get(bytes)
			assert.True(t, ok, "ok")
			assert.Equal(t, bytes, result, "value")
			if !ok {
				log.Printf(tree.root.printDbg(""))
				panic(fmt.Sprintf("!ok key: [% x]", bytes))
			}
		}
		bytes := []byte{byte(i)}
		result, ok := tree.Get(bytes)
		assert.True(t, ok, "ok")
		assert.Equal(t, bytes, result, "value")
		if !ok {
			log.Printf(tree.root.printDbg(""))
			panic(fmt.Sprintf("!ok key: [% x]", bytes))
		}
	}
	result, ok := tree.Get([]byte{})
	assert.True(t, ok, "ok")
	assert.Equal(t, []byte{}, result, "value")
}

func Test_Get_NonexistentChild_Nil(t *testing.T) {
	tree := makeTrieSimple([]byte{0x08})

	//act
	result, ok := tree.Get([]byte{0x09, 0x23})

	//assert
	assert.False(t, ok, "ok")
	assert.Nil(t, result, "value")
}

func Test_Get_KeyIsNode_Nil(t *testing.T) {
	tree := makeTrieSimple([]byte{0x08, 0x54})

	log.Printf(tree.root.printDbg(""))

	//act
	result, ok := tree.Get([]byte{0x08})

	//assert
	assert.False(t, ok, "ok")
	assert.Nil(t, result, "value")
}

func makeTrieSimpleArr(keys [][]byte) *Trie {
	ret := NilTrie()
	for i := 0; i < len(keys); i++ {
		key := keys[i]
		val := fmt.Sprintf("[% x]", key)
		ret, _ = Set(ret, key, &testItem{
			value: val,
		})
	}
	return ret
}

func makeTrieSimple(keys ...[]byte) *Trie {
	return makeTrieSimpleArr(keys)
}

func makeGiantTrie(iterations int) *Trie {
	tree := NilTrie()
	var i, j, k int
	for i = 0; i < iterations; i++ {
		log.Printf("populating %d", i)
		for j = 0; j < 256; j++ {
			for k = 0; k < 256; k++ {
				bytes := []byte{byte(i), byte(j), byte(k)}
				tree, _ = Set(tree, bytes, bytes)
			}
			bytes := []byte{byte(i), byte(j)}
			tree, _ = Set(tree, bytes, bytes)
		}
		bytes := []byte{byte(i)}
		tree, _ = Set(tree, bytes, bytes)
	}
	tree, _ = Set(tree, []byte{}, []byte{})

	return tree
}
