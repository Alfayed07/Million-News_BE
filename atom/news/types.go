package atom_news

import "time"

type NewsItem struct {
    ID         int64     `json:"id"`
    CategoryID *int64    `json:"category_id,omitempty"`
    AuthorID   *int64    `json:"author_id,omitempty"`
    Title      string    `json:"title"`
    Content    string    `json:"content"`
    ImageURL   string    `json:"image"`
    Status     string    `json:"status"`
    CreatedAt  time.Time `json:"created_at"`
    PublishedAt *time.Time `json:"published_at,omitempty"`
    Category   *string   `json:"category,omitempty"`
}

type ListResponse struct {
    Items []NewsItem `json:"items"`
    Total int64      `json:"total"`
    Page  int        `json:"page"`
    Limit int        `json:"limit"`
}

// Comments
type Comment struct {
    ID        int64   `json:"id"`
    NewsID    int64   `json:"news_id"`
    UserID    *int64  `json:"user_id,omitempty"`
    Username  *string `json:"username,omitempty"`
    Avatar    *string `json:"avatar,omitempty"`
    Content   string  `json:"content"`
    CreatedAt string  `json:"created_at"`
}

// Categories
type Category struct {
    ID   int64  `json:"id"`
    Name string `json:"name"`
}
