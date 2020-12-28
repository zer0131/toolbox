package batchprocessor

import (
	"errors"
	"time"
)

const (
	DefaultSize  = 10
	DefaultDelay = 5
)

// Operator 批量处理接口
type Operator interface {

	// BatchProcessor 是需要实现的正常批量处理接口
	BatchProcessor(isAsync bool, msg []interface{}) (err error)

	// ErrorHandler 遇到错误后需要如何解决
	ErrorHandler(err error, msg []interface{})
}

type Processor struct {
	InQueue   chan interface{}
	OpQueue   []interface{}
	Size      int
	Delay     time.Duration
	Operation Operator
	IsAsync   bool
}

func NewProcessor(inQueue chan interface{}, delay time.Duration, size int, isAsync bool, operator Operator) (*Processor, error) {
	if inQueue == nil {
		return nil, errors.New("inQueue nil")
	}
	if size == 0 {
		size = DefaultSize
	}
	if delay == 0 {
		delay = DefaultDelay
	}
	return &Processor{
		InQueue:   inQueue,
		OpQueue:   make([]interface{}, 0),
		Size:      size,
		Delay:     delay,
		Operation: operator,
		IsAsync:   isAsync,
	}, nil
}

func (s Processor) Run() {
	timer := time.NewTimer(0)
	if !timer.Stop() {
		<-timer.C
	}
	defer timer.Stop()
	for {
		select {
		case msg := <-s.InQueue:
			s.OpQueue = append(s.OpQueue, msg)
			if len(s.OpQueue) != s.Size {
				if len(s.OpQueue) == 1 {
					timer.Reset(s.Delay)
				}
				break
			}
			if err := s.Operation.BatchProcessor(s.IsAsync, s.OpQueue); err != nil {
				s.Operation.ErrorHandler(err, s.OpQueue)
			}
			if !timer.Stop() {
				<-timer.C
			}
			s.OpQueue = make([]interface{}, 0)
		case <-timer.C:
			if err := s.Operation.BatchProcessor(s.IsAsync, s.OpQueue); err != nil {
				s.Operation.ErrorHandler(err, s.OpQueue)
			}
			s.OpQueue = make([]interface{}, 0)
		}
	}
}
