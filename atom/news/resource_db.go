package atom_news

import (
	"BACKEND_SEJUTA_BERITA/config/database"
	"database/sql"
	"fmt"
	"strings"
)

func scanNewsRow(row *sql.Row) (NewsItem, error) {
    var it NewsItem
    var imgBytes []byte
    var catName sql.NullString
    err := row.Scan(&it.ID, &it.CategoryID, &it.AuthorID, &it.Title, &it.Content, &imgBytes, &it.Status, &it.CreatedAt, &it.PublishedAt, &catName)
    if catName.Valid { name := catName.String; it.Category = &name }
    if len(imgBytes) > 0 {
        // Assume image column stores URL bytes or raw; try to decode as string
        it.ImageURL = string(imgBytes)
    }
    return it, err
}

func scanNewsRows(rows *sql.Rows) ([]NewsItem, error) {
    list := []NewsItem{}
    defer rows.Close()
    for rows.Next() {
        var it NewsItem
        var imgBytes []byte
        var catName sql.NullString
        if err := rows.Scan(&it.ID, &it.CategoryID, &it.AuthorID, &it.Title, &it.Content, &imgBytes, &it.Status, &it.CreatedAt, &it.PublishedAt, &catName); err != nil {
            return list, err
        }
        if catName.Valid { name := catName.String; it.Category = &name }
        if len(imgBytes) > 0 { it.ImageURL = string(imgBytes) }
        list = append(list, it)
    }
    return list, rows.Err()
}

const baseSelect = `
SELECT n.id, n.category_id, n.author_id, n.title, n.content, n.image, n.status, n.created_at, n.published_at, c.name
FROM news n
LEFT JOIN categories c ON c.id = n.category_id
`

func listNews(category string, page, limit int) ([]NewsItem, int64, error) {
    db := database.PgOpenConnection()
    defer db.Close()
    if limit <= 0 { limit = 10 }
    if page <= 0 { page = 1 }
    offset := (page-1)*limit

    var rows *sql.Rows
    var err error
    var total int64

    if strings.TrimSpace(category) != "" {
        q := baseSelect + " WHERE LOWER(c.name) = LOWER($1) AND n.status='published'::news_status ORDER BY COALESCE(n.published_at,n.created_at) DESC LIMIT $2 OFFSET $3"
        rows, err = db.Query(q, category, limit, offset)
        if err != nil { return nil, 0, err }
        err = db.QueryRow("SELECT COUNT(1) FROM news n LEFT JOIN categories c ON c.id=n.category_id WHERE LOWER(c.name)=LOWER($1) AND n.status='published'::news_status", category).Scan(&total)
    } else {
        q := baseSelect + " WHERE n.status='published'::news_status ORDER BY COALESCE(n.published_at,n.created_at) DESC LIMIT $1 OFFSET $2"
        rows, err = db.Query(q, limit, offset)
        if err != nil { return nil, 0, err }
        err = db.QueryRow("SELECT COUNT(1) FROM news WHERE status='published'::news_status").Scan(&total)
    }
    if err != nil { return nil, 0, err }
    items, err := scanNewsRows(rows)
    return items, total, err
}

func topStories(limit int) ([]NewsItem, error) {
    db := database.PgOpenConnection(); defer db.Close()
    if limit <= 0 { limit = 5 }
    q := baseSelect + " WHERE n.status='published'::news_status ORDER BY COALESCE(n.published_at,n.created_at) DESC LIMIT $1"
    rows, err := db.Query(q, limit)
    if err != nil { return nil, err }
    return scanNewsRows(rows)
}

func trending(limit int) ([]NewsItem, error) {
    db := database.PgOpenConnection(); defer db.Close()
    if limit <= 0 { limit = 5 }
    // Trending heuristic combining recent comments (24h, 7d) and views from news_metrics
    const q = `
        WITH cm AS (
            SELECT c.news_id,
                   COUNT(*) FILTER (WHERE c.created_at >= NOW() - INTERVAL '24 hours') AS c1,
                   COUNT(*) FILTER (WHERE c.created_at >= NOW() - INTERVAL '7 days')  AS c7
            FROM comments c
            WHERE (c.is_approved = TRUE OR c.is_approved IS NULL)
            GROUP BY c.news_id
        )
        SELECT n.id, n.title
        FROM news n
        LEFT JOIN cm ON cm.news_id = n.id
        LEFT JOIN news_metrics m ON m.news_id = n.id
        WHERE n.status='published'::news_status
        ORDER BY (COALESCE(cm.c1,0)*3 + COALESCE(cm.c7,0)*1 + COALESCE(m.views,0)*0.05) DESC,
                 COALESCE(n.published_at,n.created_at) DESC
        LIMIT $1`
    rows, err := db.Query(q, limit)
    if err != nil { return nil, err }
    out := []NewsItem{}
    defer rows.Close()
    for rows.Next() {
        var id int64; var title string
        if err := rows.Scan(&id, &title); err != nil { return nil, err }
        out = append(out, NewsItem{ID: id, Title: title})
    }
    return out, rows.Err()
}

func searchNews(qs string, page, limit int) ([]NewsItem, int64, error) {
    db := database.PgOpenConnection(); defer db.Close()
    if limit <= 0 { limit = 20 }
    if page <= 0 { page = 1 }
    offset := (page-1)*limit
    term := strings.TrimSpace(qs)
    if term == "" { return []NewsItem{}, 0, nil }
    like := fmt.Sprintf("%%%s%%", term)

    // Prefer trigram similarity when available for ordering, with fallback to ILIKE; use OFFSET for paging
    qAdvPub := baseSelect + `
        WHERE n.status='published'::news_status
          AND (n.title ILIKE $1 OR n.title % $2)
        ORDER BY similarity(n.title, $2) DESC NULLS LAST,
                 COALESCE(n.published_at,n.created_at) DESC
        LIMIT $3 OFFSET $4`
    if rows, err := db.Query(qAdvPub, like, term, limit, offset); err == nil {
        if items, err2 := scanNewsRows(rows); err2 == nil && len(items) > 0 {
            // total (approx) using ILIKE only to avoid pg_trgm dependency
            var total int64
            _ = db.QueryRow("SELECT COUNT(1) FROM news n WHERE n.status='published'::news_status AND n.title ILIKE $1", like).Scan(&total)
            return items, total, nil
        }
    }

    // Simple ILIKE with status=published (paging)
    qSimplePub := baseSelect + `
        WHERE n.status='published'::news_status AND n.title ILIKE $1
        ORDER BY COALESCE(n.published_at,n.created_at) DESC
        LIMIT $2 OFFSET $3`
    if rows, err := db.Query(qSimplePub, like, limit, offset); err == nil {
        items, err2 := scanNewsRows(rows)
        if err2 != nil { return []NewsItem{}, 0, err2 }
        var total int64
        _ = db.QueryRow("SELECT COUNT(1) FROM news n WHERE n.status='published'::news_status AND n.title ILIKE $1", like).Scan(&total)
        return items, total, nil
    }

    return []NewsItem{}, 0, nil
}

func getNewsByID(id int64) (NewsItem, error) {
    db := database.PgOpenConnection(); defer db.Close()
    q := baseSelect + " WHERE n.id = $1"
    row := db.QueryRow(q, id)
    return scanNewsRow(row)
}

// Comments queries
func listCommentsForNews(newsID int64, limit int) ([]Comment, error) {
    db := database.PgOpenConnection(); defer db.Close()
    if limit <= 0 { limit = 50 }
    const q = `
        SELECT c.id, c.news_id, c.user_id, u.username, u.avatar, c.content, to_char(c.created_at,'YYYY-MM-DD"T"HH24:MI:SSZ')
        FROM comments c
        LEFT JOIN users u ON u.id = c.user_id
        WHERE c.news_id = $1 AND (c.is_approved = TRUE OR c.is_approved IS NULL)
        ORDER BY c.created_at ASC
        LIMIT $2
    `
    rows, err := db.Query(q, newsID, limit)
    if err != nil { return nil, err }
    defer rows.Close()
    out := []Comment{}
    for rows.Next() {
        var cm Comment
        err := rows.Scan(&cm.ID, &cm.NewsID, &cm.UserID, &cm.Username, &cm.Avatar, &cm.Content, &cm.CreatedAt)
        if err != nil { return nil, err }
        out = append(out, cm)
    }
    return out, rows.Err()
}

func insertCommentForNews(newsID int64, userID *int64, content string) (Comment, error) {
    db := database.PgOpenConnection(); defer db.Close()
    // mark approved directly for now
    const q = `
        INSERT INTO comments (news_id, user_id, content, is_approved)
        VALUES ($1, $2, $3, TRUE)
        RETURNING id, news_id, user_id, content, to_char(created_at,'YYYY-MM-DD"T"HH24:MI:SSZ')
    `
    var cm Comment
    cm.Username = nil; cm.Avatar = nil
    err := db.QueryRow(q, newsID, userID, strings.TrimSpace(content)).Scan(&cm.ID, &cm.NewsID, &cm.UserID, &cm.Content, &cm.CreatedAt)
    if err != nil { return Comment{}, err }
    // fill username/avatar if userID provided
    if cm.UserID != nil {
        var username, avatar sql.NullString
        _ = db.QueryRow("SELECT username, avatar FROM users WHERE id=$1", cm.UserID).Scan(&username, &avatar)
        if username.Valid { s:=username.String; cm.Username=&s }
        if avatar.Valid { s:=avatar.String; cm.Avatar=&s }
    }
    return cm, nil
}

// Categories
func listCategories() ([]Category, error) {
    db := database.PgOpenConnection(); defer db.Close()
    rows, err := db.Query("SELECT id, name FROM categories ORDER BY name ASC")
    if err != nil { return nil, err }
    defer rows.Close()
    out := []Category{}
    for rows.Next() {
        var c Category
        if err := rows.Scan(&c.ID, &c.Name); err != nil { return nil, err }
        out = append(out, c)
    }
    return out, rows.Err()
}

// Metrics
func incrementNewsView(newsID int64) error {
    db := database.PgOpenConnection(); defer db.Close()
    const q = `
        INSERT INTO news_metrics (news_id, views, last_view_at)
        VALUES ($1, 1, NOW())
        ON CONFLICT (news_id)
        DO UPDATE SET views = news_metrics.views + 1, last_view_at = NOW()`
    _, err := db.Exec(q, newsID)
    return err
}
