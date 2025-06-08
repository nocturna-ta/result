package request

type VoteResultEntry struct {
	VoteID          string `json:"vote_id"`
	VoterID         string `json:"voter_id"`
	ElectionPairID  string `json:"election_pair_id"`
	Region          string `json:"region"`
	TransactionHash string `json:"transaction_hash"`
}
