package paxos

import (
	"log"
	"net/rpc"
)

type Acceptor struct {
}

var acceptor = &Acceptor{}

func init() {
	rpc.Register(acceptor)
	rpc.HandleHTTP()
}

func GetAcceptorInstance() *Acceptor {
	return acceptor
}

// Prepare return acceptedProposal and acceptedValue
func (acceptor *Acceptor) Prepare(req *PrepareRequest, resp *PrepareResponse) error {

	log.Printf("receive Prepare at index = %d, proposalNum = %d\n", req.Index, req.ProposalNum)

	entry, err := DB.ReadLog(req.Index)
	if err != nil {
		if err == ErrorNotFound {
			resp.AcceptedProposal = 0
			resp.AcceptedValue = ""
			return nil
		}
		log.Fatal("ReadLog error ", err)
	}
	if req.ProposalNum > entry.MinProposal {
		entry.MinProposal = req.ProposalNum
		// 写回磁盘
		DB.WriteLog(req.Index, entry)
	}
	resp.AcceptedProposal = entry.AcceptedProposal
	resp.AcceptedValue = entry.AcceptedValue

	return nil
}

// Accept return minProposal
func (acceptor *Acceptor) Accept(req *AcceptRequest, resp *AcceptResponse) error {

	log.Printf("receive Accept at index = %d, proposalNum = %d, proposalValue = %s\n", req.Index, req.ProposalNum, req.ProposalValue)

	entry, err := DB.ReadLog(req.Index)
	if err != nil && err != ErrorNotFound {
		log.Fatal("ReadLog error ", err)
	}

	if err == ErrorNotFound {
		entry = &LogEntry{}
	}
	if req.ProposalNum > entry.MinProposal {

		log.Printf("accepted proposal, index = %d, proposalNum = %d, proposalValue = %s\n", req.Index, req.ProposalNum, req.ProposalValue)

		entry.MinProposal = req.ProposalNum
		entry.AcceptedProposal = req.ProposalNum
		entry.AcceptedValue = req.ProposalValue

		DB.WriteLog(req.Index, entry)
	} else {
		log.Printf("deny proposal, index = %d, proposalNum = %d, minProposal = %d\n", req.Index, req.ProposalNum, entry.MinProposal)
	}

	resp.MinProposal = entry.AcceptedProposal

	return nil
}
