CREATE SCHEMA IF NOT EXISTS url_schema;

CREATE TABLE IF NOT EXISTS url_schema.url (
                                              originalURL VARCHAR(2048) UNIQUE NOT NULL,
    shortURL VARCHAR(7) UNIQUE NOT NULL,
    count INTEGER DEFAULT 0,
    last_counter BIGINT DEFAULT 100000000000
    );