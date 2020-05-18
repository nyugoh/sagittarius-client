package exporter

import (
	"fmt"
	"github.com/nyugoh/sagittarius-client/utils"
	"os"
	"strings"
)

// Starts monitoring a certain folder
func StartMonitor() []string {
	// Check which folder it needs to watch
	res, err := utils.SendGet(os.Getenv("SAG_SERVER") + "/clients/folders")
	if err != nil {
		utils.NotifyError(err.Error())
		return nil
	}

	foldersPayload := fmt.Sprintf("%v", res["folders"])
	//config := fmt.Sprintf("%v", res["config"])

	folders := make([]string, 0)
	for _, folder := range strings.Split(foldersPayload, ",") {
		folders = append(folders, strings.TrimSpace(folder))
	}
	return folders
}
