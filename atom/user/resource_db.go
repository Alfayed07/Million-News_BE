package atom_user

import (
	"BACKEND_SEJUTA_BERITA/config/database"
	"database/sql"
	"strings"
)

func getUserByID(id int64) (Profile, error) {
    db := database.PgOpenConnection()
    defer db.Close()

    const q = `
        SELECT id, COALESCE(name,''), username, email, COALESCE(bio,''), COALESCE(avatar,''), role, is_active, registered_at
        FROM users
        WHERE id = $1
        LIMIT 1
    `
    var p Profile
    err := db.QueryRow(q, id).Scan(&p.ID, &p.Name, &p.Username, &p.Email, &p.Bio, &p.Avatar, &p.Role, &p.IsActive, &p.RegisteredAt)
    if err == sql.ErrNoRows {
        return p, err
    }
    return p, err
}

func updateUserProfile(id int64, req UpdateProfileRequest) error {
    db := database.PgOpenConnection()
    defer db.Close()

    const q = `
        UPDATE users
        SET name=$1, username=$2, email=$3, bio=$4, avatar=$5
        WHERE id=$6
    `
    _, err := db.Exec(q, strings.TrimSpace(req.Name), strings.TrimSpace(req.Username), strings.TrimSpace(req.Email), req.Bio, strings.TrimSpace(req.Avatar), id)
    return err
}

func listUsers(search string, limit, offset int) ([]UserSummary, error) {
    db := database.PgOpenConnection()
    defer db.Close()

    base := `
        SELECT id, COALESCE(name,''), username, email, role, is_active, registered_at
        FROM users
    `

    var (
        rows *sql.Rows
        err  error
    )
    if search != "" {
        like := "%" + search + "%"
        rows, err = db.Query(base+" WHERE username ILIKE $1 OR email ILIKE $1 OR COALESCE(name,'') ILIKE $1 ORDER BY registered_at DESC LIMIT $2 OFFSET $3", like, limit, offset)
    } else {
        rows, err = db.Query(base+" ORDER BY registered_at DESC LIMIT $1 OFFSET $2", limit, offset)
    }
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var items []UserSummary
    for rows.Next() {
        var u UserSummary
        if err := rows.Scan(&u.ID, &u.Name, &u.Username, &u.Email, &u.Role, &u.IsActive, &u.RegisteredAt); err != nil {
            return nil, err
        }
        items = append(items, u)
    }
    return items, rows.Err()
}

func countUsers(search string) (int64, error) {
    db := database.PgOpenConnection()
    defer db.Close()

    const base = `SELECT COUNT(*) FROM users`
    var total int64
    var err error
    if search != "" {
        like := "%" + search + "%"
        err = db.QueryRow(base+" WHERE username ILIKE $1 OR email ILIKE $1 OR COALESCE(name,'') ILIKE $1", like).Scan(&total)
    } else {
        err = db.QueryRow(base).Scan(&total)
    }
    return total, err
}

func getUserSummaryByID(id int64) (UserSummary, error) {
    db := database.PgOpenConnection()
    defer db.Close()

    const q = `
        SELECT id, COALESCE(name,''), username, email, role, is_active, registered_at
        FROM users
        WHERE id = $1
        LIMIT 1
    `
    var u UserSummary
    err := db.QueryRow(q, id).Scan(&u.ID, &u.Name, &u.Username, &u.Email, &u.Role, &u.IsActive, &u.RegisteredAt)
    if err == sql.ErrNoRows {
        return u, err
    }
    return u, err
}

func updateUserAccess(id int64, role string, isActive sql.NullBool) error {
    db := database.PgOpenConnection()
    defer db.Close()

    const q = `
        UPDATE users
        SET role = COALESCE(NULLIF($2, '')::user_role, role),
            is_active = COALESCE($3, is_active)
        WHERE id = $1
    `
    res, err := db.Exec(q, id, strings.ToLower(strings.TrimSpace(role)), isActive)
    if err != nil {
        return err
    }
    if n, _ := res.RowsAffected(); n == 0 {
        return sql.ErrNoRows
    }
    return nil
}
