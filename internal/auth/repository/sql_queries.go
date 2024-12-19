package repository

const (
	createUserQuery = `INSERT INTO users (email, password, nickname, first_name, second_name) VALUES ($1, $2, $3, $4, $5) RETURNING *`

	updateUserQuery = `UPDATE users 
						SET first_name = COALESCE(NULLIF($1, ''), first_name),
						    second_name = COALESCE(NULLIF($2, ''), second_name),
						    nickname = COALESCE(NULLIF($3, ''), nickname),
						    email = COALESCE(NULLIF($4, ''), email)
						WHERE user_id = $5
						RETURNING *
						`

	deleteUserQuery = `DELETE FROM users WHERE user_id = $1`

	getByIDQuery = `SELECT user_id, first_name, second_name, email, password
				 		FROM users 
				 		WHERE user_id = $1`

	findUserByEmailQuery = `SELECT user_id, first_name, second_name, email, password
				 		FROM users 
				 		WHERE email = $1`
)
