package repository

const (
	addMember = `INSERT INTO room_members (room_id, user_id, role, joined_at)
				VALUES ($1, $2, $3, CURRENT_TIMESTAMP)
				ON CONFLICT (room_id, user_id) DO UPDATE SET role = $3
				RETURNING *`

	updateRole = `UPDATE room_members SET role = $3 WHERE room_id = $1 AND user_id = $2`

	removeMember = `DELETE FROM room_members WHERE room_id = $1 AND user_id = $2`

	getMembersByRoomID = `SELECT rm.room_id, rm.user_id, rm.role, u.nickname, rm.joined_at
				FROM room_members rm
				JOIN users u ON rm.user_id = u.user_id
				WHERE rm.room_id = $1
				ORDER BY rm.joined_at ASC`

	getMemberRole = `SELECT role FROM room_members WHERE room_id = $1 AND user_id = $2`

	isMember = `SELECT EXISTS(SELECT 1 FROM room_members WHERE room_id = $1 AND user_id = $2)`
)
