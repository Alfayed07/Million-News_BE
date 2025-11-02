package atom_auth

import (
	"BACKEND_SEJUTA_BERITA/config/database"
	"database/sql"
	"time"
)

// getUserByUsername retrieves a user row by username
func getUserByUsername(username string) (loginRow, error) {
	db := database.PgOpenConnection()
	defer db.Close()

	const q = `
		SELECT id, username, email, password_hash, role, is_active
		FROM users
		WHERE username = $1
		LIMIT 1
	`

	var row loginRow
	err := db.QueryRow(q, username).Scan(
		&row.ID,
		&row.Username,
		&row.Email,
		&row.PasswordHash,
		&row.Role,
		&row.IsActive,
	)

	if err == sql.ErrNoRows {
		return row, err
	}
	return row, err
}

func getUserByEmail(email string) (loginRow, error) {
	db := database.PgOpenConnection()
	defer db.Close()

	const q = `
		SELECT id, username, email, password_hash, role, is_active
		FROM users
		WHERE email = $1
		LIMIT 1
	`

	var row loginRow
	err := db.QueryRow(q, email).Scan(
		&row.ID,
		&row.Username,
		&row.Email,
		&row.PasswordHash,
		&row.Role,
		&row.IsActive,
	)
	if err == sql.ErrNoRows {
		return row, err
	}
	return row, err
}

func insertUser(username, email, passwordHash string) error {
	db := database.PgOpenConnection()
	defer db.Close()

	const q = `
		INSERT INTO users (username, email, password_hash, role, is_active)
		VALUES ($1, $2, $3, 'user', TRUE)
	`
	_, err := db.Exec(q, username, email, passwordHash)
	return err
}

func createPasswordResetToken(userID int64, token string, expiresAt time.Time) error {
	db := database.PgOpenConnection()
	defer db.Close()
	const q = `
		INSERT INTO password_resets (user_id, token, expires_at)
		VALUES ($1, $2, $3)
	`
	_, err := db.Exec(q, userID, token, expiresAt)
	return err
}

type passwordResetRow struct {
	ID        int64
	UserID    int64
	Token     string
	ExpiresAt time.Time
	UsedAt    sql.NullTime
}

func getPasswordResetByToken(token string) (passwordResetRow, error) {
	db := database.PgOpenConnection()
	defer db.Close()
	const q = `
		SELECT id, user_id, token, expires_at, used_at
		FROM password_resets
		WHERE token = $1
		LIMIT 1
	`
	var r passwordResetRow
	err := db.QueryRow(q, token).Scan(&r.ID, &r.UserID, &r.Token, &r.ExpiresAt, &r.UsedAt)
	if err == sql.ErrNoRows {
		return r, err
	}
	return r, err
}

func markPasswordResetUsed(id int64) error {
	db := database.PgOpenConnection()
	defer db.Close()
	const q = `UPDATE password_resets SET used_at = NOW() WHERE id = $1`
	_, err := db.Exec(q, id)
	return err
}

func updateUserPassword(userID int64, passwordHash string) error {
	db := database.PgOpenConnection()
	defer db.Close()
	const q = `UPDATE users SET password_hash = $1 WHERE id = $2`
	_, err := db.Exec(q, passwordHash, userID)
	return err
}