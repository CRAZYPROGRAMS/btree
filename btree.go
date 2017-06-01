package btree

type IAddress interface{}

type IKey interface{}
type IValue interface{}
type IPage interface {
	Address() IAddress
	Tree() IBTree
	Leaf() bool
	SetLeaf(leaf bool)
	Write()
	Count() int
	SetCount(count int)
	MaxCount() int
	CopyItems(pos int, psource IPage, possource int, count int)
	ClearItems(pos int, count int)
	Link(pos int) IAddress
	SetLink(pos int, address IAddress)
	Key(pos int) IKey
	SetKey(pos int, key IKey)
	Value(pos int) IValue
	SetValue(pos int, value IValue)
	Insert(key IKey, value IValue, link IAddress, pos int)
}
type IBTree interface {
	Root() IPage
	NewPage() IPage
	Page(Address IAddress) IPage
	LessKey(key1 IKey, key2 IKey) bool
	EqKey(key1 IKey, key2 IKey) bool
}

func splitChild(tree IBTree, parent IPage, yPage IPage, Item int) {
	var zPage IPage
	zPage = tree.NewPage()
	zPage.SetLeaf(yPage.Leaf())
	nyNode := yPage.MaxCount() / 2
	nzNode := yPage.MaxCount() - nyNode - 1
	parent.Insert(yPage.Key(nyNode), yPage.Value(nyNode), zPage.Address(), Item)
	zPage.CopyItems(0, yPage, nyNode+1, nzNode)
	zPage.SetCount(nzNode)
	yPage.ClearItems(nyNode, nzNode+1)
	yPage.SetCount(nyNode)
	parent.Write()
	yPage.Write()
	zPage.Write()
}
func findKey(tree IBTree, page IPage, key IKey) (int, bool) {
	left := 0
	right := page.Count() - 1
	if right-left < 0 {
		return -1, false
	}
	pos := left
	for right-left >= 0 {
		pos = (right + left) / 2
		if tree.EqKey(key, page.Key(pos)) {
			return pos, true
		}
		if tree.LessKey(key, page.Key(pos)) {
			right = pos - 1
		} else {
			left = pos + 1
			pos++
		}
	}
	return pos, false
}
func insertNonFull(tree IBTree, page IPage, key IKey, value IValue) bool {
	pos, eq := findKey(tree, page, key)
	if eq {
		return false
	}
	if page.Leaf() {
		if pos == -1 {
			page.Insert(key, value, nil, 0)
		} else {
			page.Insert(key, value, nil, pos)
		}
		page.Write()
		return true
	}
	xPage := tree.Page(page.Link(pos))
	if xPage.Count() == xPage.MaxCount() {
		splitChild(tree, page, xPage, pos)
		if tree.LessKey(page.Key(pos), key) {
			pos++
			xPage = tree.Page(page.Link(pos))
		}
	}
	return insertNonFull(tree, xPage, key, value)
}
func contain(tree IBTree, page IPage, key IKey) bool {
	pos, eq := findKey(tree, page, key)
	if eq {
		return true
	}
	if page.Leaf() {
		return false
	}
	return contain(tree, tree.Page(page.Link(pos)), key)
}
func ContainKey(tree IBTree, key IKey) bool {
	return contain(tree, tree.Root(), key)
}
func Insert(tree IBTree, key IKey, value IValue) bool {
	page := tree.Root()
	if page.Count() == page.MaxCount() {
		ypage := tree.NewPage()
		ypage.SetLeaf(page.Leaf())
		ypage.CopyItems(0, page, 0, page.Count())
		ypage.SetCount(page.Count())
		page.ClearItems(0, page.Count())
		page.SetCount(0)
		page.SetLeaf(false)
		page.SetLink(0, ypage.Address())
		splitChild(tree, page, ypage, 0)
	}
	return insertNonFull(tree, page, key, value)
}
func loopPageASC(tree IBTree, page IPage, keyMin IKey, keyMax IKey, loop func(key IKey, value IValue)) int {
	var link IAddress
	num := 0
	start := 0
	end := page.Count()
	if keyMin != nil {
		start, _ = findKey(tree, page, keyMin)
	}
	if keyMax != nil {
		var eq bool
		end, eq = findKey(tree, page, keyMax)
		if eq {
			end++
		}
	}
	for i := start; i < end; i++ {
		key := page.Key(i)
		if !page.Leaf() {
			link = page.Link(i)
			num += loopPageASC(tree, tree.Page(link), keyMin, keyMax, loop)
		}
		num++
		loop(key, page.Value(i))
	}
	if !page.Leaf() {
		link = page.Link(end)
		num += loopPageASC(tree, tree.Page(link), keyMin, keyMax, loop)
	}
	return num
}
func loopPageDESC(tree IBTree, page IPage, keyMin IKey, keyMax IKey, loop func(key IKey, value IValue)) int {
	var link IAddress
	num := 0
	start := 0
	end := page.Count()
	if keyMin != nil {
		start, _ = findKey(tree, page, keyMin)
	}
	if keyMax != nil {
		var eq bool
		end, eq = findKey(tree, page, keyMax)
		if eq {
			end++
		}
	}
	for i := end - 1; i >= start; i-- {
		key := page.Key(i)
		if !page.Leaf() {
			link = page.Link(i + 1)
			num += loopPageDESC(tree, tree.Page(link), keyMin, keyMax, loop)
		}
		num++
		loop(key, page.Value(i))
	}

	if !page.Leaf() {
		link = page.Link(start)
		num += loopPageDESC(tree, tree.Page(link), keyMin, keyMax, loop)
	}
	return num
}
func Loop(tree IBTree, ASC bool, keyMin IKey, keyMax IKey, loop func(key IKey, value IValue)) int {
	if ASC {
		return loopPageASC(tree, tree.Root(), keyMin, keyMax, loop)
	}
	return loopPageDESC(tree, tree.Root(), keyMin, keyMax, loop)
}
