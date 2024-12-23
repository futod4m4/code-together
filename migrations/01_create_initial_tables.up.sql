DROP TABLE IF EXISTS users CASCADE;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS CITEXT;

CREATE TABLE users (
    user_id UUID PRIMARY KEY        DEFAULT uuid_generate_v4(),
    first_name VARCHAR(32)          NOT NULL CHECK ( first_name <> '' ),
    last_name VARCHAR(32)         NOT NULL CHECK ( last_name <> '' ),
    nickname VARCHAR (32)           NOT NULL CHECK ( nickname <> '' ),
    email VARCHAR(64)               UNIQUE NOT NULL CHECK ( email <> '' ),
    password VARCHAR(256)           NOT NULL CHECK ( octet_length(password) <> 0 )
);