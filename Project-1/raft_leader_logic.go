package raft

import (
	"math/rand"
	"time"
)

// startLeader switches this into a leader state and begins process of heartbeats.
func (this *RaftNode) startLeader() {
	this.state = "Leader"

	for _, peerId := range this.peersIds {
		this.nextIndex[peerId] = len(this.log)
		this.matchIndex[peerId] = -1
	}
	this.write_log("became Leader; term=%d, nextIndex=%v, matchIndex=%v; log=%v", this.currentTerm, this.nextIndex, this.matchIndex, this.log)

	go func() {
		ticker := time.NewTicker(1000 * time.Millisecond)
		defer ticker.Stop()

		// Send periodic heartbeats, as long as still leader.
		for {
			this.broadcastHeartbeats()
			<-ticker.C

			this.mu.Lock()
			if this.state != "Leader" {
				this.mu.Unlock()
				return
			}
			this.mu.Unlock()
		}
	}()
}

// broadcastHeartbeats sends a round of heartbeats to all peers, collects their replies and adjusts this's state.
// Since the Raft Paper uses the AppendEntries function with an EMPTY log as its heartbeat,
// we're doing the same here.
func (this *RaftNode) broadcastHeartbeats() {
	this.mu.Lock()

	if this.state != "Leader" {
		this.mu.Unlock()
		return
	}
	termWhenHeartbeatSent := this.currentTerm

	this.mu.Unlock()

	// Send a Heartbeat PER PEER.
	for _, peerId := range this.peersIds { // Peers are other nodes.

		go func(peerId int) {
			this.mu.Lock()

			currentPeer_nextIndex := this.nextIndex[peerId]
			prevLogIndex := currentPeer_nextIndex - 1
			prevLogTerm := -1
			if prevLogIndex >= 0 {
				prevLogTerm = this.log[prevLogIndex].Term
			}
			entries := this.log[currentPeer_nextIndex:] // Which entries on the leader are not there on peer?

			var aeType string
			if len(entries) > 0 {
				aeType = "AppendEntries"
			} else {
				aeType = "Heartbeat"
			}

			args := AppendEntriesArgs{
				Term:         termWhenHeartbeatSent,
				LeaderId:     this.id,
				PrevLogIndex: prevLogIndex,
				PrevLogTerm:  prevLogTerm,
				Entries:      entries,
				LeaderCommit: this.commitIndex,
				Latency:      rand.Intn(500), // Ignore Latency
			}

			this.mu.Unlock()
			if (aeType == "Heartbeat" && LogHeartbeatMessages) || aeType == "AppendEntries" {
				this.write_log("sending %s to %v: currentPeer_nextIndex=%d, args=%+v", aeType, peerId, currentPeer_nextIndex, args)
			}

			var reply AppendEntriesReply

			// Don't worry about this; this is how the RPC itself is sent.
			// Just presume that the AppendEntries went to this follower id,
			// And you now need to handle the reply.
			if err := this.server.SendRPCCallTo(peerId, "RaftNode.AppendEntries", args, &reply); err == nil {
				this.mu.Lock()
				defer this.mu.Unlock()

				if reply.Term > this.currentTerm {
					this.becomeFollower(reply.Term)
					return
				}

				if this.state == "Leader" && termWhenHeartbeatSent == reply.Term {
					if reply.Success {

						// There's changes you need to make here.
						// this.nextIndex for the received PEER (this.nextIndex[peerId]) needs to be updated.
						// So does this.matchIndex[peerId].
						// IMPLEMENT THE UPDATE LOGIC FOR THIS.
						//-------------------------------------------------------------------------------------------/
						// TODO
						//-------------------------------------------------------------------------------------------/

						if (aeType == "Heartbeat" && LogHeartbeatMessages) || aeType == "AppendEntries" {
							this.write_log("%s reply from NODE %d success: nextIndex := %v, matchIndex := %v", aeType, peerId, this.nextIndex, this.matchIndex)
						}
						oldCommitIndex := this.commitIndex

						// AppendEntries success on majority, now commit on leader (IF NOT HEARTBEAT)

						// You must update commitIndex in a specific way somewhere in this loop;
						// Figure out how and where; HINT: look for a majority of matchCounts.

						//-------------------------------------------------------------------------------------------/
						for i := this.commitIndex + 1; i < len(this.log); i++ {
							if this.log[i].Term == this.currentTerm {
								matchCount := 1 // Leader itself

								for _, peerId := range this.peersIds {
									if { // TODO  // When should you update matchCount?
										matchCount++
									}
								}

								if { // TODO  // When should you update commitIndex to i?
									this.commitIndex = i
								}
							}
						}
						//-------------------------------------------------------------------------------------------/

						// This actually applies commits. Your logic above for deciding whether or not
						// To commit needs to work succesfully in order for this to occur.
						if this.commitIndex != oldCommitIndex {
							this.write_log("leader sets commitIndex := %d", this.commitIndex)
							this.notifyToApplyCommit <- 1
						}

					} else {

						// There's changes you need to make here.
						// this.nextIndex for the received PEER (this.nextIndex[peerId]) needs to be updated.

						//-------------------------------------------------------------------------------------------/
						// TODO
						//-------------------------------------------------------------------------------------------/

						if (aeType == "Heartbeat" && LogHeartbeatMessages) || aeType == "AppendEntries" {
							this.write_log("%s reply from NODE %d was failure; Hence, decrementing its nextIndex", aeType, peerId)
						}
					}
				}
			}
		}(peerId)
	}
}
