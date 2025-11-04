package atom_user

// Profile represents the user profile returned to clients
type Profile struct {
    ID       int64  `json:"id"`
    Name     string `json:"name"`
    Username string `json:"username"`
    Email    string `json:"email"`
    Bio      string `json:"bio"`
    Avatar   string `json:"avatar"`
    Role     string `json:"role"`
    IsActive bool   `json:"is_active"`
}

// UpdateProfileRequest represents allowed fields to update
type UpdateProfileRequest struct {
    Name     string `json:"name"`
    Username string `json:"username"`
    Email    string `json:"email"`
    Bio      string `json:"bio"`
    Avatar   string `json:"avatar"`
}
