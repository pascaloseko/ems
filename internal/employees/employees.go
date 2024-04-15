package employees

import (
	"database/sql"
	"errors"
	"log"

	"github.com/pascaloseko/ems/internal/pkg/db/database"
	"golang.org/x/crypto/bcrypt"
)

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

// GetEmployeeIdByUsername check if a user exists in database by given username
func GetEmployeeIdByUsername(username string) (int64, error) {
	statement, err := database.Db.Prepare("select ID from Employees WHERE Username = ?")
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	// Execute the statement
	row := statement.QueryRow(username)

	var Id int64
	err = row.Scan(&Id)
	if err != nil {
		if err != sql.ErrNoRows {
			return 0, err
		}
		return 0, err
	}

	return Id, nil
}

func (d Department) Save() (int64, error) {
	stmt, err := database.Db.Prepare("INSERT INTO Department(Name) VALUES(?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(d.Name)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	log.Print("Row inserted!")
	return id, nil
}

func GetDepartmentIdByName(name string) (int64, error) {
	statement, err := database.Db.Prepare("select ID from Department WHERE Name = ?")
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	// Execute the statement
	row := statement.QueryRow(name)

	var Id int64
	err = row.Scan(&Id)
	if err != nil {
		if err != sql.ErrNoRows {
			return 0, err
		}
		return 0, err
	}

	return Id, nil
}

func GetDepartmentNameById(id int64) (string, error) {
	statement, err := database.Db.Prepare("select Name from Department WHERE ID = ?")
	if err != nil {
		return "", err
	}
	defer statement.Close()

	// Execute the statement
	row := statement.QueryRow(id)

	var name string
	err = row.Scan(&name)
	if err != nil {
		if err != sql.ErrNoRows {
			return "", err
		}
		return "", err
	}

	return name, nil
}

func (e Employee) Save() (int64, error) {
	departmentID, err := GetDepartmentIdByName(e.DepartmentName)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		departmentName := &Department{Name: e.DepartmentName}
		departmentID, err = departmentName.Save()
		if err != nil {
			return 0, err
		}
	} else if err != nil {
		return 0, err
	}

	stmt, err := database.Db.Prepare("INSERT INTO Employees(FirstName,LastName, Username, Password, Email, DOB, DepartmentID, Position) VALUES(?,?,?,?,?,?,?,?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(e.FirstName, e.LastName, e.Username, e.Password, e.Email, e.DOB, departmentID, e.Position)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	log.Print("Row inserted!")
	return id, nil
}

func GetAllEmployees() ([]Employee, error) {
	stmt, err := database.Db.Prepare("select ID, FirstName, LastName, Username, Email, DOB, DepartmentID, Position from Employees")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var employees []Employee
	for rows.Next() {
		var employee Employee
		err = rows.Scan(&employee.ID, &employee.FirstName, &employee.LastName, &employee.Username, &employee.Email, &employee.DOB, &employee.DepartmentID, &employee.Position)
		if err != nil {
			return nil, err
		}
		employees = append(employees, employee)
	}
	return employees, nil
}

func (user *Employee) Authenticate() bool {
	statement, err := database.Db.Prepare("select Password from Employees WHERE Username = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer statement.Close()

	// Execute the statement
	row := statement.QueryRow(user.Username)
	var hashedPassword string
	err = row.Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		} else {
			log.Fatal(err)
		}
	}

	return CheckPasswordHash(user.Password, hashedPassword)
}

// HashPassword hashes given password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPassword hash compares raw password with it's hashed values
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
