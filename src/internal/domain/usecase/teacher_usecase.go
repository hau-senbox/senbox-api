package usecase

import (
	"errors"
	"fmt"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/value"
	"sen-global-api/pkg/consulapi/gateway"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hung-senbox/senbox-cache-service/pkg/cache/caching"
)

type TeacherApplicationUseCase struct {
	TeacherRepo              *repository.TeacherApplicationRepository
	GetUserEntityUseCase     *GetUserEntityUseCase
	UserEntityRepository     *repository.UserEntityRepository
	TeacherMenuRepo          *repository.TeacherMenuRepository
	ComponentRepo            *repository.ComponentRepository
	RoleOrgRepo              *repository.RoleOrgSignUpRepository
	OrganizationRepo         *repository.OrganizationRepository
	LanguagesConfigUsecase   *LanguagesConfigUsecase
	UserImagesUsecase        *UserImagesUsecase
	LanguageSettingRepo      *repository.LanguageSettingRepository
	ProfileGateway           gateway.ProfileGateway
	UserBlockSettingUsecase  *UserBlockSettingUsecase
	GenerateOwnerCodeUseCase GenerateOwnerCodeUseCase
	CachingMainService       caching.CachingMainService
	ValuesAppCurrentUseCase  *ValuesAppCurrentUseCase
	DepartmentGateway        gateway.DepartmentGateway
}

func NewTeacherApplicationUseCase(repo *repository.TeacherApplicationRepository) *TeacherApplicationUseCase {
	return &TeacherApplicationUseCase{TeacherRepo: repo}
}

// Create a new teacher application
func (uc *TeacherApplicationUseCase) Create(teacher *entity.STeacherFormApplication) error {
	teacher.ID = uuid.New()
	teacher.CreatedAt = time.Now()
	return uc.TeacherRepo.Create(teacher)
}

// Get by ID
func (uc *TeacherApplicationUseCase) GetByID(id uuid.UUID) (*entity.STeacherFormApplication, error) {
	return uc.TeacherRepo.GetByID(id)
}

// Get all applications
func (uc *TeacherApplicationUseCase) GetAll() ([]entity.STeacherFormApplication, error) {
	return uc.TeacherRepo.GetAll()
}

// Update application
func (uc *TeacherApplicationUseCase) Update(teacher *entity.STeacherFormApplication) error {
	return uc.TeacherRepo.Update(teacher)
}

// Delete by ID
func (uc *TeacherApplicationUseCase) Delete(id uuid.UUID) error {
	return uc.TeacherRepo.Delete(id)
}

// Get by UserID
func (uc *TeacherApplicationUseCase) GetByUserIDApproved(userID string) ([]entity.STeacherFormApplication, error) {
	return uc.TeacherRepo.GetByUserIDApproved(userID)
}

// Get by OrganizationID
func (uc *TeacherApplicationUseCase) GetByOrganizationID(orgID string) ([]entity.STeacherFormApplication, error) {
	return uc.TeacherRepo.GetByOrganizationID(orgID)
}

// Get by list of OrganizationIDs
func (uc *TeacherApplicationUseCase) GetByOrganizationIDs(orgIDs []string) ([]entity.STeacherFormApplication, error) {
	return uc.TeacherRepo.GetByOrganizationIDs(orgIDs)
}

// Check if teacher belongs to one of the given organizations
func (uc *TeacherApplicationUseCase) CheckTeacherBelongsToOrganizations(teacherID uuid.UUID, orgIDs []string) (bool, error) {
	return uc.TeacherRepo.CheckTeacherBelongsToOrganizations(uc.TeacherRepo.GetDB(), teacherID, orgIDs)
}

func (uc *TeacherApplicationUseCase) GetTeacherByID4App(ctx *gin.Context, teacherID string) (*response.TeacherResponseBase, error) {
	// Lấy user hiện tại kèm danh sách tổ chức mà họ thuộc về
	user, err := uc.GetUserEntityUseCase.GetCurrentUserWithOrganizations(ctx)
	if err != nil {
		return nil, err
	}

	orgIDs := user.GetOrganizationIDsFromPreloaded()

	// Lấy thông tin giáo viên theo ID
	teacherApp, err := uc.TeacherRepo.GetByID(uuid.MustParse(teacherID))
	if err != nil {
		return nil, err
	}
	if teacherApp == nil {
		return nil, errors.New("teacher not found")
	}

	// Kiểm tra xem giáo viên có thuộc một trong các tổ chức của user không
	teacherOrgID := teacherApp.OrganizationID.String()
	isBelong := false
	for _, orgID := range orgIDs {
		if orgID == teacherOrgID {
			isBelong = true
			break
		}
	}

	if !isBelong {
		return nil, errors.New("teacher is not under your management scope")
	}

	// lay user theo user id cua teacher
	userEntity, err := uc.UserEntityRepository.GetByID(request.GetUserEntityByIDRequest{
		ID: teacherApp.UserID.String(),
	})

	if err != nil {
		return nil, err
	}

	// Trả về response
	return &response.TeacherResponseBase{
		TeacherID:   teacherID,
		TeacherName: userEntity.Username,
	}, nil
}

// GetAllTeachers4Search returns all teachers for search functionality
func (uc *TeacherApplicationUseCase) GetAllTeachers4Search(ctx *gin.Context) ([]response.TeacherResponse, error) {
	// Lấy thông tin người dùng hiện tại (kèm Organizations, Roles)
	user, err := uc.GetUserEntityUseCase.GetCurrentUserWithOrganizations(ctx)
	if err != nil {
		return nil, err
	}

	// Nếu là SuperAdmin → trả về tất cả
	if user.IsSuperAdmin() {
		apps, err := uc.TeacherRepo.GetApprovedAll()
		if err != nil {
			return nil, err
		}
		return uc.mapTeacherAppsToResponse(ctx, apps), nil
	}

	// Nếu không phải SuperAdmin → lấy orgIDs mà user đang quản lý
	orgIDs, err := user.GetManagedOrganizationIDs(uc.TeacherRepo.GetDB())
	if err != nil {
		return nil, err
	}
	if len(orgIDs) == 0 {
		return []response.TeacherResponse{}, nil
	}

	// Lấy teacher application theo các orgID
	apps, err := uc.TeacherRepo.GetByOrganizationIDsApproved(orgIDs)
	if err != nil {
		return nil, err
	}

	return uc.mapTeacherAppsToResponse(ctx, apps), nil
}

func (uc *TeacherApplicationUseCase) mapTeacherAppsToResponse(ctx *gin.Context, teachers []entity.STeacherFormApplication) []response.TeacherResponse {
	res := make([]response.TeacherResponse, 0, len(teachers))

	for _, teacher := range teachers {
		// lay user theo user id cua teacher
		userEntity, _ := uc.UserEntityRepository.GetByID(request.GetUserEntityByIDRequest{
			ID: teacher.UserID.String(),
		})
		isDeactive, _ := uc.UserBlockSettingUsecase.GetDeactive4Teacher(teacher.ID.String())
		avatar, _ := uc.UserImagesUsecase.GetAvtIsMain4Owner(teacher.ID.String(), value.OwnerRoleStaff)
		Code, _ := uc.ProfileGateway.GetTeacherCode(ctx, teacher.ID.String())
		isLoggedDevice, _ := uc.ValuesAppCurrentUseCase.GetIsLoggedDevice4Teacher(ctx, teacher.ID.String())
		res = append(res, response.TeacherResponse{
			TeacherID:        teacher.ID.String(),
			TeacherName:      userEntity.Nickname,
			CreatedIndex:     teacher.CreatedIndex,
			UserCreatedIndex: userEntity.CreatedIndex,
			IsDeactive:       isDeactive,
			Avatar:           avatar,
			LanguageKeys:     []string{"vietnamese", "english"},
			Code:             Code,
			IsLoggedDevice:   isLoggedDevice,
		})
	}
	return res
}

func (uc *TeacherApplicationUseCase) GetTeacherByID4Web(ctx *gin.Context, teacherID string) (*response.TeacherResponseBase, error) {
	// Lấy thông tin application của giáo viên
	teacherApp, err := uc.TeacherRepo.GetByID(uuid.MustParse(teacherID))
	if err != nil {
		return nil, err
	}
	if teacherApp == nil {
		return nil, errors.New("teacher not found")
	}

	// Lấy danh sách menu của giáo viên
	teacherMenus, err := uc.TeacherMenuRepo.GetByTeacherID(teacherID)
	if err != nil {
		return nil, err
	}

	// Tạo danh sách componentID để lấy Component
	componentIDs := make([]uuid.UUID, 0, len(teacherMenus))
	componentOrderMap := make(map[uuid.UUID]int)

	for _, cm := range teacherMenus {
		componentIDs = append(componentIDs, cm.ComponentID)
		componentOrderMap[cm.ComponentID] = cm.Order
	}

	// Lấy components theo ID
	components, err := uc.ComponentRepo.GetByIDs(componentIDs)
	if err != nil {
		return nil, err
	}

	// Gom components theo language_id
	menusByLang := make(map[uint][]response.ComponentResponse)
	langMap := make(map[uint]entity.LanguageSetting)

	for _, comp := range components {
		menu := response.ComponentResponse{
			ID:         comp.ID.String(),
			Name:       comp.Name,
			Type:       comp.Type.String(),
			Key:        comp.Key,
			Value:      string(comp.Value),
			Order:      componentOrderMap[comp.ID],
			LanguageID: comp.LanguageID,
		}

		menusByLang[comp.LanguageID] = append(menusByLang[comp.LanguageID], menu)

		// nếu chưa có languageID trong cache -> query DB
		if _, ok := langMap[comp.LanguageID]; !ok {
			langSetting, err := uc.LanguageSettingRepo.GetByID(comp.LanguageID)
			if err != nil {
				return nil, err
			}
			if langSetting != nil {
				langMap[comp.LanguageID] = *langSetting
			}
		}
	}

	// Build []GetMenus4Web
	getMenus := make([]response.GetMenus4Web, 0, len(menusByLang))
	for langID, comps := range menusByLang {
		getMenus = append(getMenus, response.GetMenus4Web{
			Language: langMap[langID],
			Menus:    comps,
		})
	}

	// Tạo QR cho form profile
	teacherRoleOrg, err := uc.RoleOrgRepo.GetByRoleName(string(value.RoleTeacher))
	if err != nil {
		return nil, err
	}
	formProfile := teacherRoleOrg.OrgProfile + ":" + teacherApp.ID.String()
	userEntity, _ := uc.UserEntityRepository.GetByID(request.GetUserEntityByIDRequest{
		ID: teacherApp.UserID.String(),
	})

	// get languages config
	languageConfig, _ := uc.LanguagesConfigUsecase.GetLanguagesConfigByOwnerNoCtx(teacherID, value.OwnerRoleLangTeacher)

	// get avts
	avatars, _ := uc.UserImagesUsecase.GetAvt4Owner(teacherID, value.OwnerRoleTeacher)

	// get list loged device
	logedDevices, _ := uc.ValuesAppCurrentUseCase.GetLogedDevices4Teacher(ctx, teacherID)

	return &response.TeacherResponseBase{
		TeacherID:      teacherID,
		UserID:         userEntity.ID.String(),
		TeacherName:    userEntity.Nickname,
		Avatar:         "",
		AvatarURL:      "",
		QrFormProfile:  formProfile,
		Menus:          getMenus,
		IsUserBlock:    userEntity.IsBlocked,
		LanguageConfig: languageConfig,
		Avatars:        avatars,
		CreatedIndex:   teacherApp.CreatedIndex,
		LogedDevices:   logedDevices,
	}, nil
}

func (uc *TeacherApplicationUseCase) ApproveTeacherApplication(ctx *gin.Context, applicationID string) error {
	// Tìm bản ghi hiện tại theo ID
	application, err := uc.TeacherRepo.GetByID(uuid.MustParse(applicationID))

	if err != nil {
		return err
	}

	if application == nil {
		return errors.New("application not found")
	}

	// Lấy thông tin người dùng hiện tại
	user, err := uc.GetUserEntityUseCase.GetCurrentUserWithOrganizations(ctx)
	if err != nil {
		return err
	}

	// Nếu application bị block bởi admin → chỉ SuperAdmin mới có quyền duyệt
	if application.IsAdminBlock && !user.IsSuperAdmin() {
		return fmt.Errorf("only SuperAdmin can approve an admin-blocked application")
	}

	// Cập nhật trạng thái thành Approved
	application.Status = value.Approved
	application.ApprovedAt = time.Now()
	application.IsAdminBlock = false // Reset block status when approving

	// assign teacher department group
	uc.DepartmentGateway.AssignTeacherDepartmentGroup(ctx, application.ID.String(), application.OrganizationID.String())

	// Lưu lại
	return uc.TeacherRepo.Update(application)
}

func (uc *TeacherApplicationUseCase) BlockTeacherApplication(ctx *gin.Context, applicationID string) error {
	// Tìm bản ghi hiện tại theo ID
	application, err := uc.TeacherRepo.GetByID(uuid.MustParse(applicationID))

	if err != nil {
		return err
	}

	if application == nil {
		return errors.New("application not found")
	}

	// Lấy thông tin người dùng hiện tại (kèm Organizations, Roles)
	user, err := uc.GetUserEntityUseCase.GetCurrentUserWithOrganizations(ctx)
	if err != nil {
		return err
	}

	// Nếu là SuperAdmin
	if user.IsSuperAdmin() {
		application.IsAdminBlock = true
	}

	// Cập nhật trạng thái thành Approved
	application.Status = value.Blocked

	// Lưu lại
	return uc.TeacherRepo.Update(application)
}

func (uc *TeacherApplicationUseCase) GetAllTeacherApplications(ctx *gin.Context) ([]response.TeacherFormApplicationResponse, error) {
	// Lấy thông tin người dùng hiện tại (kèm Organizations, Roles)
	user, err := uc.GetUserEntityUseCase.GetCurrentUserWithOrganizations(ctx)
	if err != nil {
		return nil, err
	}

	var apps []entity.STeacherFormApplication

	if user.IsSuperAdmin() {
		// SuperAdmin → lấy tất cả đơn
		apps, err = uc.TeacherRepo.GetAll()
		if err != nil {
			return nil, err
		}
	} else {
		// Nếu không phải SuperAdmin → lấy các orgIDs được quản lý
		orgIDs, err := user.GetManagedOrganizationIDs(uc.TeacherRepo.GetDB())
		if err != nil {
			return nil, err
		}

		// Lọc các đơn theo orgID
		apps, err = uc.TeacherRepo.GetByOrganizationIDs(orgIDs)
		if err != nil {
			return nil, err
		}
	}

	// Tạo response
	res := make([]response.TeacherFormApplicationResponse, 0, len(apps))
	for _, a := range apps {
		userEntity, _ := uc.UserEntityRepository.GetByID(request.GetUserEntityByIDRequest{
			ID: a.UserID.String(),
		})
		orgTeacher, _ := uc.OrganizationRepo.GetByID(a.OrganizationID.String())

		res = append(res, response.TeacherFormApplicationResponse{
			ID:               a.ID.String(),
			TeacherName:      userEntity.Username,
			Status:           a.Status.String(),
			ApprovedAt:       a.ApprovedAt.Format("2006-01-02 15:04:05"),
			CreatedAt:        a.CreatedAt.Format("2006-01-02 15:04:05"),
			UserID:           a.UserID.String(),
			OrganizationID:   a.OrganizationID.String(),
			OrganizationName: orgTeacher.OrganizationName,
		})
	}

	return res, nil
}

func (uc *TeacherApplicationUseCase) GetDetailTeacherApplication(ctx *gin.Context, applicationID string) (*response.TeacherFormApplicationResponse, error) {
	// Lấy thông tin application của giáo viên
	application, err := uc.TeacherRepo.GetByID(uuid.MustParse(applicationID))
	if err != nil {
		return nil, err
	}
	if application == nil {
		return nil, errors.New("teacher application not found")
	}

	userEntity, _ := uc.UserEntityRepository.GetByID(request.GetUserEntityByIDRequest{
		ID: application.UserID.String(),
	})
	orgTeacher, _ := uc.OrganizationRepo.GetByID(application.OrganizationID.String())
	return &response.TeacherFormApplicationResponse{
		ID:               application.ID.String(),
		TeacherName:      userEntity.Username,
		Status:           application.Status.String(),
		ApprovedAt:       application.ApprovedAt.Format("2006-01-02 15:04:05"),
		CreatedAt:        application.CreatedAt.Format("2006-01-02 15:04:05"),
		UserID:           application.UserID.String(),
		OrganizationID:   application.OrganizationID.String(),
		OrganizationName: orgTeacher.OrganizationName,
	}, nil
}

func (uc *TeacherApplicationUseCase) GetTeacher4Gateway(ctx *gin.Context, teacherID string) (*response.GetTeacher4Gateway, error) {
	teacher, err := uc.TeacherRepo.GetByID(uuid.MustParse(teacherID))
	if err != nil {
		return nil, err
	}
	if teacher == nil {
		return nil, errors.New("teacher not found")
	}

	userEntity, _ := uc.UserEntityRepository.GetByID(request.GetUserEntityByIDRequest{
		ID: teacher.UserID.String(),
	})

	// get avts
	avatar, _ := uc.UserImagesUsecase.GetAvtIsMain4Owner(teacherID, value.OwnerRoleTeacher)
	code, _ := uc.ProfileGateway.GetTeacherCode(ctx, teacherID)

	res := &response.GetTeacher4Gateway{
		TeacherID:      teacherID,
		OrganizationID: teacher.OrganizationID.String(),
		TeacherName:    userEntity.Nickname,
		Avatar:         avatar,
		Code:           code,
	}

	uc.CachingMainService.SetTeacherCache(ctx.Request.Context(), teacherID, res)
	return res, nil
}

func (uc *TeacherApplicationUseCase) GetTeacherByUser4Gateway(userID string) (*response.GetTeacher4Gateway, error) {
	teacher, err := uc.TeacherRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	userEntity, _ := uc.UserEntityRepository.GetByID(request.GetUserEntityByIDRequest{
		ID: userID,
	})

	// get avts
	avatar, _ := uc.UserImagesUsecase.GetAvtIsMain4Owner(teacher.ID.String(), value.OwnerRoleTeacher)

	return &response.GetTeacher4Gateway{
		TeacherID:      teacher.ID.String(),
		OrganizationID: teacher.OrganizationID.String(),
		TeacherName:    userEntity.Username,
		Avatar:         avatar,
	}, nil
}

func (uc *TeacherApplicationUseCase) GetTeachersByUser4Gateway(userID string) ([]*response.GetTeacher4Gateway, error) {
	// Lấy tất cả teacher application đã được duyệt theo userID
	teachers, err := uc.TeacherRepo.GetAllByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Lấy thông tin userEntity
	userEntity, _ := uc.UserEntityRepository.GetByID(request.GetUserEntityByIDRequest{
		ID: userID,
	})

	var result []*response.GetTeacher4Gateway
	for _, teacher := range teachers {
		// Lấy avatar cho từng teacher
		avatar, _ := uc.UserImagesUsecase.GetAvtIsMain4Owner(teacher.ID.String(), value.OwnerRoleTeacher)

		result = append(result, &response.GetTeacher4Gateway{
			TeacherID:      teacher.ID.String(),
			OrganizationID: teacher.OrganizationID.String(),
			TeacherName:    userEntity.Username,
			Avatar:         avatar,
		})
	}

	return result, nil
}

func (uc *TeacherApplicationUseCase) GetTeacherByOrgAndUser4Gateway(ctx *gin.Context, userID string, organizationID string) (*response.GetTeacher4Gateway, error) {
	teacher, err := uc.TeacherRepo.GetByUserIDAndOrgID(userID, organizationID)
	if err != nil {
		return nil, err
	}

	if teacher == nil {
		return nil, errors.New("teacher not found")
	}

	userEntity, _ := uc.UserEntityRepository.GetByID(request.GetUserEntityByIDRequest{
		ID: userID,
	})

	// get avts
	avatar, _ := uc.UserImagesUsecase.GetAvtIsMain4Owner(teacher.ID.String(), value.OwnerRoleTeacher)
	code, _ := uc.ProfileGateway.GetTeacherCode(ctx, teacher.ID.String())

	res := &response.GetTeacher4Gateway{
		TeacherID:      teacher.ID.String(),
		OrganizationID: teacher.OrganizationID.String(),
		TeacherName:    userEntity.Username,
		Avatar:         avatar,
		Code:           code,
	}

	uc.CachingMainService.SetTeacherByUserAndOrgCacheKey(ctx.Request.Context(), userID, organizationID, res)

	return res, nil
}

func (uc *TeacherApplicationUseCase) GetAllTeacherByOrg4App(ctx *gin.Context, organizationID string) ([]response.TeacherResponse, error) {

	teachers, err := uc.TeacherRepo.GetByOrganizationID(organizationID)
	if err != nil {
		return nil, err
	}

	return uc.mapTeacherAppsToResponse(ctx, teachers), nil
}

func (uc *TeacherApplicationUseCase) GenerateTeacherCode(ctx *gin.Context) {
	// get all teachers
	teachers, err := uc.TeacherRepo.GetAll()
	if err != nil {
		return
	}

	for _, tr := range teachers {
		// call profile gateway to generate teacher code
		_, _ = uc.ProfileGateway.GenerateTeacherCode(ctx, tr.ID.String(), tr.CreatedIndex)
	}
}

func (uc *TeacherApplicationUseCase) MigrateTeacherDepartmentGroup(ctx *gin.Context, organizationID string) {
	// get all teachers
	teachers, err := uc.TeacherRepo.GetByOrganizationID(organizationID)
	if err != nil {
		return
	}

	for _, tr := range teachers {
		uc.DepartmentGateway.AssignTeacherDepartmentGroup(ctx, tr.ID.String(), organizationID)
	}
}
