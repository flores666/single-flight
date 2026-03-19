package main

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

type SingleFlight struct {
	wg sync.WaitGroup
	isActive atomic.Bool
	res string
	err error
}

func NewSingleFlight() *SingleFlight {
	return &SingleFlight {}
}

func (s *SingleFlight) Do(fn func() (string, error)) (string, error) {
	s.wg.Add(1)

	go func (){
		if !s.isActive.Load() {
			s.isActive.Store(true)
			fmt.Println("Working")
			s.res, s.err = fn()
			s.isActive.Store(false)
		}

		s.wg.Done()
	}()

	s.wg.Wait()

	return s.res, s.err
}

func main() {
	requestsCount := 10
	sf := NewSingleFlight()
	wg := sync.WaitGroup{}

	for i := 0; i < requestsCount; i++ {
		wg.Add(1)

		go func() {
			fmt.Printf("Gorutine n %d started\n", i + 1)
			res, _ := sf.Do(longRequest)
			fmt.Println(res)
			wg.Done()
		}()
	}

	wg.Wait()
}

func longRequest() (string, error) {
	result := "Hello World"

	time.Sleep(1 * time.Second)

	if rand.Intn(10) + 1 <= 2 {
		return "Error", errors.New("Error")
	}

	return result, nil
}
