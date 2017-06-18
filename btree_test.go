package btree

import (
	"testing"
)

type test_page struct {
	tree    *test_btree
	address int
	leaf    bool
	count   int
	links   []int
	keys    []int
	values  []int
}

func newPage(tree *test_btree, address int, size int) *test_page {
	p := &test_page{}
	p.address = address
	p.links = make([]int, size+1, size+1)
	p.keys = make([]int, size, size)
	p.values = make([]int, size, size)
	p.count = 0
	p.leaf = false
	return p
}
func (page *test_page) Address() IAddress  { return page.address }
func (page *test_page) Tree() IBTree       { return page.tree }
func (page *test_page) Leaf() bool         { return page.leaf }
func (page *test_page) SetLeaf(leaf bool)  { page.leaf = leaf }
func (page *test_page) Write()             {}
func (page *test_page) Count() int         { return page.count }
func (page *test_page) SetCount(count int) { page.count = count }
func (page *test_page) MaxCount() int      { return len(page.keys) }
func (page *test_page) CopyItems(pos int, psource IPage, possource int, count int) {
	p2 := psource.(*test_page)
	copy(page.keys[pos:pos+count], p2.keys[possource:possource+count])
	copy(page.values[pos:pos+count], p2.values[possource:possource+count])
	copy(page.links[pos+1:pos+1+count], p2.links[possource+1:possource+1+count])
}
func (page *test_page) ClearItems(pos int, count int) {
	for index := 0; index < count; index++ {
		page.keys[pos+index] = 0
		page.values[pos+index] = 0
		page.links[pos+index+1] = 0
	}
}
func (page *test_page) Link(pos int) IAddress          { return page.links[pos] }
func (page *test_page) SetLink(pos int, link IAddress) { page.links[pos] = link.(int) }
func (page *test_page) Key(pos int) IKey               { return page.keys[pos] }
func (page *test_page) SetKey(pos int, key IKey)       { page.keys[pos] = key.(int) }
func (page *test_page) Value(pos int) IValue           { return page.values[pos] }
func (page *test_page) SetValue(pos int, value IValue) { page.values[pos] = value.(int) }

/*func (page *test_page) Insert(key IKey, value IValue, link IAddress, pos int) {
	for i := page.count - 1; i >= pos && i >= 0; i-- {
		page.keys[i+1] = page.keys[i]
		page.values[i+1] = page.values[i]
		page.links[i+1] = page.links[i]
	}
	page.count++
	page.keys[pos] = key.(int)
	page.values[pos] = value.(int)
	var ok bool
	if page.links[pos+1], ok = link.(int); !ok {
		page.links[pos+1] = 0
	}
}*/

type test_btree struct {
	IBTree
	root     IPage
	pages    []IPage
	sizepage int
}

func newBTree() *test_btree {
	tree := &test_btree{pages: make([]IPage, 0, 10)}
	tree.sizepage = 11
	tree.root = tree.NewPage().(*test_page)
	tree.root.SetLeaf(true)
	return tree
}
func (tree *test_btree) Root() IPage { return tree.root }
func (tree *test_btree) NewPage() IPage {
	page := newPage(tree, len(tree.pages), tree.sizepage)
	tree.pages = append(tree.pages, page)
	return page
}
func (tree *test_btree) FreePage(page IPage) {
	page.ClearItems(0, page.Count())
	page.SetCount(0)
	page.SetLeaf(false)
	page.Write()
}
func (tree *test_btree) Page(Address IAddress) IPage {
	a := Address.(int)
	return tree.pages[a]
}
func (tree *test_btree) LessKey(key1 IKey, key2 IKey) bool { return key1.(int) < key2.(int) }
func (tree *test_btree) EqKey(key1 IKey, key2 IKey) bool {
	v1 := key1.(int)
	v2 := key2.(int)
	return v1 == v2
}
func testLoop(t *testing.T, ASC bool, minInsert int, maxInsert int, minLoopKey IKey, maxLoopKey IKey) bool {
	tree := newBTree()
	for i := minInsert; i <= maxInsert; i++ {
		if got := Insert(tree, i, i+1000); !got {
			t.Errorf("Insert(%v) = %v, want %v", i, got, true)
			return false
		}
	}
	min := minInsert
	max := maxInsert
	if minLoopKey != nil && minLoopKey.(int) > min {
		min = minLoopKey.(int)
	}
	if maxLoopKey != nil && maxLoopKey.(int) < max {
		max = maxLoopKey.(int)
	}
	count := max - min + 1
	if count < 0 {
		count = 0
	}
	minKey := 0
	maxKey := 0
	first := true
	countKey := 0
	perv := 0
	countKey2 := Loop(tree, ASC, minLoopKey, maxLoopKey, func(key IKey, value IValue) {
		if first {
			minKey = key.(int)
			maxKey = key.(int)
			first = false
			perv = key.(int)
		} else {
			if ASC {
				if perv >= key.(int) {
					t.Errorf("ASC %v < %v", perv, key.(int))
					return
				}
			} else {
				if perv <= key.(int) {
					t.Errorf("DESC %v > %v", perv, key.(int))
					return
				}
			}
			perv = key.(int)
		}
		if minKey > key.(int) {
			minKey = key.(int)
		}
		if maxKey < key.(int) {
			maxKey = key.(int)
		}
		countKey++
	})
	if countKey2 != countKey {
		t.Errorf("count = %v, want %v", countKey2, countKey)
		return false
	}
	if min > max && count == 0 {
		return true
	}
	if minKey != min {
		t.Errorf("MinKey = %v, want %v", minKey, min)
		return false
	}
	if maxKey != max {
		t.Errorf("MaxKey = %v, want %v", maxKey, max)
		return false
	}
	if countKey != count {
		t.Errorf("Count = %v, want %v", countKey, count)
		return false
	}
	return true
}
func TestInsert(t *testing.T) {
	t.Run("Insert to back", func(t *testing.T) {
		tree := newBTree()
		for i := 1; i <= 1000; i++ {
			if got := Insert(tree, i, i+1000); !got {
				t.Errorf("Insert(%v) = %v, want %v", i, got, true)
				return
			}
		}
		for i := 1; i <= 1000; i++ {
			if got := Insert(tree, i, i+1000); got {
				t.Errorf("Dub Insert(%v) = %v, want %v", i, got, true)
				return
			}
		}
		for i := -1000; i <= 0; i++ {
			if got := ContainKey(tree, i); got {
				t.Errorf("ContainKey(%v) = %v, want %v", i, got, false)
				return
			}
		}
		for i := 1; i <= 1000; i++ {
			if got := ContainKey(tree, i); !got {
				t.Errorf("ContainKey(%v) = %v, want %v", i, got, true)
				return
			}
		}
		for i := 1001; i <= 2000; i++ {
			if got := ContainKey(tree, i); got {
				t.Errorf("ContainKey(%v) = %v, want %v", i, got, false)
				return
			}
		}
	})
	t.Run("Insert to front", func(t *testing.T) {
		tree := newBTree()
		for i := 1000; i >= 1; i-- {
			if got := Insert(tree, i, i+1000); !got {
				t.Errorf("Insert(%v) = %v, want %v", i, got, true)
				return
			}
		}
		for i := 1000; i >= 1; i-- {
			if got := Insert(tree, i, i+1000); got {
				t.Errorf("Dub Insert(%v) = %v, want %v", i, got, true)
				return
			}
		}
		for i := 0; i >= -1000; i-- {
			if got := ContainKey(tree, i); got {
				t.Errorf("ContainKey(%v) = %v, want %v", i, got, false)
				return
			}
		}
		for i := 1000; i >= 1; i-- {
			if got := ContainKey(tree, i); !got {
				t.Errorf("ContainKey(%v) = %v, want %v", i, got, true)
				return
			}
		}
		for i := 2000; i >= 1001; i-- {
			if got := ContainKey(tree, i); got {
				t.Errorf("ContainKey(%v) = %v, want %v", i, got, false)
				return
			}
		}
	})

	t.Run("Delete", func(t *testing.T) {
		tree := newBTree()
		for i := 1000; i >= 1; i-- {
			if got := Insert(tree, i, i+1000); !got {
				t.Errorf("Insert(%v) = %v, want %v", i, got, true)
				return
			}
		}
	})
}
func TestLoop(t *testing.T) {
	t.Run("Loop ASC", func(t *testing.T) {
		if !testLoop(t, true, 1, 1000, nil, nil) {
			return
		}
		if !testLoop(t, true, 1, 1000, 500, nil) {
			return
		}
		if !testLoop(t, true, 1, 1000, nil, 500) {
			return
		}
		if !testLoop(t, true, 1, 1000, 500, 500) {
			return
		}
		if !testLoop(t, true, 1, 1000, 500, 1) {
			return
		}
	})
	t.Run("Loop DESC", func(t *testing.T) {
		if !testLoop(t, false, 1, 1000, nil, nil) {
			return
		}
		if !testLoop(t, false, 1, 1000, 500, nil) {
			return
		}
		if !testLoop(t, false, 1, 1000, nil, 500) {
			return
		}
		if !testLoop(t, false, 1, 1000, 500, 500) {
			return
		}
		if !testLoop(t, false, 1, 1000, 500, 1) {
			return
		}
	})
}
func TestLoopPage(t *testing.T) {
	tree := newBTree()
	for i := 1; i <= 1000; i++ {
		if got := Insert(tree, i, i+1000); !got {
			t.Errorf("Insert(%v) = %v, want %v", i, got, true)
			return
		}
	}
	count := 0
	MaxLevel := 0
	LeafMinMaxInit := true
	LeafMinLevel := 0
	LeafMaxLevel := 0
	NumPage := 0
	p := LoopPage(tree, func(page IPage, keyMin IKey, keyMax IKey, level int) {
		NumPage++
		count += page.Count()
		if MaxLevel < level {
			MaxLevel = level
		}
		if page.Leaf() {
			if LeafMinMaxInit {
				LeafMinLevel = level
				LeafMaxLevel = level
				LeafMinMaxInit = false
			} else {
				if LeafMinLevel > level {
					LeafMinLevel = level
				}
				if LeafMaxLevel < level {
					LeafMaxLevel = level
				}
			}
		}
	})
	if MaxLevel != LeafMinLevel || LeafMinLevel != LeafMaxLevel {
		t.Errorf("MaxLevel %v, LeafMin %v, LeafMax %v", MaxLevel, LeafMinLevel, LeafMaxLevel)
	}
	if p != NumPage {
		t.Errorf("NumPage %v != %v", NumPage, p)
	}
	if count != 1000 {
		t.Errorf("Count %v, want %v", count, 1000)
	}
}
func TestGetValue(t *testing.T) {
	tree := newBTree()
	for i := 1; i <= 1000; i++ {
		if got := Insert(tree, i, i+1000); !got {
			t.Errorf("Insert(%v) = %v, want %v", i, got, true)
			return
		}
	}
	for i := -1000; i <= 0; i++ {
		if _, ok := GetValue(tree, i); ok {
			t.Errorf("GetValue(%v)  ok=%v, want %v", i, ok, true)
			return
		}
	}
	for i := 1; i <= 1000; i++ {
		if got, ok := GetValue(tree, i); !ok || (ok && got.(int) != i+1000) {
			t.Errorf("GetValue(%v) = %v, want %v, ok=%v, want %v", i, got, i+1000, ok, true)
			return
		}
	}
	for i := 1001; i <= 2000; i++ {
		if _, ok := GetValue(tree, i); ok {
			t.Errorf("GetValue(%v)  ok=%v, want %v", i, ok, true)
			return
		}
	}
}
func TestDelete(t *testing.T) {
	t.Run("Delete 1", func(t *testing.T) {
		for j := -100; j <= 1100; j++ {
			tree := newBTree()
			for i := 1; i <= 1000; i++ {
				if got := Insert(tree, i, i+1000); !got {
					t.Errorf("Insert(%v) = %v, want %v", i, got, true)
					return
				}
			}
			c0 := Loop(tree, true, nil, nil, func(key IKey, value IValue) {})
			ok := DeleteKey(tree, j)
			c1 := Loop(tree, true, nil, nil, func(key IKey, value IValue) {})
			if j < 1 || j > 1000 {
				if ok {
					t.Errorf("Delete(%v) = %v, want %v", j, ok, false)
					return
				}
				if c1 != c0 {
					t.Errorf("Before = %v, after %v", c0, c1)
					return
				}
			}
			if j >= 1 && j <= 1000 {
				if !ok {
					t.Errorf("Delete(%v) = %v, want %v", j, ok, true)
					return
				}
				if c1 != c0-1 {
					t.Errorf("Before = %v, after %v", c0, c1)
					return
				}
			}
		}
	})
	t.Run("Delete fill", func(t *testing.T) {
		num := 100
		for j := -100; j <= num+100; j++ {
			tree := newBTree()
			for i := 1; i <= num; i++ {
				if got := Insert(tree, i, i+1000); !got {
					t.Errorf("Insert(%v) = %v, want %v", i, got, true)
					return
				}
			}
			for i := -100; i <= j; i++ {
				c0 := Loop(tree, true, nil, nil, func(key IKey, value IValue) {})
				ok := DeleteKey(tree, i)
				c1 := Loop(tree, true, nil, nil, func(key IKey, value IValue) {})
				if i < 1 || i > num {
					if ok {
						t.Errorf("Delete(%v) = %v, want %v", i, ok, false)
						return
					}
					if c1 != c0 {
						t.Errorf("Before = %v, after %v", c0, c1)
						return
					}
				}
				if i >= 1 && i <= num {
					if !ok {
						t.Errorf("Delete(%v) = %v, want %v", i, ok, true)
						return
					}
					if c1 != c0-1 {
						t.Errorf("Before = %v, after %v", c0, c1)
						return
					}
				}
			}
		}
	})
	t.Run("Delete fill 2", func(t *testing.T) {
		num := 100
		for j := num + 100; j >= -100; j-- {
			tree := newBTree()
			for i := 1; i <= num; i++ {
				if got := Insert(tree, i, i+1000); !got {
					t.Errorf("Insert(%v) = %v, want %v", i, got, true)
					return
				}
			}
			for i := num + 100; i > j; i-- {
				c0 := Loop(tree, true, nil, nil, func(key IKey, value IValue) {})
				ok := DeleteKey(tree, i)
				c1 := Loop(tree, true, nil, nil, func(key IKey, value IValue) {})
				if i < 1 || i > num {
					if ok {
						t.Errorf("Delete(%v) = %v, want %v", i, ok, false)
						return
					}
					if c1 != c0 {
						t.Errorf("Before = %v, after %v", c0, c1)
						return
					}
				}
				if i >= 1 && i <= num {
					if !ok {
						t.Errorf("Delete(%v) = %v, want %v", i, ok, true)
						return
					}
					if c1 != c0-1 {
						t.Errorf("Before = %v, after %v", c0, c1)
						return
					}
				}
			}
		}
	})
}
