package atom_auth

import (
	bcryptutil "BACKEND_SEJUTA_BERITA/utils/bcrypt"
	tokenutil "BACKEND_SEJUTA_BERITA/utils/token"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"time"
)

// LoginUseCase handles the login flow: fetch user, verify password, return token
func LoginUseCase(req LoginUserRequest) (LoginUser, string, error) {
	row, err := getUserByUsername(req.Username)
	if err != nil {
		return LoginUser{}, "", errors.New("invalid username or password")
	}

	if !row.IsActive {
		return LoginUser{}, "", errors.New("account is inactive")
	}

	if err := bcryptutil.CompareHashAndPassword(row.PasswordHash, req.Password); err != nil {
		return LoginUser{}, "", errors.New("invalid username or password")
	}

	user := LoginUser{
		ID:       row.ID,
		Username: row.Username,
		Email:    row.Email,
		Role:     row.Role,
		IsActive: row.IsActive,
	}

	token, err := tokenutil.GenerateToken(fmt.Sprintf("%d", row.ID))
	if err != nil {
		return LoginUser{}, "", errors.New("failed to generate token")
	}

	return user, token, nil
}

// RegisterUseCase creates a new user (role=user, active=true) with bcrypt password
func RegisterUseCase(req RegisterRequest) error {
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return errors.New("missing fields")
	}
	if err := validatePassword(req.Password); err != nil { return err }
	if _, err := getUserByUsername(req.Username); err == nil {
		return errors.New("username already exists")
	}
	if _, err := getUserByEmail(req.Email); err == nil {
		return errors.New("email already exists")
	}
	hash, err := bcryptutil.GenerateFromPassword(req.Password)
	if err != nil { return err }
	return insertUser(req.Username, req.Email, hash)
}

// ForgotPasswordUseCase generates reset token and stores it; returns token for development
func ForgotPasswordUseCase(req ForgotPasswordRequest) (string, error) {
	if req.Email == "" { return "", errors.New("email required") }
	user, err := getUserByEmail(req.Email)
	if err != nil { return "", errors.New("user not found") }

	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil { return "", err }
	token := hex.EncodeToString(b)
	expires := time.Now().Add(1 * time.Hour)
	if err := createPasswordResetToken(user.ID, token, expires); err != nil { return "", err }
	return token, nil
}

// ResetPasswordUseCase validates token and sets new password
func ResetPasswordUseCase(req ResetPasswordRequest) error {
	if req.Token == "" || req.NewPassword == "" { return errors.New("invalid payload") }
	if err := validatePassword(req.NewPassword); err != nil { return err }
	row, err := getPasswordResetByToken(req.Token)
	if err != nil { return errors.New("invalid token") }
	if row.UsedAt.Valid { return errors.New("token already used") }
	if time.Now().After(row.ExpiresAt) { return errors.New("token expired") }
	hash, err := bcryptutil.GenerateFromPassword(req.NewPassword)
	if err != nil { return err }
	if err := updateUserPassword(row.UserID, hash); err != nil { return err }
	return markPasswordResetUsed(row.ID)
}

// validatePassword enforces: >=7 chars, at least one lower, one upper, one digit, one special
func validatePassword(p string) error {
	if len(p) < 7 { return errors.New("password must be at least 7 characters") }
	lower := regexp.MustCompile(`[a-z]`).MatchString(p)
	upper := regexp.MustCompile(`[A-Z]`).MatchString(p)
	num := regexp.MustCompile(`\d`).MatchString(p)
	special := regexp.MustCompile(`[^A-Za-z0-9]`).MatchString(p)
	if !(lower && upper) { return errors.New("password must contain lowercase and uppercase letters") }
	if !num { return errors.New("password must contain at least one digit") }
	if !special { return errors.New("password must contain at least one special character") }
	return nil
}