create table city_ranking
(
    election_pair_id   String,
    city_name          String,
    confirmed_votes    UInt64,
    unique_voters      UInt64,
    participation_rate Float64,
    city_rank          UInt64,
    last_updated       DateTime,
    updated_at         DateTime
)
    engine = ReplacingMergeTree(updated_at)
        PARTITION BY toYYYYMM(last_updated)
        ORDER BY (election_pair_id, confirmed_votes)
        SETTINGS index_granularity = 8192;

