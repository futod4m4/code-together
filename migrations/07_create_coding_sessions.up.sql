CREATE TABLE IF NOT EXISTS coding_sessions (
    session_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    room_id UUID REFERENCES rooms(room_id) ON DELETE CASCADE,
    title VARCHAR(200) NOT NULL DEFAULT 'Live Session',
    started_by UUID REFERENCES users(user_id),
    started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ended_at TIMESTAMP,
    is_active BOOLEAN DEFAULT true,
    max_viewers INT DEFAULT 0
);

CREATE TABLE IF NOT EXISTS session_snapshots (
    snapshot_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    session_id UUID REFERENCES coding_sessions(session_id) ON DELETE CASCADE,
    code TEXT NOT NULL,
    language VARCHAR(50),
    filename VARCHAR(255) DEFAULT 'main',
    timestamp_ms BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_coding_sessions_room ON coding_sessions(room_id);
CREATE INDEX idx_session_snapshots_session ON session_snapshots(session_id);
