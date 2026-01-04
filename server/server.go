package server

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"time"
	//"tcp/test/constants"
	//"tcp/test/constants"
)

const (
	MaxWorkers = 10 
	MaxQueue   = 100 
	MaxConn    = 50  
)

type Job struct {
	Buf  []byte
	Addr string
}

var (
	pendingJobs   int64
	batchStart    time.Time
	batchMutex    sync.Mutex
	metricsLog    *log.Logger
	jobQueue      chan Job
)

func workerProcess(id int) {
	for job := range jobQueue {
		Process(job.Buf, job.Addr)

		batchMutex.Lock()
		pendingJobs--
		if pendingJobs == 0 {
			duration := time.Since(batchStart)
			msg := fmt.Sprintf("Batch FINISHED! Duration: %v", duration)
			fmt.Println(msg)
			if metricsLog != nil {
				metricsLog.Println(msg)
			}
		}
		batchMutex.Unlock()
	}
}

func work(srv net.Conn) {

	defer srv.Close()
	clientAddr := srv.RemoteAddr().String()

	for {
		var hdr [4]byte
		if _, err := io.ReadFull(srv, hdr[:]); err != nil {
			if err == io.EOF {
				return
			}
			log.Println("read size error:", err)
			return
		}

		size := int(binary.BigEndian.Uint32(hdr[:]))
		buf := make([]byte, size)
		if _, err := io.ReadFull(srv, buf); err != nil {
			log.Println("read data error:", err)
			return
		}

		batchMutex.Lock()
		if pendingJobs == 0 {
			batchStart = time.Now()
			msg := fmt.Sprintf("Batch STARTED at %v", batchStart)
			fmt.Println(msg)
			if metricsLog != nil {
				metricsLog.Println(msg)
			}
		}
		pendingJobs++
		batchMutex.Unlock()
		jobQueue <- Job{Buf: buf, Addr: clientAddr}
	}
}

func Server() {
	f, err := os.OpenFile("metrics.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	metricsLog = log.New(f, "", log.LstdFlags)
	jobQueue = make(chan Job, MaxQueue)
	for i := 0; i < MaxWorkers; i++ {
		go workerProcess(i)
	}
	connLimiter := make(chan struct{}, MaxConn)
	conn, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Print(err)
		return
	}
	defer conn.Close()
	log.Printf("Server started on :8082 (Workers: %d, MaxConn: %d)", MaxWorkers, MaxConn)
	for {
		connLimiter <- struct{}{}
		srv, err := conn.Accept()
		if err != nil {
			log.Fatal(err)
			<-connLimiter
			continue
		}
		go work(srv)
		<-connLimiter
	}
}
