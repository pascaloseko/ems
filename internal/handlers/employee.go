package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/pascaloseko/ems/graph"
	"github.com/pascaloseko/ems/graph/model"
)

// LoginHandler handles employee authentication
func LoginHandler(w http.ResponseWriter, r *http.Request) {
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
		FirstName:    credentials.FirstName,
		LastName:     credentials.LastName,
		Email:        credentials.Email,
		Username:     credentials.Username,
		Password:     credentials.Password,
		Dob:          credentials.Dob,
		DepartmentID: credentials.DepartmentID,
		Position:     credentials.Position,
	}

	resolver := graph.Resolver{}
	token, err := resolver.Mutation().CreateEmployee(r.Context(), newEmployee)
	if err != nil {
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
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
func GetAllEmployeesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	resolver := graph.Resolver{}
	employees, err := resolver.Query().Employees(r.Context())
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Return the employees
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(employees)

}
