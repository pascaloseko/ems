package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/pascaloseko/ems/graph"
	"github.com/pascaloseko/ems/graph/model"
)

type Handlers struct {
	resolver *graph.Resolver
}

func NewHandlers(resolver *graph.Resolver) *Handlers {
	return &Handlers{
		resolver: resolver,
	}
}

// LoginHandler handles employee authentication
func (h *Handlers) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var credentials model.Employee
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	newEmployee := model.NewEmployee{
		FirstName: credentials.FirstName,
		LastName:  credentials.LastName,
		Email:     credentials.Email,
		Username:  credentials.Username,
		Password:  credentials.Password,
		Dob:       credentials.Dob,
		DepartmentID: credentials.DepartmentID,
		Position: credentials.Position,
	}

	token, err := h.resolver.Mutation().CreateEmployee(r.Context(), newEmployee)
	if err != nil {
		log.Println("ERROR", err)
		http.Error(w, "Failed to create employee", http.StatusInternalServerError)
		return
	}

	if token == nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Return the JWT token
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": *token})
}

// GetEmployees handles employees
func (h *Handlers) GetAllEmployeesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	employees, err := h.resolver.Query().Employees(r.Context())
	if err != nil {
		if errors.Is(err, graph.ErrAccessDenied) {
			http.Error(w, "access denied", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Error fetching employees", http.StatusInternalServerError)
		return
	}

	if employees == nil {
		http.Error(w, "No employees found", http.StatusNotFound)
		return
	}

	// Return the employees
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(employees)

}
