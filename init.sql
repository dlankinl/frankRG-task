CREATE TABLE IF NOT EXISTS Files(
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    size BIGINT NOT NULL,
    modtime timestamptz NOT NULL,
    isdirectory BOOLEAN NOT NULL,
    content TEXT NOT NULL,
    parentid INT NOT NULL
);