-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users
(
    id       UUID PRIMARY KEY,
    email    VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(100)        NOT NULL,
    type     VARCHAR(100)        NOT NULL
);

CREATE TABLE IF NOT EXISTS houses
(
    id         SERIAL PRIMARY KEY,
    address    TEXT NOT NULL,
    year       INT      NOT NULL,
    developer  TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    update_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

CREATE TABLE IF NOT EXISTS flats
(
    id        SERIAL NOT NULL,
    house_id  INT      NOT NULL,
    price     INT      NOT NULL,
    rooms     INT      NOT NULL,
    status VARCHAR(100) NOT NULL default 'created',
    foreign key (house_id) references houses(id) on delete set null
--         references houses(id) on delete set null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS flats;
DROP TABLE IF EXISTS houses;
-- +goose StatementEnd