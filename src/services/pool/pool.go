package pool

import (
	"sync"
)

type Pool struct {
	Channel chan *Task
	Tasks   []*Task

	concurrency int
	tasksChan   chan *Task
	wg          sync.WaitGroup
	done        int
}

type Communication struct {
}

func NewPool(concurrency int) *Pool {
	return &Pool{
		Channel:     make(chan *Task, concurrency),
		concurrency: concurrency,
		tasksChan:   make(chan *Task),
	}
}

func (p *Pool) AddTask(job *Task) {
	p.Tasks = append(p.Tasks, job)
}

func (p *Pool) Run() {
	for i := 0; i < p.concurrency; i++ {
		go p.work()
	}

	p.wg.Add(len(p.Tasks))
	for _, task := range p.Tasks {
		p.tasksChan <- task
	}

	// all workers return
	close(p.tasksChan)

	p.wg.Wait()
}

func (p *Pool) ReceiveAnswer() {
	p.done++

	if p.done == p.concurrency {
		close(p.Channel)
	}
}

func (p *Pool) work() {
	for task := range p.tasksChan {
		task.Run(&p.wg)
	}
}
