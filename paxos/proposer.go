package paxos

import (
	"errors"
	"fmt"
	"log"
	"net/rpc"
	"time"

	"github.com/fudute/GoPaxos/sm"
)

type Proposer struct {
	ServerID int
	LogIndex int // 记录最小的没有Chosen的logIndex
	Clients  []*rpc.Client
}

var proposer Proposer

type PrepareRequest struct {
	Index       int
	ProposalNum int
}

type PrepareResponse struct {
	AcceptedProposal int
	AcceptedValue    string
}

type AcceptRequest struct {
	Index         int
	ProposalNum   int
	ProposalValue string
}
type AcceptResponse struct {
	MinProposal int
}

func InitNetwork() {
	time.Sleep(time.Second * 3)
	peers := GetPeerList()
	for _, peer := range peers {
		addr := fmt.Sprintf("%s:1234", peer)

		client, err := rpc.DialHTTP("tcp", addr)
		if err != nil {
			log.Println("rpc dial error, addr =", addr)
		}
		proposer.Clients = append(proposer.Clients, client)
	}
}

// Prepare starts a Paxos round sending
// a prepare request to all the Paxos
// peers including itself
// 这里执行multiPaxos中逐渐一次向后选择logIndex的逻辑，具体的数据传输在doPrepare中完成
// oper取值范围为SET DELETE
func StartNewInstance(oper int, key string, value string) error {

	var command string
	if oper == SET {
		command = "SET " + key + " " + value
	} else if oper == DELETE {
		command = "DELETE " + key
	} else {
		return ErrorUnkonwCommand
	}

	log.Println("command :", command)

	// 循环获得第一个没有被Chosen的index，直到成功Prepare
	isCommited := false
	for !isCommited {
		var err error
		le, err := DB.ReadLog(proposer.LogIndex)
		for err == nil && le.isCommited {
			proposer.LogIndex++
			le, err = DB.ReadLog(proposer.LogIndex)
		}

		if err != nil && err != ErrorNotFound {
			log.Fatal("read error", err)
		}

		isCommited, err = DoPrepare(proposer.LogIndex, command, 0)
		if err != nil {
			return err
		}
		proposer.LogIndex++
	}
	return nil
}

// DoPrepare可以确定index位置的值
// 这里的value格式为 [SET key value]或者[DELETE key]
// 如果成功提交当前value，返回true，否则返回false
func DoPrepare(index int, value string, minProposal int) (bool, error) {

	log.Printf("start prepare at index = %d\n", index)

	proposalNum := GenerateProposalNum(minProposal, proposer.ServerID)

	curValue := value   // 记录当前index的value，有可能之后会变更
	curMaxProposal := 0 // 记录当前看到的最大的accpetedProposal
	preparedPeersCount := 0
	majorityPeersCount := len(proposer.Clients)/2 + 1

	isCommited := true

	for _, client := range proposer.Clients {

		req := &PrepareRequest{
			Index:       proposer.LogIndex,
			ProposalNum: proposalNum,
		}
		resp := &PrepareResponse{}
		err := client.Call("Acceptor.Prepare", req, resp)
		if err != nil {
			log.Println("rpc failed", err)
			continue
		}
		preparedPeersCount++

		if resp.AcceptedValue != "" && resp.AcceptedProposal > curMaxProposal {
			curMaxProposal = resp.AcceptedProposal
			curValue = resp.AcceptedValue
			isCommited = false
		}
		// Break when majorityPeersCount reached
		if preparedPeersCount >= majorityPeersCount {
			DoAccept(index, proposalNum, curValue)
			break
		}
	}
	if preparedPeersCount < majorityPeersCount {
		return false, errors.New("majority consensus not obtained")
	}
	return isCommited, nil
}

// DoAccept starts the accept phase sending
// an accept request to all the Paxos
// peers including itself
func DoAccept(index, proposalNum int, proposalValue string) error {

	log.Printf("start Accept at index = %d, proposalNum = %d, value = %s\n", index, proposalNum, proposalValue)

	acceptedPeersCount := 0
	majorityPeersCount := len(proposer.Clients)/2 + 1

	for _, client := range proposer.Clients {
		req := &AcceptRequest{
			Index:         index,
			ProposalNum:   proposalNum,
			ProposalValue: proposalValue,
		}
		resp := &AcceptResponse{}

		err := client.Call("Acceptor.Accept", req, resp)
		if err != nil {
			log.Println("rpc failed", err)
			continue
		}
		acceptedPeersCount++

		if resp.MinProposal > proposalNum {
			// 从新prepare，选择更大的proposalNum
			DoPrepare(index, proposalValue, resp.MinProposal)
			return nil
		}

		// Break when majorityPeersCount reached
		if acceptedPeersCount >= majorityPeersCount {
			// 到这里就可以确定已经被大多是接受了，那么可以直接提交到状态机中运行
			le, err := DB.ReadLog(index)
			if err != nil && err != ErrorNotFound {
				log.Fatal("log read error ", err)
			}
			if err == ErrorNotFound {
				le = &LogEntry{
					isCommited:       true,
					AcceptedProposal: proposalNum,
					AcceptedValue:    proposalValue,
					MinProposal:      proposalNum,
				}
			}
			le.isCommited = true
			DB.WriteLog(index, le)
			sm.GetKVStatMachineInstance().Execute(proposalValue)
			break
		}
	}

	if acceptedPeersCount < majorityPeersCount {
		return errors.New("majority consensus not obtained")
	}

	return nil
}
