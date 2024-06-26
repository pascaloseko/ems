// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Department struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Employee struct {
	ID           string `json:"id"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Email        string `json:"email"`
	Dob          string `json:"dob"`
	DepartmentID int    `json:"departmentID"`
	Position     string `json:"position"`
}

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Mutation struct {
}

type NewEmployee struct {
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Email        string `json:"email"`
	Dob          string `json:"dob"`
	DepartmentID int    `json:"departmentID"`
	Position     string `json:"position"`
}

type Query struct {
}

type RefreshTokenInput struct {
	Token string `json:"token"`
}
