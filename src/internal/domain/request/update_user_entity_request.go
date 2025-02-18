package request

import "regexp"

type UpdateUserEntityRequest struct {
	ID         string    `json:"id" binding:"required"`
	Username   string    `json:"username" binding:"required"`
	Fullname   *string   `json:"fullname"`
	Phone      *string   `json:"phone"`
	Email      *string   `json:"email"`
	UserConfig *uint     `json:"user_config"`
	Guardians  *[]string `json:"guardians"`
	Roles      *[]string `json:"roles"`
	Policies   *[]uint   `json:"policies"`
	Devices    *[]string `json:"devices"`
}

// Email pattern to validate email format
func (req *UpdateUserEntityRequest) ValidateEmail() bool {
	// Regular expression for validating email format
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(*req.Email)
}

// Phone pattern to validate phone format (simple example for international numbers)
func (req *UpdateUserEntityRequest) ValidatePhone() bool {
	// Regular expression to check phone number format (e.g., +1234567890)
	re := regexp.MustCompile(`^\+?[0-9]{10,15}$`)
	return re.MatchString(*req.Phone)
}
