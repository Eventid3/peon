package main

import (
	"fmt"
	"peon"
	"time"
)

func main() {
	worker := peon.NewWorker()
	worker.AddRepeatableJob("Counter Job", time.Second*5, &CounterJob{numbers: []int{1, 2, 3}})

	worker.StartWorker()
	time.Sleep(time.Second * 14)
	worker.PauseWorker()
	time.Sleep(time.Second * 2)
	worker.ResumeWorker()
	time.Sleep(time.Second * 7)
	worker.StopWorker()
}

type CounterJob struct {
	numbers []int
}

func (c *CounterJob) Execute() error {
	for _, num := range c.numbers {
		fmt.Println("Num", num)
		time.Sleep(time.Millisecond * 300)
	}
	return nil
}
