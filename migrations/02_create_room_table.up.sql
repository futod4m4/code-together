DROP TABLE IF EXISTS rooms CASCADE;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS CITEXT;


CREATE TABLE rooms (
    room_id UUID PRIMARY KEY default uuid_generate_v4(),
    room_name VARCHAR(50) NOT NULL,
    join_code VARCHAR(12) UNIQUE NOT NULL,
    language VARCHAR(50) NOT NULL,
    owner_id UUID REFERENCES users(user_id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx_rooms_join_code ON rooms(join_code);

CREATE TABLE room_code (
    room_code_id UUID PRIMARY KEY default uuid_generate_v4(),
    room_id UUID REFERENCES rooms(room_id) ON DELETE CASCADE,
    code TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);



