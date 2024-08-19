package config

import (
	"encoding/json"
	"os"

	"github.com/nazzarr03/notification/models"
)

func WriteToFile(emailMsg models.Message) error {
	file, err := os.OpenFile("emails.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := json.Marshal(emailMsg)
	if err != nil {
		return err
	}

	_, err = file.Write(append(data, '\n'))
	return err
}
