// Пул обработчиков для паралельных задач
package workflow

import (
	"sync"
	"time"
)

// Pool - структура, нам потребуется Мутекс, для гарантий атомарности изменений самого объекта
// Канал входящих задач
// Канал отмены, для завершения работы
// WaitGroup для контроля завершнеия работ
type pool struct {
	size      int
	limitPool int
	tasks     chan Task
	kill      chan struct{}
	wg        sync.WaitGroup
}

// NewPool Создаем пул воркеров указанного размера
func NewPool(LimitCh, LimitPool int) *pool {
	p := &pool{
		limitPool: LimitPool,
		// Канал задач - буферизированный, чтобы основная программа не блокировалась при постановке задач
		tasks: make(chan Task, LimitCh),
		// Канал kill для убийства "лишних воркеров"
		kill: make(chan struct{}),
	}
	p.size++
	p.wg.Add(2)
	go p.worker()
	go p.resize()
	return p
}

// Жизненный цикл воркера
func (p *pool) worker() {
	defer func() {
		p.size--
		p.wg.Done()
	}()
	for {
		select {
		// Если есть задача, то ее нужно обработать
		// Блокируется пока канал не будет закрыт, либо не поступит новая задача
		case task, ok := <-p.tasks:
			if ok {
				task.Execute()
			} else {
				return
			}
			// Если пришел сигнал умирать, выходим
		case <-p.kill:
			return
		}
	}
}

func (p *pool) resize() {
	defer p.wg.Done()
	for 0 < p.size {
		step := cap(p.tasks) / 20
		if step*p.size < len(p.tasks) && p.size < p.limitPool {
			p.size++
			p.wg.Add(1)
			go p.worker()
		} else if 1 < p.size && len(p.tasks) <= step*(p.size-1) {
			p.kill <- struct{}{}
		}
		time.Sleep(time.Second * 1)
	}
}

// TaskAdd Добавляем задачу в пул
func (p *pool) TaskAdd(task Task) {
	p.tasks <- task
}

// Wait Завершаем работу пула
func (p *pool) Wait() {
	close(p.tasks)
	p.wg.Wait()
}
