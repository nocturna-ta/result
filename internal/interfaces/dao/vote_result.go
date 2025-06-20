package dao

import (
	"context"
	"fmt"
	"github.com/nocturna-ta/golib/database/sql"
	"github.com/nocturna-ta/golib/log"
	"github.com/nocturna-ta/golib/tracing"
	"github.com/nocturna-ta/golib/txmanager/utils"
	"github.com/nocturna-ta/result/internal/domain/model"
	"github.com/nocturna-ta/result/internal/domain/repository"
)

type VoteResultRepository struct {
	db *sql.Store
}

type OptsVoteResultRepository struct {
	DB *sql.Store
}

func NewVoteResultRepository(opts *OptsVoteResultRepository) repository.VoteResultRepository {
	return &VoteResultRepository{
		db: opts.DB,
	}
}

const (
	insertVoteResultQuery = `
		INSERT INTO vote_results (
			id, voter_id, election_pair_id, region, status, 
			transaction_hash, error_message, voted_at, 
			processed_at, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	selectVoteResultQuery = `SELECT %s FROM vote_results %s WHERE TRUE %s `
	updateVoteResultQuery = `ALTER TABLE vote_results UPDATE %s WHERE TRUE %s `
)

func (v *VoteResultRepository) InsertVoteResult(ctx context.Context, result *model.VoteResult) error {
	span, ctx := tracing.StartSpanFromContext(ctx, "VoteResultRepository.InsertVoteResult")
	defer span.End()

	sqlTrx := utils.GetSqlTx(ctx)
	var (
		err error
	)

	if sqlTrx != nil {
		_, err = sqlTrx.ExecContext(ctx, insertVoteResultQuery, result.ID, result.VoterID, result.ElectionPairID,
			result.Region, result.Status, result.TransactionHash, result.ErrorMessage,
			result.VotedAt, result.ProcessedAt, result.CreatedAt, result.UpdatedAt)
	} else {
		_, err = v.db.GetMaster().ExecContext(ctx, insertVoteResultQuery, result.ID, result.VoterID, result.ElectionPairID,
			result.Region, result.Status, result.TransactionHash, result.ErrorMessage,
			result.VotedAt, result.ProcessedAt, result.CreatedAt, result.UpdatedAt)
	}

	if err != nil {
		log.WithFields(log.Fields{
			"error":  err,
			"result": result,
		}).ErrorWithCtx(ctx, "[VoteResultRepository.InsertVoteResult] failed to insert vote result")
		return err
	}

	return nil
}

func (v *VoteResultRepository) UpdateVoteResult(ctx context.Context, result *model.VoteResult) error {
	span, ctx := tracing.StartSpanFromContext(ctx, "VoteResultRepository.UpdateVoteResult")
	defer span.End()

	var (
		args []any
		err  error
	)

	sqlTrx := utils.GetSqlTx(ctx)

	setQuery := `status = ?, transaction_hash = ?, error_message = ?, processed_at = ?, updated_at = ?`
	whereQuery := ` AND id = ?`
	args = append(args, result.Status, result.TransactionHash, result.ErrorMessage, result.ProcessedAt, result.UpdatedAt, result.ID)

	query := fmt.Sprintf(updateVoteResultQuery, setQuery, whereQuery)

	if sqlTrx != nil {
		_, err = sqlTrx.ExecContext(ctx, query, args...)
	} else {
		_, err = v.db.GetMaster().ExecContext(ctx, query, args...)
	}

	if err != nil {
		log.WithFields(log.Fields{
			"error":  err,
			"result": result,
		}).ErrorWithCtx(ctx, "[VoteResultRepository.UpdateVoteResult] failed to update vote result")
		return err
	}

	return nil
}

func (v *VoteResultRepository) GetVoteResultByID(ctx context.Context, id string) (*model.VoteResult, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "VoteResultRepository.GetVoteResultByID")
	defer span.End()

	var (
		result model.VoteResult
		err    error
		args   []any
	)

	sqlTrx := utils.GetSqlTx(ctx)

	selectQuery := `id, voter_id, election_pair_id, region, status, transaction_hash, error_message, voted_at, processed_at, created_at, updated_at`
	whereQuery := `AND id = ?`
	args = append(args, id)
	query := fmt.Sprintf(selectVoteResultQuery, selectQuery, "", whereQuery)

	if sqlTrx != nil {
		err = sqlTrx.GetContext(ctx, &result, query, args...)
	} else {
		err = v.db.GetMaster().GetContext(ctx, &result, query, args...)
	}

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"id":    id,
		}).ErrorWithCtx(ctx, "[VoteResultRepository.GetVoteResultByID] failed to get vote result by ID")
		return nil, err
	}

	return &result, nil
}
