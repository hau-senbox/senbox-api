package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/mapper"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/value"

	"github.com/gin-gonic/gin"
)

type OwnerAssignUseCase struct {
	TeacherRepo       *repository.TeacherApplicationRepository
	StaffRepo         *repository.StaffApplicationRepository
	UserEntityRepo    *repository.UserEntityRepository
	UserImagesUsecase *UserImagesUsecase
}

func (uc *OwnerAssignUseCase) GetListOwner2Assign(
	ctx *gin.Context,
	organizationID string,
) (*response.ListOwnerAssignResponse, error) {

	// get list teachers by org
	teachers, err := uc.TeacherRepo.GetByOrganizationID(organizationID)
	if err != nil {
		return nil, err
	}

	// get list staffs by org
	staffs, err := uc.StaffRepo.GetByOrganizationID(organizationID)
	if err != nil {
		return nil, err
	}

	listResp := &response.ListOwnerAssignResponse{
		Teachers: []*response.OwnerAssignResponse{},
		Staffs:   []*response.OwnerAssignResponse{},
	}

	// loop teachers
	for _, t := range teachers {
		// get name
		user, _ := uc.UserEntityRepo.GetByID(request.GetUserEntityByIDRequest{
			ID: t.UserID.String(),
		})
		// get avatar key & url
		avatar, _ := uc.UserImagesUsecase.GetAvtIsMain4Owner(t.UserID.String(), value.OwnerRoleTeacher)
		listResp.Teachers = append(listResp.Teachers,
			mapper.MapTeacherToOwnerAssignResponse(&t, user.Nickname, avatar.ImageKey, avatar.ImageUrl),
		)
	}

	// loop staffs
	for _, s := range staffs {
		// get name
		user, _ := uc.UserEntityRepo.GetByID(request.GetUserEntityByIDRequest{
			ID: s.UserID.String(),
		})
		// get avatar key & url
		avatar, _ := uc.UserImagesUsecase.GetAvtIsMain4Owner(s.UserID.String(), value.OwnerRoleStaff)
		listResp.Staffs = append(listResp.Staffs,
			mapper.MapStaffToOwnerAssignResponse(&s, user.Nickname, avatar.ImageKey, avatar.ImageUrl),
		)
	}

	return listResp, nil
}
