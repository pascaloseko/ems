package employees

type Employee struct {
	ID             int64  `json:"id"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	Email          string `json:"email"`
	DOB            string `json:"dob"`
	DepartmentName string `json:"department_name"`
	DepartmentID   int64  `json:"department_id"`
	Position       string `json:"position"`
}

type Department struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

