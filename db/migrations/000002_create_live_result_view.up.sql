CREATE MATERIALIZED VIEW election_live_results_mv
ENGINE = SummingMergeTree()
PARTITION BY toYYYYMM(date)
ORDER BY (election_pair_id, status, date)
AS SELECT
    election_pair_id,
    status,
    toDate(created_at) as date,
    count() as vote_count,
    max(updated_at) as last_updated
   FROM vote_results
   WHERE status = 'confirmed'
   GROUP BY election_pair_id, status, toDate(created_at);