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

// GetTeachersDbHandler retrieves all teachers with optional filters and sorting
func GetTeachersDbHandler(teachers []models.Teacher, r *http.Request) ([]models.Teacher, error) {
	db, err := ConnectDB()
	if err != nil {
		return nil, utils.ErrorHandler(err, "error connecting to database")
	}
	defer db.Close()

	query := "SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE 1=1"
	var args []interface{}

	// Add filters from query params
	params := map[string]string{
		"first_name": "first_name",
		"last_name":  "last_name",
		"email":      "email",
		"class":      "class",
		"subject":    "subject",
	}
	for param, dbField := range params {
		value := r.URL.Query().Get(param)
		if value != "" {
			query += " AND " + dbField + " = ?"
			args = append(args, value)
		}
	}

	// Add sorting
	sortParams := r.URL.Query()["sortby"]
	if len(sortParams) > 0 {
		orderBys := []string{}
		for _, param := range sortParams {
			parts := strings.Split(param, ":")
			if len(parts) == 2 {
				field := parts[0]
				order := strings.ToUpper(parts[1])
				// Validate field and order
				validFields := map[string]bool{
					"first_name": true,
					"last_name":  true,
					"email":      true,
					"class":      true,
					"subject":    true,
				}
				if validFields[field] && (order == "ASC" || order == "DESC") {
					orderBys = append(orderBys, field+" "+order)
				}
			}
		}
		if len(orderBys) > 0 {
			query += " ORDER BY " + strings.Join(orderBys, ", ")
		}
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, utils.ErrorHandler(err, "error querying database")
	}
	defer rows.Close()

	for rows.Next() {
		var teacher models.Teacher
		err := rows.Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
		if err != nil {
			return nil, utils.ErrorHandler(err, "error scanning row")
		}
		teachers = append(teachers, teacher)
	}
	return teachers, nil
}

// GetTeacherByID retrieves a single teacher by ID
func GetTeacherByID(id int) (models.Teacher, error) {
	db, err := ConnectDB()
	if err != nil {
		return models.Teacher{}, utils.ErrorHandler(err, "error connecting to database")
	}
	defer db.Close()

	var teacher models.Teacher
	err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", id).
		Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)

	if err == sql.ErrNoRows {
		return models.Teacher{}, utils.ErrorHandler(err, "teacher not found")
	} else if err != nil {
		return models.Teacher{}, utils.ErrorHandler(err, "error querying database")
	}
	return teacher, nil
}

// AddTeachersDBHandler adds multiple teachers in batch
func AddTeachersDBHandler(newTeachers []models.Teacher) ([]models.Teacher, error) {
	db, err := ConnectDB()
	if err != nil {
		return nil, utils.ErrorHandler(err, "error connecting to database")
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO teachers (first_name, last_name, email, class, subject) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return nil, utils.ErrorHandler(err, "error preparing statement")
	}
	defer stmt.Close()

	addedTeachers := make([]models.Teacher, 0, len(newTeachers))
	for _, newTeacher := range newTeachers {
		res, err := stmt.Exec(newTeacher.FirstName, newTeacher.LastName, newTeacher.Email, newTeacher.Class, newTeacher.Subject)
		if err != nil {
			return nil, utils.ErrorHandler(err, "error inserting teacher")
		}
		lastID, err := res.LastInsertId()
		if err != nil {
			return nil, utils.ErrorHandler(err, "error getting last insert id")
		}
		newTeacher.ID = int(lastID)
		addedTeachers = append(addedTeachers, newTeacher)
	}
	return addedTeachers, nil
}

// UpdateTeacher performs a full update of a teacher (all fields required)
func UpdateTeacher(id int, updatedTeacher models.Teacher) (models.Teacher, error) {
	db, err := ConnectDB()
	if err != nil {
		return models.Teacher{}, utils.ErrorHandler(err, "error connecting to database")
	}
	defer db.Close()

	// Check if teacher exists
	var existingTeacher models.Teacher
	err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", id).
		Scan(&existingTeacher.ID, &existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Email, &existingTeacher.Class, &existingTeacher.Subject)

	if err == sql.ErrNoRows {
		return models.Teacher{}, utils.ErrorHandler(err, "teacher not found")
	} else if err != nil {
		return models.Teacher{}, utils.ErrorHandler(err, "error querying database")
	}

	// Update the teacher
	updatedTeacher.ID = existingTeacher.ID
	_, err = db.Exec("UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?",
		updatedTeacher.FirstName, updatedTeacher.LastName, updatedTeacher.Email, updatedTeacher.Class, updatedTeacher.Subject, updatedTeacher.ID)
	if err != nil {
		return models.Teacher{}, utils.ErrorHandler(err, "error updating teacher")
	}
	return updatedTeacher, nil
}

// PatchOneTeacher performs a partial update of a single teacher
func PatchOneTeacher(id int, updates map[string]interface{}) (models.Teacher, error) {
	db, err := ConnectDB()
	if err != nil {
		return models.Teacher{}, utils.ErrorHandler(err, "error connecting to database")
	}
	defer db.Close()

	// Check if teacher exists and get current data
	var existingTeacher models.Teacher
	err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", id).
		Scan(&existingTeacher.ID, &existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Email, &existingTeacher.Class, &existingTeacher.Subject)

	if err == sql.ErrNoRows {
		return models.Teacher{}, utils.ErrorHandler(err, "teacher not found")
	} else if err != nil {
		return models.Teacher{}, utils.ErrorHandler(err, "error querying database")
	}

	// Map JSON keys to struct fields
	fieldMap := map[string]*string{
		"firstName": &existingTeacher.FirstName,
		"lastName":  &existingTeacher.LastName,
		"email":     &existingTeacher.Email,
		"class":     &existingTeacher.Class,
		"subject":   &existingTeacher.Subject,
	}

	// Apply updates
	for key, value := range updates {
		if fieldPtr, ok := fieldMap[key]; ok {
			if strVal, ok := value.(string); ok {
				*fieldPtr = strVal
			}
		}
	}

	// Update in database
	_, err = db.Exec("UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?",
		existingTeacher.FirstName, existingTeacher.LastName, existingTeacher.Email, existingTeacher.Class, existingTeacher.Subject, existingTeacher.ID)
	if err != nil {
		return models.Teacher{}, utils.ErrorHandler(err, "error updating teacher")
	}
	return existingTeacher, nil
}

// PatchTeachers performs a batch partial update of multiple teachers
func PatchTeachers(updates []map[string]interface{}) ([]models.Teacher, error) {
	db, err := ConnectDB()
	if err != nil {
		return nil, utils.ErrorHandler(err, "error connecting to database")
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return nil, utils.ErrorHandler(err, "error beginning transaction")
	}

	updatedTeachers := []models.Teacher{}

	for i, update := range updates {
		// Extract ID
		var id int
		idVal, exists := update["id"]
		if !exists {
			tx.Rollback()
			return nil, utils.ErrorHandler(fmt.Errorf("missing id"), "missing id in update #"+strconv.Itoa(i+1))
		}

		switch v := idVal.(type) {
		case float64:
			id = int(v)
		case string:
			id, err = strconv.Atoi(v)
			if err != nil {
				tx.Rollback()
				return nil, utils.ErrorHandler(err, "invalid ID in update #"+strconv.Itoa(i+1))
			}
		default:
			tx.Rollback()
			return nil, utils.ErrorHandler(fmt.Errorf("invalid ID type"), "invalid ID type in update #"+strconv.Itoa(i+1))
		}

		// Get existing teacher
		var teacherFromDb models.Teacher
		err = tx.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", id).
			Scan(&teacherFromDb.ID, &teacherFromDb.FirstName, &teacherFromDb.LastName, &teacherFromDb.Email, &teacherFromDb.Class, &teacherFromDb.Subject)

		if err == sql.ErrNoRows {
			tx.Rollback()
			return nil, utils.ErrorHandler(err, "teacher with ID "+strconv.Itoa(id)+" not found")
		} else if err != nil {
			tx.Rollback()
			return nil, utils.ErrorHandler(err, "error querying teacher")
		}

		// Map JSON keys to struct fields
		fieldMap := map[string]*string{
			"firstName": &teacherFromDb.FirstName,
			"lastName":  &teacherFromDb.LastName,
			"email":     &teacherFromDb.Email,
			"class":     &teacherFromDb.Class,
			"subject":   &teacherFromDb.Subject,
		}

		// Apply updates (skip 'id' field)
		for key, value := range update {
			if key == "id" {
				continue
			}
			if fieldPtr, ok := fieldMap[key]; ok {
				if strVal, ok := value.(string); ok {
					*fieldPtr = strVal
				}
			}
		}

		// Update in database
		_, err = tx.Exec("UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?",
			teacherFromDb.FirstName, teacherFromDb.LastName, teacherFromDb.Email, teacherFromDb.Class, teacherFromDb.Subject, teacherFromDb.ID)
		if err != nil {
			tx.Rollback()
			return nil, utils.ErrorHandler(err, "error updating teacher #"+strconv.Itoa(i+1))
		}

		updatedTeachers = append(updatedTeachers, teacherFromDb)
	}

	err = tx.Commit()
	if err != nil {
		return nil, utils.ErrorHandler(err, "error committing transaction")
	}
	return updatedTeachers, nil
}

// DeleteOneTeacher deletes a single teacher by ID
func DeleteOneTeacher(id int) error {
	db, err := ConnectDB()
	if err != nil {
		return utils.ErrorHandler(err, "error connecting to database")
	}
	defer db.Close()

	result, err := db.Exec("DELETE FROM teachers WHERE id = ?", id)
	if err != nil {
		return utils.ErrorHandler(err, "error deleting teacher")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return utils.ErrorHandler(err, "error getting rows affected")
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
		return nil, utils.ErrorHandler(err, "error connecting to database")
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return nil, utils.ErrorHandler(err, "error beginning transaction")
	}

	stmt, err := tx.Prepare("DELETE FROM teachers WHERE id = ?")
	if err != nil {
		tx.Rollback()
		return nil, utils.ErrorHandler(err, "error preparing statement")
	}
	defer stmt.Close()

	deletedIds := []int{}

	for _, id := range ids {
		result, err := stmt.Exec(id)
		if err != nil {
			tx.Rollback()
			return nil, utils.ErrorHandler(err, "error deleting teacher")
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			tx.Rollback()
			return nil, utils.ErrorHandler(err, "error getting rows affected")
		}

		if rowsAffected > 0 {
			deletedIds = append(deletedIds, id)
		} else {
			tx.Rollback()
			return nil, utils.ErrorHandler(fmt.Errorf("teacher not found"), "teacher with ID "+strconv.Itoa(id)+" not found")
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, utils.ErrorHandler(err, "error committing transaction")
	}

	if len(deletedIds) < 1 {
		return nil, utils.ErrorHandler(fmt.Errorf("no teachers deleted"), "no teachers were deleted")
	}
	return deletedIds, nil
}
