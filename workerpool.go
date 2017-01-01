package trabago

import "sync"

// TODO:
// - change callback paradigm to be "optional" (in the case of a user wanting to create a "pipeline")
// 	 include a queue of responses as well, for serial handling

type WorkerPool struct {
	mu         *sync.Mutex
	size       int
	running    bool
	workerfunc func(v interface{}) interface{}
	work       chan interface{}
	callback   chan interface{}
	queue      []interface{}
	kill       chan struct{}
	wg         *sync.WaitGroup
}

func New(size uint16, workerfunc func(v interface{}) interface{}, callback chan interface{}) *WorkerPool {
	return &WorkerPool{
		mu:         new(sync.Mutex),
		size:       int(size),
		running:    false,
		workerfunc: workerfunc,
		work:       make(chan interface{}),
		callback:   callback,
		queue:      make([]interface{}, 0),
		kill:       make(chan struct{}, size),
		wg:         new(sync.WaitGroup),
	}
}

func (wp *WorkerPool) Run() {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if wp.running {
		return
	}

	for i := 0; i < wp.size; i++ {
		wp.wg.Add(1)
		go func(wp *WorkerPool) {
			for {
				select {
				case w := <-wp.work:
					if v := wp.workerfunc(w); wp.callback != nil && v != nil {
						wp.callback <- v
					}
				case <-wp.kill:
					wp.wg.Done()
					return
				}
			}
		}(wp)
	}

	// clear any queued work
	for i := range wp.queue {
		wp.work <- wp.queue[i]
	}

	wp.running = true
}

func (wp *WorkerPool) Stop() {
	wp.mu.Lock()
	if wp.running {
		for i := 0; i < wp.size; i++ {
			wp.kill <- struct{}{}
		}
		wp.wg.Wait()

		close(wp.work)
		close(wp.callback)
	}
	wp.running = false
	wp.mu.Unlock()
}

func (wp *WorkerPool) DoWork(v interface{}) {
	if wp.running {
		wp.work <- v
	} else {
		wp.mu.Lock()
		wp.queue = append(wp.queue, v)
		wp.mu.Unlock()
	}
}
