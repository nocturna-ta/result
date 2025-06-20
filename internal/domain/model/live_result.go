package model

import (
	"database/sql"
	"time"
)

type LiveElectionResult struct {
	ElectionPairID    string    `db:"election_pair_id"`
	Region            string    `db:"region"`
	ConfirmedVotes    uint64    `db:"confirmed_votes"`
	TotalVotes        uint64    `db:"total_votes"`
	SuccessPercentage float64   `db:"success_percentage"`
	LastUpdated       time.Time `db:"last_updated"`
}

type LiveCityResult struct {
	CityName                  string          `db:"city_name"`
	ElectionPairID            string          `db:"election_pair_id"`
	TotalUniqueVoters         uint64          `db:"total_unique_voters"`
	ConfirmedVotes            uint64          `db:"confirmed_votes"`
	VoteSuccessRate           sql.NullFloat64 `db:"vote_success_rate"`
	CandidatePercentageInCity sql.NullFloat64 `db:"candidate_percentage_in_city"`
	LastUpdated               time.Time       `db:"last_updated"`
}

type ElectionSummary struct {
	ElectionPairID      string    `db:"election_pair_id"`
	TotalUniqueVoters   uint64    `db:"total_unique_voters"`
	TotalRegions        uint64    `db:"total_regions"`
	TotalConfirmedVotes uint64    `db:"total_confirmed_votes"`
	OverallSuccessRate  float64   `db:"overall_success_rate"`
	LastUpdated         time.Time `db:"last_updated"`
}

type CityRanking struct {
	ElectionPairID    string    `db:"election_pair_id"`
	CityName          string    `db:"city_name"`
	ConfirmedVotes    uint64    `db:"confirmed_votes"`
	UniqueVoters      uint64    `db:"unique_voters"`
	ParticipationRate float64   `db:"participation_rate"`
	CityRank          uint32    `db:"city_rank"`
	LastUpdated       time.Time `db:"last_updated"`
}

func (l *LiveCityResult) GetVoteSuccessRate() float64 {
	if l.VoteSuccessRate.Valid {
		return l.VoteSuccessRate.Float64
	}
	return 0.0
}

func (l *LiveCityResult) GetCandidatePercentageInCity() float64 {
	if l.CandidatePercentageInCity.Valid {
		return l.CandidatePercentageInCity.Float64
	}
	return 0.0
}
