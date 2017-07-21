package btree

import (
	"fmt"
)

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
}
type IBTree interface {
	Root() IPage
	NewPage() IPage
	FreePage(page IPage)
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
	insertItem(tree, parent, yPage.Key(nyNode), yPage.Value(nyNode), zPage.Address(), Item)
	zPage.SetLink(0, yPage.Link(nyNode+1))
	zPage.CopyItems(0, yPage, nyNode+1, nzNode)
	zPage.SetCount(nzNode)
	yPage.ClearItems(nyNode, nzNode+1)
	yPage.SetCount(nyNode)
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
func insertItem(tree IBTree, page IPage, key IKey, value IValue, link IAddress, pos int) {
	count := page.Count()
	page.CopyItems(pos+1, page, pos, count-pos)
	page.SetLink(pos+1, page.Link(pos))
	page.SetCount(count + 1)
	page.SetKey(pos, key)
	page.SetValue(pos, value)
	if link != nil {
		page.SetLink(pos+1, link)
	}
	page.Write()
}
func insertNonFull(tree IBTree, page IPage, key IKey, value IValue) bool {
	pos, eq := findKey(tree, page, key)
	if eq {
		return false
	}
	if page.Leaf() {
		if pos == -1 {
			insertItem(tree, page, key, value, nil, 0)
		} else {
			insertItem(tree, page, key, value, nil, pos)
		}
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
		ypage.SetLink(0, page.Link(0))
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
func loopPageASC(tree IBTree, page IPage, keyMin IKey, keyMax IKey, loop func(key IKey, value IValue) (next bool)) (int, bool) {
	var link IAddress
	var n int
	var next bool
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
			if n, next = loopPageASC(tree, tree.Page(link), keyMin, keyMax, loop); !next {
				return num + n, false
			}
			num += n
		}
		num++
		if !loop(key, page.Value(i)) {
			return num, false
		}
	}
	if !page.Leaf() {
		link = page.Link(end)
		if n, next = loopPageASC(tree, tree.Page(link), keyMin, keyMax, loop); !next {
			return num + n, false
		}
		num += n
	}
	return num, true
}
func loopPageDESC(tree IBTree, page IPage, keyMin IKey, keyMax IKey, loop func(key IKey, value IValue) (next bool)) (int, bool) {
	var link IAddress
	var n int
	var next bool
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
			if n, next = loopPageDESC(tree, tree.Page(link), keyMin, keyMax, loop); !next {
				return num + n, false
			}
			num += n
		}
		num++
		if !loop(key, page.Value(i)) {
			return num, false
		}
	}

	if !page.Leaf() {
		link = page.Link(start)
		if n, next = loopPageDESC(tree, tree.Page(link), keyMin, keyMax, loop); !next {
			return num + n, false
		}
		num += n
	}
	return num, true
}
func Loop(tree IBTree, ASC bool, keyMin IKey, keyMax IKey, loop func(key IKey, value IValue) (next bool)) (int, bool) {
	if ASC {
		return loopPageASC(tree, tree.Root(), keyMin, keyMax, loop)
	}
	return loopPageDESC(tree, tree.Root(), keyMin, keyMax, loop)
}
func loopPage(tree IBTree, page IPage, keyMin IKey, keyMax IKey, level int, loop func(page IPage, keyMin IKey, keyMax IKey, level int) (next bool)) (int, bool) {
	var next bool
	var n int
	num := 1
	if !loop(page, keyMin, keyMax, level) {
		return num, false
	}
	if !page.Leaf() {
		count := page.Count()
		if count > 0 {
			ValueLeft := page.Value(0)
			if n, next = loopPage(tree, tree.Page(page.Link(0)), nil, ValueLeft, level+1, loop); !next {
				return num + n, false
			}
			num = num + n
			for i := 1; i < count; i++ {
				ValueRight := page.Value(0)
				if n, next = loopPage(tree, tree.Page(page.Link(i)), ValueLeft, ValueRight, level+1, loop); !next {
					return num + n, false
				}
				num = num + n
				ValueLeft = ValueRight
			}
			if n, next = loopPage(tree, tree.Page(page.Link(count)), page.Value(count-1), nil, level+1, loop); !next {
				return num + n, false
			}
			num = num + n
		}
	}
	return num, true
}
func LoopPage(tree IBTree, loop func(page IPage, keyMin IKey, keyMax IKey, level int) (next bool)) (int, bool) {
	return loopPage(tree, tree.Root(), nil, nil, 0, loop)
}
func GetValue(tree IBTree, key IKey) (IValue, bool) {
	page := tree.Root()
	for !page.Leaf() {
		pos, eq := findKey(tree, page, key)
		if eq {
			return page.Value(pos), true
		}
		page = tree.Page(page.Link(pos))
	}
	pos, eq := findKey(tree, page, key)
	if eq {
		return page.Value(pos), true
	}
	return nil, false
}
func SetValue(tree IBTree, key IKey, value IValue) bool {
	page := tree.Root()
	for !page.Leaf() {
		pos, eq := findKey(tree, page, key)
		if eq {
			page.SetValue(pos, value)
			return true
		}
		page = tree.Page(page.Link(pos))
	}
	pos, eq := findKey(tree, page, key)
	if eq {
		page.SetValue(pos, value)
		return true
	}
	return false
}
func findLeaf(tree IBTree, page IPage, positions []int, pages []IPage, Left bool) IPage {
	leaf := false
	for !leaf {
		var pos int
		if Left {
			pos = 0
		} else {
			pos = page.Count() - 1
		}
		positions = append(positions, pos)
		pages = append(pages, page)
		leaf = page.Leaf()
		if leaf {
			break
		}
		page = tree.Page(page.Link(pos))
	}
	return page
}
func deleteItem(tree IBTree, page IPage, pos int) {
	count := page.Count()
	page.ClearItems(pos, 1)
	page.CopyItems(pos, page, pos+1, count-pos-1)
	count--
	page.SetCount(count)
	page.ClearItems(count, page.MaxCount()-count)
	page.Write()
}

func balans(tree IBTree, page IPage, pos int) {
	LeftMove := true
	LeftUnion := true
	RightMove := true
	RightUnion := true
	var page_count int
	var p IPage
	var p_count int
	var p_maxcount int
	var lp IPage
	var lp_count int
	var lp_maxcount int
	var rp IPage
	var rp_count int
	var rp_maxcount int
	p = tree.Page(page.Link(pos))
	p_count = p.Count()
	p_maxcount = p.MaxCount()
	page_count = page.Count()
	if pos == 0 {
		LeftMove = false
		LeftUnion = false
	} else if pos == page_count {
		RightMove = false
		RightUnion = false
	}
	if LeftUnion || LeftMove {
		lp = tree.Page(page.Link(pos - 1))
		lp_count = lp.Count()
		lp_maxcount = lp.MaxCount()
	}
	if RightMove || RightUnion {
		rp = tree.Page(page.Link(pos + 1))
		rp_count = rp.Count()
		rp_maxcount = rp.MaxCount()
	}
	if LeftUnion {
		if lp_count+1+p_count > lp_maxcount {
			LeftUnion = false
		}
	}
	if RightUnion {
		if p_count+1+rp_count > p_maxcount {
			RightUnion = false
		}
	}
	if LeftMove {
		if lp_count-1 < lp_maxcount/2 {
			LeftMove = false
		}
	}
	if RightMove {
		if rp_count-1 < rp_maxcount/2 {
			RightMove = false
		}
	}
	if LeftMove {
		//insert p[0]
		p.CopyItems(1, p, 0, p_count)
		p.SetLink(1, p.Link(0))
		p.SetLink(0, lp.Link(lp_count-1))
		p.SetKey(0, page.Key(pos-1))
		p.SetValue(0, page.Value(pos-1))
		p.SetCount(p_count + 1)
		//page[pos-1]=lp[lp_count-1]
		page.SetKey(pos-1, lp.Key(lp_count-1))
		page.SetValue(pos-1, lp.Value(lp_count-1))
		//delete lp[lp_count-1]
		lp.SetCount(lp_count - 1)
		lp.ClearItems(lp_count-1, 1)
		lp.Write()
		p.Write()
		page.Write()
	} else if RightMove {
		//insert
		p.SetKey(p_count, page.Key(pos))
		p.SetValue(p_count, page.Value(pos))
		p.SetLink(p_count+1, rp.Link(0))
		p.SetCount(p_count + 1)
		//set
		page.SetKey(pos, rp.Key(0))
		page.SetValue(pos, rp.Value(0))
		//delete
		rp.SetLink(0, rp.Link(1))
		rp.CopyItems(0, rp, 1, rp_count-1)
		rp.ClearItems(rp_count-1, 1)
		rp.SetCount(rp_count - 1)

		p.Write()
		rp.Write()
		page.Write()

	} else if LeftUnion {
		// lp = join lp[0:count]+page[pos-1]+p[0:count]
		lp.SetKey(lp_count, page.Key(pos-1))
		lp.SetValue(lp_count, page.Value(pos-1))
		lp.SetLink(lp_count+1, p.Link(0))
		lp.CopyItems(lp_count+1, p, 0, p_count)
		lp.SetCount(lp_count + 1 + p_count)
		// delete page[pos-1]
		deleteItem(tree, page, pos-1)
		// write edits
		lp.Write()
		tree.FreePage(p)
	} else if RightUnion {
		// p = join p[0:count]+page[pos]+rp[0:count]
		p.SetKey(p_count, page.Key(pos))
		p.SetValue(p_count, page.Value(pos))
		p.SetLink(p_count+1, rp.Link(0))
		p.CopyItems(p_count+1, rp, 0, rp_count)
		p.SetCount(p_count + 1 + rp_count)
		// delete page[pos]
		deleteItem(tree, page, pos)
		// write edits
		p.Write()
		tree.FreePage(rp)
	} else {
		fmt.Println("CRACK")
		//перенос из верхнего уровня
	}
}
func DeleteKey(tree IBTree, key IKey) bool {
	positions := make([]int, 0, 10)
	pages := make([]IPage, 0, 10)
	page := tree.Root()
	leaf := false
	pos := 0
	for !leaf {
		var eq bool
		pos, eq = findKey(tree, page, key)
		positions = append(positions, pos)
		pages = append(pages, page)
		leaf = page.Leaf()
		if eq {
			break
		}
		if leaf && !eq {
			return false
		}
		page = tree.Page(page.Link(pos))
	}
	if leaf {
		deleteItem(tree, page, pos)
	} else {
		leftpage := tree.Page(page.Link(pos))
		rightpage := tree.Page(page.Link(pos + 1))
		if leftpage.Count() > rightpage.Count() {
			lpage := findLeaf(tree, leftpage, positions, pages, false)
			lpos := lpage.Count() - 1
			page.SetValue(pos, lpage.Value(lpos))
			page.SetKey(pos, lpage.Key(lpos))
			page.Write()
			deleteItem(tree, lpage, lpos)
		} else {
			rpage := findLeaf(tree, rightpage, positions, pages, true)
			rpos := 0
			page.SetValue(pos, rpage.Value(rpos))
			page.SetKey(pos, rpage.Key(rpos))
			page.Write()
			deleteItem(tree, rpage, rpos)
		}
	}
	for i := len(pages) - 1; i >= 1; i-- {
		page = pages[i]
		num := page.Count()
		max := page.MaxCount()
		if num < (max-1)/2 {
			balans(tree, pages[i-1], positions[i-1])
			//fmt.Println(num, max)
		}
	}
	if pages[0].Count() == 0 && len(pages) > 1 {
		p1 := tree.Page(pages[0].Link(0))
		p1count := p1.Count()
		pages[0].SetLink(0, p1.Link(0))
		pages[0].CopyItems(0, p1, 0, p1count)
		pages[0].SetCount(p1count)
		pages[0].SetLeaf(p1.Leaf())
		p1.ClearItems(0, p1count)
		tree.FreePage(p1)
	}
	return true
}
func CreateRoot(tree IBTree, NullKey IKey, NullValue IValue, NullLink IAddress) {
	root := tree.Root()
	max := root.MaxCount()
	root.SetLeaf(true)
	root.SetCount(max)
	for index := 0; index < max; index++ {
		root.SetKey(index, NullKey)
		root.SetValue(index, NullValue)
		root.SetLink(index, NullLink)
	}
	root.SetLink(max, NullLink)
	root.SetCount(0)
	root.Write()
}
