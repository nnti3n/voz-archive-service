-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE TABLE boxes(id INTEGER PRIMARY KEY);

CREATE TABLE threads(id INTEGER PRIMARY KEY, title TEXT not null, source TEXT, page_count INT not null, post_count INT not null, view_count INT not null, user_id_starter INT, user_name_starter VARCHAR(40), last_updated TIMESTAMPTZ not null, box_id INTEGER REFERENCES boxes not null);

CREATE TABLE posts(id INTEGER PRIMARY KEY, user_id INTEGER not null, number INTEGER, user_name VARCHAR(40), content TEXT, time TIMESTAMPTZ not null, thread_id INTEGER REFERENCES threads not null);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

DROP TABLE posts;
DROP TABLE threads;
DROP TABLE boxes;