package repository

const (
	createRoom = `INSERT INTO rooms (room_name, join_code, language, owner_id, description, is_private, created_at, updated_at)
					VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
					RETURNING *`

	updateRoom = `UPDATE rooms
					SET room_name = COALESCE(NULLIF($1, ''), room_name),
					    language = COALESCE(NULLIF($2, ''), language),
					    owner_id = COALESCE(NULLIF($3, '00000000-0000-0000-0000-000000000000'), owner_id),
					    description = $5,
					    is_private = $6,
					    updated_at = now()
					WHERE room_id = $4
					RETURNING *`

	getRoomByID = `SELECT room_id, room_name, join_code, owner_id, language, description, is_private, created_at, updated_at
					FROM rooms
					WHERE room_id = $1`

	getRoomByJoinCode = `SELECT room_id, room_name, join_code, owner_id, language, description, is_private, created_at, updated_at
					FROM rooms
					WHERE join_code = $1`

	deleteRoomByID = `DELETE FROM rooms WHERE room_id = $1`

	getRoomsByOwnerID = `SELECT room_id, room_name, join_code, owner_id, language, description, is_private, created_at, updated_at
					FROM rooms
					WHERE owner_id = $1
					ORDER BY updated_at DESC`
)
