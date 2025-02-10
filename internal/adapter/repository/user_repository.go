package repository

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type UserRepository struct {
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) ValidateToken(token string, ownerId string) (bool, error) {
	req, err := http.NewRequest("POST", "http://svc-user-app/token", nil)
	if err != nil {
		fmt.Println("Error creating request", err)
		return false, err
	}

	req.Header.Add("Authorization", token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request", err)
		return false, err
	}
	defer resp.Body.Close()

	var response map[string]interface{}

	err = json.NewDecoder(resp.Body).Decode(&response)

	if err != nil {
		fmt.Println("Error decoding response", err)
		return false, err
	}

	id, err := strconv.ParseFloat(ownerId, 64)
	if err != nil {
		fmt.Println("Error parsing owner id", err)
		return false, err
	}

	return response["id"].(float64) == id, nil
}
