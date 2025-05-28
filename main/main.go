package main

import (
	"context"
	"fmt"
	"peon"
	"time"
)

func main() {
	ctx := context.Background()

	worker := peon.NewWorker()
	counter := CounterJob{[]int{1, 2, 3, 4, 5}}

	// TODO this fails at a lower timeout duration than the job takes to execute
	worker.AddRepeatableJob("Counter Job", time.Second, time.Second*3, &counter)

	worker.StartWorker(ctx)
	time.Sleep(time.Second * 10)
	worker.StopWorker()
}

type CounterJob struct {
	numbers []int
}

func (c *CounterJob) Execute(ctx context.Context) error {
	for _, num := range c.numbers {
		select {
		case <-ctx.Done():
			fmt.Printf("Context timed out at number %v: %v", num, ctx.Err())
			return ctx.Err()
		default:
			fmt.Println("Num", num)
			time.Sleep(time.Second * 1)
		}
	}
	return nil
}
