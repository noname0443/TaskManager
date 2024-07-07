package integration

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/noname0443/task_manager/env"
	"github.com/noname0443/task_manager/models"
)

func GetPeopleInfo(passport string) (*models.User, error) {
	strs := strings.Split(passport, " ")
	if len(strs) != 2 {
		return nil, fmt.Errorf("incorrect format of passportNumber: %s", passport)
	}
	passportSerie, passportNumber := strs[0], strs[1]

	client := &http.Client{}
	req, err := http.NewRequest("GET", os.Getenv(env.EXTERNAL_WEBSERIVCE_URL), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("passportSerie", passportSerie)
	q.Add("passportNumber", passportNumber)
	req.URL.RawQuery = q.Encode()
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	user := models.User{}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(b, &user)
	user.PassportNumber = passport
	return &user, nil
}
