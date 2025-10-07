package mapper

import (
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/value"
)

func MapTeacherToOwnerAssignResponse(teacher *entity.STeacherFormApplication, name string, avatarKey string, avatarUrl string, createdIndex int, userCreatedIndex int) *response.OwnerAssignResponse {
	return &response.OwnerAssignResponse{
		OwnerID:          teacher.ID.String(),
		OwnerRole:        value.OwnerRoleTeacher,
		Name:             name,
		AvatarKey:        avatarKey,
		AvatarUrl:        avatarUrl,
		CreatedIndex:     createdIndex,
		UserCreatedIndex: userCreatedIndex,
	}
}

func MapStaffToOwnerAssignResponse(staff *entity.SStaffFormApplication, name string, avatarKey string, avatarUrl string, createdIndex int, userCreatedIndex int) *response.OwnerAssignResponse {
	return &response.OwnerAssignResponse{
		OwnerID:          staff.ID.String(),
		OwnerRole:        value.OwnerRoleStaff,
		Name:             name,
		AvatarKey:        avatarKey,
		AvatarUrl:        avatarUrl,
		CreatedIndex:     createdIndex,
		UserCreatedIndex: userCreatedIndex,
	}
}

func MapStudentToOwnerAssignResponse(student *entity.SStudentFormApplication, name string, avatarKey string, avatarUrl string, createdIndex int) *response.OwnerAssignResponse {
	return &response.OwnerAssignResponse{
		OwnerID:      student.ID.String(),
		OwnerRole:    value.OwnerRoleStudent,
		Name:         name,
		AvatarKey:    avatarKey,
		AvatarUrl:    avatarUrl,
		CreatedIndex: createdIndex,
	}
}
