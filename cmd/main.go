package main

import (
	"context"
	"fmt"
	"time"

	peon "github.com/Eventid3/peon/pkg"
)

type SomeJob struct {
	numbers []int
}

func NewSomeJob() SomeJob {
	return SomeJob{[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}}
}

func (j *SomeJob) Exec(_ context.Context) error {
	for i := range j.numbers {
		j.numbers[i] = j.numbers[i] * 2
	}
	return nil
}

func main() {
	p := peon.NewPeon(1)
	job := NewSomeJob()
	err := p.RegisterJob("SomeJob", &job, time.Second*1)
	if err != nil {
		fmt.Println("Error in main.go: %w", err)
	}

	fmt.Println("State of SomeJob before starting peon:")
	fmt.Println(job)
	fmt.Println("Starting peon...")
	p.Start()
	time.Sleep(time.Second * 10)
	fmt.Println("Stopping peon...")
	fmt.Println("State of SomeJob after peon:")
	fmt.Println(job)
	p.Stop()
	fmt.Println("Terminating...")
}
