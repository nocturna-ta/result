create table live_city_result
(
    city_name                    String,
    election_pair_id             String,
    total_unique_voters          UInt64,
    confirmed_votes              UInt64,
    total_vote_attempts          UInt64,
    vote_success_rate            Float64,
    candidate_percentage_in_city Float64,
    last_updated                 DateTime,
    updated_at                   DateTime
)
    engine = ReplacingMergeTree(updated_at)
        PARTITION BY toYYYYMM(last_updated)
        ORDER BY (city_name, election_pair_id)
        SETTINGS index_granularity = 8192;