CREATE MATERIALIZED VIEW city_live_results_mv
ENGINE = SummingMergeTree()
PARTITION BY toYYYYMM(date)
ORDER BY (city, election_pair_id, status, date)
AS SELECT
              region as city,
              election_pair_id,
              status,
              toDate(created_at) as date,
    count() as vote_count,
    uniq(voter_id) as unique_voters,
    max(updated_at) as last_updated
   FROM vote_results
   WHERE status = 'confirmed'
   GROUP BY region, election_pair_id, status, toDate(created_at);