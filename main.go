package main

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Call struct {
	wg sync.WaitGroup
	res string
	err error
}

type CallsMap struct {
	m map[string]*Call
	mutex sync.Mutex
}

func NewCallsMap() CallsMap {
	return CallsMap{
		m: make(map[string]*Call),
	}
}

func (c *CallsMap) GetOrCreate(key string) (*Call, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if call, ok := c.m[key]; ok {
		return call, true
	}

	call := &Call{}
	c.m[key] = call

	return call, false
}

func (c *CallsMap) Remove(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.m, key)
}

type SingleFlight struct {
	calls CallsMap
}

func NewSingleFlight() *SingleFlight {
	return &SingleFlight {
		calls: NewCallsMap(),
	}
}

func (s *SingleFlight) Do(key string, fn func() (string, error)) (string, error) {
	call, ok := s.calls.GetOrCreate(key)

	if !ok {
		call.wg.Add(1)

		go func (){
			defer call.wg.Done()

			fmt.Println("Working")
			call.res, call.err = fn()
		}()

		defer func() {
			s.calls.Remove(key)
		}()
	}

	call.wg.Wait()

	return call.res, call.err
}

func main() {
	requestsCount := 10
	sf := NewSingleFlight()
	wg := sync.WaitGroup{}

	for i := 0; i < requestsCount; i++ {
		wg.Add(1)

		go func() {
			res, _ := sf.Do("123", longRequest)
			fmt.Printf("Gorutine %d finished with result = %s\n", i + 1, res)
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
