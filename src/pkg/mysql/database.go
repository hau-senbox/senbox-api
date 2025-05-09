package mysql

import (
	"fmt"
	"log"
	"os"
	"sen-global-api/config"
	"time"

	_ "time/tzdata"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func Establish(appConfig config.AppConfig) (*gorm.DB, error) {
	return connect(appConfig)
}

func connect(appConfig config.AppConfig) (*gorm.DB, error) {
	fmt.Println("Connect to database.....")
	var err error
	var (
		devHostName = appConfig.Config.Host
		devDbName   = appConfig.Config.Database.Database
		devUser     = appConfig.Config.User
		devDbPort   = appConfig.Config.Database.Port
		devPassword = appConfig.Config.Password
	)

	dsn := devUser + ":" + devPassword + "@tcp(" + devHostName + ":" + devDbPort + ")/" + devDbName + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt: false,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             time.Second, // Slow SQL threshold
				LogLevel:                  logger.Info, // Log level
				IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
				Colorful:                  false,       // Disable color
			},
		),
	})
	if err != nil {
		panic(err)
	}
	//db.LogMode(appConfig.Config.Env == "development")
	//db.DB().SetMaxIdleConns(appConfig.Config.Database.MaxIdleConn)
	//db.DB().SetMaxOpenConns(appConfig.Config.Database.MaxConn)
	//db.DB().SetConnMaxLifetime(time.Duration(appConfig.Config.Database.MaxLifetime))
	logrus.Info("Established database!")

	return db, nil
}
