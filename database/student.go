package database

import (
	"context"
	"database/sql"
	"fmt"

	"goProject-2024/models"

	uuid "github.com/satori/go.uuid"
)

// student model
type StudentRow struct {
	ID          string
	FName       sql.NullString
	LName       sql.NullString
	DateOfBirth sql.NullTime
	Email       sql.NullString
	Address     sql.NullString
	Gender      sql.NullString
	Age         sql.NullInt32
	CreatedBy   sql.NullString
	CreatedOn   sql.NullTime
	UpdatedBy   sql.NullString
	UpdatedOn   sql.NullTime
}

func convertStudentRowToStudent(s StudentRow) models.Student {
	return models.Student{
		ID:          s.ID,
		FName:       s.FName.String,
		LName:       s.LName.String,
		DateOfBirth: s.DateOfBirth.Time,
		Email:       s.Email.String,
		Address:     s.Address.String,
		Gender:      s.Gender.String,
		Age:         int(s.Age.Int32),
		CreatedBy:   s.CreatedBy.String,
		CreatedOn:   s.CreatedOn.Time,
		UpdatedBy:   s.UpdatedBy.String,
		UpdatedOn:   s.UpdatedOn.Time,
	}
}

func (d *Database) GetStudent(ctx context.Context, id string) (models.Student, error) {
	var studentRow StudentRow
	row := d.Client.QueryRowContext(
		ctx,
		`SELECT id, f_name, l_name, date_of_birth, email, address, gender, age, 
		created_by, created_on, updated_by, updated_on 
		FROM students 
		WHERE id = ?`,
		id,
	)
	err := row.Scan(
		&studentRow.ID, &studentRow.FName, &studentRow.LName, &studentRow.DateOfBirth,
		&studentRow.Email, &studentRow.Address, &studentRow.Gender, &studentRow.Age,
		&studentRow.CreatedBy, &studentRow.CreatedOn, &studentRow.UpdatedBy, &studentRow.UpdatedOn,
	)
	if err != nil {
		return models.Student{}, fmt.Errorf("an error occurred fetching a student by ID: %w", err)
	}
	return convertStudentRowToStudent(studentRow), nil
}

func (d *Database) PostStudent(ctx context.Context, student models.Student) (models.Student, error) {
	student.ID = uuid.NewV4().String()
	studentRow := StudentRow{
		ID:          student.ID,
		FName:       sql.NullString{String: student.FName, Valid: true},
		LName:       sql.NullString{String: student.LName, Valid: true},
		DateOfBirth: sql.NullTime{Time: student.DateOfBirth, Valid: true},
		Email:       sql.NullString{String: student.Email, Valid: true},
		Address:     sql.NullString{String: student.Address, Valid: true},
		Gender:      sql.NullString{String: student.Gender, Valid: true},
		Age:         sql.NullInt32{Int32: int32(student.Age), Valid: true},
		CreatedBy:   sql.NullString{String: student.CreatedBy, Valid: true},
		CreatedOn:   sql.NullTime{Time: student.CreatedOn, Valid: true},
	}

	_, err := d.Client.ExecContext(
		ctx,
		`INSERT INTO students 
		(id, f_name, l_name, date_of_birth, email, address, gender, age, 
		created_by, created_on) VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		studentRow.ID, studentRow.FName, studentRow.LName, studentRow.DateOfBirth,
		studentRow.Email, studentRow.Address, studentRow.Gender, studentRow.Age,
		studentRow.CreatedBy, studentRow.CreatedOn,
	)
	if err != nil {
		return models.Student{}, fmt.Errorf("failed to insert student: %w", err)
	}
	return student, nil
}

func (d *Database) UpdateStudent(ctx context.Context, id string, student models.Student) (models.Student, error) {
	studentRow := StudentRow{
		ID:          id,
		FName:       sql.NullString{String: student.FName, Valid: true},
		LName:       sql.NullString{String: student.LName, Valid: true},
		DateOfBirth: sql.NullTime{Time: student.DateOfBirth, Valid: true},
		Email:       sql.NullString{String: student.Email, Valid: true},
		Address:     sql.NullString{String: student.Address, Valid: true},
		Gender:      sql.NullString{String: student.Gender, Valid: true},
		Age:         sql.NullInt32{Int32: int32(student.Age), Valid: true},
		UpdatedBy:   sql.NullString{String: student.UpdatedBy, Valid: true},
		UpdatedOn:   sql.NullTime{Time: student.UpdatedOn, Valid: true},
	}

	_, err := d.Client.ExecContext(
		ctx,
		`UPDATE students SET
		f_name = ?, l_name = ?, date_of_birth = ?,
		email = ?, address = ?, gender = ?, age = ?,
		updated_by = ?, updated_on = ?
		WHERE id = ?`,
		studentRow.FName, studentRow.LName, studentRow.DateOfBirth,
		studentRow.Email, studentRow.Address, studentRow.Gender, studentRow.Age,
		studentRow.UpdatedBy, studentRow.UpdatedOn,
		studentRow.ID,
	)
	if err != nil {
		return models.Student{}, fmt.Errorf("failed to update student: %w", err)
	}

	return convertStudentRowToStudent(studentRow), nil
}

func (d *Database) DeleteStudent(ctx context.Context, id string) error {
	_, err := d.Client.ExecContext(
		ctx,
		`DELETE FROM students where id = ?`,
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to delete student from the database: %w", err)
	}
	return nil
}
