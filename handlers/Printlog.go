package handlers

import (
	"fmt"

	"github.com/fudute/GoPaxos/paxos"
	"github.com/gin-gonic/gin"
)

func PrintLog(c *gin.Context) {
	fileName := c.Param("path")
	if fileName == "" {
		fileName = "paxos_log_1.txt"
	}
	fmt.Printf("try to write log to file %v\n", fileName)
	paxos.DB.PrintLog(fileName)

}
