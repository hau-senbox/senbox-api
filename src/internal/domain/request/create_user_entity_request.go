package request

import (
	"errors"
	"strings"
	"time"
)

type CreateUserEntityRequest struct {
	Username string `json:"username" binding:"required"`
	Nickname string `json:"nickname" default:""`
	// Fullname *string `json:"fullname"`
	// Phone    *string `json:"phone"`
	Email      *string `json:"email"`
	Birthday   string  `json:"birthday" binding:"required"`
	Password   string  `json:"password" binding:"required"`
	Role       *string  `json:"role"`
	DeviceUUID string  `json:"device_uuid" binding:"required"`
}

// IsOver18 validates if the user is over 18 years old
func (r *CreateUserEntityRequest) IsOver18() error {
	// Parse the birthday string to time.Time object
	birthday, err := time.Parse("2006-01-02", r.Birthday)
	if err != nil {
		return err
	}

	// Calculate age
	now := time.Now()
	age := now.Year() - birthday.Year()
	if now.YearDay() < birthday.YearDay() {
		age--
	}

	// If age is less than 18, return an error
	if age < 18 {
		return errors.New("user must be at least 18 years old")
	}

	return nil
}

// IsUsernameValid checks if the username contains any spaces
func (r *CreateUserEntityRequest) IsUsernameValid() error {
	if strings.Contains(r.Username, " ") {
		return errors.New("username should not contain spaces")
	}
	return nil
}
