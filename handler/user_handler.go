package handler

import (
	"encoding/json"
	"go-api/dto"
	"go-api/service"

	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	UserService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{UserService: userService}
}

// helper to send JSON responses
func jsonResponse(w http.ResponseWriter, status int, resp dto.APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(resp)
}

// GET /users
func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.UserService.GetAllUsers()
	if err != nil {
		jsonResponse(w, 500, dto.APIResponse{Success: false, Message: err.Error()})
		return
	}
	jsonResponse(w, 200, dto.APIResponse{Success: true, Message: "Users fetched", Data: users})
}

// GET /users/{id}
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		jsonResponse(w, 400, dto.APIResponse{Success: false, Message: "Invalid ID"})
		return
	}

	user, err := h.UserService.GetUserByID(id)
	if err != nil {
		jsonResponse(w, 404, dto.APIResponse{Success: false, Message: err.Error()})
		return
	}
	jsonResponse(w, 200, dto.APIResponse{Success: true, Message: "User fetched", Data: user})
}

// POST /users
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonResponse(w, 400, dto.APIResponse{Success: false, Message: "Invalid request body"})
		return
	}

	user, err := h.UserService.CreateUser(req)
	if err != nil {
		jsonResponse(w, 400, dto.APIResponse{Success: false, Message: err.Error()})
		return
	}
	jsonResponse(w, 201, dto.APIResponse{Success: true, Message: "User created", Data: user})
}

// PUT /users/{id}
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		jsonResponse(w, 400, dto.APIResponse{Success: false, Message: "Invalid ID"})
		return
	}

	var req dto.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonResponse(w, 400, dto.APIResponse{Success: false, Message: "Invalid request body"})
		return
	}

	user, err := h.UserService.UpdateUser(id, req)
	if err != nil {
		jsonResponse(w, 400, dto.APIResponse{Success: false, Message: err.Error()})
		return
	}
	jsonResponse(w, 200, dto.APIResponse{Success: true, Message: "User updated", Data: user})
}

// DELETE /users/{id}
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		jsonResponse(w, 400, dto.APIResponse{Success: false, Message: "Invalid ID"})
		return
	}

	if err := h.UserService.DeleteUser(id); err != nil {
		jsonResponse(w, 404, dto.APIResponse{Success: false, Message: err.Error()})
		return
	}
	jsonResponse(w, 200, dto.APIResponse{Success: true, Message: "User deleted"})
}
