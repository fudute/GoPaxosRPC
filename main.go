package main

// The following implements the main Go
// package starting up the paxos server

import (
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/fudute/GoPaxos/handlers"
	"github.com/fudute/GoPaxos/paxos"
	log "github.com/sirupsen/logrus"
)

const (
	// PORT defines the port value
	// for the Paxos Server service
	PORT     = "8080"
	RPC_PORT = ":1234"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func main() {
	r := handlers.Router()
	defer paxos.DB.Close()

	log.WithFields(log.Fields{
		"port": PORT,
	}).Info("starting paxos server")

	var wg sync.WaitGroup
	wg.Add(2)

	// RESTful API
	go func() {
		http.ListenAndServe(":"+PORT, r)
		wg.Done()
	}()

	// RPC
	go func() {
		l, e := net.Listen("tcp", RPC_PORT)
		if e != nil {
			log.Fatal("listen error:", e)
		}
		http.Serve(l, nil)
		wg.Done()
	}()

	go func() {
		for {
			time.Sleep(time.Second)

			req := paxos.Request{
				Oper: paxos.NOP,
				Done: make(chan error),
			}
			paxos.GetProposerInstance().In <- req

			err := <-req.Done
			if err != nil {
				log.Printf("NOP error: ", err)
			}
		}
	}()
	paxos.InitNetwork()
	wg.Wait()
}
