

CREATE MATERIALIZED VIEW election_summary_mv
REFRESH EVERY 1 MINUTE
TO election_summary
AS SELECT
              election_pair_id,
              uniq(voter_id) as total_unique_voters,
              uniq(region) as total_regions,
              countIf(status = 'confirmed') as total_confirmed_votes,
              count() as total_vote_attempts,
              (total_confirmed_votes * 100.0) / nullIf(total_vote_attempts, 0) as overall_success_rate,
              max(updated_at) as last_updated,
              now() as updated_at
   FROM vote_results
   GROUP BY election_pair_id;