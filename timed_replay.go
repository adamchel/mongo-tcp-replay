package main

import (
	"fmt"
	"sync"
	"time"
	"strconv"
)

var wg sync.WaitGroup

func replay(payload string) {
	defer wg.Done()
	const layout = "Jan 2, 2006 at 3:04:01pm (MST)"
	fmt.Println(payload + strconv.FormatInt(time.Now().UnixNano(), 10) + "before")
	timer := time.NewTimer(time.Second)
	<- timer.C
	fmt.Println(payload + strconv.FormatInt(time.Now().UnixNano(), 10))
}

func executor() {
	wg.Add(3)
	go replay("test1")
	go replay("test2")
	go replay("test3")
	wg.Wait()
}