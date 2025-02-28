CREATE TABLE urls (
    id SERIAL,
    from_url VARCHAR(255) NOT NULL,
    to_url VARCHAR(255) NOT NULL,
    hit_count INT DEFAULT 0
)