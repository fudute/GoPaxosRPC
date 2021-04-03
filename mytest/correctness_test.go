package mytest

import (
	"net/http"
	"testing"
)

func TestPrintLog(t *testing.T) {

	logFileName := "log2"
	for _, addr := range serverAddrs {
		resp, err := http.Get(addr + optPrint + logFileName)
		if err != nil || resp.StatusCode != http.StatusOK {
			t.Errorf("request for printing log failed with msg: %v \n", err)
		}
	}
}
