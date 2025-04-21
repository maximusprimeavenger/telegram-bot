package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)


func NotifierIdTaking(email string) (string, error) {
	url := fmt.Sprintf("http://user-auth-service:8081/notifier/%s", email)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request:%v", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request:%v", err)
	}
	defer resp.Body.Close()
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
