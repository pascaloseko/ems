package graph

import "github.com/pascaloseko/ems/internal/employees"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct{
	emp employees.Store
}


func NewResolver(emp employees.Store) *Resolver {
	return &Resolver{
		emp: emp,
	}
}