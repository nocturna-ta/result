package repository

import (
	"context"
	"github.com/nocturna-ta/result/internal/domain/model"
)

type LiveResultRepository interface {
	GetLiveElectionResults(ctx context.Context, electionPairID string) ([]*model.LiveElectionResult, error)
	GetLiveCityResults(ctx context.Context, cityName string) ([]*model.LiveCityResult, error)
	GetElectionSummary(ctx context.Context, electionPairID string) (*model.ElectionSummary, error)
	GetCityRankings(ctx context.Context, electionPairID string, limit int) ([]*model.CityRanking, error)
	GetAllElectionsSummary(ctx context.Context) ([]*model.ElectionSummary, error)
	GetLiveResultsByCityAndElection(ctx context.Context, cityName, electionPairID string) (*model.LiveCityResult, error)
}
