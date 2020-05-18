package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	uuid2 "github.com/google/uuid"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func SendError(c *gin.Context, msg string) {
	LogError(msg)
	c.JSON(http.StatusBadRequest, gin.H{
		"error": msg,
	})
}

func SendJson(c *gin.Context, payload gin.H) {
	c.JSON(http.StatusOK, payload)
}

func RespondWithError(w http.ResponseWriter, code int, msg string) {
	RespondWithJSON(w, code, map[string]interface{}{"error": msg})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	Log(fmt.Sprintf("RESPONSE:: Status:%d Payload: %v", code, payload))
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func CurrentTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func ValidateEmail(email string) (bool, error) {
	var rxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if len(email) > 254 || !rxEmail.MatchString(email) {
		return false, errors.New("email is invalid")
	}
	return true, nil
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		Log("Hit an auth required endpoint...")
		session := sessions.Default(c)
		userId := session.Get("userId")
		username := session.Get("username")
		if username == nil || userId == nil {
			//c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Redirect(302, "/auth/login")
			return
		}
		Log("Request made by user Id:", userId, "Username:", username)
		c.Set("userId", userId)
		c.Set("username", username)
		c.Next()
	}
}

func MetricsMonitor() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		// before request
		c.Next()
		// after request
		latency := time.Since(t)
		// access the status we are sending
		status := c.Writer.Status()
		if !strings.Contains(c.Request.URL.Path, "socket.io") { // Avoid logging for /socket.io
			log.Println("Request: ", c.Request.URL.Path, "took:", latency, "Status:", status)
		}
	}
}

func ExtractUser(c *gin.Context) (username string, userId int) {
	if name, ok := c.Get("username"); ok {
		username = fmt.Sprintf("%v", name)
	} else {
		username = "Unknown"
	}

	if id, ok := c.Get("userId"); ok {
		userId, _ = strconv.Atoi(fmt.Sprintf("%v", id))
	} else {
		userId = 0
	}
	return
}

func ExitApp(code int) {
	Log("Exiting app...")
	os.Exit(code)
}

func GenerateUUID() string {
	uuid := uuid2.New()
	return uuid.String()
}

func SendPost(url string, payload map[string]interface{}) (result map[string]interface{}, err error) {
	bytesRepresentation, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	res, err := http.Post(url, "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		return nil, err
	}

	json.NewDecoder(res.Body).Decode(&result)

	Log(result)
	if res.StatusCode == 200 {
		return result, nil
	} else {
		errMsg := fmt.Sprintf("%s", result["error"])
		LogError(errMsg)
		NotifyError(errMsg)
		return nil, errors.New(errMsg)
	}
}

func SendGet(url string) (result map[string]interface{}, err error)  {
	client := http.Client{}
	request, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer([]byte("")))
	name := os.Getenv("APP_NAME")
	request.Header.Set("Authorization", os.Getenv(name+"_JWT"))
	Log("Sending a GET request: Authorization:", os.Getenv(name+"_JWT"))

	res, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	json.NewDecoder(res.Body).Decode(&result)

	Log(result)
	if res.StatusCode == 200 {
		return result, nil
	} else {
		errMsg := fmt.Sprintf("%s", result["error"])
		LogError(errMsg)
		NotifyError(errMsg)
		return nil, errors.New(errMsg)
	}
}

func NotifyError(message string)  {
	_, err := SendMail("admin@quebasetech.co.ke", "Joe Nyugoh","joenyugoh@gmail.com", "Error occurred in "+os.Getenv("APP_NAME"), `
<div>
    <p>Error on a sagittarius client</p>
    <br>
	<p>`+ message+ `</p>
    <br>
    <p>Regards</p>
    <p>Que Base Tech</p>
    <br>
</div>

`)
	if err != nil {
		LogError(err.Error())
	}
}