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

type StudentApplicationUseCase struct {
	StudentAppRepo                *repository.StudentApplicationRepository
	StudentMenuRepo               *repository.StudentMenuRepository
	ComponentRepo                 *repository.ComponentRepository
	RoleOrgRepo                   *repository.RoleOrgSignUpRepository
	GetUserEntityUseCase          *GetUserEntityUseCase
	OrganizationRepo              *repository.OrganizationRepository
	DeviceRepo                    *repository.DeviceRepository
	StudentBlockSettingUsecase    *StudentBlockSettingUsecase
	LanguagesConfigUsecase        *LanguagesConfigUsecase
	UserImagesUsecase             *UserImagesUsecase
	LanguageSettingRepo           *repository.LanguageSettingRepository
	ProfileGateway                gateway.ProfileGateway
	StudentBlockSettingRepository *repository.StudentBlockSettingRepository
	GenerateOwnerCodeUseCase      GenerateOwnerCodeUseCase
	CachingMainService            caching.CachingMainService
	ValuesAppCurrentUseCase       *ValuesAppCurrentUseCase
}

func NewStudentApplicationUseCase(
	studentRepo *repository.StudentApplicationRepository,
	menuRepo *repository.StudentMenuRepository,
	componentRepo *repository.ComponentRepository,
	roleOrgRepo *repository.RoleOrgSignUpRepository,
	getUserEntityUseCase *GetUserEntityUseCase,
	organizationRepo *repository.OrganizationRepository,
	studentBlockSettingUsecase *StudentBlockSettingUsecase,
	languagesConfigUsecase *LanguagesConfigUsecase,
	userImagesUsecase *UserImagesUsecase,
	languageSettingRepo *repository.LanguageSettingRepository,
	studentBlockRepo *repository.StudentBlockSettingRepository,
	profileGw gateway.ProfileGateway,
	generateOwnerCodeUseCase GenerateOwnerCodeUseCase,
	cachingMainService caching.CachingMainService,
	valuesAppCurrentUseCase *ValuesAppCurrentUseCase,
) *StudentApplicationUseCase {
	return &StudentApplicationUseCase{
		StudentAppRepo:                studentRepo,
		StudentMenuRepo:               menuRepo,
		ComponentRepo:                 componentRepo,
		RoleOrgRepo:                   roleOrgRepo,
		GetUserEntityUseCase:          getUserEntityUseCase,
		OrganizationRepo:              organizationRepo,
		StudentBlockSettingUsecase:    studentBlockSettingUsecase,
		LanguagesConfigUsecase:        languagesConfigUsecase,
		UserImagesUsecase:             userImagesUsecase,
		LanguageSettingRepo:           languageSettingRepo,
		StudentBlockSettingRepository: studentBlockRepo,
		ProfileGateway:                profileGw,
		GenerateOwnerCodeUseCase:      generateOwnerCodeUseCase,
		CachingMainService:            cachingMainService,
		ValuesAppCurrentUseCase:       valuesAppCurrentUseCase,
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

// func (uc *StudentApplicationUseCase) GetStudentByID(studentID string) (*response.StudentResponseBase, error) {
// 	studentApp, err := uc.StudentAppRepo.GetByID(uuid.MustParse(studentID))
// 	if err != nil {
// 		return nil, err
// 	}
// 	if studentApp == nil {
// 		return nil, errors.New("student not found")
// 	}

// 	// Lấy danh sách ChildMenu
// 	studentMenus, err := uc.StudentMenuRepo.GetByStudentID(studentID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Tạo danh sách componentID để lấy Component
// 	componentIDs := make([]uuid.UUID, 0, len(studentMenus))
// 	componentOrderMap := make(map[uuid.UUID]int)
// 	componentIsShowMap := make(map[uuid.UUID]bool)

// 	for _, cm := range studentMenus {
// 		componentIDs = append(componentIDs, cm.ComponentID)
// 		componentOrderMap[cm.ComponentID] = cm.Order
// 		componentIsShowMap[cm.ComponentID] = cm.IsShow
// 	}

// 	// Lấy tất cả components theo danh sách ID
// 	components, err := uc.ComponentRepo.GetByIDs(componentIDs)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Build danh sách ComponentChildResponse
// 	menus := make([]response.ComponentResponse, 0)
// 	for _, comp := range components {
// 		menu := response.ComponentResponse{
// 			ID:       comp.ID.String(),
// 			Name:     comp.Name,
// 			Type:     comp.Type.String(),
// 			Key:      comp.Key,
// 			Value:    string(comp.Value),
// 			Order:    componentOrderMap[comp.ID],
// 			IsShow:   componentIsShowMap[comp.ID],
// 			Language: comp.Language,
// 		}
// 		menus = append(menus, menu)
// 	}
// 	// lay qr profile form
// 	studentRoleOrg, err := uc.RoleOrgRepo.GetByRoleName(string(value.RoleStudent))
// 	if err != nil {
// 		return nil, err
// 	}
// 	formProfile := studentRoleOrg.OrgProfile + ":" + studentApp.ID.String()

// 	// get student block setting
// 	studentBlockSetting, _ := uc.StudentBlockSettingUsecase.GetByStudentID(studentID)

// 	// get languages config
// 	languageConfig, _ := uc.LanguagesConfigUsecase.GetLanguagesConfigByOwnerNoCtx(studentID, value.OwnerRoleLangStudent)

// 	// get avts
// 	avatars, _ := uc.UserImagesUsecase.GetAvt4Owner(studentID, value.OwnerRoleStudent)

// 	return &response.StudentResponseBase{
// 		StudentID:      studentID,
// 		StudentName:    studentApp.StudentName,
// 		Avatar:         "",
// 		AvatarURL:      "",
// 		QrFormProfile:  formProfile,
// 		Menus:          menus,
// 		CustomID:       studentApp.CustomID,
// 		StudentBlock:   studentBlockSetting,
// 		LanguageConfig: languageConfig,
// 		Avatars:        avatars,
// 		CreatedIndex:   studentApp.CreatedIndex,
// 	}, nil
// }

func (uc *StudentApplicationUseCase) GetByID4WebAdmin(ctx *gin.Context, studentID string) (*response.StudentResponseBase, error) {
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
			IsShow:     componentIsShowMap[comp.ID],
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

	// lay qr profile form
	studentRoleOrg, err := uc.RoleOrgRepo.GetByRoleName(string(value.RoleStudent))
	if err != nil {
		return nil, err
	}
	formProfile := studentRoleOrg.OrgProfile + ":" + studentApp.ID.String()

	// get student block setting
	studentBlockSetting, _ := uc.StudentBlockSettingUsecase.GetByStudentID(studentID)

	// get languages config
	languageConfig, _ := uc.LanguagesConfigUsecase.GetLanguagesConfigByOwnerNoCtx(studentID, value.OwnerRoleLangStudent)

	// get avts
	avatars, _ := uc.UserImagesUsecase.GetAvt4Owner(studentID, value.OwnerRoleStudent)

	// get list loged device
	logedDevices, _ := uc.ValuesAppCurrentUseCase.GetLogedDevices4Student(ctx, studentID)

	return &response.StudentResponseBase{
		StudentID:      studentID,
		StudentName:    studentApp.StudentName,
		Avatar:         "",
		AvatarURL:      "",
		QrFormProfile:  formProfile,
		Menus:          getMenus,
		CustomID:       studentApp.CustomID,
		StudentBlock:   studentBlockSetting,
		LanguageConfig: languageConfig,
		Avatars:        avatars,
		CreatedIndex:   studentApp.CreatedIndex,
		LogedDevices:   logedDevices,
	}, nil
}

func (uc *StudentApplicationUseCase) GetStudent4Gateway(ctx *gin.Context, studentID string) (*response.GetStudent4Gateway, error) {
	studentApp, err := uc.StudentAppRepo.GetByID(uuid.MustParse(studentID))
	if err != nil {
		return nil, err
	}
	if studentApp == nil {
		return nil, errors.New("student not found")
	}
	// get avts
	avatar, _ := uc.UserImagesUsecase.GetAvtIsMain4Owner(studentID, value.OwnerRoleStudent)
	code, _ := uc.ProfileGateway.GetStudentCode(ctx, studentID)

	res := &response.GetStudent4Gateway{
		StudentID:      studentID,
		StudentName:    studentApp.StudentName,
		OrganizationID: studentApp.OrganizationID.String(),
		Avatar:         avatar,
		Code:           code,
	}

	uc.CachingMainService.SetStudentCache(ctx.Request.Context(), studentID, res)

	return res, nil
}

func (uc *StudentApplicationUseCase) GetStudentByID4App(ctx *gin.Context, studentID string, deviceID string) (*response.StudentResponseBase, error) {

	studentApp, err := uc.StudentAppRepo.GetByID(uuid.MustParse(studentID))
	if err != nil {
		return nil, err
	}
	if studentApp == nil {
		return nil, errors.New("student not found")
	}

	// kiem tra student org va device org
	// deviceOrgIds, _ := uc.DeviceRepo.GetOrgIDsByDeviceID(deviceID)

	// found := false
	// for _, orgID := range deviceOrgIds {
	// 	if orgID == studentApp.OrganizationID {
	// 		found = true
	// 		break
	// 	}
	// }

	// if !found {
	// 	return nil, errors.New("device not associated with student's organization")
	// }

	// get student main avt url
	avatar, _ := uc.UserImagesUsecase.GetAvtIsMain4Owner(studentID, value.OwnerRoleStudent)

	return &response.StudentResponseBase{
		StudentID:   studentID,
		StudentName: studentApp.StudentName,
		CustomID:    studentApp.CustomID,
		AvatarURL:   avatar.ImageUrl,
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
		apps, err := uc.StudentAppRepo.GetApprovedAll()
		if err != nil {
			return nil, err
		}
		return uc.mapStudentAppsToResponse(ctx, apps), nil
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
	apps, err := uc.StudentAppRepo.GetByOrganizationIDsApproved(orgIDs)
	if err != nil {
		return nil, err
	}

	return uc.mapStudentAppsToResponse(ctx, apps), nil
}

func (uc *StudentApplicationUseCase) mapStudentAppsToResponse(ctx *gin.Context, students []entity.SStudentFormApplication) []response.StudentResponse {
	res := make([]response.StudentResponse, 0, len(students))
	for _, std := range students {
		isDeactive, _ := uc.StudentBlockSettingRepository.GetIsDeactiveByStudentID(std.ID.String())
		avatar, _ := uc.UserImagesUsecase.GetAvtIsMain4Owner(std.ID.String(), value.OwnerRoleStudent)
		code, _ := uc.ProfileGateway.GetStudentCode(ctx, std.ID.String())
		res = append(res, response.StudentResponse{
			StudentID:    std.ID.String(),
			StudentName:  std.StudentName,
			CreatedIndex: std.CreatedIndex,
			IsDeactive:   isDeactive,
			Avatar:       avatar,
			LanguageKeys: []string{"vietnamese", "english"},

			Code: code,
		})
	}
	return res
}

func (uc *StudentApplicationUseCase) ApproveStudentApplication(ctx *gin.Context, applicationID string) error {
	// Tìm bản ghi hiện tại theo ID
	application, err := uc.StudentAppRepo.GetByID(uuid.MustParse(applicationID))

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

	// Lưu lại
	return uc.StudentAppRepo.Update(application)
}

func (uc *StudentApplicationUseCase) BlockStudentApplication(ctx *gin.Context, applicationID string) error {
	// Tìm bản ghi hiện tại theo ID
	application, err := uc.StudentAppRepo.GetByID(uuid.MustParse(applicationID))

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
	return uc.StudentAppRepo.Update(application)
}

func (uc *StudentApplicationUseCase) GetAllStudentApplications(ctx *gin.Context) ([]response.StudentFormApplicationResponse, error) {
	// Lấy thông tin người dùng hiện tại (kèm Organizations, Roles)
	user, err := uc.GetUserEntityUseCase.GetCurrentUserWithOrganizations(ctx)
	if err != nil {
		return nil, err
	}

	var apps []entity.SStudentFormApplication

	if user.IsSuperAdmin() {
		// SuperAdmin → lấy tất cả đơn
		apps, err = uc.StudentAppRepo.GetAll()
		if err != nil {
			return nil, err
		}
	} else {
		// Nếu không phải SuperAdmin → lấy các orgIDs được quản lý
		orgIDs, err := user.GetManagedOrganizationIDs(uc.StudentAppRepo.GetDB())
		if err != nil {
			return nil, err
		}

		// Lọc các đơn theo orgID
		apps, err = uc.StudentAppRepo.GetByOrganizationIDs(orgIDs)
		if err != nil {
			return nil, err
		}
	}

	// Tạo response
	res := make([]response.StudentFormApplicationResponse, 0, len(apps))
	for _, a := range apps {
		orgStaff, _ := uc.OrganizationRepo.GetByID(a.OrganizationID.String())

		res = append(res, response.StudentFormApplicationResponse{
			ID:               a.ID.String(),
			StudentName:      a.StudentName,
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

func (uc *StudentApplicationUseCase) GetDetailStudentApplication(ctx *gin.Context, applicationID string) (*response.StudentFormApplicationResponse, error) {
	// Tìm bản ghi hiện tại theo ID
	application, err := uc.StudentAppRepo.GetByID(uuid.MustParse(applicationID))
	if err != nil {
		return nil, err
	}

	if application == nil {
		return nil, errors.New("application not found")
	}

	orgStudent, _ := uc.OrganizationRepo.GetByID(application.OrganizationID.String())
	return &response.StudentFormApplicationResponse{
		ID:               application.ID.String(),
		StudentName:      application.StudentName,
		Status:           application.Status.String(),
		ApprovedAt:       application.ApprovedAt.Format("2006-01-02 15:04:05"),
		CreatedAt:        application.CreatedAt.Format("2006-01-02 15:04:05"),
		UserID:           application.UserID.String(),
		OrganizationID:   application.OrganizationID.String(),
		OrganizationName: orgStudent.OrganizationName,
	}, nil
}

func (uc *StudentApplicationUseCase) AddCustomID(req request.AddCustomId2StudentRequest) error {
	// Tìm bản ghi hiện tại theo ID
	student := &entity.SStudentFormApplication{}
	err := uc.StudentAppRepo.DB.
		Where("id = ?", req.StudentID).
		First(student).Error
	if err != nil {
		return err
	}

	// Cập nhật tên
	student.CustomID = req.CustomID

	// Lưu lại
	return uc.StudentAppRepo.Update(student)
}

func (uc *StudentApplicationUseCase) GetStudentsByUser(userID string) ([]entity.SStudentFormApplication, error) {
	return uc.StudentAppRepo.GetByUserIDApproved(userID)
}

func (uc *StudentApplicationUseCase) GetStudentOrganizationsByUser(userID string) ([]response.StudentOrganization, error) {
	students, err := uc.GetStudentsByUser(userID)
	if err != nil {
		return nil, err
	}

	res := make([]response.StudentOrganization, 0, len(students))
	seen := make(map[string]bool) // track org.ID đã thêm

	for _, s := range students {
		org, _ := uc.OrganizationRepo.GetByID(s.OrganizationID.String())
		if org == nil {
			continue
		}

		orgID := org.ID.String()
		if seen[orgID] {
			continue
		}
		seen[orgID] = true

		res = append(res, response.StudentOrganization{
			ID:               orgID,
			OrganizationName: org.OrganizationName,
			AvatarURL:        org.AvatarURL,
		})
	}

	return res, nil
}

func (uc *StudentApplicationUseCase) GetAllByOrganizationID(organizationID string) ([]entity.SStudentFormApplication, error) {
	return uc.StudentAppRepo.GetByOrganizationID(organizationID)
}

// call profile gateway to generate student code
func (uc *StudentApplicationUseCase) GenerateStudentCode(ctx *gin.Context) {
	// get all students
	students, err := uc.StudentAppRepo.GetAll()
	if err != nil {
		return
	}

	for _, student := range students {
		// call profile gateway to generate student code
		_, _ = uc.ProfileGateway.GenerateStudentCode(ctx, student.ID.String(), student.CreatedIndex)
	}
}
