package usecase

import (
	"errors"
	"fmt"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/response"

	"github.com/google/uuid"
)

type StudentApplicationUseCase struct {
	StudentAppRepo  *repository.StudentApplicationRepository
	StudentMenuRepo *repository.StudentMenuRepository
	ComponentRepo   *repository.ComponentRepository
}

func NewStudentApplicationUseCase(repo *repository.StudentApplicationRepository) *StudentApplicationUseCase {
	return &StudentApplicationUseCase{
		StudentAppRepo: repo,
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
			StudentID:   fmt.Sprintf("%d", a.ID),
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
	menus := make([]response.ComponentStudentResponse, 0)
	for _, comp := range components {
		menu := response.ComponentStudentResponse{
			ID:    comp.ID.String(),
			Name:  comp.Name,
			Type:  comp.Type.String(),
			Key:   comp.Key,
			Value: string(comp.Value),
			Order: componentOrderMap[comp.ID],
			Ishow: componentIsShowMap[comp.ID],
		}
		menus = append(menus, menu)
	}

	return &response.StudentResponseBase{
		StudentID:   studentID,
		StudentName: studentApp.StudentName,
		Avatar:      "", // Nếu bạn có trường Avatar trong DB thì lấy thêm ở đây
		AvatarURL:   "", // Có thể generate từ link
		Menus:       menus,
	}, nil
}
