-- +goose Up
-- +goose StatementBegin
CREATE TABLE events (
   id SERIAL PRIMARY KEY,
   event_type TEXT NOT NULL,
   payload TEXT NOT NULL,
   status TEXT NOT NULL DEFAULT 'new' CHECK ( status in ('new', 'done')),
   created_at TIMESTAMPTZ DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE events;
-- +goose StatementEnd
