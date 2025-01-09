package repository

const (
	createRoomCode = `INSERT INTO room_code (room_id, code, created_at, updated_at)
					VALUES ($1, '', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
					RETURNING *`

	updateRoomCode = `UPDATE room_code SET 
    					room_id = COALESCE(NULLIF($1::UUID, '00000000-0000-0000-0000-000000000000'::UUID), room_id),
    					code = COALESCE(NULLIF($2, ''), code),
    					updated_at = now()
						WHERE room_code_id = $3
						RETURNING *;
`

	getRoomCodeByID = `SELECT room_code_id, room_id, code, created_at, updated_at
					FROM room_code
					WHERE room_code_id = $1`

	getRoomCodeByRoomID = `SELECT room_code_id FROM room_code WHERE room_id = $1`

	deleteRoomCodeByID = `DELETE FROM room_code WHERE room_code_id = $1`
)
