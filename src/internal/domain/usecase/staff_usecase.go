package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
)

type StaffApplicationUseCase struct {
	StaffAppRepo         *repository.StaffApplicationRepository
	StaffMenuRepo        *repository.StaffMenuRepository
	ComponentRepo        *repository.ComponentRepository
	RoleOrgRepo          *repository.RoleOrgSignUpRepository
	GetUserEntityUseCase *GetUserEntityUseCase
}

func NewStaffApplicationUseCase(
	staffRepo *repository.StaffApplicationRepository,
	menuRepo *repository.StaffMenuRepository,
	componentRepo *repository.ComponentRepository,
	roleOrgRepo *repository.RoleOrgSignUpRepository,
	getUserEntityUseCase *GetUserEntityUseCase,
) *StaffApplicationUseCase {
	return &StaffApplicationUseCase{
		StaffAppRepo:         staffRepo,
		StaffMenuRepo:        menuRepo,
		ComponentRepo:        componentRepo,
		RoleOrgRepo:          roleOrgRepo,
		GetUserEntityUseCase: getUserEntityUseCase,
	}
}

// GetAllStaffApplications retrieves all staff applications
func (uc *StaffApplicationUseCase) GetAllStaffApplications() ([]response.StaffFormApplicationResponse, error) {
	apps, err := uc.StaffAppRepo.GetAll()
	if err != nil {
		return nil, err
	}

	res := make([]response.StaffFormApplicationResponse, 0, len(apps))
	for _, a := range apps {
		// get user details for the staff application
		userStaff, _ := uc.GetUserEntityUseCase.GetUserByID(request.GetUserEntityByIDRequest{
			ID: a.UserID.String(),
		})

		res = append(res, response.StaffFormApplicationResponse{
			ID:         a.ID.String(),
			StaffName:  userStaff.Username,
			Status:     a.Status.String(),
			ApprovedAt: a.ApprovedAt.Format("2006-01-02 15:04:05"),
			CreatedAt:  a.CreatedAt.Format("2006-01-02 15:04:05"),
			UserID:     a.UserID.String(),
		})
	}
	return res, nil
}
