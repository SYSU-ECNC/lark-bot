package db

import (
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var db *gorm.DB

func init() {
	var err error

	db, err = gorm.Open(sqlite.Open("lark.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&User{})
}

// GetUserByNetID returns User by given NetID.
func GetUserByNetID(netID string) *User {
	user := &User{}
	result := db.Where(&User{NetID: netID}).First(user)
	if result.Error != nil {
		return nil
	}
	return user
}

func GetUsers() *[]User {
	var users []User
	result := db.Find(&users)
	if result.Error != nil {
		return nil
	}
	return &users
}

// UpdateUser updates user record in database.
func UpdateUser(user *User) {
	db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(user)
}

type UpdateUserTokenParams struct {
	AccessToken            string
	AccessTokenExpireTime  *time.Time
	RefreshToken           string
	RefreshTokenExpireTime *time.Time
}

func UpdateUserTokenByOpenID(openID string, params *UpdateUserTokenParams) {
	user := User{}
	db.FirstOrInit(&user, User{
		OpenID: openID,
	})
	user.AccessToken = params.AccessToken
	user.AccessTokenExpireTime = params.AccessTokenExpireTime
	user.RefreshToken = params.RefreshToken
	user.RefreshTokenExpireTime = params.RefreshTokenExpireTime
	db.Save(&user)
}
