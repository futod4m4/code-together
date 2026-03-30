package repository

const (
	createMessage = `INSERT INTO room_messages (room_id, user_id, nickname, content, created_at)
					VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP)
					RETURNING *`

	getMessagesByRoomID = `SELECT message_id, room_id, user_id, nickname, content, created_at
					FROM room_messages
					WHERE room_id = $1
					ORDER BY created_at ASC
					LIMIT $2 OFFSET $3`
)
