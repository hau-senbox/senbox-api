package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/value"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type StaffApplicationUseCase struct {
	StaffAppRepo         *repository.StaffApplicationRepository
	StaffMenuRepo        *repository.StaffMenuRepository
	ComponentRepo        *repository.ComponentRepository
	RoleOrgRepo          *repository.RoleOrgSignUpRepository
	OrganizationRepo     *repository.OrganizationRepository
	GetUserEntityUseCase *GetUserEntityUseCase
}

func NewStaffApplicationUseCase(
	staffRepo *repository.StaffApplicationRepository,
	menuRepo *repository.StaffMenuRepository,
	componentRepo *repository.ComponentRepository,
	roleOrgRepo *repository.RoleOrgSignUpRepository,
	organizationRepo *repository.OrganizationRepository,
	getUserEntityUseCase *GetUserEntityUseCase,
) *StaffApplicationUseCase {
	return &StaffApplicationUseCase{
		StaffAppRepo:         staffRepo,
		StaffMenuRepo:        menuRepo,
		ComponentRepo:        componentRepo,
		RoleOrgRepo:          roleOrgRepo,
		OrganizationRepo:     organizationRepo,
		GetUserEntityUseCase: getUserEntityUseCase,
	}
}

func (uc *StaffApplicationUseCase) GetAllStaffApplications(ctx *gin.Context) ([]response.StaffFormApplicationResponse, error) {
	// Lấy thông tin người dùng hiện tại (kèm Organizations, Roles)
	user, err := uc.GetUserEntityUseCase.GetCurrentUserWithOrganizations(ctx)
	if err != nil {
		return nil, err
	}

	var apps []entity.SStaffFormApplication

	if user.IsSuperAdmin() {
		// SuperAdmin → lấy tất cả đơn
		apps, err = uc.StaffAppRepo.GetAll()
		if err != nil {
			return nil, err
		}
	} else {
		// Nếu không phải SuperAdmin → lấy các orgIDs được quản lý
		orgIDs, err := user.GetManagedOrganizationIDs(uc.StaffAppRepo.GetDB())
		if err != nil {
			return nil, err
		}

		// Lọc các đơn theo orgID
		apps, err = uc.StaffAppRepo.GetByOrganizationIDs(orgIDs)
		if err != nil {
			return nil, err
		}
	}

	// Tạo response
	res := make([]response.StaffFormApplicationResponse, 0, len(apps))
	for _, a := range apps {
		userStaff, _ := uc.GetUserEntityUseCase.GetUserByID(request.GetUserEntityByIDRequest{
			ID: a.UserID.String(),
		})
		orgStaff, _ := uc.OrganizationRepo.GetByID(a.OrganizationID.String())

		res = append(res, response.StaffFormApplicationResponse{
			ID:               a.ID.String(),
			StaffName:        userStaff.Username,
			Status:           a.Status.String(),
			ApprovedAt:       a.ApprovedAt.Format("2006-01-02 15:04:05"),
			CreatedAt:        a.CreatedAt.Format("2006-01-02 15:04:05"),
			UserID:           a.UserID.String(),
			OrganizationID:   a.OrganizationID.String(),
			OrganizationName: orgStaff.OrganizationName,
		})
	}

	return res, nil
}

func (uc *StaffApplicationUseCase) ApproveStaffApplication(applicationID string) error {
	// Tìm bản ghi hiện tại theo ID
	application, err := uc.StaffAppRepo.GetByID(uuid.MustParse(applicationID))

	if err != nil {
		return err
	}

	// Cập nhật trạng thái thành Approved
	application.Status = value.Approved

	// Lưu lại
	return uc.StaffAppRepo.Update(application)
}

func (uc *StaffApplicationUseCase) BlockStaffApplication(applicationID string) error {
	// Tìm bản ghi hiện tại theo ID
	application, err := uc.StaffAppRepo.GetByID(uuid.MustParse(applicationID))

	if err != nil {
		return err
	}

	// Cập nhật trạng thái thành Approved
	application.Status = value.Blocked

	// Lưu lại
	return uc.StaffAppRepo.Update(application)
}
