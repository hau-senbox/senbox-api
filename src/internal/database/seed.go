package database

import (
	"bufio"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/entity/components"
	"sen-global-api/internal/domain/entity/menu"
	"sen-global-api/pkg/common"
	"time"

	"gorm.io/gorm/clause"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	_, b, _, _ = runtime.Caller(0)

	// Root folder of this project
	Root = filepath.Join(filepath.Dir(b), "../..")
)

func Seed(db *gorm.DB, config *common.Config, seedSQLFile string) error {

	// Migrate the schema
	err := db.AutoMigrate(
		&entity.SAppKey{},
		&entity.SQuestion{},
		&entity.SForm{},
		&entity.SRedirectUrl{},
		&entity.SSetting{},
		&entity.SToDo{},
		&entity.SSubmission{},
		&entity.SFormQuestion{},
		&entity.SMobileDevice{},
		&entity.SCodeCounting{},
		&entity.SDevice{},
		&entity.SOrganization{},
		&entity.SDeviceComponentValues{},
		&entity.SRole{},
		&entity.SUserEntity{},
		&entity.SUserRoles{},
		&entity.SFunctionClaim{},
		&entity.SFunctionClaimPermission{},
		&entity.SUserDevices{},
		&entity.SImage{},
		&entity.SVideo{},
		&entity.SAudio{},
		&entity.SPdf{},
		&entity.SUserFCMToken{},
		&entity.SUserFunctionAuthorize{},
		&entity.SUserOrg{},
		&entity.SOrgFormApplication{},
		&entity.SPreRegister{},
		&components.Component{},
		&menu.SuperAdminMenu{},
		&menu.OrgMenu{},
		&menu.UserMenu{},
		&menu.DeviceMenu{},
		&entity.PublicImage{},
		&entity.SStaffFormApplication{},
		&entity.SStudentFormApplication{},
		&entity.STeacherFormApplication{},
		&entity.SOrgDevices{},
		&entity.MemoryComponentValue{},
		&entity.SUserParentChild{},
		&entity.SAnswer{},
		&entity.SRoleOrgSignUp{},
		&entity.SChild{},
		&entity.ChildMenu{},
		&entity.StudentMenu{},
		&entity.TeacherMenu{},
		&entity.StaffMenu{},
		&entity.OrganizationMenuTemplate{},
		&entity.SyncQueue{},
		&entity.UserBlockSetting{},
		&entity.SDeviceMenuV2{},
		&entity.ParentMenu{},
		&entity.StudentBlockSetting{},
		&entity.OrganizationSetting{},
		&entity.LanguagesConfig{},
		&entity.OrganizationNewsSetting{},
		&entity.UserImages{},
		&entity.AppConfig{},
		&entity.TeacherMenuOrganization{},
		&entity.UserDevicesLogin{},
		&entity.UserSetting{},
		&entity.DepartmentMenu{},
		&entity.DepartmentMenuOrganization{},
		//&entity.ClassroomMenu{},
		&entity.ValuesAppCurrent{},
		&entity.AccountsLog{},
		&entity.SuperAdminEmergencyMenu{},
		&entity.OrganizationEmergencyMenu{},
		&entity.LanguageSetting{},
		&entity.DataLog{},
		&entity.OrganizationSettingMenu{},
		&entity.MessageLanguage{},
		&entity.SParent{},
		&entity.SParentChilds{},
		&entity.StudentMenuOrganization{},
		&entity.ValuesAppHistories{},
	)

	// Seed
	log.Debug("Seeding database...")
	if err != nil {
		return err
	}

	//Seeding data
	file, err := os.Open(Root + seedSQLFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		query := scanner.Text()
		if query != "" {
			err = db.Exec(scanner.Text()).Error
			if err != nil {
				return err
			}
		}
	}

	// var superAdmin = &entity.SUserEntity{
	// 	RoleID:             uuid.New(),
	// 	Username:       "root",
	// 	Fullname:       "Senbox Super Administrator",
	// 	Birthday:       time.Time{},
	// 	Phone:          "",
	// 	Email:          "admin@senbox.vn",
	// 	Password:       "SEN@box",
	// 	OrganizationID: 1,
	// }

	var appKey = &entity.SAppKey{
		ID:     1,
		AppKey: "8c30f8d5-f430-4079-bc7f-cfea3d61704d",
	}

	// var users = []*entity.SUserEntity{superAdmin}
	// err = db.Clauses(
	// 	clause.OnConflict{
	// 		Columns:   []clause.Column{{Name: "id"}},
	// 		DoUpdates: clause.AssignmentColumns([]string{"username"}),
	// 	}).Create(&users).Error
	// if err != nil {
	// 	log.Error(err)
	// }

	err = db.Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"id"}),
		}).Create(&appKey).Error
	if err != nil {
		log.Error(err)
	}

	log.Debug("Seeding database done")
	//for i := 0; i < 24; i++ {
	//	randString := RandString(10)
	//	var code = &entity.SCodeCounting{
	//		Token:        randString,
	//		CurrentValue: i,
	//	}
	//
	//	err = db.Clauses(
	//		clause.OnConflict{
	//			Columns:   []clause.Column{{Name: "token"}},
	//			DoUpdates: clause.AssignmentColumns([]string{"current_value"}),
	//		}).Create(&code).Error
	//	if err != nil {
	//		log.Error(err)
	//	}
	//}

	return nil
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
