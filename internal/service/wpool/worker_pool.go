package wpool

import (
	"container/heap"
	"context"
	"errors"
	"fmt"
	"sync"
)

var (
	ErrBadNumOfGoRoutines = errors.New("the number of goroutines cannot be like that")
)

const (
	maxGoRoutines = 100
	maxQueueSize  = 100
)

type WorkerPool struct {
	ctx          context.Context
	cancel       context.CancelFunc
	wg           sync.WaitGroup
	mu           sync.Mutex
	numWorkers   int
	notification chan string
	jobQueue     PriorityQueue
	stopChan     chan struct{}
	semaphore    chan struct{}
}

func NewWorkerPool(ctx context.Context, numWorkers int, notification chan string) (*WorkerPool, error) {
	if numWorkers <= 0 {
		numWorkers = 1
	}
	if numWorkers > maxGoRoutines {
		return nil, ErrBadNumOfGoRoutines
	}
	ctx, cancel := context.WithCancel(ctx)
	wp := &WorkerPool{
		ctx:          ctx,
		cancel:       cancel,
		numWorkers:   numWorkers,
		notification: notification,
		jobQueue:     make(PriorityQueue, 0),
		stopChan:     make(chan struct{}),
		semaphore:    make(chan struct{}, maxQueueSize),
	}
	heap.Init(&wp.jobQueue)
	return wp, nil
}

func (wp *WorkerPool) worker() {
	defer wp.wg.Done()
	for {
		select {
		case <-wp.stopChan:
			return
		case <-wp.semaphore:
			wp.mu.Lock()
			if wp.jobQueue.Len() > 0 {
				job := heap.Pop(&wp.jobQueue).(PriorityJob)
				wp.mu.Unlock()

				select {
				case <-job.ctx.Done():
					continue
				default:
					wp.notification <- fmt.Sprintf("work by name %s started", job.name)
					job.task()
					wp.notification <- fmt.Sprintf("work by name %s finished", job.name)
				}
			} else {
				wp.mu.Unlock()
			}
		}
	}
}

func (wp *WorkerPool) Start() {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	for i := 0; i < wp.numWorkers; i++ {
		wp.wg.Add(1)
		go wp.worker()
	}
}

func (wp *WorkerPool) Stop() {
	close(wp.stopChan)
	wp.wg.Wait()
	for len(wp.semaphore) > 0 {
		<-wp.semaphore
	}
	close(wp.notification)
}

func (wp *WorkerPool) AddJob(ctx context.Context, task func(), name string, priority int) {

	wp.semaphore <- struct{}{}

	wp.mu.Lock()
	defer wp.mu.Unlock()

	job := PriorityJob{
		task:     task,
		ctx:      ctx,
		name:     name,
		priority: priority,
	}

	heap.Push(&wp.jobQueue, job)
}

func (wp *WorkerPool) AddWorker(n int) error {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if wp.numWorkers+n > maxGoRoutines {
		return ErrBadNumOfGoRoutines
	}

	for i := 1; i <= n; i++ {
		wp.wg.Add(1)
		go wp.worker()
	}
	wp.numWorkers += n
	return nil
}

func (wp *WorkerPool) RemoveWorkers(n int) error {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if wp.numWorkers-n < 1 {
		return ErrBadNumOfGoRoutines
	}

	for i := 0; i < n; i++ {
		wp.stopChan <- struct{}{}
	}

	wp.numWorkers -= n
	return nil
}

func (wp *WorkerPool) PrintWorkers() {
	wp.mu.Lock()
	defer wp.mu.Unlock()
	fmt.Println("\nNum of workers: ", wp.numWorkers)
}

func (wp *WorkerPool) ChangeNumOfWorkers(numOfWorkers int) error {
	if numOfWorkers > 0 {
		return wp.AddWorker(numOfWorkers)
	}
	if numOfWorkers < 0 {
		numOfWorkers *= -1
		return wp.RemoveWorkers(numOfWorkers)
	}
	return nil
}
