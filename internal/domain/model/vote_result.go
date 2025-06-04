package model

import (
	"github.com/google/uuid"
	"time"
)

type VoteStatus string

const (
	VoteStatusPending   VoteStatus = "pending"
	VoteStatusConfirmed VoteStatus = "confirmed"
	VoteStatusRejected  VoteStatus = "rejected"
	VoteStatusError     VoteStatus = "error"
	VoteStatusQueued    VoteStatus = "queued"
	VoteStatusRetrying  VoteStatus = "retrying"
)

type VoteResult struct {
	ID              uuid.UUID  `db:"id" json:"id"`
	VoteID          uuid.UUID  `db:"vote_id" json:"vote_id"`
	VoterID         string     `db:"voter_id" json:"voter_id"`
	ElectionPairID  uuid.UUID  `db:"election_pair_id" json:"election_pair_id"`
	Region          string     `db:"region" json:"region"`
	Status          string     `db:"status" json:"status"`
	TransactionHash string     `db:"transaction_hash" json:"transaction_hash"`
	ErrorMessage    string     `db:"error_message" json:"error_message,omitempty"`
	VotedAt         time.Time  `db:"voted_at" json:"voted_at"`
	ProcessedAt     *time.Time `db:"processed_at" json:"processed_at,omitempty"`
	CreatedAt       time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at" json:"updated_at"`
}

type ElectionResult struct {
	ElectionPairID string    `db:"election_pair_id" json:"election_pair_id"`
	Region         string    `db:"region" json:"region"`
	TotalVotes     uint64    `db:"total_votes" json:"total_votes"`
	ConfirmedVotes uint64    `db:"confirmed_votes" json:"confirmed_votes"`
	PendingVotes   uint64    `db:"pending_votes" json:"pending_votes"`
	ErrorVotes     uint64    `db:"error_votes" json:"error_votes"`
	LastUpdated    time.Time `db:"last_updated" json:"last_updated"`
}

type RegionResult struct {
	Region         string    `db:"region" json:"region"`
	TotalVotes     uint64    `db:"total_votes" json:"total_votes"`
	ConfirmedVotes uint64    `db:"confirmed_votes" json:"confirmed_votes"`
	PendingVotes   uint64    `db:"pending_votes" json:"pending_votes"`
	ErrorVotes     uint64    `db:"error_votes" json:"error_votes"`
	LastUpdated    time.Time `db:"last_updated" json:"last_updated"`
}

type VoteStatistics struct {
	TotalVotes     uint64    `json:"total_votes"`
	ConfirmedVotes uint64    `json:"confirmed_votes"`
	PendingVotes   uint64    `json:"pending_votes"`
	ErrorVotes     uint64    `json:"error_votes"`
	SuccessRate    float64   `json:"success_rate"`
	LastUpdated    time.Time `json:"last_updated"`
}

func (vs *VoteStatistics) CalculateSuccessRate() {
	if vs.TotalVotes > 0 {
		vs.SuccessRate = float64(vs.ConfirmedVotes) / float64(vs.TotalVotes) * 100
	} else {
		vs.SuccessRate = 0
	}

}
