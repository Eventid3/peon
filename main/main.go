package main

import (
	"fmt"
	"github.com/Eventid3/peon"
	"time"
)

func main() {
	cj1 := &CounterJob{numbers: []int{1, 2, 3}}
	cj2 := &CounterJob{numbers: []int{4, 5, 6}}
	workforce := peon.NewWorkforce("Counters")
	workforce.AddRepeatableJob("Counter Job 1", time.Second*5, cj1.Execute)
	workforce.AddRepeatableJob("Counter Job 2", time.Second*10, cj2.Execute)
	workforce.StartAllWorkers()

	time.Sleep(time.Second * 30)
}

type CounterJob struct {
	numbers []int
}

func (c *CounterJob) Execute() error {
	for _, num := range c.numbers {
		if num == 6 {
			return fmt.Errorf("number %d is not allowed", num)
		}
		fmt.Println("Num", num)
		time.Sleep(time.Millisecond * 300)
	}
	return nil
}
