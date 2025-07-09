create table election_summary
(
    election_pair_id      String,
    total_unique_voters   UInt64,
    total_regions         UInt64,
    total_confirmed_votes UInt64,
    total_vote_attempts   UInt64,
    overall_success_rate  Float64,
    last_updated          DateTime,
    updated_at            DateTime
)
    engine = ReplacingMergeTree(updated_at)
        PARTITION BY toYYYYMM(last_updated)
        ORDER BY election_pair_id
        SETTINGS index_granularity = 8192;
