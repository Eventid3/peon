package peon

import (
	"testing"
	"time"
)

func TestWorker_AddRepeatableJob(t *testing.T) {
	uut := NewWorker()

	mockJob := &MockJob{false, 0}
	uut.AddRepeatableJob("MockJob", time.Second*1, mockJob)
	uut.StartWorker()

	time.Sleep(time.Second * 4)

	uut.StopWorker()

	if !mockJob.WasExecuted {
		t.Error("Expected job to be executed, but it wasn't")
	}

	t.Logf("Job executed %v times.", mockJob.NumExecuted)
}

type MockJob struct {
	WasExecuted bool
	NumExecuted int
}

func (j *MockJob) Execute() {
	j.WasExecuted = true
	j.NumExecuted++
}
