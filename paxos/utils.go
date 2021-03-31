package paxos

import (
	"os"
	"strings"
)

const SHIFT = 3

// GenerateProposalNum Generates a proposal number out of minProposalNum and Server ID
func GenerateProposalNum(minProposalNum, ID int) int {
	return (minProposalNum>>SHIFT+1)<<SHIFT + ID
}

// GetPeerList Obtains Peer List
// From Environment Variable
func GetPeerList() []string {
	return strings.Split(os.Getenv("PEERS"), ",")
}

// GetNetwork Obtains Network
// From Environment Variable
func GetNetwork() string {
	return os.Getenv("NETWORK") + ":1234"
}
