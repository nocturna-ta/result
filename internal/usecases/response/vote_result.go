package response

import "time"

type VoteResultResponse struct {
	ID              string     `json:"id"`
	VoterID         string     `json:"voter_id"`
	ElectionPairID  string     `json:"election_pair_id"`
	Region          string     `json:"region"`
	Status          string     `json:"status"`
	TransactionHash string     `json:"transaction_hash"`
	ErrorMessage    string     `json:"error_message,omitempty"`
	VotedAt         time.Time  `json:"voted_at"`
	ProcessedAt     *time.Time `json:"processed_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type ElectionVoteResultResponse struct {
	ElectionPairID string    `json:"election_pair_id"`
	Region         string    `json:"region"`
	TotalVotes     uint64    `json:"total_votes"`
	ConfirmedVotes uint64    `json:"confirmed_votes"`
	PendingVotes   uint64    `json:"pending_votes"`
	ErrorVotes     uint64    `json:"error_votes"`
	LastUpdated    time.Time `json:"last_updated"`
}

type RegionVoteResultResponse struct {
	Region         string    `json:"region"`
	TotalVotes     uint64    `json:"total_votes"`
	ConfirmedVotes uint64    `json:"confirmed_votes"`
	PendingVotes   uint64    `json:"pending_votes"`
	ErrorVotes     uint64    `json:"error_votes"`
	LastUpdated    time.Time `json:"last_updated"`
}

type VoteStatisticsResponse struct {
	TotalVotes     uint64    `json:"total_votes"`
	ConfirmedVotes uint64    `json:"confirmed_votes"`
	PendingVotes   uint64    `json:"pending_votes"`
	ErrorVotes     uint64    `json:"error_votes"`
	SuccessRate    float64   `json:"success_rate"`
	LastUpdated    time.Time `json:"last_updated"`
}
