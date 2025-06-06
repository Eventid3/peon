package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(time.Second * 10))

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select{
			case <- ticker.C:
				fmt.Println("Tick!")
			case <-ctx.Done():
				fmt.Println("Done!")
				return
		}
	}
}

type CounterJob struct {
	numbers []int
}

func (c *CounterJob) Execute(ctx context.Context) error {
	for _, num := range c.numbers {
		select {
		case <-ctx.Done():
			fmt.Printf("Context timed out at number %v: %v\n", num, ctx.Err())
			return ctx.Err()
		default:
			fmt.Println("Num", num)
			time.Sleep(time.Second * 1)
		}
	}
	return nil
}
