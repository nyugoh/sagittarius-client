package api

import (
	"github.com/gin-gonic/gin"
	"github.com/nyugoh/sagittarius-client/app/models"
	"github.com/nyugoh/sagittarius-client/utils"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strings"
)

func (app *App) ListFolders(c *gin.Context) {
	logs := make([]map[string][]models.LogFile, 0)
	utils.Log(app.Folders)
	for _, folder := range app.Folders {
		folder = strings.TrimSpace(folder)
		logs = append(logs, map[string][]models.LogFile{
			folder: ListDir(folder, ".log"),
		})
	}

	// Current app logs
	appName := strings.TrimSpace(strings.ToLower(os.Getenv("APP_NAME")))
	logs = append(logs, map[string][]models.LogFile{appName: ListLogs()})

	utils.SendJson(c, gin.H{
		"logs": logs,
	})
}

func (app *App) ReadLog(c *gin.Context) {
	logFile := c.Request.URL.Query()["log"]
	utils.Log("Reading log:", logFile)
	if !strings.Contains(logFile[0], ".log") && !strings.Contains(logFile[0], ".sql") {
		_, err := utils.SendMail("admin@quebasetech.co.ke", "Joe Nyugoh", "joenyugoh@gmail.com", "Server Notification", "<p>Someone is trying to access files outside logs folder</p><p>Folder::"+logFile[0]+"</p>")
		if err != nil {
			utils.LogError(err.Error())
		}
		utils.SendError(c, "You are trying to access restricted file: E-mail sent to admin")
		return
	}
	content, err := ioutil.ReadFile(logFile[0])
	if err != nil {
		utils.SendError(c, "Unable to read log file:"+err.Error())
		return
	}
	utils.SendJson(c, gin.H{
		"status":  "success",
		"payload": string(content),
	})
}

func ListLogs() []models.LogFile {
	appName := strings.TrimSpace(strings.ToLower(os.Getenv("APP_NAME")))
	if len(appName) == 0 {
		appName = "app-logs" // Default log file name
	}
	logFolder := os.Getenv(appName + "_LOG_FOLDER")
	logs := ListDir(logFolder, ".log")
	return logs
}

func ListDir(dirPath, fileExt string) []models.LogFile {
	utils.Log("Trying to read ", dirPath, " for ", fileExt, "files")

	files := make([]models.LogFile, 0)
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			utils.LogError(err.Error())
			return err
		}
		if info.IsDir() || filepath.Ext(path) != fileExt {
			return nil
		}
		utils.Log("Path:", path, "Info: Size:", info.Size(), "Name:", info.Name())
		file := models.LogFile{
			Path: path,
			Size: toFixed(float64(info.Size())/1024576.00, 2),
			Date: strings.Split(info.Name(), ".")[1],
			Name: info.Name(),
		}
		files = append(files, file)
		return nil
	})
	if err != nil {
		utils.LogError(err.Error())
		return files
	}
	return files
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}
