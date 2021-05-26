package app

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sysu-ecnc/lark-bot/internal/pkg/db"
	"github.com/sysu-ecnc/lark-bot/internal/pkg/lark"
)

// Start a HTTP Server.
func Start() {
	r := gin.Default()

	r.GET("/authorize", beginAuthorize)
	r.GET("/callback", authorizeCallback)
	r.GET("/select", selectAll)
	r.GET("/files", listAllFiles)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func beginAuthorize(c *gin.Context) {
	c.Redirect(302, lark.GetAuthorizeURL("http://127.0.0.1:8080/callback", ""))
}

func authorizeCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(400, gin.H{
			"message": "missing parameter",
		})
		return
	}

	ret, err := lark.GetUserInfo(code)

	if err != nil {
		c.JSON(500, err)
		return
	}

	accessTokenExpireTime := time.Now().Add(time.Duration(ret.ExpiresIn) * time.Second)
	refreshTokenExpireTime := time.Now().Add(time.Duration(ret.RefreshExpiresIn) * time.Second)

	db.UpdateUserTokenByOpenID(ret.OpenId, &db.UpdateUserTokenParams{
		AccessToken:            ret.AccessToken,
		AccessTokenExpireTime:  &accessTokenExpireTime,
		RefreshToken:           ret.RefreshToken,
		RefreshTokenExpireTime: &refreshTokenExpireTime,
	})
	c.JSON(200, ret)
}

func selectAll(c *gin.Context) {
	c.JSON(200, db.GetUsers())
}

func listAllFiles(c *gin.Context) {
	user := db.GetUserByNetID("")
	ret, err := lark.ListFiles(user.AccessToken, "fldcnSod1sJbqmUJ1udYzj7ZEEd")

	if err != nil {
		c.JSON(500, err)
	}
	c.JSON(200, ret)
}
