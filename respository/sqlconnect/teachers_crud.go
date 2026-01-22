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

// getTeacherFieldMapping returns the JSON-to-struct field mapping for Teacher
func getTeacherFieldMapping(teacher *models.Teacher) map[string]*string {
	return map[string]*string{
		"firstName": &teacher.FirstName,
		"lastName":  &teacher.LastName,
		"email":     &teacher.Email,
		"class":     &teacher.Class,
		"subject":   &teacher.Subject,
	}
}

// applyTeacherUpdates applies map updates to teacher fields
func applyTeacherUpdates(teacher *models.Teacher, updates map[string]interface{}) {
	fieldMap := getTeacherFieldMapping(teacher)
	for key, value := range updates {
		if key == "id" {
			continue // Skip ID field
		}
		if fieldPtr, ok := fieldMap[key]; ok {
			if strVal, ok := value.(string); ok {
				*fieldPtr = strVal
			}
		}
	}
}

// ============================================================================
// READ OPERATIONS
// ============================================================================

// GetTeachersDbHandler retrieves all teachers with optional filters and sorting
func GetTeachersDbHandler(teachers []models.Teacher, r *http.Request) ([]models.Teacher, error) {
	db, err := ConnectDB()
	if err != nil {
		return nil, utils.ErrorHandler(err, "database connection failed")
	}

	query := "SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE 1=1"
	args := []interface{}{}

	// Add filters
	filters := map[string]string{
		"first_name": "first_name",
		"last_name":  "last_name",
		"email":      "email",
		"class":      "class",
		"subject":    "subject",
	}
	for param, dbField := range filters {
		if value := r.URL.Query().Get(param); value != "" {
			query += " AND " + dbField + " = ?"
			args = append(args, value)
		}
	}

	// Add sorting (format: ?sortby=first_name:ASC&sortby=class:DESC)
	if sortParams := r.URL.Query()["sortby"]; len(sortParams) > 0 {
		validFields := map[string]bool{
			"first_name": true, "last_name": true, "email": true,
			"class": true, "subject": true,
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
		var teacher models.Teacher
		if err := rows.Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName,
			&teacher.Email, &teacher.Class, &teacher.Subject); err != nil {
			return nil, utils.ErrorHandler(err, "scan failed")
		}
		teachers = append(teachers, teacher)
	}
	return teachers, nil
}

// GetTeacherByID retrieves a single teacher by ID
func GetTeacherByID(id int) (models.Teacher, error) {
	db, err := ConnectDB()
	if err != nil {
		return models.Teacher{}, utils.ErrorHandler(err, "database connection failed")
	}

	var teacher models.Teacher
	query := "SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?"
	err = db.QueryRow(query, id).Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName,
		&teacher.Email, &teacher.Class, &teacher.Subject)

	if err == sql.ErrNoRows {
		return models.Teacher{}, utils.ErrorHandler(err, "teacher not found")
	}
	if err != nil {
		return models.Teacher{}, utils.ErrorHandler(err, "query failed")
	}
	return teacher, nil
}

// ============================================================================
// CREATE OPERATIONS
// ============================================================================

// AddTeachersDBHandler adds multiple teachers in batch
func AddTeachersDBHandler(newTeachers []models.Teacher) ([]models.Teacher, error) {
	db, err := ConnectDB()
	if err != nil {
		return nil, utils.ErrorHandler(err, "database connection failed")
	}

	stmt, err := db.Prepare("INSERT INTO teachers (first_name, last_name, email, class, subject) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return nil, utils.ErrorHandler(err, "prepare statement failed")
	}
	defer stmt.Close()

	addedTeachers := make([]models.Teacher, 0, len(newTeachers))
	for _, teacher := range newTeachers {
		res, err := stmt.Exec(teacher.FirstName, teacher.LastName, teacher.Email,
			teacher.Class, teacher.Subject)
		if err != nil {
			return nil, utils.ErrorHandler(err, "insert failed")
		}

		lastID, err := res.LastInsertId()
		if err != nil {
			return nil, utils.ErrorHandler(err, "failed to get insert ID")
		}
		teacher.ID = int(lastID)
		addedTeachers = append(addedTeachers, teacher)
	}
	return addedTeachers, nil
}

// ============================================================================
// UPDATE OPERATIONS
// ============================================================================

// UpdateTeacher performs a full update of a teacher
func UpdateTeacher(id int, updatedTeacher models.Teacher) (models.Teacher, error) {
	db, err := ConnectDB()
	if err != nil {
		return models.Teacher{}, utils.ErrorHandler(err, "database connection failed")
	}

	// Verify teacher exists
	var exists int
	err = db.QueryRow("SELECT 1 FROM teachers WHERE id = ?", id).Scan(&exists)
	if err == sql.ErrNoRows {
		return models.Teacher{}, utils.ErrorHandler(err, "teacher not found")
	}
	if err != nil {
		return models.Teacher{}, utils.ErrorHandler(err, "query failed")
	}

	// Update teacher
	updatedTeacher.ID = id
	_, err = db.Exec("UPDATE teachers SET first_name=?, last_name=?, email=?, class=?, subject=? WHERE id=?",
		updatedTeacher.FirstName, updatedTeacher.LastName, updatedTeacher.Email,
		updatedTeacher.Class, updatedTeacher.Subject, id)
	if err != nil {
		return models.Teacher{}, utils.ErrorHandler(err, "update failed")
	}
	return updatedTeacher, nil
}

// PatchOneTeacher performs a partial update of a single teacher
func PatchOneTeacher(id int, updates map[string]interface{}) (models.Teacher, error) {
	db, err := ConnectDB()
	if err != nil {
		return models.Teacher{}, utils.ErrorHandler(err, "database connection failed")
	}

	// Get current teacher
	teacher, err := GetTeacherByID(id)
	if err != nil {
		return models.Teacher{}, err
	}

	// Apply updates
	applyTeacherUpdates(&teacher, updates)

	// Save to database
	_, err = db.Exec("UPDATE teachers SET first_name=?, last_name=?, email=?, class=?, subject=? WHERE id=?",
		teacher.FirstName, teacher.LastName, teacher.Email, teacher.Class, teacher.Subject, teacher.ID)
	if err != nil {
		return models.Teacher{}, utils.ErrorHandler(err, "update failed")
	}
	return teacher, nil
}

// PatchTeachers performs batch partial updates in a transaction
func PatchTeachers(updates []map[string]interface{}) ([]models.Teacher, error) {
	db, err := ConnectDB()
	if err != nil {
		return nil, utils.ErrorHandler(err, "database connection failed")
	}

	tx, err := db.Begin()
	if err != nil {
		return nil, utils.ErrorHandler(err, "transaction start failed")
	}

	updatedTeachers := []models.Teacher{}

	for i, update := range updates {
		// Extract and validate ID
		id, err := extractTeacherID(update, i+1)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		// Get existing teacher
		var teacher models.Teacher
		query := "SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?"
		err = tx.QueryRow(query, id).Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName,
			&teacher.Email, &teacher.Class, &teacher.Subject)

		if err == sql.ErrNoRows {
			tx.Rollback()
			return nil, utils.ErrorHandler(err, fmt.Sprintf("teacher #%d (ID %d) not found", i+1, id))
		}
		if err != nil {
			tx.Rollback()
			return nil, utils.ErrorHandler(err, fmt.Sprintf("query failed for teacher #%d", i+1))
		}

		// Apply updates
		applyTeacherUpdates(&teacher, update)

		// Update in database
		_, err = tx.Exec("UPDATE teachers SET first_name=?, last_name=?, email=?, class=?, subject=? WHERE id=?",
			teacher.FirstName, teacher.LastName, teacher.Email, teacher.Class, teacher.Subject, teacher.ID)
		if err != nil {
			tx.Rollback()
			return nil, utils.ErrorHandler(err, fmt.Sprintf("update failed for teacher #%d", i+1))
		}

		updatedTeachers = append(updatedTeachers, teacher)
	}

	if err := tx.Commit(); err != nil {
		return nil, utils.ErrorHandler(err, "transaction commit failed")
	}
	return updatedTeachers, nil
}

// extractTeacherID extracts and validates ID from update map
func extractTeacherID(update map[string]interface{}, index int) (int, error) {
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
// DELETE OPERATIONS
// ============================================================================

// DeleteOneTeacher deletes a single teacher by ID
func DeleteOneTeacher(id int) error {
	db, err := ConnectDB()
	if err != nil {
		return utils.ErrorHandler(err, "database connection failed")
	}

	result, err := db.Exec("DELETE FROM teachers WHERE id = ?", id)
	if err != nil {
		return utils.ErrorHandler(err, "delete failed")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return utils.ErrorHandler(err, "failed to get rows affected")
	}
	if rowsAffected == 0 {
		return utils.ErrorHandler(fmt.Errorf("no rows affected"), "teacher not found")
	}
	return nil
}

// DeleteTeachers deletes multiple teachers in a transaction
func DeleteTeachers(ids []int) ([]int, error) {
	db, err := ConnectDB()
	if err != nil {
		return nil, utils.ErrorHandler(err, "database connection failed")
	}

	tx, err := db.Begin()
	if err != nil {
		return nil, utils.ErrorHandler(err, "transaction start failed")
	}

	stmt, err := tx.Prepare("DELETE FROM teachers WHERE id = ?")
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
			return nil, utils.ErrorHandler(err, fmt.Sprintf("delete failed for ID %d", id))
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			tx.Rollback()
			return nil, utils.ErrorHandler(err, "failed to get rows affected")
		}
		if rowsAffected == 0 {
			tx.Rollback()
			return nil, utils.ErrorHandler(fmt.Errorf("not found"), fmt.Sprintf("teacher with ID %d not found", id))
		}
		deletedIds = append(deletedIds, id)
	}

	if err := tx.Commit(); err != nil {
		return nil, utils.ErrorHandler(err, "transaction commit failed")
	}
	return deletedIds, nil
}

func GetStudentsByTeacherIdFromDb(teacherId string, students []models.Student) ([]models.Student, error) {
	db, err := ConnectDB()
	if err != nil {
		return nil, utils.ErrorHandler(err, "database connection failed")
	}
	defer db.Close()
	query := `SELECT id, first_name, last_name, email, class FROM students WHERE class = (SELECT class from teachers WHERE id = ?)`
	rows, err := db.Query(query, teacherId)
	if err != nil {
		return nil, utils.ErrorHandler(err, "error retrieving data")
	}
	defer rows.Close()

	for rows.Next() {
		var student models.Student
		err = rows.Scan(&student.ID, &student.FirstName, &student.LastName, &student.Email, &student.Class)
		if err != nil {
			return nil, utils.ErrorHandler(err, "error retrieving data")
		}
		students = append(students, student)
	}
	err = rows.Err()
	if err != nil {
		return nil, utils.ErrorHandler(err, "error retrieving data")
	}
	return students, nil
}

func GetStudentCountByTeacherIDFromDB(teacherID string) (int, error) {
	db, err := ConnectDB()
	if err != nil {
		return 0, utils.ErrorHandler(err, "error retrieving data")
	}

	defer db.Close()

	query := `SELECT COUNT(*) FROM students WHERE class = (SELECT class FROM teachers WHERE id = ?)`
	var studentCount int
	err = db.QueryRow(query, teacherID).Scan(&studentCount)
	if err != nil {
		return 0, utils.ErrorHandler(err, "error retrieving data")
	}
	return studentCount, nil
}
