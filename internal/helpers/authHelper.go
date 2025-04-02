package helpers

import (
	"fmt"
	"net/http"

	"encoding/json"
)

func NotidierIdTaking(email string) (string, error) {
	url := fmt.Sprintf("http://user-auth-service:8081/users/%s", email)
	resp, err := http.Get(url)
	if err != nil {
		return "", ErrorHelper(err, "error taking id from user service")
	}
	var data map[string]string
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return "", ErrorHelper(err, "Error while decoding response")
	}
	userId, ok := data["user_id"]
	if !ok || userId == "" {
		return "", fmt.Errorf("error parsing data from map")
	}
	return userId, nil
}
