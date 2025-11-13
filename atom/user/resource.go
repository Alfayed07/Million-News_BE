package atom_user

import (
	"database/sql"
	"errors"
	"net/mail"
	"strings"
)

var allowedRoles = map[string]struct{}{
    "admin":  {},
    "editor": {},
    "user":   {},
}

var (
    ErrInvalidRole         = errors.New("invalid role")
    ErrNoUpdateFields      = errors.New("no update fields provided")
    ErrUserNotFound        = errors.New("user not found")
    ErrCannotDeactivateSelf = errors.New("cannot deactivate own account")
    ErrCannotDowngradeSelf = errors.New("cannot change own role")
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

func ListUsersUseCase(search string, page, limit int) (ListUsersResponse, error) {
    search = strings.TrimSpace(search)
    if limit <= 0 {
        limit = 10
    }
    if limit > 100 {
        limit = 100
    }
    if page <= 0 {
        page = 1
    }
    offset := (page - 1) * limit

    items, err := listUsers(search, limit, offset)
    if err != nil {
        return ListUsersResponse{}, err
    }
    total, err := countUsers(search)
    if err != nil {
        return ListUsersResponse{}, err
    }
    pages := 0
    if limit > 0 {
        pages = int((total + int64(limit) - 1) / int64(limit))
    }
    return ListUsersResponse{
        Items:  items,
        Page:   page,
        Limit:  limit,
        Total:  total,
        Pages:  pages,
        Search: search,
    }, nil
}

func UpdateUserAccessUseCase(targetUserID, actorID int64, req UpdateUserAccessRequest) (UserSummary, error) {
    desiredRole := strings.ToLower(strings.TrimSpace(req.Role))
    if desiredRole != "" {
        if _, ok := allowedRoles[desiredRole]; !ok {
            return UserSummary{}, ErrInvalidRole
        }
    }
    if desiredRole == "" && req.IsActive == nil {
        return UserSummary{}, ErrNoUpdateFields
    }

    if targetUserID == actorID {
        if req.IsActive != nil && !*req.IsActive {
            return UserSummary{}, ErrCannotDeactivateSelf
        }
        if desiredRole != "" && desiredRole != "admin" {
            return UserSummary{}, ErrCannotDowngradeSelf
        }
    }

    var active sql.NullBool
    if req.IsActive != nil {
        active = sql.NullBool{Bool: *req.IsActive, Valid: true}
    }
    if err := updateUserAccess(targetUserID, desiredRole, active); err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return UserSummary{}, ErrUserNotFound
        }
        return UserSummary{}, err
    }
    return getUserSummaryByID(targetUserID)
}
