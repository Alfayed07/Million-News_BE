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
        SELECT id, COALESCE(name,''), username, email, COALESCE(bio,''), COALESCE(avatar,''), role, is_active
        FROM users
        WHERE id = $1
        LIMIT 1
    `
    var p Profile
    err := db.QueryRow(q, id).Scan(&p.ID, &p.Name, &p.Username, &p.Email, &p.Bio, &p.Avatar, &p.Role, &p.IsActive)
    if err == sql.ErrNoRows { return p, err }
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
