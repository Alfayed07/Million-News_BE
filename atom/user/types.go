package atom_user

import "time"

// Profile represents the user profile returned to clients
type Profile struct {
    ID          int64     `json:"id"`
    Name        string    `json:"name"`
    Username    string    `json:"username"`
    Email       string    `json:"email"`
    Bio         string    `json:"bio"`
    Avatar      string    `json:"avatar"`
    Role        string    `json:"role"`
    IsActive    bool      `json:"is_active"`
    RegisteredAt time.Time `json:"registered_at"`
}

// UpdateProfileRequest represents allowed fields to update
type UpdateProfileRequest struct {
    Name     string `json:"name"`
    Username string `json:"username"`
    Email    string `json:"email"`
    Bio      string `json:"bio"`
    Avatar   string `json:"avatar"`
}

// UserSummary is a lightweight view used for admin management tables
type UserSummary struct {
    ID          int64     `json:"id"`
    Name        string    `json:"name"`
    Username    string    `json:"username"`
    Email       string    `json:"email"`
    Role        string    `json:"role"`
    IsActive    bool      `json:"is_active"`
    RegisteredAt time.Time `json:"registered_at"`
}

// ListUsersResponse wraps paginated results for admin listings
type ListUsersResponse struct {
    Items []UserSummary `json:"items"`
    Page  int           `json:"page"`
    Limit int           `json:"limit"`
    Total int64         `json:"total"`
    Pages int           `json:"pages"`
    Search string       `json:"search"`
}

// UpdateUserAccessRequest allows administrators to modify role/active flags
type UpdateUserAccessRequest struct {
    Role     string `json:"role"`
    IsActive *bool  `json:"is_active"`
}
