CREATE TABLE IF NOT EXISTS banned_members (
    ban_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    room_id UUID REFERENCES rooms(room_id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(user_id) ON DELETE CASCADE,
    banned_by UUID REFERENCES users(user_id),
    ip_address VARCHAR(45),
    reason VARCHAR(500) DEFAULT '',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(room_id, user_id)
);

CREATE INDEX idx_banned_members_room ON banned_members(room_id);
CREATE INDEX idx_banned_members_ip ON banned_members(room_id, ip_address);
