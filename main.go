package main

import (
	"fmt"
	"sync"
	"time"
)

type Ttype struct {
	id         int
	cT         string // время создания
	fT         string // время выполнения
	taskRESULT []byte
}

func main() {
	taskCreator := func(a chan<- Ttype) {
		defer close(a)
		for {
			ft := time.Now().Format(time.RFC3339)
			if time.Now().Nanosecond()%2 > 0 { 
				ft = "Some error occurred"
			}
			a <- Ttype{cT: ft, id: int(time.Now().Unix())}
		}
	}

	superChan := make(chan Ttype, 10)
	var wg sync.WaitGroup

	workerPool := func(a <-chan Ttype, done chan<- Ttype, errChan chan<- error, wg *sync.WaitGroup) {
		defer wg.Done()
		for t := range a {
			t = taskWorker(t)
			if string(t.taskRESULT) == "task has been succeeded" {
				done <- t
			} else {
				errChan <- fmt.Errorf("Task id %d time %s, error %s", t.id, t.cT, string(t.taskRESULT))
			}
		}
	}

	doneTasks := make(chan Ttype)
	undoneTasks := make(chan error)

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go workerPool(superChan, doneTasks, undoneTasks, &wg)
	}

	go taskCreator(superChan)

	go func() {
		wg.Wait()
		close(doneTasks)
		close(undoneTasks)
	}()

	fmt.Println("Errors:")
	for err := range undoneTasks {
		fmt.Println(err)
	}

	fmt.Println("Done tasks:")
	for result := range doneTasks {
		fmt.Println(result.id)
	}
}

func taskWorker(t Ttype) Ttype {
	tt, _ := time.Parse(time.RFC3339, t.cT)
	if tt.After(time.Now().Add(-20 * time.Second)) {
		t.taskRESULT = []byte("task has been succeeded")
	} else {
		t.taskRESULT = []byte("something went wrong")
	}
	t.fT = time.Now().Format(time.RFC3339Nano)
	time.Sleep(time.Millisecond * 150)
	return t
}
