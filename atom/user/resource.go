package atom_user

import (
	"errors"
	"net/mail"
	"strings"
)

func GetProfileUseCase(userID int64) (Profile, error) {
    return getUserByID(userID)
}

func UpdateProfileUseCase(userID int64, req UpdateProfileRequest) (Profile, error) {
    // Basic validations
    req.Name = strings.TrimSpace(req.Name)
    req.Username = strings.TrimSpace(req.Username)
    req.Email = strings.TrimSpace(req.Email)
    req.Avatar = strings.TrimSpace(req.Avatar)

    if req.Username == "" || req.Email == "" {
        return Profile{}, errors.New("username and email are required")
    }
    if _, err := mail.ParseAddress(req.Email); err != nil {
        return Profile{}, errors.New("invalid email format")
    }
    if err := updateUserProfile(userID, req); err != nil {
        // Surface unique constraint issues
        if strings.Contains(err.Error(), "duplicate key value") {
            if strings.Contains(err.Error(), "users_username_key") {
                return Profile{}, errors.New("username already exists")
            }
            if strings.Contains(err.Error(), "users_email_key") {
                return Profile{}, errors.New("email already exists")
            }
            return Profile{}, errors.New("duplicate value")
        }
        return Profile{}, err
    }
    return getUserByID(userID)
}
