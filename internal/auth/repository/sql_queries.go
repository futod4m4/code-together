package repository

const (
	createUserQuery = `INSERT INTO users (email, password, nickname, first_name, last_name, avatar_url, github_url, bio) VALUES ($1, $2, $3, $4, $5, '', '', '') RETURNING *`

	updateUserQuery = `UPDATE users 
						SET first_name = COALESCE(NULLIF($1, ''), first_name),
							last_name = COALESCE(NULLIF($2, ''),last_name),
						    nickname = COALESCE(NULLIF($3, ''), nickname),
						    email = COALESCE(NULLIF($4, ''), email),
						WHERE user_id = $5
						RETURNING *
						`

	deleteUserQuery = `DELETE FROM users WHERE user_id = $1`

	getByIDQuery = `SELECT user_id, first_name, last_name, email, password, nickname, avatar_url, github_url, bio
				 		FROM users
				 		WHERE user_id = $1`

	findUserByEmailQuery = `SELECT user_id, first_name, last_name, email, password, nickname, avatar_url, github_url, bio FROM users WHERE email = $1`

	updateProfileQuery = `UPDATE users
SET avatar_url = COALESCE(NULLIF($1, ''), avatar_url),
    github_url = COALESCE(NULLIF($2, ''), github_url),
    bio = COALESCE(NULLIF($3, ''), bio)
WHERE user_id = $4
RETURNING user_id, first_name, last_name, email, nickname, avatar_url, github_url, bio`
)
