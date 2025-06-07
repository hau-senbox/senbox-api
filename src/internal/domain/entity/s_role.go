package entity

import (
	"github.com/pkg/errors"
	"time"
)

type Role string

const (
	Student    Role = "Student"
	Guardian   Role = "Guardian"
	Teacher    Role = "Teacher"
	Staff      Role = "Staff"
	Admin      Role = "Admin"
	Doctor     Role = "Doctor"
	Nanny      Role = "Nanny"
	Parent     Role = "Parent"
	SuperAdmin Role = "SuperAdmin"
)

func (r Role) String() string {
	return string(r)
}

func RoleFromString(role string) (Role, error) {
	switch role {
	case "Student":
		return Student, nil
	case "Guardian":
		return Guardian, nil
	case "Teacher":
		return Teacher, nil
	case "Staff":
		return Staff, nil
	case "Admin":
		return Admin, nil
	case "Doctor":
		return Doctor, nil
	case "Nanny":
		return Nanny, nil
	case "Parent":
		return Parent, nil
	case "SuperAdmin":
		return SuperAdmin, nil
	default:
		return "", errors.New("invalid role")
	}
}

type SRole struct {
	ID        int64     `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	Role      Role      `gorm:"type:varchar(16);not null;"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP;not null"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP;not null"`
}
