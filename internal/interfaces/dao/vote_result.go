package dao

import (
	"context"
	"github.com/nocturna-ta/golib/database/nosql/clickhouse"
	"github.com/nocturna-ta/golib/log"
	"github.com/nocturna-ta/golib/tracing"
	"github.com/nocturna-ta/result/internal/domain/model"
	"github.com/nocturna-ta/result/internal/domain/repository"
	"time"
)

type VoteResultRepository struct {
	db clickhouse.Client
}

type OptsVoteResultRepository struct {
	DB clickhouse.Client
}

func NewVoteResultRepository(opts *OptsVoteResultRepository) repository.VoteResultRepository {
	return &VoteResultRepository{
		db: opts.DB,
	}
}

const (
	insertVoteResultQuery = `
		INSERT INTO vote_results (
			id, vote_id, voter_id, election_pair_id, region, status, 
			transaction_hash, error_message, voted_at, 
			processed_at, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
)

func (v *VoteResultRepository) InsertVoteResult(ctx context.Context, result *model.VoteResult) error {
	span, ctx := tracing.StartSpanFromContext(ctx, "VoteResultRepository.InsertVoteResult")
	defer span.End()

	err := v.db.Exec(ctx, insertVoteResultQuery, result.ID.String(), result.VoteID.String(), result.VoterID, result.ElectionPairID.String(),
		result.Region, result.Status, result.TransactionHash, result.ErrorMessage,
		result.VotedAt, result.ProcessedAt, result.CreatedAt, result.UpdatedAt)

	if err != nil {
		log.WithFields(log.Fields{
			"error":  err,
			"result": result,
		}).ErrorWithCtx(ctx, "[VoteResultRepository.InsertVoteResult] failed to insert vote result")
		return err
	}

	return nil
}

func (v *VoteResultRepository) InsertVoteResultBatch(ctx context.Context, results []*model.VoteResult) error {
	//TODO implement me
	panic("implement me")
}

func (v *VoteResultRepository) UpdateVoteResult(ctx context.Context, result *model.VoteResult) error {
	//TODO implement me
	panic("implement me")
}

func (v *VoteResultRepository) GetVoteResultByID(ctx context.Context, id string) (*model.VoteResult, error) {
	//TODO implement me
	panic("implement me")
}

func (v *VoteResultRepository) GetVoteResultsByElectionPair(ctx context.Context, electionPairID string, limit, offset int) ([]*model.VoteResult, error) {
	//TODO implement me
	panic("implement me")
}

func (v *VoteResultRepository) GetVoteResultsByRegion(ctx context.Context, region string, limit, offset int) ([]*model.VoteResult, error) {
	//TODO implement me
	panic("implement me")
}

func (v *VoteResultRepository) GetVoteResultsByStatus(ctx context.Context, status string, limit, offset int) ([]*model.VoteResult, error) {
	//TODO implement me
	panic("implement me")
}

func (v *VoteResultRepository) GetVoteResultsByVoter(ctx context.Context, voterID string, limit, offset int) ([]*model.VoteResult, error) {
	//TODO implement me
	panic("implement me")
}

func (v *VoteResultRepository) GetElectionResults(ctx context.Context, electionPairID string) (*model.ElectionResult, error) {
	//TODO implement me
	panic("implement me")
}

func (v *VoteResultRepository) GetRegionResults(ctx context.Context, region string) (*model.RegionResult, error) {
	//TODO implement me
	panic("implement me")
}

func (v *VoteResultRepository) GetOverallStatistics(ctx context.Context) (*model.VoteStatistics, error) {
	//TODO implement me
	panic("implement me")
}

func (v *VoteResultRepository) GetVoteResultsByDateRange(ctx context.Context, startDate, endDate time.Time, limit, offset int) ([]*model.VoteResult, error) {
	//TODO implement me
	panic("implement me")
}

func (v *VoteResultRepository) GetElectionResultsByRegion(ctx context.Context, region string) ([]*model.ElectionResult, error) {
	//TODO implement me
	panic("implement me")
}

func (v *VoteResultRepository) GetRegionStatistics(ctx context.Context) ([]*model.RegionResult, error) {
	//TODO implement me
	panic("implement me")
}

func (v *VoteResultRepository) CountVotesByStatus(ctx context.Context, status string) (uint64, error) {
	//TODO implement me
	panic("implement me")
}

func (v *VoteResultRepository) CountVotesByElectionPair(ctx context.Context, electionPairID string) (uint64, error) {
	//TODO implement me
	panic("implement me")
}

func (v *VoteResultRepository) CountVotesByRegion(ctx context.Context, region string) (uint64, error) {
	//TODO implement me
	panic("implement me")
}

func (v *VoteResultRepository) GetVoteResultsByHour(ctx context.Context, date time.Time) ([]*model.VoteResult, error) {
	//TODO implement me
	panic("implement me")
}

func (v *VoteResultRepository) GetVoteResultsByDay(ctx context.Context, date time.Time) ([]*model.VoteResult, error) {
	//TODO implement me
	panic("implement me")
}

func (v *VoteResultRepository) GetDailyStatistics(ctx context.Context, startDate, endDate time.Time) ([]*model.VoteStatistics, error) {
	//TODO implement me
	panic("implement me")
}
