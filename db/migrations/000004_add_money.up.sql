CREATE TABLE money_records (
    id SERIAL PRIMARY KEY, user_id integer NOT NULL REFERENCES users ON DELETE CASCADE, reference VARCHAR(50) UNIQUE NOT NULL, status VARCHAR(50) NOT NULL, amount bigint NOT NULL
);