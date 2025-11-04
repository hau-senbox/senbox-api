package mapper

import (
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/value"
)

func MapTeacherToOwnerAssignResponse(teacher *entity.STeacherFormApplication, name string, avatarKey string, avatarUrl string, createdIndex int, userCreatedIndex int, code string) *response.OwnerAssignResponse {
	return &response.OwnerAssignResponse{
		OwnerID:          teacher.ID.String(),
		OwnerRole:        value.OwnerRoleTeacher,
		Name:             name,
		AvatarKey:        avatarKey,
		AvatarUrl:        avatarUrl,
		CreatedIndex:     createdIndex,
		UserCreatedIndex: userCreatedIndex,
		Code:             code,
		LanguageKeys:     []string{"vietnamese", "english"},
	}
}

func MapStaffToOwnerAssignResponse(staff *entity.SStaffFormApplication, name string, avatarKey string, avatarUrl string, createdIndex int, userCreatedIndex int, code string) *response.OwnerAssignResponse {
	return &response.OwnerAssignResponse{
		OwnerID:          staff.ID.String(),
		OwnerRole:        value.OwnerRoleStaff,
		Name:             name,
		AvatarKey:        avatarKey,
		AvatarUrl:        avatarUrl,
		CreatedIndex:     createdIndex,
		UserCreatedIndex: userCreatedIndex,
		Code:             code,
		LanguageKeys:     []string{"vietnamese", "english"},
	}
}

func MapStudentToOwnerAssignResponse(student *entity.SStudentFormApplication, name string, avatarKey string, avatarUrl string, createdIndex int, code string) *response.OwnerAssignResponse {
	return &response.OwnerAssignResponse{
		OwnerID:      student.ID.String(),
		OwnerRole:    value.OwnerRoleStudent,
		Name:         name,
		AvatarKey:    avatarKey,
		AvatarUrl:    avatarUrl,
		CreatedIndex: createdIndex,
		Code:         code,
		LanguageKeys: []string{"vietnamese", "english"},
	}
}

func MapParentToOwnerAssignResponse(parent *entity.SParent, name string, avatarKey string, avatarUrl string, createdIndex int, userCreatedIndex int, code string) *response.OwnerAssignResponse {
	return &response.OwnerAssignResponse{
		OwnerID:          parent.ID.String(),
		OwnerRole:        value.OwnerRoleParent,
		Name:             name,
		AvatarKey:        avatarKey,
		AvatarUrl:        avatarUrl,
		CreatedIndex:     createdIndex,
		UserCreatedIndex: userCreatedIndex,
		Code:             code,
		LanguageKeys:     []string{"vietnamese", "english"},
	}
}
