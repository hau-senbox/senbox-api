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
	UserEntityRepository *repository.UserEntityRepository
}

func NewStaffApplicationUseCase(
	staffRepo *repository.StaffApplicationRepository,
	menuRepo *repository.StaffMenuRepository,
	componentRepo *repository.ComponentRepository,
	roleOrgRepo *repository.RoleOrgSignUpRepository,
	organizationRepo *repository.OrganizationRepository,
	getUserEntityUseCase *GetUserEntityUseCase,
	userEntityResitory *repository.UserEntityRepository,
) *StaffApplicationUseCase {
	return &StaffApplicationUseCase{
		StaffAppRepo:         staffRepo,
		StaffMenuRepo:        menuRepo,
		ComponentRepo:        componentRepo,
		RoleOrgRepo:          roleOrgRepo,
		OrganizationRepo:     organizationRepo,
		GetUserEntityUseCase: getUserEntityUseCase,
		UserEntityRepository: userEntityResitory,
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

// GetAllStaff4Search returns all staff for search functionality
func (uc *StaffApplicationUseCase) GetAllStaff4Search(ctx *gin.Context) ([]response.StaffResponse, error) {
	// Lấy thông tin người dùng hiện tại (kèm Organizations, Roles)
	user, err := uc.GetUserEntityUseCase.GetCurrentUserWithOrganizations(ctx)
	if err != nil {
		return nil, err
	}

	// Nếu là SuperAdmin → trả về tất cả
	if user.IsSuperAdmin() {
		apps, err := uc.StaffAppRepo.GetAll()
		if err != nil {
			return nil, err
		}
		return mapStaffAppsToResponse(apps, uc), nil
	}

	// Nếu không phải SuperAdmin → lấy orgIDs mà user đang quản lý
	orgIDs, err := user.GetManagedOrganizationIDs(uc.StaffAppRepo.GetDB())
	if err != nil {
		return nil, err
	}
	if len(orgIDs) == 0 {
		return []response.StaffResponse{}, nil
	}

	// 4. Lấy student application theo các orgID
	apps, err := uc.StaffAppRepo.GetByOrganizationIDs(orgIDs)
	if err != nil {
		return nil, err
	}

	return mapStaffAppsToResponse(apps, uc), nil
}

func mapStaffAppsToResponse(apps []entity.SStaffFormApplication, uc *StaffApplicationUseCase) []response.StaffResponse {
	res := make([]response.StaffResponse, 0, len(apps))
	for _, a := range apps {

		// lay user theo user id cua teacher
		userEntity, _ := uc.UserEntityRepository.GetByID(request.GetUserEntityByIDRequest{
			ID: a.UserID.String(),
		})
		res = append(res, response.StaffResponse{
			StaffID:   a.ID.String(),
			StaffName: userEntity.Username,
		})
	}
	return res
}
