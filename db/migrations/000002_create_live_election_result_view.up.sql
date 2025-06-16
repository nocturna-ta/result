CREATE MATERIALIZED VIEW live_election_results_mv
ENGINE = ReplacingMergeTree(updated_at)
PARTITION BY toYYYYMM(last_updated)
ORDER BY (election_pair_id, region)
AS SELECT
              election_pair_id,
              region,
              countIf(status = 'confirmed') as confirmed_votes,
              countIf(status = 'pending') as pending_votes,
              countIf(status = 'error') as error_votes,
              count() as total_votes,
              (confirmed_votes * 100.0) / nullIf(total_votes, 0) as success_percentage,
              max(updated_at) as last_updated,
              now() as updated_at
   FROM vote_results
   GROUP BY election_pair_id, region;