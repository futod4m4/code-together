package repository

const (
	createFile = `INSERT INTO room_files (room_id, filename, language, content, is_entry_point, created_at, updated_at)
				VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
				RETURNING *`

	updateFile = `UPDATE room_files
				SET content = $1, language = $2, filename = COALESCE(NULLIF($3, ''), filename), updated_at = now()
				WHERE file_id = $4
				RETURNING *`

	deleteFile = `DELETE FROM room_files WHERE file_id = $1`

	getFileByID = `SELECT file_id, room_id, filename, language, content, is_entry_point, created_at, updated_at
				FROM room_files WHERE file_id = $1`

	getFilesByRoomID = `SELECT file_id, room_id, filename, language, content, is_entry_point, created_at, updated_at
				FROM room_files WHERE room_id = $1 ORDER BY is_entry_point DESC, filename ASC`

	countFilesByRoomID = `SELECT COUNT(*) FROM room_files WHERE room_id = $1`
)
