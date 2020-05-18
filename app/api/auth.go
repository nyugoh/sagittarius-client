package api

import (
	"fmt"
	"github.com/nyugoh/sagittarius-client/utils"
	"os"
)

func (app *App) RequestJWT() (res map[string]interface{}, err error) {
	appName := os.Getenv("CLIENT_NAME")
	appHash := os.Getenv("APP_HASH")
	appPort := os.Getenv("APP_PORT")
	appIp := os.Getenv("APP_IP")

	res, err = utils.SendPost(os.Getenv("SAG_SERVER")+"/clients/auth", map[string]interface{}{
		"appName": appName,
		"appHash": appHash,
		"appPort": appPort,
		"appIp": appIp,
	})
	return
}

func (app *App) Login() {
	res, err := app.RequestJWT()

	if err != nil {
		utils.NotifyError(err.Error())
		utils.LogError(err.Error())
		return
	}
	utils.Log("Token::", res["payload"])

	name := os.Getenv("APP_NAME")
	os.Setenv(name+"_JWT", fmt.Sprintf("%v", res["payload"]))
}
