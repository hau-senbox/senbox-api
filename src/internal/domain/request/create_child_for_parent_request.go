package request

import (
	"errors"
	"strings"
)

type CreateChildForParentRequest struct {
	Username string `json:"username" binding:"required"`
	Nickname string `json:"nickname" default:""`
	Fullname string `json:"fullname" default:""`
	Birthday string `json:"birthday" binding:"required"`
}

// IsUsernameValid checks if the username contains any spaces
func (r *CreateChildForParentRequest) IsUsernameValid() error {
	if strings.Contains(r.Username, " ") {
		return errors.New("username should not contain spaces")
	}
	return nil
}
