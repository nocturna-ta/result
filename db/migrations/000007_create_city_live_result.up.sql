CREATE MATERIALIZED VIEW live_city_results_mv
    REFRESH EVERY 1 MINUTE
    TO live_city_result
AS SELECT
              region as city_name,
              election_pair_id,
              uniq(voter_id) as total_unique_voters,
              countIf(status = 'confirmed') as confirmed_votes,
              count() as total_vote_attempts,
              (confirmed_votes * 100.0) / nullIf(total_vote_attempts, 0) as vote_success_rate,
              -- Fixed: Use window function to ensure percentages sum to 100% per city
              (confirmed_votes * 100.0) / nullIf(
                      sum(confirmed_votes) OVER (PARTITION BY region), 0
                                          ) as candidate_percentage_in_city,
              max(updated_at) as last_updated,
              now() as updated_at
   FROM vote_results
   GROUP BY region, election_pair_id;
