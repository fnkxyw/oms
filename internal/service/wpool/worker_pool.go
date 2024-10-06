package wpool

import (
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
)

type job struct {
	task func()
	ctx  context.Context
	name string
}

type WorkerPool struct {
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
	mu         sync.Mutex
	numWorkers int
	//numJobs      atomic.Int64
	notification chan string
	stopChan     chan struct{}
	jobChan      chan job
}

func NewWorkerPool(ctx context.Context, numWorkers int, notification chan string) (*WorkerPool, error) {
	if numWorkers <= 0 {
		numWorkers = 1
	}
	if numWorkers > maxGoRoutines {
		return nil, ErrBadNumOfGoRoutines
	}
	ctx, cancel := context.WithCancel(ctx)
	return &WorkerPool{
		ctx:          ctx,
		cancel:       cancel,
		numWorkers:   numWorkers,
		notification: notification,
		stopChan:     make(chan struct{}),
		jobChan:      make(chan job, 100),
	}, nil
}

func (wp *WorkerPool) worker() {
	defer wp.wg.Done()
	for {
		select {
		case <-wp.stopChan:
			return
		case job, ok := <-wp.jobChan:
			if !ok {
				return
			}

			select {
			case <-job.ctx.Done():
				continue
			default:
			}

			wp.notification <- fmt.Sprintf("work by name %s started", job.name)
			job.task()
			wp.notification <- fmt.Sprintf("work by name %s finished", job.name)
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
	wp.cancel()
	close(wp.stopChan)
	wp.wg.Wait()
	close(wp.notification)
	close(wp.jobChan)
}

func (wp *WorkerPool) AddJob(ctx context.Context, task func(), name string) {
	job := job{
		task: task,
		ctx:  ctx,
		name: name,
	}

	wp.jobChan <- job
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
	//time.Sleep(10 * time.Millisecond)
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
