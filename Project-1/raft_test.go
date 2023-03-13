package raft

import (
	"testing"
)

func Test1(t *testing.T) { // Simple Leader Election

	cluster := NewCluster(t, 5)
	defer cluster.Shutdown()

	sleepMs(3000) // Wait for a leader to be elected

	firstLeaderId := cluster.getClusterLeader()
	cluster.DisconnectPeer(firstLeaderId)

	secondLeaderId := cluster.getClusterLeader()
	cluster.DisconnectPeer(secondLeaderId)

	thirdLeaderId := cluster.getClusterLeader()
	cluster.DisconnectPeer(thirdLeaderId)

	sleepMs(3000)

	// Fails, no leader present
	cluster.getClusterLeader()
	sleepMs(3000)

}

func Test2(t *testing.T) {
	/* Replication failure scenario: Leader drops after committing, comes back later*/

	cluster := NewCluster(t, 5)
	defer cluster.Shutdown()

	// ReceiveClientCommand a couple of values to a fully connected nodes.
	origLeaderId := cluster.getClusterLeader()
	cluster.SubmitClientCommand(origLeaderId, "Set X = 5")
	cluster.SubmitClientCommand(origLeaderId, "Set X = 1000")

	sleepMs(3000)

	// Leader disconnected...
	cluster.DisconnectPeer(origLeaderId)

	// ReceiveClientCommand 7 to original leader, even though it's disconnected. Should not reflect.
	cluster.SubmitClientCommand(origLeaderId, "Set X = X-5")

	newLeaderId := cluster.getClusterLeader()

	// ReceiveClientCommand 8.. to new leader.
	cluster.SubmitClientCommand(newLeaderId, "Set X = X+10")
	cluster.SubmitClientCommand(newLeaderId, "Set X = X+1")
	cluster.SubmitClientCommand(newLeaderId, "Set Y = 5")
	cluster.SubmitClientCommand(newLeaderId, "Set Y = X+Y")
	cluster.SubmitClientCommand(newLeaderId, "Set Y = Y+3")
	cluster.SubmitClientCommand(newLeaderId, "Set Z = -1")
	sleepMs(3000)

	// ReceiveClientCommand 9 and check it's fully committed.
	cluster.SubmitClientCommand(newLeaderId, "Set Z = 3")
	sleepMs(3000)

	cluster.ReconnectPeer(origLeaderId)
	sleepMs(15000)
}
