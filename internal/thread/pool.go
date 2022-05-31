package thread

import (
	"sync"
)

// Pool is a thread pool that takes a bunch of tasks and uses the threads in the pool to execute them concurrently
type Pool struct {
	tasks       Tasks
	threadCount int
	tasksChan   chan *Task
	waitGroup   sync.WaitGroup
}

// NewPool returns a new thread pool with the given thread count and tasks to execute
func NewPool(tasks Tasks, threadCount int) *Pool {
	return &Pool{
		tasks:       tasks,
		threadCount: threadCount,
		tasksChan:   make(chan *Task),
	}
}

// Run starts executing the tasks concurrently and blocks until all the tasks have been completed
func (p *Pool) Run() {
	p.waitGroup.Add(len(p.tasks))

	for i := 0; i < p.threadCount; i++ {
		go p.work()
	}

	for _, task := range p.tasks {
		p.tasksChan <- task
	}

	close(p.tasksChan)
	p.waitGroup.Wait()
}

// work is the work function of a thread. All threads from the pool call this function and wait to pick a task from the
// tasksChan queue. Once a thread gets a task it executes it, finishes it, and waits again to pick a new task
// from the queue until there are no more task remaining
func (p *Pool) work() {
	for task := range p.tasksChan {
		task.Run(&p.waitGroup)
	}
}
