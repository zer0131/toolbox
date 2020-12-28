package cron

import (
	"context"
	"errors"
)

type Worker interface {
	Name() string
	Run(ctx context.Context) error
}

var workerList []Worker

func Register(worker Worker) {
	workerList = append(workerList, worker)
}

func WorkerList() []Worker {
	return workerList
}

func Run(ctx context.Context, name string, body []byte) error {
	// TODO body暂时没有用，包含task_id
	for _, w := range workerList {
		if w.Name() == name {
			return w.Run(ctx)
		}
	}
	return errors.New("Can not find worker called " + name)
}
