package golockfree

/*
The Compare and Swap lock free algorithm (CAS) implementation
for the stack linked data structure
*/
import (
	"sync/atomic"
)

// TODO: move the ListNode into a utils module (it will be used by other implementations)
// ListNode used for the stack implementation
type ListNode[T any] struct {
	val  T
	next *ListNode[T]
}

// Lock Free stack LinkedList
type LockFreeStackLinked[T any] struct {
	top  atomic.Pointer[ListNode[T]]
	size atomic.Uint64
}

// Put a value of type T to the stack
func (stack *LockFreeStackLinked[T]) Put(val T) {
	for {
		var topList *ListNode[T] = stack.top.Load()
		var node *ListNode[T] = &ListNode[T]{val: val, next: topList}
		if stack.top.CompareAndSwap(topList, node) {
			// succeed on adding the node and updating the topList
			// atomic.AddUint64(&stack.size, 1)
			stack.size.Add(1)
			return
		}
	}
}

// Popout the top value from the stack
func (stack *LockFreeStackLinked[T]) Pop() (T, bool) {
	for {
		var topList *ListNode[T] = stack.top.Load()
		if topList == nil {
			// https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md#the-zero-value
			var zero T
			// empty response
			return zero, false
		}
		if stack.top.CompareAndSwap(topList, topList.next) {
			// atomic.AddUint64(&stack.size, ^uint64(0)) // to decrement the size
			stack.size.Add(^uint64(0))
			return topList.val, true
		}
	}
}

// Returns the size of the stack
func (stack *LockFreeStackLinked[T]) Size() uint64 {
	return stack.size.Load()
}

// Check if the stack is empty or not
func (stack *LockFreeStackLinked[T]) Empty() bool {
	return stack.size.Load() == 0
}

// Peek the top element in the stack
func (stack *LockFreeStackLinked[T]) Peek() (T, bool) {
	var top *ListNode[T] = stack.top.Load()
	if top == nil {
		var zero T
		return zero, false
	}
	return top.val, true
}
