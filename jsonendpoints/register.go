package jsonendpoints

import (
	"encoding/json"
	"io/ioutil"
	storage "lena/inmemory"
	"lena/models"
	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		handleRegisterErrorResponse(w, -1, "Method not allowed")
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		handleRegisterErrorResponse(w, -1, "Error reading request body")
		return
	}
	defer r.Body.Close()
	var request models.RegisterRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		handleRegisterErrorResponse(w, -1, "Request body is malformed")
		return
	}
	if storage.UserDoesExist(request.Name) {
		handleRegisterErrorResponse(w, -1, "Name is already used")
		return
	}
	user := models.GenerateUser(request.Name)
	err = storage.AddUser(user)
	if err != nil {
		handleRegisterErrorResponse(w, -1, "Unable to add user")
		return
	}
	authentication := generateAuthentication(user.Name)
	err = storage.AddAuthentication(authentication)
	if err != nil {
		handleRegisterErrorResponse(w, -1, "Unable to add authentication")
		return
	}
	response := models.RegisterResponse{
		Okay:              true,
		ErrorMessage:      "",
		ErrorCode:         0,
		AuthenticationKey: authentication.Key,
	}
	jsonData, err := json.Marshal(response)
	if err != nil {
		handleRegisterErrorResponse(w, -1, "Unable to generate response")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func handleRegisterErrorResponse(w http.ResponseWriter, code int, message string) {
	response := models.RegisterResponse{
		Okay:              false,
		ErrorCode:         code,
		ErrorMessage:      message,
		AuthenticationKey: "",
	}
	jsonData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
