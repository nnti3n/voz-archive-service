
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE INDEX thread_posts ON posts (thread_id);


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP INDEX thread_posts;