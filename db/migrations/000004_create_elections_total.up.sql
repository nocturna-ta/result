CREATE MATERIALIZED VIEW election_totals_mv
ENGINE = SummingMergeTree()
PARTITION BY toYYYYMM(date)
ORDER BY (election_pair_id, date)
AS SELECT
              election_pair_id,
              toDate(created_at) as date,
    count() as total_votes,
    uniq(voter_id) as total_voters,
    max(updated_at) as last_updated
   FROM vote_results
   WHERE status = 'confirmed'
   GROUP BY election_pair_id, toDate(created_at);