create table live_election_results
(
    election_pair_id   String,
    region             String,
    confirmed_votes    UInt64,
    pending_votes      UInt64,
    error_votes        UInt64,
    total_votes        UInt64,
    success_percentage Float64,
    last_updated       DateTime,
    updated_at         DateTime
)
    engine = ReplacingMergeTree(updated_at)
        PARTITION BY toYYYYMM(last_updated)
        ORDER BY (election_pair_id, region)
        SETTINGS index_granularity = 8192;