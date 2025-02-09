package repository

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type UserRepository struct {
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) ValidateToken(token string) (string, error) {
	req, err := http.NewRequest("GET", "https://api.com/validate", nil)
	if err != nil {
		fmt.Println("Error creating request", err)
		return "", err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request", err)
		return "",
			err
	}
	defer resp.Body.Close()

	var response map[string]interface{}

	err = json.NewDecoder(resp.Body).Decode(&response)

	if err != nil {
		fmt.Println("Error decoding response", err)
		return "", err
	}

	return response["user_id"].(string), nil

}
