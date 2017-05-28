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
	for index := 0; index < count; index++ {
		page.keys[pos+index] = p2.keys[possource+index]
		page.values[pos+index] = p2.values[possource+index]
		page.links[pos+index] = p2.links[possource+index]
	}
	page.links[pos+count] = p2.links[possource+count]
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
func (page *test_page) Value(pos int) IValue           { return page.keys[pos] }
func (page *test_page) SetValue(pos int, value IValue) { page.values[pos] = value.(int) }
func (page *test_page) Insert(key IKey, value IValue, link IAddress, pos int) {
	for i := page.count - 1; i >= pos && i >= 0; i-- {
		page.keys[i+1] = page.keys[i]
		page.values[i+1] = page.values[i]
		page.keys[i+1] = page.keys[i]
	}
	page.count++
	page.keys[pos] = key.(int)
	page.values[pos] = value.(int)
	var ok bool
	if page.links[pos+1], ok = link.(int); !ok {
		page.links[pos+1] = 0
	}
}

type test_btree struct {
	IBTree
	pageid   int
	root     *test_page
	pages    map[int]*test_page
	sizepage int
}

func newBTree() *test_btree {
	tree := &test_btree{pages: map[int]*test_page{}, pageid: 1}
	tree.sizepage = 10
	tree.root = tree.NewPage().(*test_page)
	tree.root.SetLeaf(true)
	return tree
}
func (tree *test_btree) Root() IPage { return tree.root }
func (tree *test_btree) NewPage() IPage {
	page := newPage(tree, tree.pageid, tree.sizepage)
	tree.pageid++
	tree.pages[page.address] = page
	return page
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
func TestInsert(t *testing.T) {
	tree := newBTree()
	type args struct {
		tree  IBTree
		key   IKey
		value IValue
	}
	type args2 struct {
		tree IBTree
		key  IKey
	}
	tests1 := []struct {
		name string
		args args
		want bool
	}{
		{name: "insert", args: args{key: 1, value: 1001, tree: tree}, want: true},
		{name: "insert", args: args{key: 3, value: 1002, tree: tree}, want: true},
		{name: "insert", args: args{key: 5, value: 1003, tree: tree}, want: true},
		{name: "insert", args: args{key: 7, value: 1003, tree: tree}, want: true},
		{name: "insert", args: args{key: 9, value: 1004, tree: tree}, want: true},
		{name: "insert", args: args{key: 11, value: 1004, tree: tree}, want: true},
		{name: "insert", args: args{key: 2, value: 1005, tree: tree}, want: true},
		{name: "insert", args: args{key: 4, value: 1006, tree: tree}, want: true},
		{name: "insert", args: args{key: 6, value: 1007, tree: tree}, want: true},
		{name: "insert", args: args{key: 8, value: 1008, tree: tree}, want: true},
		{name: "insert", args: args{key: 10, value: 1008, tree: tree}, want: true},
		{name: "insert", args: args{key: 12, value: 1008, tree: tree}, want: true},
		// TODO: Add test cases.
	}
	tests2 := []struct {
		name string
		args args2
		want bool
	}{
		{name: "contain", args: args2{key: 0, tree: tree}, want: false},
		{name: "contain", args: args2{key: 1, tree: tree}, want: true},
		{name: "contain", args: args2{key: 2, tree: tree}, want: true},
		{name: "contain", args: args2{key: 3, tree: tree}, want: true},
		{name: "contain", args: args2{key: 4, tree: tree}, want: true},
		{name: "contain", args: args2{key: 5, tree: tree}, want: true},
		{name: "contain", args: args2{key: 6, tree: tree}, want: true},
		{name: "contain", args: args2{key: 7, tree: tree}, want: true},
		{name: "contain", args: args2{key: 8, tree: tree}, want: true},
		{name: "contain", args: args2{key: 9, tree: tree}, want: true},
		{name: "contain", args: args2{key: 10, tree: tree}, want: true},
		{name: "contain", args: args2{key: 11, tree: tree}, want: true},
		{name: "contain", args: args2{key: 12, tree: tree}, want: true},
		{name: "contain", args: args2{key: 13, tree: tree}, want: false},
		// TODO: Add test cases.
	}
	tests3 := []struct {
		name string
		args args
		want bool
	}{
		{name: "dub", args: args{key: 1, value: 2001, tree: tree}, want: false},
		{name: "dub", args: args{key: 2, value: 2002, tree: tree}, want: false},
		{name: "dub", args: args{key: 3, value: 2003, tree: tree}, want: false},
		{name: "dub", args: args{key: 4, value: 2004, tree: tree}, want: false},
		{name: "dub", args: args{key: 5, value: 2005, tree: tree}, want: false},
		{name: "dub", args: args{key: 6, value: 2006, tree: tree}, want: false},
		{name: "dub", args: args{key: 7, value: 2007, tree: tree}, want: false},
		{name: "dub", args: args{key: 8, value: 2008, tree: tree}, want: false},
		{name: "dub", args: args{key: 9, value: 2009, tree: tree}, want: false},
		{name: "dub", args: args{key: 100, value: 2010, tree: tree}, want: true},
		// TODO: Add test cases.
	}
	for _, tt := range tests1 {
		//fmt.Println(tt.name, tt.args.tree, tt.args.key, tt.args.value)
		t.Run(tt.name, func(t *testing.T) {
			if got := Insert(tt.args.tree, tt.args.key, tt.args.value); got != tt.want {
				t.Errorf("Insert() = %v, want %v", got, tt.want)
			}
		})
	}
	for _, tt := range tests2 {
		//fmt.Println(tt.name, tt.args.tree, tt.args.key)
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainKey(tt.args.tree, tt.args.key); got != tt.want {
				t.Errorf("ContainKey() = %v, want %v", got, tt.want)
			}
		})
	}
	for _, tt := range tests3 {
		//fmt.Println(tt.name, tt.args.tree, tt.args.key, tt.args.value)
		t.Run(tt.name, func(t *testing.T) {
			if got := Insert(tt.args.tree, tt.args.key, tt.args.value); got != tt.want {
				t.Errorf("Insert() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestInsert2(t *testing.T) {
	tree := newBTree()
	for i := 0; i < 1000; i++ {
		//fmt.Println("TestInsert2", i)
		t.Run("TestInsert2", func(t *testing.T) {
			if got := Insert(tree, i, i+1000); !got {
				t.Errorf("Insert() = %v, want %v", got, true)
			}
		})
	}
}
