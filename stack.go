package trabago

import (
	"bytes"
	"fmt"
	"sync"
)

type Stack struct {
	lock *sync.RWMutex
	head *node
	size int
}

type node struct {
	next *node
	val  interface{}
}

func NewStack() *Stack {
	return &Stack{
		lock: new(sync.RWMutex),
		head: nil,
		size: 0,
	}
}

func (s *Stack) Push(v interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()

	n := &node{val: v, next: new(node)}

	if !s.isEmpty() {
		*n.next = *s.head
	}

	s.head = n
	s.size++
}

func (s *Stack) Pop() interface{} {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.isEmpty() {
		return nil
	}

	val := s.head.val
	s.head = s.head.next
	s.size--

	return val
}

func (s *Stack) Peak() interface{} {
	return s.head.val
}

func (s *Stack) Clear() {
	s.lock.Lock()
	s.head = nil
	s.size = 0
	s.lock.Unlock()
}

func (s *Stack) Size() int {
	return s.size
}

func (s *Stack) String() string {
	s.lock.RLock()
	defer s.lock.RUnlock()

	buf := bytes.NewBufferString("[")
	for n := s.head; n != nil; n = n.next {
		buf.WriteString(fmt.Sprintf(" %#v", n.val))
	}
	buf.WriteString(" ]")

	return buf.String()
}

func (s *Stack) isEmpty() bool {
	return s.head == nil
}
