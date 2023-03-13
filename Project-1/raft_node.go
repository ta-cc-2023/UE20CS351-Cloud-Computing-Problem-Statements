package raft

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

const LogHeartbeatMessages = false
const LogVoteRequestMessages = true

type LogEntry struct {
	Command interface{}
	Term    int
}

// Main Raft Data Structure
type RaftNode struct {
	mu sync.Mutex

	id       int
	peersIds []int

	// Persistent state on all servers
	currentTerm int
	votedFor    int
	log         []LogEntry

	// Volatile state on all servers
	commitIndex int
	lastApplied int

	// Volatile Raft state on leaders
	nextIndex  map[int]int
	matchIndex map[int]int

	// Utility States
	state                        string
	lastElectionTimerStartedTime time.Time
	notifyToApplyCommit          chan int
	LOG_ENTRIES                  bool
	filePath                     string

	// Networking Component, do NOT worry about this whatsoever.
	server *Server
}

// Constructor for RaftNodes
func NewRaftNode(id int, peersIds []int, server *Server, ready <-chan interface{}) *RaftNode {
	this := new(RaftNode)

	this.server = server
	this.notifyToApplyCommit = make(chan int, 16)

	this.id = id
	this.peersIds = peersIds

	this.votedFor = -1
	this.currentTerm = 0

	this.commitIndex = -1
	this.lastApplied = -1

	this.nextIndex = make(map[int]int)
	this.matchIndex = make(map[int]int)

	this.state = "Follower"

	this.LOG_ENTRIES = true

	this.filePath = "NodeLogs/" + strconv.Itoa(this.id)
	f, _ := os.Create(this.filePath)
	f.Close()

	go func() {
		// Signalled when all servers are up and running, ready to receive RPCs;
		// Again, this is code you don't need to worry about.
		<-ready

		this.mu.Lock()
		this.lastElectionTimerStartedTime = time.Now()
		this.mu.Unlock()

		this.startElectionTimer()
	}()

	go this.applyCommitedLogEntries() // Fire off a watcher to apply any committed entries

	return this
}

// This function implements the 'application' of a query to the leader
// This is the function that also writes queries accepted by the leader to files
// to observe as output
func (this *RaftNode) applyCommitedLogEntries() {
	for range this.notifyToApplyCommit {
		this.mu.Lock()

		var entriesToApply []LogEntry

		if this.commitIndex > this.lastApplied {
			entriesToApply = this.log[this.lastApplied+1 : this.commitIndex+1]
		}

		f, _ := os.OpenFile(this.filePath, os.O_APPEND|os.O_WRONLY, 0644)

		defer f.Close()

		for i, entry := range entriesToApply {
			strentry := fmt.Sprintf("%s; T:[%d]; I:[%d]", entry.Command, this.currentTerm, this.commitIndex+i)
			f.WriteString(strentry)
			f.WriteString("\n")
		}

		this.lastApplied = this.commitIndex
		this.mu.Unlock()
	}

	this.write_log("applyCommitedLogEntries done")
}

/* UTILITY FUNCTIONS */

// GetNodeState reports the state of this RN.
func (this *RaftNode) GetNodeState() (id int, term int, isLeader bool) {
	this.mu.Lock()
	defer this.mu.Unlock()
	return this.id, this.currentTerm, this.state == "Leader"
}

// Kills a RaftNode and sets its state to Dead
func (this *RaftNode) KillNode() {
	this.mu.Lock()
	defer this.mu.Unlock()
	this.state = "Dead"
	this.write_log("KILLED")
	close(this.notifyToApplyCommit)
}

// This function logs all messages to the terminal
func (this *RaftNode) write_log(format string, args ...interface{}) {
	if this.LOG_ENTRIES {
		format = fmt.Sprintf("AT NODE %d: ", this.id) + format
		log.Printf(format, args...)
	}
}
