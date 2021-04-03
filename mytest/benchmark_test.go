package mytest

import (
	"testing"
)

const (
	smallValLen = 10
	largeValLen = 1000
	serverCnt   = 3
	optSet      = "/store/set"
	optGet      = "/store/get"
	optPrint    = "/log/print/"
	optIndex    = "/"
	serverIP    = "http://127.0.0.1"
)

var serverAddrs []string = []string{
	serverIP + ":8000",
	serverIP + ":8001",
	serverIP + ":8002",
}

var serverSetUrls []string = []string{
	serverAddrs[0] + optSet,
	serverAddrs[1] + optSet,
	serverAddrs[2] + optSet,
}

func BenchmarkClusterSmallValueAsync(b *testing.B) {

	if !asyncRequest(smallValLen, serverSetUrls[0:serverCnt], b) {
		b.Error("async cluster benchmark with small value length failed")
	}
}

func BenchmarkClusterSmallValue(b *testing.B) {
	if !request(smallValLen, serverSetUrls[0:serverCnt], b) {
		b.Error("cluster benchmark with small value length failed")
	}
}

func BenchmarkSmallValueAsync(b *testing.B) {

	if !asyncRequest(smallValLen, serverSetUrls[0:1], b) {
		b.Error("async benchmark for signal node with small value length failed")
	}
}

func BenchmarkSmallValue(b *testing.B) {
	if !request(smallValLen, serverSetUrls[0:1], b) {
		b.Error("benchmark for signal node with small value length failed")
	}
}
func BenchmarkLargeValueAsync(b *testing.B) {
	if !asyncRequest(largeValLen, serverSetUrls[0:1], b) {
		b.Error("async benchmark for signal node with large value length failed")
	}
}

func BenchmarkLargeValue(b *testing.B) {
	if !request(largeValLen, serverSetUrls[0:1], b) {
		b.Error("benchmark for signal node with large value length failed")
	}
}
