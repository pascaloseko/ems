package employees

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type Store interface {
	GetEmployeeIdByUsername(ctx context.Context, username string) (int64, error)
	GetDepartmentIdByName(ctx context.Context, name string) (int64, error)
	GetDepartmentNameById(ctx context.Context, id int64) (string, error)
	GetAllEmployees(ctx context.Context) ([]Employee, error)
	Save(ctx context.Context, emp Employee) (int64, error)
	Authenticate(ctx context.Context, emp Employee) bool
	SaveDepartment(ctx context.Context, dept Department) (int64, error)
}

type EmployeeStore struct {
	store *sql.DB
}

func NewEmployeeStore(db *sql.DB) Store {
	return &EmployeeStore{
		store: db,
	}
}

// SaveDepartment implements Store.
func (e *EmployeeStore) SaveDepartment(ctx context.Context, d Department) (int64, error) {
	stmt, err := e.store.PrepareContext(ctx, "INSERT INTO Department(Name) VALUES(?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, d.Name)
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

// Authenticate implements Store.
func (e *EmployeeStore) Authenticate(ctx context.Context, user Employee) bool {
	row := e.store.QueryRowContext(ctx, "SELECT Password FROM Employees WHERE Username = ?", user.Username)
	var hashedPassword string
	err := row.Scan(&hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false
		}
		log.Fatal(err)
	}

	return CheckPasswordHash(user.Password, hashedPassword)
}

// GetAllEmployees implements Store.
func (e *EmployeeStore) GetAllEmployees(ctx context.Context) ([]Employee, error) {
	rows, err := e.store.QueryContext(ctx, "SELECT ID, FirstName, LastName, Username, Email, DOB, DepartmentID, Position FROM Employees")
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

// GetDepartmentIdByName implements Store.
func (e *EmployeeStore) GetDepartmentIdByName(ctx context.Context, name string) (int64, error) {
	row := e.store.QueryRowContext(ctx, "SELECT ID FROM Department WHERE Name = ?", name)
	var id int64
	err := row.Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}

	return id, nil
}

// GetDepartmentNameById implements Store.
func (e *EmployeeStore) GetDepartmentNameById(ctx context.Context, id int64) (string, error) {
	row := e.store.QueryRowContext(ctx, "SELECT Name FROM Department WHERE ID = ?", id)
	var name string
	err := row.Scan(&name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", err
	}

	return name, nil
}

// GetEmployeeIdByUsername implements Store.
func (e *EmployeeStore) GetEmployeeIdByUsername(ctx context.Context, username string) (int64, error) {
	row := e.store.QueryRowContext(ctx, "SELECT ID FROM Employees WHERE Username = ?", username)
	var id int64
	err := row.Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}

	return id, nil
}

// Save implements Store.
func (e *EmployeeStore) Save(ctx context.Context, emp Employee) (int64, error) {
	// Check if database is alive.
	err := e.store.PingContext(ctx)
	if err != nil {
		return -1, err
	}

	// if there is no department with the name provided go ahead and create the department
	departmentID, err := e.GetDepartmentIdByName(ctx, emp.DepartmentName)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		departmentName := Department{Name: emp.DepartmentName}
		departmentID, err = e.SaveDepartment(ctx, departmentName)
		if err != nil {
			return 0, err
		}
	} else if err != nil {
		return 0, err
	}

	tsql := `
	INSERT INTO EmployeesSchema.Employees(FirstName,LastName, Username, Password, Email, DOB, DepartmentID, Position) 
	VALUES(@FirstName,@LastName,@Username,@Password,@Email,@DOB,@DepartmentID,@Position)
	`

	stmt, err := e.store.Prepare(tsql)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, emp.FirstName, emp.LastName, emp.Username, emp.Password, emp.Email, emp.DOB, departmentID, emp.Position)
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
