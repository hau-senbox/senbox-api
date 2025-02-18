package database

import (
	"bufio"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"sen-global-api/internal/domain/entity"
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
		&entity.SUser{},
		&entity.SQuestion{},
		&entity.SForm{},
		&entity.SRedirectUrl{},
		&entity.SSetting{},
		&entity.SToDo{},
		&entity.SDeviceFormDataset{},
		&entity.SSubmission{},
		&entity.SFormQuestion{},
		&entity.SMobileDevice{},
		&entity.SCodeCounting{},
		&entity.SDevice{},
		&entity.SCompany{},
		&entity.SDeviceComponentValues{},
		&entity.SRole{},
		&entity.SUserEntity{},
		&entity.SRolePolicy{},
		&entity.SUserGuardians{},
		&entity.SUserRoles{},
		&entity.SRoleClaim{},
		&entity.SRolePolicyRoles{},
		&entity.SRolePolicyClaims{},
		&entity.SUserPolicies{},
		&entity.SUserDevices{},
		&entity.SUserConfig{},
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

	var bot = &entity.SUser{
		UserId:      "F9046EC9-5703-42F0-97D8-750138D52E4C",
		Username:    "techbot",
		Fullname:    "Bot",
		Birthday:    time.Time{},
		Phone:       "",
		Email:       "techbot@senbox.vn",
		Address:     "",
		Job:         "",
		CountryCode: "",
		Password:    "techbot",
		Role:        28,
	}

	var admin = &entity.SUser{
		UserId:      "663688d9-639c-4691-bc3b-612fcb3b6b48",
		Username:    "admin",
		Fullname:    "SEN Admin",
		Birthday:    time.Time{},
		Phone:       "",
		Email:       "admin@senbox.vn",
		Address:     "",
		Job:         "",
		CountryCode: "",
		Password:    "SEN@box",
		Role:        28,
	}

	var appKey = &entity.SAppKey{
		ID:     1,
		AppKey: "8c30f8d5-f430-4079-bc7f-cfea3d61704d",
	}

	var users = []*entity.SUser{bot, admin}
	err = db.Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"username"}),
		}).Create(&users).Error
	if err != nil {
		log.Error(err)
	}

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
