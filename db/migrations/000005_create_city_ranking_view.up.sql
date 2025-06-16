CREATE MATERIALIZED VIEW city_ranking_mv
ENGINE = ReplacingMergeTree(updated_at)
PARTITION BY toYYYYMM(last_updated)
ORDER BY (election_pair_id, confirmed_votes)
AS SELECT
    election_pair_id,
    region as city_name,
    countIf(status = 'confirmed') as confirmed_votes,
    uniq(voter_id) as unique_voters,
    (confirmed_votes * 100.0) / nullIf(unique_voters, 0) as participation_rate,
    rank() OVER (PARTITION BY election_pair_id ORDER BY confirmed_votes DESC) as city_rank,
    max(updated_at) as last_updated,
    now() as updated_at
FROM vote_results
GROUP BY election_pair_id, region;