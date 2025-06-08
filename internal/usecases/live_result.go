package usecases

import (
	"context"
	"time"
)

type LiveResultUsecases interface {
	// Broadcast individual updates
	BroadcastVoteUpdate(ctx context.Context, voteID string) error
	BroadcastElectionUpdate(ctx context.Context, electionPairID string) error
	BroadcastRegionUpdate(ctx context.Context, region string) error
	BroadcastStatisticsUpdate(ctx context.Context) error

	// Broadcast multiple updates at once
	BroadcastAllUpdates(ctx context.Context, electionPairID, region string) error

	// Management functions
	GetConnectedClients(ctx context.Context) int
	StartPeriodicBroadcast(ctx context.Context, interval time.Duration)
}
