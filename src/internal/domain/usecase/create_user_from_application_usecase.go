package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/value"

	"github.com/google/uuid"
)

type CreateUserFormApplicationUseCase struct {
	*repository.UserEntityRepository
	*repository.RoleOrgSignUpRepository
	*repository.ComponentRepository
	*repository.StudentMenuRepository
}

func (receiver *CreateUserFormApplicationUseCase) CreateTeacherFormApplication(req request.CreateTeacherFormApplicationRequest) error {
	_, err := receiver.UserEntityRepository.GetByID(request.GetUserEntityByIDRequest{ID: req.UserID})
	if err != nil {
		return err
	}

	return receiver.UserEntityRepository.CreateTeacherFormApplication(req)
}

func (receiver *CreateUserFormApplicationUseCase) CreateStaffFormApplication(req request.CreateStaffFormApplicationRequest) error {
	_, err := receiver.UserEntityRepository.GetByID(request.GetUserEntityByIDRequest{ID: req.UserID})
	if err != nil {
		return err
	}

	return receiver.UserEntityRepository.CreateStaffFormApplication(req)
}

func (receiver *CreateUserFormApplicationUseCase) CreateStudentFormApplication(req request.CreateStudentFormApplicationRequest) error {
	_, err := receiver.UserEntityRepository.GetByID(request.GetUserEntityByIDRequest{ID: req.UserID})
	if err != nil {
		return err
	}

	studentID := uuid.New()

	err = receiver.UserEntityRepository.CreateStudentFormApplication(&entity.SStudentFormApplication{
		ID:             studentID,
		StudentName:    req.StudentName,
		ChildID:        uuid.MustParse(req.ChildID),
		UserID:         uuid.MustParse(req.UserID),
		OrganizationID: uuid.MustParse(req.OrganizationID),
	})

	if err == nil {
		// tao student menu
		studentRoleOrg, _ := receiver.RoleOrgSignUpRepository.GetByRoleName(string(value.RoleStudent))
		if studentRoleOrg != nil {
			components, _ := receiver.ComponentRepository.GetBySectionID(studentRoleOrg.ID.String())

			for index, component := range components {
				err := receiver.StudentMenuRepository.Create(&entity.StudentMenu{
					ID:          uuid.New(),
					StudentID:   studentID,
					ComponentID: component.ID,
					Order:       index,
					IsShow:      true,
					Visible:     true,
				})
				if err != nil {
					continue
				}
			}
		}
	}

	return err
}
