package list

import (
	"container/list"
	"sync"
)

type List struct {
	mux  sync.Mutex
	list *list.List
	size int
}

func NewList(size int) *List {
	sl := &List{
		list: list.New(),
		size: size,
	}
	if sl.size == 0 {
		sl.size = 1000
	}

	return sl
}

func (sl *List) PushFront(item interface{}) {
	sl.mux.Lock()
	defer sl.mux.Unlock()
	sl.list.PushFront(item)
	if sl.list.Len() > sl.size {
		sl.list.Remove(sl.list.Back())

	}
}

func (sl *List) PushBack(item interface{}) {
	sl.mux.Lock()
	defer sl.mux.Unlock()
	sl.list.PushBack(item)
	if sl.list.Len() > sl.size {
		sl.list.Remove(sl.list.Front())

	}
}
func (sl *List) RemoveBack() interface{} {
	sl.mux.Lock()
	defer sl.mux.Unlock()
	if sl.list.Len() > 0 {
		return sl.list.Remove(sl.list.Back())
	}
	return nil
}

func (sl *List) RemoveFront() interface{} {
	sl.mux.Lock()
	defer sl.mux.Unlock()
	if sl.list.Len() > 0 {
		return sl.list.Remove(sl.list.Front())
	}
	return nil
}

func (sl *List) RemoveItemFromFront(key interface{}, condition func(key, value interface{}) bool) interface{} {
	sl.mux.Lock()
	defer sl.mux.Unlock()
	for e := sl.list.Front(); e != nil; e = e.Next() {
		if condition(key, e.Value) {
			return sl.list.Remove(e)

		}

	}
	return nil
}
func (sl *List) RemoveItemFromBack(key interface{}, condition func(key, value interface{}) bool) interface{} {
	sl.mux.Lock()
	defer sl.mux.Unlock()
	for e := sl.list.Back(); e != nil; e = e.Prev() {
		if condition(key, e.Value) {
			return sl.list.Remove(e)

		}

	}
	return nil
}

func (sl *List) RemoveAllItems(key interface{}, condition func(key, value interface{}) bool) {
	sl.mux.Lock()
	defer sl.mux.Unlock()
	tmpList := []*list.Element{}
	for e := sl.list.Front(); e != nil; e = e.Next() {
		if condition(key, e.Value) {
			tmpList = append(tmpList, e)

		}
	}
	for i := 0; i < len(tmpList); i++ {
		sl.list.Remove(tmpList[i])
	}

}

func (sl *List) SearchFront(key interface{}, size int, condition func(key, value interface{}) bool) (result []interface{}) {

	sl.mux.Lock()
	defer sl.mux.Unlock()

	for e := sl.list.Front(); e != nil; e = e.Next() {
		if condition(key, e.Value) {
			result = append(result, e.Value)
		}
		if size > 0 {
			if len(result) >= size {
				return
			}
		}
	}

	return
}

func (sl *List) SearchBack(key interface{}, size int, condition func(key, value interface{}) bool) (result []interface{}) {
	sl.mux.Lock()
	defer sl.mux.Unlock()

	for e := sl.list.Back(); e != nil; e = e.Prev() {

		if condition(key, e.Value) {
			result = append(result, e.Value)
		}
		if size > 0 {
			if len(result) >= size {
				return
			}
		}
	}

	return
}

func (sl *List) GetListSize() int {
	sl.mux.Lock()
	defer sl.mux.Unlock()
	return sl.list.Len()
}
