package usecase

import (
	"errors"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/value"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type StudentApplicationUseCase struct {
	StudentAppRepo       *repository.StudentApplicationRepository
	StudentMenuRepo      *repository.StudentMenuRepository
	ComponentRepo        *repository.ComponentRepository
	RoleOrgRepo          *repository.RoleOrgSignUpRepository
	GetUserEntityUseCase *GetUserEntityUseCase
}

func NewStudentApplicationUseCase(
	studentRepo *repository.StudentApplicationRepository,
	menuRepo *repository.StudentMenuRepository,
	componentRepo *repository.ComponentRepository,
	roleOrgRepo *repository.RoleOrgSignUpRepository,
	getUserEntityUseCase *GetUserEntityUseCase,
) *StudentApplicationUseCase {
	return &StudentApplicationUseCase{
		StudentAppRepo:       studentRepo,
		StudentMenuRepo:      menuRepo,
		ComponentRepo:        componentRepo,
		RoleOrgRepo:          roleOrgRepo,
		GetUserEntityUseCase: getUserEntityUseCase,
	}
}

// Get all students
func (uc *StudentApplicationUseCase) GetAllStudents() ([]response.StudentResponse, error) {
	apps, err := uc.StudentAppRepo.GetAll()
	if err != nil {
		return nil, err
	}

	res := make([]response.StudentResponse, 0, len(apps))
	for _, a := range apps {
		res = append(res, response.StudentResponse{
			StudentID:   a.ID.String(),
			StudentName: a.StudentName,
		})
	}
	return res, nil
}

func (uc *StudentApplicationUseCase) GetStudentByID(studentID string) (*response.StudentResponseBase, error) {
	studentApp, err := uc.StudentAppRepo.GetByID(uuid.MustParse(studentID))
	if err != nil {
		return nil, err
	}
	if studentApp == nil {
		return nil, errors.New("student not found")
	}

	// Lấy danh sách ChildMenu
	studentMenus, err := uc.StudentMenuRepo.GetByStudentID(studentID)
	if err != nil {
		return nil, err
	}

	// Tạo danh sách componentID để lấy Component
	componentIDs := make([]uuid.UUID, 0, len(studentMenus))
	componentOrderMap := make(map[uuid.UUID]int)
	componentIsShowMap := make(map[uuid.UUID]bool)

	for _, cm := range studentMenus {
		componentIDs = append(componentIDs, cm.ComponentID)
		componentOrderMap[cm.ComponentID] = cm.Order
		componentIsShowMap[cm.ComponentID] = cm.IsShow
	}

	// Lấy tất cả components theo danh sách ID
	components, err := uc.ComponentRepo.GetByIDs(componentIDs)
	if err != nil {
		return nil, err
	}

	// Build danh sách ComponentChildResponse
	menus := make([]response.ComponentResponse, 0)
	for _, comp := range components {
		menu := response.ComponentResponse{
			ID:     comp.ID.String(),
			Name:   comp.Name,
			Type:   comp.Type.String(),
			Key:    comp.Key,
			Value:  string(comp.Value),
			Order:  componentOrderMap[comp.ID],
			IsShow: componentIsShowMap[comp.ID],
		}
		menus = append(menus, menu)
	}
	// lay qr profile form
	studentRoleOrg, err := uc.RoleOrgRepo.GetByRoleName(string(value.RoleStudent))
	if err != nil {
		return nil, err
	}
	formProfile := studentRoleOrg.OrgProfile + ":" + studentApp.ID.String()

	return &response.StudentResponseBase{
		StudentID:     studentID,
		StudentName:   studentApp.StudentName,
		Avatar:        "",
		AvatarURL:     "",
		QrFormProfile: formProfile,
		Menus:         menus,
	}, nil
}

func (uc *StudentApplicationUseCase) GetStudentByID4App(ctx *gin.Context, studentID string) (*response.StudentResponseBase, error) {
	user, err := uc.GetUserEntityUseCase.GetCurrentUserWithOrganizations(ctx)
	if err != nil {
		return nil, err
	}

	orgIDs := user.GetOrganizationIDsFromPreloaded()

	studentApp, err := uc.StudentAppRepo.GetByID(uuid.MustParse(studentID))
	if err != nil {
		return nil, err
	}
	if studentApp == nil {
		return nil, errors.New("student not found")
	}

	// Kiểm tra student có thuộc 1 trong các tổ chức mà user quản lý không
	studentOrgID := studentApp.OrganizationID.String()
	isBelong := false
	for _, orgID := range orgIDs {
		if orgID == studentOrgID {
			isBelong = true
			break
		}
	}

	if !isBelong {
		return nil, errors.New("student is not under your management scope")
	}

	return &response.StudentResponseBase{
		StudentID:   studentID,
		StudentName: studentApp.StudentName,
	}, nil
}

// usecase/student_application_usecase.go
func (uc *StudentApplicationUseCase) UpdateStudentName(req request.UpdateStudentRequest) error {
	// Tìm bản ghi hiện tại theo ID
	student := &entity.SStudentFormApplication{}
	err := uc.StudentAppRepo.DB.
		Where("id = ?", req.StudentID).
		First(student).Error
	if err != nil {
		return err
	}

	// Cập nhật tên
	student.StudentName = req.StudentName

	// Lưu lại
	return uc.StudentAppRepo.Update(student)
}

// GetAllStudents4Search returns all students for search functionality
func (uc *StudentApplicationUseCase) GetAllStudents4Search(ctx *gin.Context) ([]response.StudentResponse, error) {
	// Lấy thông tin người dùng hiện tại (kèm Organizations, Roles)
	user, err := uc.GetUserEntityUseCase.GetCurrentUserWithOrganizations(ctx)
	if err != nil {
		return nil, err
	}

	// Nếu là SuperAdmin → trả về tất cả
	if user.IsSuperAdmin() {
		apps, err := uc.StudentAppRepo.GetAll()
		if err != nil {
			return nil, err
		}
		return mapStudentAppsToResponse(apps), nil
	}

	// Nếu không phải SuperAdmin → lấy orgIDs mà user đang quản lý
	orgIDs, err := user.GetManagedOrganizationIDs(uc.StudentAppRepo.GetDB())
	if err != nil {
		return nil, err
	}
	if len(orgIDs) == 0 {
		return []response.StudentResponse{}, nil
	}

	// 4. Lấy student application theo các orgID
	apps, err := uc.StudentAppRepo.GetByOrganizationIDs(orgIDs)
	if err != nil {
		return nil, err
	}

	return mapStudentAppsToResponse(apps), nil
}

func mapStudentAppsToResponse(apps []entity.SStudentFormApplication) []response.StudentResponse {
	res := make([]response.StudentResponse, 0, len(apps))
	for _, a := range apps {
		res = append(res, response.StudentResponse{
			StudentID:   a.ID.String(),
			StudentName: a.StudentName,
		})
	}
	return res
}
