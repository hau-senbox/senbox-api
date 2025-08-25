package router

import (
	"sen-global-api/config"
	"sen-global-api/internal/controller"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/usecase"
	"sen-global-api/internal/middleware"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func setupGatewayRoutes(r *gin.Engine, dbConn *gorm.DB, appCfg config.AppConfig) {
	// init repository + usecase
	sessionRepository := repository.SessionRepository{
		OrganizationRepository: &repository.OrganizationRepository{DBConn: dbConn},
		AuthorizeEncryptKey:    appCfg.AuthorizeEncryptKey,

		TokenExpireTimeInHour: time.Duration(appCfg.TokenExpireDurationInHour),
	}
	secureMiddleware := middleware.SecuredMiddleware{SessionRepository: sessionRepository}

	userEntityRepository := &repository.UserEntityRepository{DBConn: dbConn}

	// student
	studentRepo := &repository.StudentApplicationRepository{DB: dbConn}
	studentUsecase := &usecase.StudentApplicationUseCase{StudentAppRepo: studentRepo}

	// teacher
	teacherRepo := &repository.TeacherApplicationRepository{DBConn: dbConn}
	teacherUsecase := &usecase.TeacherApplicationUseCase{
		TeacherRepo:          teacherRepo,
		UserEntityRepository: userEntityRepository,
	}

	// staff
	staffRepo := &repository.StaffApplicationRepository{DBConn: dbConn}
	staffUsecase := &usecase.StaffApplicationUseCase{
		StaffAppRepo:         staffRepo,
		UserEntityRepository: userEntityRepository,
	}

	userEntityCtrl := &controller.UserEntityController{
		StudentApplicationUseCase: studentUsecase,
		TeacherApplicationUseCase: teacherUsecase,
		StaffApplicationUseCase:   staffUsecase,
	}

	api := r.Group("/v1/gateway", secureMiddleware.Secured())
	{
		student := api.Group("/students")
		{
			student.GET("/:student_id", userEntityCtrl.GetStudent4Gateway)
		}

		teacher := api.Group("/teachers")
		{
			teacher.GET("/:teacher_id", userEntityCtrl.GetTeacher4Gateway)
		}

		staff := api.Group("/staffs")
		{
			staff.GET("/:staff_id", userEntityCtrl.GetStaff4Gateway)
		}
	}
}
