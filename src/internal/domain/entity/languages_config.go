package entity

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"sen-global-api/internal/domain/value"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LanguagesConfig struct {
	ID         uuid.UUID                  `gorm:"type:char(36);primaryKey"`
	OwnerID    string                     `json:"owner_id" gorm:"type:varchar(255);not null"`
	OwnerRole  value.OwnerRole4LangConfig `json:"owner_role" gorm:"type:varchar(50);not null"`
	SpokenLang LanguageConfigList         `json:"spoken_lang" gorm:"type:json"`
	StudyLang  LanguageConfigList         `json:"study_lang" gorm:"type:json"`
	CreatedAt  time.Time                  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time                  `json:"updated_at" gorm:"autoUpdateTime"`
}

type LanguageConfig struct {
	Order    int    `json:"order"`
	Language string `json:"language"`
	Origin   string `json:"origin"`
	Percent  int    `json:"percent"`
	Note     string `json:"note"`
}

type LanguageConfigList []LanguageConfig

// driver.Valuer → serialize slice thành JSON khi lưu MySQL
func (lc LanguageConfigList) Value() (driver.Value, error) {
	return json.Marshal(lc)
}

// sql.Scanner → parse JSON từ MySQL thành slice struct
func (lc *LanguageConfigList) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("LanguageConfigList: Scan source is not []byte")
	}
	return json.Unmarshal(bytes, lc)
}

// GORM hook: tự tạo UUID nếu chưa có
func (lc *LanguagesConfig) BeforeCreate(tx *gorm.DB) (err error) {
	if lc.ID == uuid.Nil {
		lc.ID = uuid.New()
	}
	return
}
