package main

import "fmt"

func (n *Node) StartElection() {
	fmt.Println("Node", n.ID, "starting election")
	var reply bool
	n.CallNextNeighbor("RPCHandler.Election", n.ID, &reply)
}

func (h *RPCHandler) Election(candidateID int, reply *bool) error {
	n := h.node
	*reply = true

	if candidateID > n.ID {
		fmt.Println("Node", n.ID, "forwarding Election message (ID:", candidateID, ")")
		go n.CallNextNeighbor("RPCHandler.Election", candidateID, reply)
	} else if candidateID < n.ID {
		fmt.Println("Node", n.ID, "replacing Election ID", candidateID, "with My ID", n.ID)
		go n.CallNextNeighbor("RPCHandler.Election", n.ID, reply)
	} else if candidateID == n.ID {
		fmt.Println("Node", n.ID, "Election message returned! Becoming LEADER.")
		n.BecomeLeader()
	}

	return nil
}

func (n *Node) BecomeLeader() {
	n.mu.Lock()
	n.LeaderID = n.ID
	n.mu.Unlock()

	fmt.Println("Node", n.ID, "is now the COORDINATOR")
	
	// Inform everyone in the ring
	var ack bool
	n.CallNextNeighbor("RPCHandler.Coordinator", n.ID, &ack)
}

func (h *RPCHandler) Coordinator(leaderID int, ack *bool) error {
	n := h.node
	*ack = true

	n.mu.Lock()
	if n.LeaderID == leaderID {
		n.mu.Unlock()
		return nil // Already informed
	}
	n.LeaderID = leaderID
	n.mu.Unlock()

	fmt.Println("Node", n.ID, "recognized Node", leaderID, "as the new LEADER")

	return nil
}
