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

type LiveResultRepository struct {
	db *sql.Store
}

type OptsLiveResultRepository struct {
	DB *sql.Store
}

func NewLiveResultRepository(opts *OptsLiveResultRepository) repository.LiveResultRepository {
	return &LiveResultRepository{
		db: opts.DB,
	}
}

const (
	selectLiveElectionResults = `SELECT %s FROM live_election_results_mv %s WHERE %s`
	selectLiveCityResults     = `SELECT %s FROM live_city_results_mv %s WHERE %s`
	selectElectionSummary     = `SELECT %s FROM election_summary_mv %s WHERE %s`
	selectCityRankings        = `SELECT %s FROM city_rankings_mv %s WHERE %s`
)

func (l *LiveResultRepository) GetLiveElectionResults(ctx context.Context, electionPairID string) ([]*model.LiveElectionResult, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "LiveResultRepository.GetLiveElectionResults")
	defer span.End()

	var (
		results []*model.LiveElectionResult
		err     error
		args    []any
	)

	sqlTrx := utils.GetSqlTx(ctx)

	selectQuery := "election_pair_id, region, confirmed_votes, total_votes, success_percentage, last_updated"
	whereQuery := "election_pair_id = ? ORDER BY confirmed_votes DESC"

	args = append(args, electionPairID)

	query := fmt.Sprintf(selectLiveElectionResults, selectQuery, "", whereQuery)

	if sqlTrx != nil {
		err = sqlTrx.SelectContext(ctx, &results, query, args...)
	} else {
		err = l.db.GetMaster().SelectContext(ctx, &results, query, args...)
	}

	if err != nil {
		log.WithFields(log.Fields{
			"error":            err,
			"election_pair_id": electionPairID,
		}).ErrorWithCtx(ctx, "[LiveResultRepository.GetLiveElectionResults] Failed to get live election results")
		return nil, err
	}

	return results, nil
}

func (l *LiveResultRepository) GetLiveCityResults(ctx context.Context, cityName string) ([]*model.LiveCityResult, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "LiveResultRepository.GetLiveCityResults")
	defer span.End()

	var (
		results []*model.LiveCityResult
		err     error
		args    []any
	)
	sqlTrx := utils.GetSqlTx(ctx)

	selectQuery := "city_name, election_pair_id, confirmed_votes, vote_success_rate, candidate_percentage_in_city, last_updated"
	whereQuery := "city_name = ? ORDER BY confirmed_votes DESC"
	args = append(args, cityName)
	query := fmt.Sprintf(selectLiveCityResults, selectQuery, "", whereQuery)

	if sqlTrx != nil {
		err = sqlTrx.SelectContext(ctx, &results, query, args...)
	} else {
		err = l.db.GetMaster().SelectContext(ctx, &results, query, args...)
	}

	if err != nil {
		log.WithFields(log.Fields{
			"error":     err,
			"city_name": cityName,
		}).ErrorWithCtx(ctx, "[LiveResultRepository.GetLiveCityResults] Failed to get live city results")
		return nil, err
	}

	return results, nil
}

func (l *LiveResultRepository) GetElectionSummary(ctx context.Context, electionPairID string) (*model.ElectionSummary, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "LiveResultRepository.GetElectionSummary")
	defer span.End()

	var (
		summary *model.ElectionSummary
		err     error
		args    []any
	)

	sqlTrx := utils.GetSqlTx(ctx)

	selectQuery := "election_pair_id, total_unique_voters, total_regions, total_confirmed_votes, overallsuccess_rate, last_updated"
	whereQuery := "election_pair_id = ?"

	args = append(args, electionPairID)
	query := fmt.Sprintf(selectElectionSummary, selectQuery, "", whereQuery)

	if sqlTrx != nil {
		err = sqlTrx.GetContext(ctx, &summary, query, args...)
	} else {
		err = l.db.GetMaster().GetContext(ctx, &summary, query, args...)
	}

	if err != nil {
		log.WithFields(log.Fields{
			"error":            err,
			"election_pair_id": electionPairID,
		}).ErrorWithCtx(ctx, "[LiveResultRepository.GetElectionSummary] Failed to get election summary")
		return nil, err
	}

	return summary, nil
}

func (l *LiveResultRepository) GetCityRankings(ctx context.Context, electionPairID string, limit int) ([]*model.CityRanking, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "LiveResultRepository.GetCityRankings")
	defer span.End()

	if limit <= 0 {
		limit = 20
	}

	var (
		rankings []*model.CityRanking
		err      error
		args     []any
	)

	sqlTrx := utils.GetSqlTx(ctx)
	selectQuery := "election_pair_id, city_name, confirmed_votes, unique_voters, participation_rate, city_rank, last_updated"
	whereQuery := "election_pair_id = ? ORDER BY city_rank ASC LIMIT ?"

	args = append(args, electionPairID, limit)

	query := fmt.Sprintf(selectCityRankings, selectQuery, "", whereQuery)
	if sqlTrx != nil {
		err = sqlTrx.SelectContext(ctx, &rankings, query, args...)
	} else {
		err = l.db.GetMaster().SelectContext(ctx, &rankings, query, args...)
	}

	if err != nil {
		log.WithFields(log.Fields{
			"error":            err,
			"election_pair_id": electionPairID,
		}).ErrorWithCtx(ctx, "[LiveResultRepository.GetCityRankings] Failed to get city rankings")
		return nil, err
	}

	return rankings, nil
}

func (l *LiveResultRepository) GetAllElectionsSummary(ctx context.Context) ([]*model.ElectionSummary, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "LiveResultRepository.GetAllElectionsSummary")
	defer span.End()

	var (
		summaries []*model.ElectionSummary
		err       error
	)

	sqlTrx := utils.GetSqlTx(ctx)
	selectQuery := "election_pair_id, total_unique_voters, total_regions, total_confirmed_votes, overall_success_rate, last_updated"
	whereQuery := "ORDER BY total_confirmed_votes DESC"

	query := fmt.Sprintf(selectElectionSummary, selectQuery, "", whereQuery)
	if sqlTrx != nil {
		err = sqlTrx.SelectContext(ctx, &summaries, query)
	} else {
		err = l.db.GetMaster().SelectContext(ctx, &summaries, query)
	}

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).ErrorWithCtx(ctx, "[LiveResultRepository.GetAllElectionsSummary] Failed to get all elections summary")
		return nil, err
	}

	return summaries, nil
}

func (l *LiveResultRepository) GetLiveResultsByCityAndElection(ctx context.Context, cityName, electionPairID string) (*model.LiveCityResult, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "LiveResultRepository.GetLiveResultsByCityAndElection")
	defer span.End()

	var (
		result *model.LiveCityResult
		err    error
		args   []any
	)

	sqlTrx := utils.GetSqlTx(ctx)
	selectQuery := "city_name, election_pair_id, total_unique_voters, confirmed_votes, vote_success_rate, candidate_percentage_in_city, last_updated"

	whereQuery := "city_name = ? AND election_pair_id = ?"
	args = append(args, cityName, electionPairID)

	query := fmt.Sprintf(selectLiveCityResults, selectQuery, "", whereQuery)
	if sqlTrx != nil {
		err = sqlTrx.GetContext(ctx, &result, query, args...)
	} else {
		err = l.db.GetMaster().GetContext(ctx, &result, query, args...)
	}

	if err != nil {
		log.WithFields(log.Fields{
			"error":            err,
			"city_name":        cityName,
			"election_pair_id": electionPairID,
		}).ErrorWithCtx(ctx, "[LiveResultRepository.GetLiveResultsByCityAndElection] Failed to get live results by city and election")
		return nil, err
	}

	return result, nil
}
