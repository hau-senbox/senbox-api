package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/request"
	"sen-global-api/pkg/consulapi/gateway"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GenerateOwnerCodeUseCase interface {
	GenerateUserCode(ctx *gin.Context, userID string) (*string, error)
	GenerateStudentCode(ctx *gin.Context, studentID string) (*string, error)
	GenerateTeacherCode(ctx *gin.Context, teacherID string) (*string, error)
	GenerateStaffCode(ctx *gin.Context, staffID string) (*string, error)
	GenerateParentCode(ctx *gin.Context, parentID string) (*string, error)
	GenerateChildCode(ctx *gin.Context, childID string) (*string, error)
}

type generateOwnerCodeUseCase struct {
	userRepo       *repository.UserEntityRepository
	teacherRepo    *repository.TeacherApplicationRepository
	studentRepo    *repository.StudentApplicationRepository
	staffRepo      *repository.StaffApplicationRepository
	childRepo      *repository.ChildRepository
	parentRepo     *repository.ParentRepository
	profileGateway gateway.ProfileGateway
}

func NewGenerateOwnerCodeUseCase(
	userRepo *repository.UserEntityRepository,
	teacherRepo *repository.TeacherApplicationRepository,
	studentRepo *repository.StudentApplicationRepository,
	staffRepo *repository.StaffApplicationRepository,
	childRepo *repository.ChildRepository,
	parentRepo *repository.ParentRepository,
	profileGateway gateway.ProfileGateway,
) GenerateOwnerCodeUseCase {
	return &generateOwnerCodeUseCase{
		userRepo:       userRepo,
		teacherRepo:    teacherRepo,
		studentRepo:    studentRepo,
		staffRepo:      staffRepo,
		childRepo:      childRepo,
		parentRepo:     parentRepo,
		profileGateway: profileGateway,
	}
}

func (u *generateOwnerCodeUseCase) GenerateUserCode(ctx *gin.Context, userID string) (*string, error) {
	user, _ := u.userRepo.GetByID(request.GetUserEntityByIDRequest{ID: userID})
	if user == nil {
		return nil, nil
	}
	u.profileGateway.GenerateUserCode(ctx, userID, user.CreatedIndex)
	return nil, nil
}

func (u *generateOwnerCodeUseCase) GenerateStudentCode(ctx *gin.Context, studentID string) (*string, error) {
	student, _ := u.studentRepo.GetByID(uuid.MustParse(studentID))
	if student == nil {
		return nil, nil
	}
	u.profileGateway.GenerateStudentCode(ctx, studentID, student.CreatedIndex)
	return nil, nil
}

func (u *generateOwnerCodeUseCase) GenerateTeacherCode(ctx *gin.Context, teacherID string) (*string, error) {
	teacher, _ := u.teacherRepo.GetByID(uuid.MustParse(teacherID))
	if teacher == nil {
		return nil, nil
	}
	u.profileGateway.GenerateTeacherCode(ctx, teacherID, teacher.CreatedIndex)
	return nil, nil
}

func (u *generateOwnerCodeUseCase) GenerateStaffCode(ctx *gin.Context, staffID string) (*string, error) {
	staff, _ := u.staffRepo.GetByID(uuid.MustParse(staffID))
	if staff == nil {
		return nil, nil
	}
	u.profileGateway.GenerateStaffCode(ctx, staffID, staff.CreatedIndex)
	return nil, nil
}

func (u *generateOwnerCodeUseCase) GenerateParentCode(ctx *gin.Context, parentID string) (*string, error) {
	parent, _ := u.parentRepo.GetByID(ctx, parentID)
	if parent == nil {
		return nil, nil
	}
	u.profileGateway.GenerateParentCode(ctx, parentID, parent.CreatedIndex)
	return nil, nil
}

func (u *generateOwnerCodeUseCase) GenerateChildCode(ctx *gin.Context, childID string) (*string, error) {
	child, _ := u.childRepo.GetByID(childID)
	if child == nil {
		return nil, nil
	}
	u.profileGateway.GenerateChildCode(ctx, childID, child.CreatedIndex)
	return nil, nil
}
