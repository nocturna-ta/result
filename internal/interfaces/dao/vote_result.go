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
	"time"
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

func (v *VoteResultRepository) GetVoteResultsByElectionPair(ctx context.Context, electionPairID string, limit, offset int) ([]*model.VoteResult, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "VoteResultRepository.GetVoteResultsByElectionPair")
	defer span.End()

	var (
		results []*model.VoteResult
		err     error
		args    []any
	)

	sqlTrx := utils.GetSqlTx(ctx)

	selectQuery := `id, voter_id, election_pair_id, region, status, transaction_hash, error_message, voted_at, processed_at, created_at, updated_at`
	whereQuery := `AND election_pair_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`
	args = append(args, electionPairID, limit, offset)
	query := fmt.Sprintf(selectVoteResultQuery, selectQuery, "", whereQuery)

	if sqlTrx != nil {
		err = sqlTrx.SelectContext(ctx, &results, query, args...)
	} else {
		err = v.db.GetMaster().SelectContext(ctx, &results, query, args...)
	}
	if err != nil {
		log.WithFields(log.Fields{
			"error":         err,
			"election_pair": electionPairID,
			"limit":         limit,
			"offset":        offset,
		}).ErrorWithCtx(ctx, "[VoteResultRepository.GetVoteResultsByElectionPair] failed to get vote results by election pair")
		return nil, err
	}

	return results, nil
}

func (v *VoteResultRepository) GetVoteResultsByRegion(ctx context.Context, region string, limit, offset int) ([]*model.VoteResult, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "VoteResultRepository.GetVoteResultsByRegion")
	defer span.End()

	var (
		results []*model.VoteResult
		err     error
		args    []any
	)

	sqlTrx := utils.GetSqlTx(ctx)

	selectQuery := `id, voter_id, election_pair_id, region, status, transaction_hash, error_message, voted_at, processed_at, created_at, updated_at`
	whereQuery := `AND region = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`
	args = append(args, region, limit, offset)
	query := fmt.Sprintf(selectVoteResultQuery, selectQuery, "", whereQuery)

	if sqlTrx != nil {
		err = sqlTrx.SelectContext(ctx, &results, query, args...)
	} else {
		err = v.db.GetMaster().SelectContext(ctx, &results, query, args...)
	}

	if err != nil {
		log.WithFields(log.Fields{
			"error":  err,
			"region": region,
			"limit":  limit,
			"offset": offset,
		}).ErrorWithCtx(ctx, "[VoteResultRepository.GetVoteResultsByRegion] failed to get vote results by region")
		return nil, err
	}

	return results, nil
}

func (v *VoteResultRepository) GetVoteResultsByStatus(ctx context.Context, status string, limit, offset int) ([]*model.VoteResult, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "VoteResultRepository.GetVoteResultsByStatus")
	defer span.End()

	var (
		results []*model.VoteResult
		err     error
		args    []any
	)

	sqlTrx := utils.GetSqlTx(ctx)

	selectQuery := `id, voter_id, election_pair_id, region, status, transaction_hash, error_message, voted_at, processed_at, created_at, updated_at`
	whereQuery := `AND status = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`
	args = append(args, status, limit, offset)
	query := fmt.Sprintf(selectVoteResultQuery, selectQuery, "", whereQuery)

	if sqlTrx != nil {
		err = sqlTrx.SelectContext(ctx, &results, query, args...)
	} else {
		err = v.db.GetMaster().SelectContext(ctx, &results, query, args...)
	}
	if err != nil {
		log.WithFields(log.Fields{
			"error":  err,
			"status": status,
			"limit":  limit,
			"offset": offset,
		}).ErrorWithCtx(ctx, "[VoteResultRepository.GetVoteResultsByStatus] failed to get vote results by status")
		return nil, err
	}

	return results, nil
}

func (v *VoteResultRepository) GetElectionResults(ctx context.Context, electionPairID string) (*model.ElectionResult, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "VoteResultRepository.GetElectionResults")
	defer span.End()

	var (
		result model.ElectionResult
		err    error
		args   []any
	)

	sqlTrx := utils.GetSqlTx(ctx)

	selectQuery := `election_pair_id, region, count(*) as total_votes, countIf(status = 'confirmed') as confirmed_votes, countIf(status = 'pending') as pending_votes, countIf(status = 'error') as error_votes, max(updated_at) as last_updated`
	whereQuery := ` AND election_pair_id = ? GROUP BY election_pair_id, region`
	args = append(args, electionPairID)
	query := fmt.Sprintf(selectVoteResultQuery, selectQuery, "", whereQuery)

	if sqlTrx != nil {
		err = sqlTrx.GetContext(ctx, &result, query, args...)
	} else {
		err = v.db.GetMaster().GetContext(ctx, &result, query, args...)
	}
	if err != nil {
		log.WithFields(log.Fields{
			"error":         err,
			"election_pair": electionPairID,
		}).ErrorWithCtx(ctx, "[VoteResultRepository.GetElectionResults] failed to get election results")
		return nil, err
	}

	return &result, nil
}

func (v *VoteResultRepository) GetRegionResults(ctx context.Context, region string) (*model.RegionResult, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "VoteResultRepository.GetRegionResults")
	defer span.End()

	var (
		result model.RegionResult
		err    error
		args   []any
	)

	sqlTrx := utils.GetSqlTx(ctx)

	selectQuery := `region, count(*) as total_votes, countIf(status = 'confirmed') as confirmed_votes,
	countIf(status = 'pending') as pending_votes, countIf(status = 'error') as error_votes, max(updated_at) as last_updated`
	whereQuery := `AND region = ? GROUP BY region`
	args = append(args, region)
	query := fmt.Sprintf(selectVoteResultQuery, selectQuery, "", whereQuery)

	if sqlTrx != nil {
		err = sqlTrx.GetContext(ctx, &result, query, args...)
	} else {
		err = v.db.GetMaster().GetContext(ctx, &result, query, args...)
	}
	if err != nil {
		log.WithFields(log.Fields{
			"error":  err,
			"region": region,
		}).ErrorWithCtx(ctx, "[VoteResultRepository.GetRegionResults] failed to get region results")
		return nil, err
	}

	return &result, nil
}

func (v *VoteResultRepository) GetOverallStatistics(ctx context.Context) (*model.VoteStatistics, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "VoteResultRepository.GetOverallStatistics")
	defer span.End()

	var (
		result model.VoteStatistics
		err    error
	)
	sqlTrx := utils.GetSqlTx(ctx)

	selectQuery := `count(*) as total_votes, countIf(status = 'confirmed') as confirmed_votes,
	countIf(status = 'pending') as pending_votes, countIf(status = 'error') as error_votes, max(updated_at) as last_updated`
	query := fmt.Sprintf(selectVoteResultQuery, selectQuery, "", "")

	if sqlTrx != nil {
		err = sqlTrx.GetContext(ctx, &result, query)
	} else {
		err = v.db.GetMaster().GetContext(ctx, &result, query)
	}

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).ErrorWithCtx(ctx, "[VoteResultRepository.GetOverallStatistics] failed to get overall statistics")
		return nil, err
	}

	return &result, nil
}

func (v *VoteResultRepository) GetVoteResultsByDateRange(ctx context.Context, startDate, endDate time.Time, limit, offset int) ([]*model.VoteResult, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "VoteResultRepository.GetVoteResultsByDateRange")
	defer span.End()

	var (
		results []*model.VoteResult
		err     error
		args    []any
	)

	sqlTrx := utils.GetSqlTx(ctx)

	selectQuery := `id, voter_id, election_pair_id, region, status, transaction_hash, error_message, voted_at, processed_at, created_at, updated_at`
	whereQuery := ` AND created_at >= ? AND created_at <= ? ORDER BY created_at DESC LIMIT ? OFFSET ?`
	args = append(args, startDate, endDate, limit, offset)
	query := fmt.Sprintf(selectVoteResultQuery, selectQuery, "", whereQuery)

	if sqlTrx != nil {
		err = sqlTrx.SelectContext(ctx, &results, query, args...)
	} else {
		err = v.db.GetMaster().SelectContext(ctx, &results, query, args...)
	}

	if err != nil {
		log.WithFields(log.Fields{
			"error":     err,
			"startDate": startDate,
			"endDate":   endDate,
			"limit":     limit,
			"offset":    offset,
		}).ErrorWithCtx(ctx, "[VoteResultRepository.GetVoteResultsByDateRange] failed to get vote results by date range")
		return nil, err
	}

	return results, nil
}

func (v *VoteResultRepository) GetElectionResultsByRegion(ctx context.Context, region string) ([]*model.ElectionResult, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "VoteResultRepository.GetElectionResultsByRegion")
	defer span.End()

	var (
		results []*model.ElectionResult
		err     error
		args    []any
	)

	sqlTrx := utils.GetSqlTx(ctx)

	selectQuery := `election_pair_id, region, count(*) as total_votes, countIf(status = 'confirmed') as confirmed_votes,
            countIf(status = 'pending') as pending_votes, countIf(status = 'error') as error_votes, max(updated_at) as last_updated`
	whereQuery := `AND region = ? GROUP BY election_pair_id, region ORDER BY total_votes DESC`
	args = append(args, region)
	query := fmt.Sprintf(selectVoteResultQuery, selectQuery, "", whereQuery)

	if sqlTrx != nil {
		err = sqlTrx.SelectContext(ctx, &results, query, args...)
	} else {
		err = v.db.GetMaster().SelectContext(ctx, &results, query, args...)
	}

	if err != nil {
		log.WithFields(log.Fields{
			"error":  err,
			"region": region,
		}).ErrorWithCtx(ctx, "[VoteResultRepository.GetElectionResultsByRegion] failed to get election results by region")
		return nil, err
	}

	return results, nil
}

func (v *VoteResultRepository) GetRegionStatistics(ctx context.Context) ([]*model.RegionResult, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "VoteResultRepository.GetRegionStatistics")
	defer span.End()

	var (
		results []*model.RegionResult
		err     error
	)

	sqlTrx := utils.GetSqlTx(ctx)

	selectQuery := `region, count(*) as total_votes, countIf(status = 'confirmed') as confirmed_votes,
            countIf(status = 'pending') as pending_votes, countIf(status = 'error') as error_votes, max(updated_at) as last_updated`
	whereQuery := `GROUP BY region ORDER BY total_votes DESC`
	query := fmt.Sprintf(selectVoteResultQuery, selectQuery, "", whereQuery)

	if sqlTrx != nil {
		err = sqlTrx.SelectContext(ctx, &results, query)
	} else {
		err = v.db.GetMaster().SelectContext(ctx, &results, query)
	}

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).ErrorWithCtx(ctx, "[VoteResultRepository.GetRegionStatistics] failed to get region statistics")
		return nil, err
	}

	return results, nil
}

func (v *VoteResultRepository) CountVotesByStatus(ctx context.Context, status string) (uint64, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "VoteResultRepository.CountVotesByStatus")
	defer span.End()

	var (
		count uint64
		err   error
		args  []any
	)

	sqlTrx := utils.GetSqlTx(ctx)

	selectQuery := `count(*)`
	whereQuery := `AND status = ?`
	args = append(args, status)
	query := fmt.Sprintf(selectVoteResultQuery, selectQuery, "", whereQuery)

	if sqlTrx != nil {
		err = sqlTrx.GetContext(ctx, &count, query, args...)
	} else {
		err = v.db.GetMaster().GetContext(ctx, &count, query, args...)
	}

	if err != nil {
		log.WithFields(log.Fields{
			"error":  err,
			"status": status,
		}).ErrorWithCtx(ctx, "[VoteResultRepository.CountVotesByStatus] failed to count votes by status")
		return 0, err
	}

	return count, nil
}

func (v *VoteResultRepository) CountVotesByElectionPair(ctx context.Context, electionPairID string) (uint64, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "VoteResultRepository.CountVotesByElectionPair")
	defer span.End()

	var (
		count uint64
		err   error
		args  []any
	)

	sqlTrx := utils.GetSqlTx(ctx)

	selectQuery := `count(*)`
	whereQuery := `AND election_pair_id = ?`
	args = append(args, electionPairID)
	query := fmt.Sprintf(selectVoteResultQuery, selectQuery, "", whereQuery)

	if sqlTrx != nil {
		err = sqlTrx.GetContext(ctx, &count, query, args...)
	} else {
		err = v.db.GetMaster().GetContext(ctx, &count, query, args...)
	}

	if err != nil {
		log.WithFields(log.Fields{
			"error":         err,
			"election_pair": electionPairID,
		}).ErrorWithCtx(ctx, "[VoteResultRepository.CountVotesByElectionPair] failed to count votes by election pair")
		return 0, err
	}
	return count, nil
}

func (v *VoteResultRepository) CountVotesByRegion(ctx context.Context, region string) (uint64, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "VoteResultRepository.CountVotesByRegion")
	defer span.End()

	var (
		count uint64
		err   error
		args  []any
	)

	sqlTrx := utils.GetSqlTx(ctx)

	selectQuery := `count(*)`
	whereQuery := `AND region = ?`
	args = append(args, region)
	query := fmt.Sprintf(selectVoteResultQuery, selectQuery, "", whereQuery)

	if sqlTrx != nil {
		err = sqlTrx.GetContext(ctx, &count, query, args...)
	} else {
		err = v.db.GetMaster().GetContext(ctx, &count, query, args...)
	}

	if err != nil {
		log.WithFields(log.Fields{
			"error":  err,
			"region": region,
		}).ErrorWithCtx(ctx, "[VoteResultRepository.CountVotesByRegion] failed to count votes by region")
		return 0, err
	}
	return count, nil
}

func (v *VoteResultRepository) GetVoteResultsByHour(ctx context.Context, date time.Time) ([]*model.VoteResult, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "VoteResultRepository.GetVoteResultsByHour")
	defer span.End()

	var (
		results []*model.VoteResult
		err     error
		args    []any
	)

	sqlTrx := utils.GetSqlTx(ctx)

	startHour := time.Date(date.Year(), date.Month(), date.Day(), date.Hour(), 0, 0, 0, date.Location())
	endHour := startHour.Add(time.Hour)

	selectQuery := `id, voter_id, election_pair_id, region, status, transaction_hash, error_message, voted_at, processed_at, created_at, updated_at`
	whereQuery := `AND created_at >= ? AND created_at < ? ORDER BY created_at DESC`
	args = append(args, startHour, endHour)
	query := fmt.Sprintf(selectVoteResultQuery, selectQuery, "", whereQuery)

	if sqlTrx != nil {
		err = sqlTrx.SelectContext(ctx, &results, query, args...)
	} else {
		err = v.db.GetMaster().SelectContext(ctx, &results, query, args...)
	}
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"date":  date,
		}).ErrorWithCtx(ctx, "[VoteResultRepository.GetVoteResultsByHour] failed to get vote results by hour")
		return nil, err
	}

	return results, nil
}

func (v *VoteResultRepository) GetVoteResultsByDay(ctx context.Context, date time.Time) ([]*model.VoteResult, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "VoteResultRepository.GetVoteResultsByDay")
	defer span.End()

	var (
		results []*model.VoteResult
		err     error
		args    []any
	)

	sqlTrx := utils.GetSqlTx(ctx)

	startDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endDay := startDay.Add(24 * time.Hour)

	selectQuery := `id, voter_id, election_pair_id, region, status, transaction_hash, error_message, voted_at, processed_at, created_at, updated_at`
	whereQuery := `AND created_at >= ? AND created_at < ? ORDER BY created_at DESC`
	args = append(args, startDay, endDay)
	query := fmt.Sprintf(selectVoteResultQuery, selectQuery, "", whereQuery)

	if sqlTrx != nil {
		err = sqlTrx.SelectContext(ctx, &results, query, args...)
	} else {
		err = v.db.GetMaster().SelectContext(ctx, &results, query, args...)
	}
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"date":  date,
		}).ErrorWithCtx(ctx, "[VoteResultRepository.GetVoteResultsByDay] failed to get vote results by day")
		return nil, err
	}

	return results, nil
}

func (v *VoteResultRepository) GetDailyStatistics(ctx context.Context, startDate, endDate time.Time) ([]*model.VoteStatistics, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "VoteResultRepository.GetDailyStatistics")
	defer span.End()

	var (
		results []*model.VoteStatistics
		err     error
		args    []any
	)

	sqlTrx := utils.GetSqlTx(ctx)

	selectQuery := `toDate(created_at) as date, count(*) as total_votes, countIf(status = 'confirmed') as confirmed_votes,
                 countIf(status = 'pending') as pending_votes, countIf(status = 'error') as error_votes, max(updated_at) as last_updated`
	whereQuery := `AND created_at >= ? AND created_at <= ? GROUP BY toDate(created_at) ORDER BY date DESC`
	args = append(args, startDate, endDate)
	query := fmt.Sprintf(selectVoteResultQuery, selectQuery, "", whereQuery)

	if sqlTrx != nil {
		err = sqlTrx.SelectContext(ctx, &results, query, args...)
	} else {
		err = v.db.GetMaster().SelectContext(ctx, &results, query, args...)
	}

	if err != nil {
		log.WithFields(log.Fields{
			"error":     err,
			"startDate": startDate,
			"endDate":   endDate,
		}).ErrorWithCtx(ctx, "[VoteResultRepository.GetDailyStatistics] failed to get daily statistics")
		return nil, err
	}

	for _, stat := range results {
		stat.CalculateSuccessRate()
	}

	return results, nil
}
