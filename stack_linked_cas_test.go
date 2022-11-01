package golockfree

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConcurrent_Put_Then_Pop(t *testing.T) {
	var size uint64 = 1000
	var wait sync.WaitGroup
	stack := CasLockFreeStackLinked[int]{}
	for i := 0; i < int(size); i++ {
		wait.Add(1)
		go func(i int) {
			stack.Put(i)
			stack.Peek()
			wait.Done()
		}(i)
	}
	wait.Wait()
	assert.Equal(t, stack.Size(), size)
	assert.Equal(t, stack.Empty(), false)
	for !stack.Empty() {
		wait.Add(1)
		go func() {
			stack.Peek()
			stack.Pop()
			wait.Done()
		}()
	}
	wait.Wait()
	assert.Equal(t, stack.Size(), uint64(0))
	assert.Equal(t, stack.Empty(), true)
}

func TestConcurrent_Put_Peek_Pop(t *testing.T) {
	var size uint64 = 1000
	stack := CasLockFreeStackLinked[int]{}
	var wait sync.WaitGroup
	for i := 0; i < int(size); i++ {
		wait.Add(1)
		go func(i int) {
			stack.Put(i)
			stack.Peek()
			stack.Pop()
			wait.Done()
		}(i)
	}
	wait.Wait()
	assert.Equal(t, stack.Size(), uint64(0))
	assert.Equal(t, stack.Empty(), true)
}
