package raft

import (
	"log"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func init() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)
	rand.Seed(time.Now().UnixNano())
}

type Cluster struct {
	mu sync.Mutex

	// nodes is a list of all the raft servers participating in a nodes.
	nodes []*Server

	// Maintains whether server is partioned or not
	connected []bool

	n int

	t *testing.T
}

func NewCluster(t *testing.T, n int) *Cluster {
	ns := make([]*Server, n)
	connected := make([]bool, n)
	ready := make(chan interface{})

	// Create all Servers in this nodes, assign ids and peer ids.
	for i := 0; i < n; i++ {
		peersIds := make([]int, 0)
		for p := 0; p < n; p++ {
			if p != i {
				peersIds = append(peersIds, p)
			}
		}

		// if i == 2 {
		// 	ns[i] = NewServer(i, peersIds, ready, 100)
		// } else {
		// 	ns[i] = NewServer(i, peersIds, ready, 20)
		// }

		ns[i] = NewServer(i, peersIds, ready, 20)
		ns[i].Serve()
	}

	// Connect all peers to each other.
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if i != j {
				ns[i].ConnectToPeer(j, ns[j].GetCurrentAddress())
			}
		}
		connected[i] = true
	}
	close(ready) // Channel!

	this := &Cluster{
		nodes:     ns,
		connected: connected,
		n:         n,
		t:         t,
	}
	return this
}

func (this *Cluster) Shutdown() {
	for i := 0; i < this.n; i++ {
		this.nodes[i].DisconnectAll()
		this.connected[i] = false
	}
	for i := 0; i < this.n; i++ {
		this.nodes[i].Shutdown()
	}
}

// DisconnectPeer disconnects a server from all other servers in the nodes.
func (this *Cluster) DisconnectPeer(id int) {
	testing_log("Disconnecting %d", id)
	this.nodes[id].DisconnectAll()
	for j := 0; j < this.n; j++ {
		if j != id {
			this.nodes[j].DisconnectPeer(id)
		}
	}
	this.connected[id] = false

	this.nodes[id].raftLogic.mu.Lock()
	this.nodes[id].raftLogic.LOG_ENTRIES = false
	this.nodes[id].raftLogic.mu.Unlock()
}

// ReconnectPeer connects a server to all other servers in the nodes.
func (this *Cluster) ReconnectPeer(id int) {
	testing_log("Reconnecting %d", id)
	for j := 0; j < this.n; j++ {
		if j != id && this.connected[j] {

			if err := this.nodes[id].ConnectToPeer(j, this.nodes[j].GetCurrentAddress()); err != nil {
				this.t.Fatal(err)
			}
			if err := this.nodes[j].ConnectToPeer(id, this.nodes[id].GetCurrentAddress()); err != nil {
				this.t.Fatal(err)
			}
		}
	}
	this.connected[id] = true

	this.nodes[id].raftLogic.mu.Lock()
	this.nodes[id].raftLogic.LOG_ENTRIES = true
	this.nodes[id].raftLogic.mu.Unlock()
}

/* getClusterLeader checks that only a single server thinks it's the leader.
Returns the leader's id and term. It retries several times if no leader is
identified yet. */
func (this *Cluster) getClusterLeader() int {
	for r := 0; r < 20; r++ {
		leaderId := -1
		for i := 0; i < this.n; i++ {
			if this.connected[i] {
				_, _, isLeader := this.nodes[i].raftLogic.GetNodeState()
				if isLeader {
					if leaderId < 0 {
						leaderId = i
					} else {
						this.t.Fatalf("Somehow have more than one leader!!!!!")
					}
				}
			}
		}
		if leaderId >= 0 {
			return leaderId
		}
		sleepMs(750)
	}

	this.t.Fatalf("leader not found")
	return -1
}

// SubmitClientCommand submits the command to serverId.
func (this *Cluster) SubmitClientCommand(serverId int, cmd interface{}) bool {
	return this.nodes[serverId].raftLogic.ReceiveClientCommand(cmd)
}

func testing_log(format string, a ...interface{}) {
	format = "[ACTION] " + format
	log.Printf(format, a...)
}

func sleepMs(n int) {
	time.Sleep(time.Duration(n) * time.Millisecond)
}
