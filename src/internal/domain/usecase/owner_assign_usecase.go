package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/mapper"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/value"
	"sen-global-api/pkg/consulapi/gateway"

	"github.com/gin-gonic/gin"
)

type OwnerAssignUseCase struct {
	TeacherRepo       *repository.TeacherApplicationRepository
	StaffRepo         *repository.StaffApplicationRepository
	StudentRepo       *repository.StudentApplicationRepository
	UserEntityRepo    *repository.UserEntityRepository
	UserImagesUsecase *UserImagesUsecase
	ProfileGateway    gateway.ProfileGateway
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

	// get list students by org
	students, err := uc.StudentRepo.GetByOrganizationID(organizationID)
	if err != nil {
		return nil, err
	}

	listResp := &response.ListOwnerAssignResponse{
		Teachers: []*response.OwnerAssignResponse{},
		Staffs:   []*response.OwnerAssignResponse{},
		Students: []*response.OwnerAssignResponse{},
	}

	// loop teachers
	for _, t := range teachers {
		// get name
		user, _ := uc.UserEntityRepo.GetByID(request.GetUserEntityByIDRequest{
			ID: t.UserID.String(),
		})
		// get avatar key & url
		avatar, _ := uc.UserImagesUsecase.GetAvtIsMain4Owner(t.ID.String(), value.OwnerRoleTeacher)
		code, _ := uc.ProfileGateway.GetTeacherCode(ctx, t.ID.String())

		// get user created index
		listResp.Teachers = append(listResp.Teachers,
			mapper.MapTeacherToOwnerAssignResponse(&t, user.Nickname, avatar.ImageKey, avatar.ImageUrl, t.CreatedIndex, user.CreatedIndex, code),
		)
	}

	// loop staffs
	for _, s := range staffs {
		// get name
		user, _ := uc.UserEntityRepo.GetByID(request.GetUserEntityByIDRequest{
			ID: s.UserID.String(),
		})
		// get avatar key & url
		avatar, _ := uc.UserImagesUsecase.GetAvtIsMain4Owner(s.ID.String(), value.OwnerRoleStaff)
		code, _ := uc.ProfileGateway.GetTeacherCode(ctx, s.ID.String())

		listResp.Staffs = append(listResp.Staffs,
			mapper.MapStaffToOwnerAssignResponse(&s, user.Nickname, avatar.ImageKey, avatar.ImageUrl, s.CreatedIndex, user.CreatedIndex, code),
		)
	}

	// loop students
	for _, s := range students {
		// get avatar key & url
		avatar, _ := uc.UserImagesUsecase.GetAvtIsMain4Owner(s.ID.String(), value.OwnerRoleStudent)
		code, _ := uc.ProfileGateway.GetTeacherCode(ctx, s.ID.String())
		listResp.Students = append(listResp.Students,
			mapper.MapStudentToOwnerAssignResponse(&s, s.StudentName, avatar.ImageKey, avatar.ImageUrl, s.CreatedIndex, code),
		)
	}

	return listResp, nil
}
