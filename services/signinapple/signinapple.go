package signinapple

import (
	"context"
	"dongchamao/global"
	"errors"
	"github.com/Timothylock/go-signin-with-apple/apple"
	"github.com/astaxie/beego/logs"
)

var secretString = `-----BEGIN PRIVATE KEY-----
MIGTAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBHkwdwIBAQQgLGo1mDA959OvZ2Vn
3fPYiaC2vFCCMf9Y0fWTCWS5IeigCgYIKoZIzj0DAQehRANCAATXjPxyw8Ijq2at
Tqc/jfJ07dgEJY3D4OWcL8r93iYlbLIKaS1LSdPI0S+ep23BQi9hLl9tAcIEjlr5
9tGI4r+k
-----END PRIVATE KEY-----`

func GetUniqueId(token string) (unique string, err error) {
	// Generate the client secret used to authenticate with Apple's validation servers
	teamID := global.Cfg.String("apple_team_id")
	clientID := global.Cfg.String("apple_client_id")
	keyID := global.Cfg.String("apple_key_id")
	secret, err := apple.GenerateClientSecret(secretString, teamID, clientID, keyID)
	if err != nil {
		logs.Error("error generating secret: " + err.Error())
		return
	}

	// Generate a new validation client
	client := apple.New()

	vReq := apple.AppValidationTokenRequest{
		ClientID:     clientID,
		ClientSecret: secret,
		Code:         token,
	}

	var resp apple.ValidationResponse

	// Do the verification
	err = client.VerifyAppToken(context.Background(), vReq, &resp)
	if err != nil {
		logs.Error("error verifying: " + err.Error())
		return
	}

	if resp.Error != "" {
		logs.Error("apple returned an error: " + resp.Error)
		err = errors.New(resp.Error)
		return
	}

	// Get the unique user ID
	unique, err = apple.GetUniqueID(resp.IDToken)
	if err != nil {
		logs.Error("failed to get unique ID: " + err.Error())
		return
	}
	return
}
