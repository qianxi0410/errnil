package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

var count = 0

var workerCount = 0
var maxWorkers = runtime.NumCPU() * 2
var searchReq = make(chan string)
var workDone = make(chan bool)
var caculateCh = make(chan int)

func caculate(s string) int {
	return strings.Count(s, "err != nil") + strings.Count(s, "err == nil")
}

func serch(master bool, paths ...string) {
	for _, path := range paths {
		file, err := os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) || os.ErrPermission == err {
				workDone <- true
				continue
			}
			log.Fatalf("%s", err)
		}

		if file.IsDir() {
			fs, err := ioutil.ReadDir(path)
			if err != nil {
				log.Fatalf("%s", err)
			}

			for _, f := range fs {
				if workerCount < maxWorkers {
					searchReq <- (path + "/" + f.Name())
				} else {
					serch(false, path+"/"+f.Name())
				}
			}
		} else {
			if strings.HasSuffix(file.Name(), ".go") {
				s, err := ioutil.ReadFile(path)
				if err != nil {
					log.Fatalf("%s", err)
				}
				caculateCh <- caculate(string(s))
			}
		}

		if master {
			workDone <- true
		}
	}
}

func wait() {
	for {
		select {
		case path := <-searchReq:
			workerCount++
			go serch(true, path)
		case <-workDone:
			workerCount--
			if workerCount == 0 {
				return
			}
		case cnt := <-caculateCh:
			count += cnt
		}
	}
}

func main() {
	args := os.Args[1:]

	start := time.Now()

	workerCount = 1
	go serch(true, args...)
	wait()

	end := time.Now()
	fmt.Println("errnil count: ", count)
	fmt.Println("cost: ", end.Sub(start))
}
