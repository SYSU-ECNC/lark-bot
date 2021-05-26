package lark

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/larksuite/oapi-sdk-go/api/core/request"
	"github.com/larksuite/oapi-sdk-go/api/core/response"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/config"
	"github.com/larksuite/oapi-sdk-go/core/constants"
	"github.com/larksuite/oapi-sdk-go/core/log"
	"github.com/larksuite/oapi-sdk-go/core/tools"
	authen "github.com/larksuite/oapi-sdk-go/service/authen/v1"
	driveExplorer "github.com/larksuite/oapi-sdk-go/service/drive_explorer/v2"
)

var conf *config.Config

// InitConfig should be invoked before any API Requests.
func InitConfig(appID string, appSecret string, verificationToken string, encryptKey string) error {
	if appID == "" || appSecret == "" || verificationToken == "" || encryptKey == "" {
		return errors.New("Missing required parameters")
	}

	// 企业自建应用的配置
	// AppID、AppSecret: "开发者后台" -> "凭证与基础信息" -> 应用凭证（App ID、App Secret）
	// VerificationToken、EncryptKey："开发者后台" -> "事件订阅" -> 事件订阅（Verification Token、Encrypt Key）。
	appSetting := config.NewInternalAppSettings(appID, appSecret, verificationToken, encryptKey)

	// 当前访问的是飞书，使用默认存储、默认日志（Debug级别），更多可选配置，请看：README.zh.md->高级使用->如何构建整体配置（Config）。
	conf = config.NewConfigWithDefaultStore(constants.DomainFeiShu, appSetting, log.NewDefaultLogger(), log.LevelInfo)

	return nil
}

func GetAuthorizeURL(redirectUrl, state string) string {
	return "https://open.feishu.cn/open-apis/authen/v1/index?redirect_uri=" +
		url.QueryEscape(redirectUrl) +
		"&app_id=" +
		conf.GetAppSettings().AppID +
		"&state=" +
		url.QueryEscape(state)
}

func GetUserInfo(code string) (*authen.UserAccessTokenInfo, error) {
	// body := map[string]interface{}{
	// 	"grant_type": "authorization_code",
	// 	"code": code,
	// }
	// ret := make(map[string]interface{})
	// req := request.NewRequestWithNative("access_token", "GET", request.AccessTokenTypeApp, body, &ret)
	// coreCtx := core.WrapContext(context.Background())
	// api.Send(coreCtx, conf, req)
	service := authen.NewService(conf)
	coreCtx := core.WrapContext(context.Background())

	body := &authen.AuthenAccessTokenReqBody{
		GrantType: "authorization_code",
		Code:      code,
	}
	req := service.Authens.AccessToken(coreCtx, body)
	ret, err := req.Do()

	fmt.Println(coreCtx.GetRequestID())
	// 打印请求的响应状态吗
	fmt.Println(coreCtx.GetHTTPStatusCode())

	if err != nil {
		e := err.(*response.Error)
		fmt.Println(e.Code)
		fmt.Println(e.Msg)
		fmt.Println(tools.Prettify(err))
		return nil, err
	}
	fmt.Println(tools.Prettify(ret))

	return ret, nil
}

func ListFiles(accessToken string, folderToken string) (*driveExplorer.FolderChildrenResult, error) {
	service := driveExplorer.NewService(conf)
	coreCtx := core.WrapContext(context.Background())

	req := service.Folders.Children(coreCtx, request.SetUserAccessToken(accessToken))
	req.SetFolderToken(folderToken) // "fldcnSod1sJbqmUJ1udYzj7ZEEd"

	ret, err := req.Do()

	fmt.Println(coreCtx.GetRequestID())
	// 打印请求的响应状态吗
	fmt.Println(coreCtx.GetHTTPStatusCode())

	if err != nil {
		e := err.(*response.Error)
		fmt.Println(e.Code)
		fmt.Println(e.Msg)
		fmt.Println(tools.Prettify(err))
		return nil, err
	}
	fmt.Println(tools.Prettify(ret))

	return ret, nil
}
