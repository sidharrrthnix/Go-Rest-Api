package sqlconnect

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"http-api.com/models"
	"http-api.com/utils"
)

// ============================================================================
// HELPERS
// ============================================================================

func getStudentFieldMapping(student *models.Student) map[string]*string {
	return map[string]*string{
		"firstName": &student.FirstName,
		"lastName":  &student.LastName,
		"email":     &student.Email,
		"class":     &student.Class,
	}
}

func applyStudentUpdate(student *models.Student, updates map[string]interface{}) {
	fieldMap := getStudentFieldMapping(student)
	for key, value := range updates {
		if key == "id" {
			continue
		}
		if fieldPtr, ok := fieldMap[key]; ok {
			if strVal, ok := value.(string); ok {
				*fieldPtr = strVal
			}
		}
	}
}

func extractStudentID(update map[string]interface{}, index int) (int, error) {
	idVal, exists := update["id"]
	if !exists {
		return 0, utils.ErrorHandler(fmt.Errorf("missing id"), fmt.Sprintf("missing id in update #%d", index))
	}
	switch v := idVal.(type) {
	case float64:
		return int(v), nil
	case int:
		return v, nil
	case string:
		id, err := strconv.Atoi(v)
		if err != nil {
			return 0, utils.ErrorHandler(err, fmt.Sprintf("invalid id in update #%d", index))
		}
		return id, nil
	default:
		return 0, utils.ErrorHandler(fmt.Errorf("invalid id type"), fmt.Sprintf("invalid id type in update #%d", index))
	}
}

// ============================================================================
// READ OPERATIONS
// ============================================================================

func GetStudentsDbHandler(students []models.Student, r *http.Request) ([]models.Student, error) {
	db, err := ConnectDB()
	if err != nil {
		return nil, utils.ErrorHandler(err, "database connection failed")
	}

	query := "SELECT id, first_name, last_name, email, class FROM students WHERE 1=1"
	args := []interface{}{}

	// Add filters
	filters := map[string]string{
		"first_name": "first_name",
		"last_name":  "last_name",
		"email":      "email",
		"class":      "class",
	}
	for param, dbField := range filters {
		if value := r.URL.Query().Get(param); value != "" {
			query += " AND " + dbField + " = ?"
			args = append(args, value)
		}
	}

	// Add sorting
	if sortParams := r.URL.Query()["sortby"]; len(sortParams) > 0 {
		validFields := map[string]bool{
			"first_name": true, "last_name": true,
			"email": true, "class": true,
		}
		orderBys := []string{}

		for _, param := range sortParams {
			parts := strings.Split(param, ":")
			if len(parts) == 2 {
				field, order := parts[0], strings.ToUpper(parts[1])
				if validFields[field] && (order == "ASC" || order == "DESC") {
					orderBys = append(orderBys, field+" "+order)
				}
			}
		}
		if len(orderBys) > 0 {
			query += " ORDER BY " + strings.Join(orderBys, ", ")
		}
	}

	// Execute query
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, utils.ErrorHandler(err, "query failed")
	}
	defer rows.Close()

	// Scan results
	for rows.Next() {
		var student models.Student
		if err := rows.Scan(&student.ID, &student.FirstName, &student.LastName, &student.Email, &student.Class); err != nil {
			return nil, utils.ErrorHandler(err, "scan failed")
		}
		students = append(students, student)
	}
	return students, nil
}

func GetStudentByID(id int) (models.Student, error) {
	db, err := ConnectDB()
	if err != nil {
		return models.Student{}, utils.ErrorHandler(err, "database connection failed")
	}

	var student models.Student
	query := "SELECT id, first_name, last_name, email, class FROM students WHERE id = ?"
	err = db.QueryRow(query, id).Scan(&student.ID, &student.FirstName, &student.LastName, &student.Email, &student.Class)

	if err == sql.ErrNoRows {
		return models.Student{}, utils.ErrorHandler(err, "student not found")
	}
	if err != nil {
		return models.Student{}, utils.ErrorHandler(err, "query failed")
	}
	return student, nil
}

// ============================================================================
// CREATE OPERATIONS
// ============================================================================

func AddStudentsDBHandler(newStudents []models.Student) ([]models.Student, error) {
	db, err := ConnectDB()
	if err != nil {
		return nil, utils.ErrorHandler(err, "database connection failed")
	}

	stmt, err := db.Prepare("INSERT INTO students (first_name, last_name, email, class) VALUES (?, ?, ?, ?)")
	if err != nil {
		return nil, utils.ErrorHandler(err, "prepare statement failed")
	}
	defer stmt.Close()

	addedStudents := make([]models.Student, 0, len(newStudents))
	for _, student := range newStudents {
		res, err := stmt.Exec(student.FirstName, student.LastName, student.Email, student.Class)
		if err != nil {
			return nil, utils.ErrorHandler(err, "insert failed")
		}
		lastId, err := res.LastInsertId()
		if err != nil {
			return nil, utils.ErrorHandler(err, "failed to get insert ID")
		}
		student.ID = int(lastId)
		addedStudents = append(addedStudents, student)
	}
	return addedStudents, nil
}

// ============================================================================
// UPDATE OPERATIONS
// ============================================================================

func UpdateStudent(id int, updatedStudent models.Student) (models.Student, error) {
	db, err := ConnectDB()
	if err != nil {
		return models.Student{}, utils.ErrorHandler(err, "database connection failed")
	}

	// Verify student exists
	var exists int
	err = db.QueryRow("SELECT 1 FROM students WHERE id = ?", id).Scan(&exists)
	if err == sql.ErrNoRows {
		return models.Student{}, utils.ErrorHandler(err, "student not found")
	}
	if err != nil {
		return models.Student{}, utils.ErrorHandler(err, "query failed")
	}

	// Update student
	updatedStudent.ID = id
	_, err = db.Exec("UPDATE students SET first_name=?, last_name=?, email=?, class=? WHERE id=?",
		updatedStudent.FirstName, updatedStudent.LastName, updatedStudent.Email, updatedStudent.Class, id)
	if err != nil {
		return models.Student{}, utils.ErrorHandler(err, "update failed")
	}
	return updatedStudent, nil
}

func PatchOneStudent(id int, updates map[string]interface{}) (models.Student, error) {
	db, err := ConnectDB()
	if err != nil {
		return models.Student{}, utils.ErrorHandler(err, "database connection failed")
	}

	// Get current student
	student, err := GetStudentByID(id)
	if err != nil {
		return models.Student{}, err
	}

	// Apply updates
	applyStudentUpdate(&student, updates)

	// Save to database
	_, err = db.Exec("UPDATE students SET first_name=?, last_name=?, email=?, class=? WHERE id=?",
		student.FirstName, student.LastName, student.Email, student.Class, student.ID)
	if err != nil {
		return models.Student{}, utils.ErrorHandler(err, "update failed")
	}
	return student, nil
}

func PatchStudents(updates []map[string]interface{}) ([]models.Student, error) {
	db, err := ConnectDB()
	if err != nil {
		return nil, utils.ErrorHandler(err, "database connection failed")
	}

	tx, err := db.Begin()
	if err != nil {
		return nil, utils.ErrorHandler(err, "transaction start failed")
	}

	updatedStudents := []models.Student{}

	for i, update := range updates {
		// Extract and validate ID
		id, err := extractStudentID(update, i+1)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		// Get existing student
		var student models.Student
		query := "SELECT id, first_name, last_name, email, class FROM students WHERE id = ?"
		err = tx.QueryRow(query, id).Scan(&student.ID, &student.FirstName, &student.LastName, &student.Email, &student.Class)

		if err == sql.ErrNoRows {
			tx.Rollback()
			return nil, utils.ErrorHandler(err, fmt.Sprintf("student #%d (ID %d) not found", i+1, id))
		}
		if err != nil {
			tx.Rollback()
			return nil, utils.ErrorHandler(err, fmt.Sprintf("query failed for student #%d", i+1))
		}

		// Apply updates
		applyStudentUpdate(&student, update)

		// Update in database
		_, err = tx.Exec("UPDATE students SET first_name=?, last_name=?, email=?, class=? WHERE id=?",
			student.FirstName, student.LastName, student.Email, student.Class, student.ID)
		if err != nil {
			tx.Rollback()
			return nil, utils.ErrorHandler(err, fmt.Sprintf("update failed for student #%d", i+1))
		}

		updatedStudents = append(updatedStudents, student)
	}

	if err := tx.Commit(); err != nil {
		return nil, utils.ErrorHandler(err, "transaction commit failed")
	}
	return updatedStudents, nil
}

// ============================================================================
// DELETE OPERATIONS
// ============================================================================

func DeleteOneStudent(id int) error {
	db, err := ConnectDB()
	if err != nil {
		return utils.ErrorHandler(err, "database connection failed")
	}

	result, err := db.Exec("DELETE FROM students WHERE id = ?", id)
	if err != nil {
		return utils.ErrorHandler(err, "delete failed")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return utils.ErrorHandler(err, "failed to get rows affected")
	}
	if rowsAffected == 0 {
		return utils.ErrorHandler(fmt.Errorf("no rows affected"), "student not found")
	}
	return nil
}

func DeleteStudents(ids []int) ([]int, error) {
	db, err := ConnectDB()
	if err != nil {
		return nil, utils.ErrorHandler(err, "database connection failed")
	}

	tx, err := db.Begin()
	if err != nil {
		return nil, utils.ErrorHandler(err, "transaction start failed")
	}

	stmt, err := tx.Prepare("DELETE FROM students WHERE id = ?")
	if err != nil {
		tx.Rollback()
		return nil, utils.ErrorHandler(err, "prepare statement failed")
	}
	defer stmt.Close()

	deletedIds := []int{}
	for _, id := range ids {
		result, err := stmt.Exec(id)
		if err != nil {
			tx.Rollback()
			return nil, utils.ErrorHandler(err, fmt.Sprintf("delete failed for ID %d", id)) // ✅ Fixed
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			tx.Rollback()
			return nil, utils.ErrorHandler(err, "failed to get rows affected")
		}
		if rowsAffected == 0 {
			tx.Rollback()
			return nil, utils.ErrorHandler(fmt.Errorf("not found"), fmt.Sprintf("student with ID %d not found", id))
		}
		deletedIds = append(deletedIds, id)
	}

	if err := tx.Commit(); err != nil {
		return nil, utils.ErrorHandler(err, "transaction commit failed")
	}
	return deletedIds, nil
}
