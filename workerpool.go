package trabago

import "sync"

type WorkerPool struct {
	mu         *sync.Mutex
	wg         *sync.WaitGroup
	size       int
	running    bool
	workerfunc func(v interface{}) interface{}
	work       chan interface{}
	results    chan interface{}
	queue      []interface{}
}

func New(size uint16, results bool, workerfunc func(v interface{}) interface{}) *WorkerPool {
	wp := WorkerPool{
		mu:         new(sync.Mutex),
		size:       int(size),
		running:    false,
		workerfunc: workerfunc,
		work:       make(chan interface{}, size),
		queue:      make([]interface{}, 0),
		wg:         new(sync.WaitGroup),
	}

	if results {
		wp.results = make(chan interface{})
	}

	return &wp
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
			for w := range wp.work {
				if v := wp.workerfunc(w); wp.results != nil && v != nil {
					wp.results <- v
				}
			}
			wp.wg.Done()
			return
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
		close(wp.work)
		wp.wg.Wait()
		close(wp.results)
	}
	wp.running = false
	wp.mu.Unlock()
}

func (wp *WorkerPool) DoWork(v interface{}) {
	wp.mu.Lock()
	if wp.running {
		wp.work <- v
	} else {
		wp.queue = append(wp.queue, v)
	}
	wp.mu.Unlock()
}

func (wp *WorkerPool) Results() chan interface{} {
	return wp.results
}
