-- +migrate Up
CREATE TABLE IF NOT EXISTS animals (
    id INT NOT NULL PRIMARY KEY,
    kind VARCHAR(128),
    age INT
);
-- +migrate Down
DROP TABLE IF EXISTS animals;