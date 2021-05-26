package db

import (
	"time"

	"gorm.io/gorm"
)

// User model.
type User struct {
	gorm.Model
	NetID                  string `gorm:"index"`
	OpenID                 string `gorm:"uniqueIndex"`
	UnionID                string `gorm:"uniqueIndex"`
	AccessToken            string
	AccessTokenExpireTime  *time.Time
	RefreshToken           string
	RefreshTokenExpireTime *time.Time
}
