package handlers

import (
	"fmt"
	"net/http"

	"github.com/fudute/GoPaxos/paxos"
	"github.com/gin-gonic/gin"
)

type Pair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// SetValue is the HTTP handler to process the incoming write message
// It starts the paxos round in the cluster and makes the incoming node as the leader
func SetValue(c *gin.Context) {
	var p Pair
	c.Bind(&p)
	fmt.Println(p)
	// 这里启动一个新的Instance
	err := paxos.StartNewInstance(paxos.SET, p.Key, p.Value)

	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
	}
	c.JSON(http.StatusOK, nil)
}