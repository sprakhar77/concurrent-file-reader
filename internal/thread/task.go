package thread

import (
	"sync"
)

// Task is a simple wrapper around a function object that can be executed concurrently in a thread
type Task struct {
	job func() error
	err error
}

// NewTask returns a new task that will run the given function
func NewTask(job func() error) *Task {
	return &Task{job: job}
}

// Run executes the task and after completion it informs the calling thread via the WaitGroup provided
func (t *Task) Run(waitGroup *sync.WaitGroup) {
	t.err = t.job()
	waitGroup.Done()
}

// Tasks is an array of Task pointers
type Tasks []*Task

// Errors returns the list of errors that were encountered while executing the tasks
func (ts Tasks) Errors() []error {
	var errors []error
	for _, t := range ts {
		if t.err == nil {
			continue
		}
		errors = append(errors, t.err)
	}
	return errors
}
