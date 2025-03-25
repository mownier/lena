package jsonendpoints

import (
	"encoding/json"
	"io/ioutil"
	storage "lena/inmemory"
	"lena/models"
	"net/http"
)

func SignOut(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		handleSignOutErrorResponse(w, -1, "Method not allowed")
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		handleSignOutErrorResponse(w, -1, "Error reading request body")
		return
	}
	defer r.Body.Close()
	var request models.SignOutRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		handleSignOutErrorResponse(w, -1, "Request body is malformed")
		return
	}
	if storage.AuthenticationDoesNotExist(request.AuthenticationKey) {
		handleSignOutErrorResponse(w, -1, "Authentication does not exist")
		return
	}
	archived, err := storage.AuthenticationIsArchived(request.AuthenticationKey)
	if err != nil {
		handleSignOutErrorResponse(w, -1, "Problem determining validity of authentication")
		return
	}
	if archived {
		handleSignOutErrorResponse(w, -1, "Authentication is already invalidated")
		return
	}
	err = storage.ArchiveAuthentication(request.AuthenticationKey)
	if err != nil {
		handleSignOutErrorResponse(w, -1, "Problem invalidating authentication")
		return
	}
	response := models.SignOutResponse{
		Okay:         true,
		ErrorCode:    0,
		ErrorMessage: "",
	}
	jsonData, err := json.Marshal(response)
	if err != nil {
		handleSignOutErrorResponse(w, -1, "Problem generating response")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func handleSignOutErrorResponse(w http.ResponseWriter, code int, message string) {
	response := models.SignOutResponse{
		Okay:         false,
		ErrorCode:    code,
		ErrorMessage: message,
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
