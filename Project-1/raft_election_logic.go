package raft

import (
	"math/rand"
	"time"
)

/* startElectionTimer implements an election timer. It should be launched whenever
we want to start a timer towards becoming a candidate in a new election.
This function runs as a go routine */
func (this *RaftNode) startElectionTimer() {
	timeoutDuration := time.Duration(3000+rand.Intn(3000)) * time.Millisecond
	this.mu.Lock()
	termStarted := this.currentTerm
	this.mu.Unlock()
	this.write_log("Election timer started: %v, with term=%d", timeoutDuration, termStarted)

	// Keep checking for a resolution
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()
	for {
		<-ticker.C

		this.mu.Lock()

		// if node has become a leader
		if this.state != "Candidate" && this.state != "Follower" {
			this.mu.Unlock()
			return
		}

		// if node received requestVote or appendEntries of a higher term and updated itself
		if termStarted != this.currentTerm {
			this.mu.Unlock()
			return
		}

		// Start an election if we haven't heard from a leader or haven't voted for someone for the duration of the timeout.
		if elapsed := time.Since(this.lastElectionTimerStartedTime); elapsed >= timeoutDuration {
			this.startElection()
			this.mu.Unlock()
			return
		}
		this.mu.Unlock()
	}
}

// startElection starts a new election with this RN as a candidate.
func (this *RaftNode) startElection() {
	this.state = "Candidate"
	this.currentTerm += 1
	termWhenVoteRequested := this.currentTerm
	this.lastElectionTimerStartedTime = time.Now()
	this.votedFor = this.id
	this.write_log("became Candidate with term=%d;", termWhenVoteRequested)

	votesReceived := 1

	// Send RequestVote RPCs to all other servers concurrently.
	for _, peerId := range this.peersIds {
		go func(peerId int) {
			this.mu.Lock()
			var LastLogIndexWhenVoteRequested, LastLogTermWhenVoteRequested int

			if len(this.log) > 0 {
				lastIndex := len(this.log) - 1
				LastLogIndexWhenVoteRequested, LastLogTermWhenVoteRequested = lastIndex, this.log[lastIndex].Term
			} else {
				LastLogIndexWhenVoteRequested, LastLogTermWhenVoteRequested = -1, -1
			}
			this.mu.Unlock()

			args := RequestVoteArgs{
				Term:         termWhenVoteRequested,
				CandidateId:  this.id,
				LastLogIndex: LastLogIndexWhenVoteRequested,
				LastLogTerm:  LastLogTermWhenVoteRequested,

				Latency: rand.Intn(500), // Ignore Latency.
			}

			if LogVoteRequestMessages {
				this.write_log("sending RequestVote to %d: %+v", peerId, args)
			}

			var reply RequestVoteReply
			if err := this.server.SendRPCCallTo(peerId, "RaftNode.RequestVote", args, &reply); err == nil {
				this.mu.Lock()
				defer this.mu.Unlock()
				if LogVoteRequestMessages {
					this.write_log("received RequestVoteReply from %d: %+v", peerId, reply)
				}
				if this.state != "Candidate" {
					this.write_log("State changed from Candidate to %s", this.state)
					return
				}

				// IMPLEMENT HANDLING THE VOTEREQUEST's REPLY;
				// You probably need to have implemented becomeFollower before this.

				//-------------------------------------------------------------------------------------------/
				if reply.Term > {
					// TODO
				} else if reply.Term ==  {
					// TODO
				}
				//-------------------------------------------------------------------------------------------/

			}
		}(peerId)
	}

	// Run another election timer, in case this election is not successful.
	go this.startElectionTimer()
}

// becomeFollower sets a node to be a follower and resets its state.
func (this *RaftNode) becomeFollower(term int) {
	this.write_log("became Follower with term=%d; log=%v", term, this.log)

	// IMPLEMENT becomeFollower; do you need to start a goroutine here, maybe?
	//-------------------------------------------------------------------------------------------/
	// TODO
	//-------------------------------------------------------------------------------------------/
}
