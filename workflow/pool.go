// Пул обработчиков для паралельных задач
package workflow

import (
	"sync"
	"time"
)

// Task Задача
type Task interface {
	Execute()
}

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
	self := &pool{
		limitPool: LimitPool,
		// Канал задач - буферизированный, чтобы основная программа не блокировалась при постановке задач
		tasks: make(chan Task, LimitCh),
		// Канал kill для убийства "лишних воркеров"
		kill: make(chan struct{}),
	}
	self.size++
	self.wg.Add(2)
	go self.worker()
	go self.resize()
	return self
}

// Жизненный цикл воркера
func (self *pool) worker() {
	defer func() {
		self.size--
		self.wg.Done()
	}()
	for {
		select {
		// Если есть задача, то ее нужно обработать
		// Блокируется пока канал не будет закрыт, либо не поступит новая задача
		case task, ok := <-self.tasks:
			if ok {
				task.Execute()
			} else {
				return
			}
			// Если пришел сигнал умирать, выходим
		case <-self.kill:
			return
		}
	}
}

func (self *pool) resize() {
	defer self.wg.Done()
	for 0 < self.size {
		step := cap(self.tasks) / 20
		if step*self.size < len(self.tasks) && self.size < self.limitPool {
			self.size++
			self.wg.Add(1)
			go self.worker()
		} else if 1 < self.size && len(self.tasks) <= step*(self.size-1) {
			self.kill <- struct{}{}
		}
		time.Sleep(time.Second * 1)
	}
}

// TaskAdd Добавляем задачу в пул
func (self *pool) TaskAdd(task Task) {
	self.tasks <- task
}

// Wait Завершаем работу пула
func (self *pool) Wait() {
	close(self.tasks)
	self.wg.Wait()
}
