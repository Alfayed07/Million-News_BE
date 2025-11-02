package atom_auth

// LoginUserRequest represents the payload sent from client to login
type LoginUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginUser represents a subset of user data returned after login
type LoginUser struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	IsActive  bool   `json:"is_active"`
}

// internal db scan model (includes password hash)
type loginRow struct {
	ID           int64
	Username     string
	Email        string
	PasswordHash string
	Role         string
	IsActive     bool
}

// Register
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Forgot/Reset Password
type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}