CREATE TABLE IF NOT EXISTS vote_results (
    id String,
    voter_id String,
    election_pair_id String,
    region String,
    status String,
    transaction_hash String,
    error_message String,
    voted_at DateTime,
    processed_at Nullable(DateTime),
    created_at DateTime DEFAULT now(),
    updated_at DateTime DEFAULT now()
    ) ENGINE = MergeTree()
    ORDER BY (election_pair_id, region, created_at)
    PARTITION BY toYYYYMM(created_at)
    SETTINGS index_granularity = 8192;
