package main

import (
	"fmt"
	"sync"
	"time"
	"strconv"
)

var wg sync.WaitGroup

func rply(payload string) {
	fmt.Println("Before wg" + payload)
	wg.Done()
	fmt.Println("After wg" + payload)
	// const layout = "Jan 2, 2006 at 3:04:01pm (MST)"
	// startTime := time.Now().UnixNano()
	// fmt.Println(payload + strconv.FormatInt(startTime, 10) + "before")
	// timer := time.NewTimer(time.Second)
	// <- timer.C
	// endTime := time.Now().UnixNano()
	// fmt.Println(payload + strconv.FormatInt(endTime, 10))
	// fmt.Println(strconv.FormatInt(endTime - startTime, 10))
}

func executor() {
	wg.Add(3)
	go rply("test1")
	go rply("test2")
	go rply("test3")
	wg.Wait()
}