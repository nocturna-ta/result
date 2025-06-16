CREATE MATERIALIZED VIEW live_city_results_mv
ENGINE = ReplacingMergeTree(updated_at)
PARTITION BY toYYYYMM(last_updated)
ORDER BY (city_name, election_pair_id)
AS SELECT
              region as city_name,
              election_pair_id,
              uniq(voter_id) as total_unique_voters,
              countIf(status = 'confirmed') as confirmed_votes,
              count() as total_vote_attempts,
              (confirmed_votes * 100.0) / nullIf(total_vote_attempts, 0) as vote_success_rate,
              (confirmed_votes * 100.0) / nullIf(
                      (SELECT countIf(status = 'confirmed') FROM vote_results WHERE region = city_name), 0
                                          ) as candidate_percentage_in_city,
              max(updated_at) as last_updated,
              now() as updated_at
   FROM vote_results
   GROUP BY city_name, election_pair_id;