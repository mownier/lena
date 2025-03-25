package jsonendpoints

import (
	"encoding/json"
	"io/ioutil"
	storage "lena/inmemory"
	"lena/models"
	"net/http"
)

func SignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		handleSignInErrorResponse(w, -1, "Method not allowed")
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		handleSignInErrorResponse(w, -1, "Error reading request body")
		return
	}
	defer r.Body.Close()
	var request models.SignInRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		handleSignInErrorResponse(w, -1, "Request body is malformed")
		return
	}
	if storage.UserDoesNotExist(request.Name) {
		handleSignInErrorResponse(w, -1, "User does not exist")
		return
	}
	if storage.UserDoesHaveAuthentication(request.Name) {
		handleSignInErrorResponse(w, -1, "User does have already an authentication")
		return
	}
	authentication := generateAuthentication(request.Name)
	err = storage.AddAuthentication(authentication)
	if err != nil {
		handleSignInErrorResponse(w, -1, "Unable to add authentication")
		return
	}
	response := models.SignInResponse{
		Okay:              true,
		ErrorCode:         0,
		ErrorMessage:      "",
		AuthenticationKey: authentication.Key,
	}
	jsonData, err := json.Marshal(response)
	if err != nil {
		handleSignInErrorResponse(w, -1, "Unable to generate response")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func handleSignInErrorResponse(w http.ResponseWriter, code int, message string) {
	response := models.SignInResponse{
		Okay:              false,
		ErrorCode:         code,
		ErrorMessage:      message,
		AuthenticationKey: "",
	}
	jsonData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
