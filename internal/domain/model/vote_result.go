package model

import (
	"github.com/nocturna-ta/common-model/models/event"
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
	ID              string     `db:"id" `
	VoterID         string     `db:"voter_id"`
	ElectionPairID  string     `db:"election_pair_id"`
	Region          string     `db:"region" `
	Status          string     `db:"status"`
	TransactionHash string     `db:"transaction_hash"`
	ErrorMessage    string     `db:"error_message"`
	VotedAt         time.Time  `db:"voted_at" `
	ProcessedAt     *time.Time `db:"processed_at" `
	CreatedAt       time.Time  `db:"created_at" `
	UpdatedAt       time.Time  `db:"updated_at"`
}

func FromVoteProcessedMessage(msg *event.VoteProcessedMessage) *VoteResult {
	return &VoteResult{
		ID:              msg.VoteID,
		VoterID:         msg.VoterID,
		Status:          msg.Status,
		TransactionHash: msg.TransactionHash,
		ErrorMessage:    msg.ErrorMessage,
		ProcessedAt:     &msg.ProcessedAt,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

func FromVoteSubmitMessage(msg *event.VoteSubmitMessage) *VoteResult {
	return &VoteResult{
		ID:             msg.VoteID,
		VoterID:        msg.VoterID,
		ElectionPairID: msg.ElectionPairID,
		Region:         msg.Region,
		Status:         "pending",
		VotedAt:        msg.SubmittedAt,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}
