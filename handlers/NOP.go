package handlers

import (
	"log"
	"net/http"

	"github.com/fudute/GoPaxos/paxos"
	"github.com/gin-gonic/gin"
)

func SendNop(c *gin.Context) {
	req := paxos.Request{
		Oper: paxos.NOP,
		Done: make(chan error),
	}
	paxos.GetProposerInstance().In <- req

	err := <-req.Done
	if err != nil {
		log.Printf("NOP error: ", err)
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
	}
	c.JSON(http.StatusOK, nil)
}
