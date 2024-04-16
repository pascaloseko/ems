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
	HashPassword(password string) string
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
	tsql := `
	INSERT INTO Department_Entities(Name) VALUES(@Name);
	SELECT ID = convert(bigint, SCOPE_IDENTITY());
	`
	stmt, err := e.store.Prepare(tsql)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	row := stmt.QueryRowContext(ctx, sql.Named("Name", d.Name))
	if err != nil {
		return 0, err
	}
	var newID int64
	err = row.Scan(&newID)
	if err != nil {
		return 0, err
	}
	log.Print("Row inserted!")
	return newID, nil
}

// Authenticate implements Store.
func (e *EmployeeStore) Authenticate(ctx context.Context, user Employee) bool {
	row := e.store.QueryRowContext(ctx, "SELECT Password FROM Employee_Entities WHERE Username = @Username", sql.Named("Username", user.Username))
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
	tsql := `
	SELECT ID FROM Department_Entities WHERE Name = @Name;
	`
	row := e.store.QueryRowContext(ctx, tsql, sql.Named("Name", name))
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
	row := e.store.QueryRowContext(ctx, "SELECT ID FROM Employee_Entities WHERE Username = @Username", sql.Named("Username", username))
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
	// if there is no department with the name provided go ahead and create the department
	departmentID, err := e.GetDepartmentIdByName(ctx, emp.DepartmentName)
	if errors.Is(err, sql.ErrNoRows) {
		departmentName := Department{Name: emp.DepartmentName}
		departmentID, err = e.SaveDepartment(ctx, departmentName)
		if err != nil {
			return 0, err
		}
	} else if err != nil {
		return 0, err
	}

	tsql := `
	INSERT INTO Employee_Entities (First_Name, Last_Name, Username, Password, Email, DOB, Department_Id, Position)
	VALUES (@First_Name, @Last_Name, @Username, @Password, @Email, @DOB, @Department_Id, @Position);
	SELECT ID = convert(bigint, SCOPE_IDENTITY());
	`

	stmt, err := e.store.Prepare(tsql)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	row := stmt.QueryRowContext(
		ctx,
		sql.Named("First_Name", emp.FirstName),
		sql.Named("Last_Name", emp.LastName),
		sql.Named("Username", emp.Username),
		sql.Named("Password", emp.Password),
		sql.Named("Email", emp.Email),
		sql.Named("DOB", emp.DOB),
		sql.Named("Department_Id", departmentID),
		sql.Named("Position", emp.Position))
	if err != nil {
		return 0, err
	}
	var newID int64
	err = row.Scan(&newID)
	if err != nil {
		return 0, err
	}
	log.Print("Row inserted!")
	return newID, nil
}

// HashPassword hashes given password
func (e *EmployeeStore) HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Fatal(err)
	}
	return string(bytes)
}

// CheckPassword hash compares raw password with it's hashed values
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
