package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	head  *ListItem
	tail  *ListItem
	count int
}

func (l list) Len() int {
	return l.count
}

func (l *list) Front() *ListItem {
	return l.head
}

func (l *list) Back() *ListItem {
	return l.tail
}

func (l *list) PushBack(v interface{}) *ListItem {
	i := &ListItem{
		Value: v,
		Next:  nil,
		Prev:  nil,
	}
	if l.head == nil {
		l.head = i
		l.tail = i
	} else {
		l.tail.Next = i
		i.Prev = l.tail
		l.tail = i
	}
	l.count++
	return i
}

func (l *list) PushFront(v interface{}) *ListItem {
	i := &ListItem{
		Value: v,
		Next:  nil,
		Prev:  nil,
	}
	if l.head == nil {
		l.head = i
		l.tail = i
	} else {
		l.head.Prev = i
		i.Next = l.head
		l.head = i
	}
	l.count++
	return i
}

func (l *list) Remove(i *ListItem) {
	if i == nil {
		return
	}
	if i == l.head {
		l.head = l.head.Next
		if l.head != nil {
			l.head.Prev = nil
		}
	}
	if i == l.tail {
		l.tail = l.tail.Prev
		if l.tail != nil {
			l.tail.Next = nil
		}
	}
	if i.Prev != nil {
		i.Prev.Next = i.Next
		if i.Next != nil {
			i.Next.Prev = i.Prev
		}
	}
	l.count--
}

func (l *list) MoveToFront(i *ListItem) {
	l.Remove(i)
	l.PushFront(i.Value)
}

func NewList() List {
	return new(list)
}
