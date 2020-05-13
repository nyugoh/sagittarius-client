package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
)

func SendMail(from, toName, to, subject, body string) (ok bool, err error)  {
	requestBody := map[string]interface{}{
		"to": to,
		"name": toName,
		"from": from,
		"subject": subject,
		"body": body,
		"sentBy": os.Getenv("APP_NAME"),
		"feedbackEmail": "joenyugoh@gmail.com",
	}

	bytesRepresentation, err := json.Marshal(requestBody)
	if err != nil {
		return false, err
	}

	res , err := http.Post(os.Getenv("Q_MAILER_URL"), "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		return false, err
	}
	var result map[string]interface{}

	json.NewDecoder(res.Body).Decode(&result)

	Log(result)
	if res.StatusCode == 200 {
		return true, nil
	} else {
		LogError(result["error"])
		return false, errors.New(fmt.Sprintf("%s", result["error"]))
	}
}

