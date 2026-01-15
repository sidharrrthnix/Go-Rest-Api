package sqlconnect

import (
	"database/sql"

	"http-api.com/models"
	"http-api.com/utils"
)

func getFieldMappings(student *models.Student) map[string]*string {
	return map[string]*string{
		"first_name": &student.FirstName,
		"last_name":  &student.LastName,
		"email":      &student.Email,
		"class":      &student.Class,
	}
}

func applyUpdate(student *models.Student, updates map[string]interface{}) {
	fieldMap := getFieldMappings(student)
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

func GetStudentsById(id int) (models.Student, error) {
	db, err := ConnectDB()
	if err != nil {
		return models.Student{}, utils.ErrorHandler(err, "database connection failed")
	}
	defer db.Close()
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

func AddStudentsDBHandler(newStudents []models.Student) ([]models.Student, error) {
	db, err := ConnectDB()
	if err != nil {
		return nil, utils.ErrorHandler(err, "database connection failed")
	}
	defer db.Close()
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
		lastId, _ := res.LastInsertId()
		student.ID = int(lastId)
		addedStudents = append(addedStudents, student)
	}
	return addedStudents, nil
}

func UpdateStudent(id int, updatedStudent models.Student) (models.Student, error) {
	db, err := ConnectDB()
	if err != nil {
		return models.Student{}, utils.ErrorHandler(err, "database connection failed")
	}
	defer db.Close()

	var exists int
	err = db.QueryRow("SELECT 1 FROM students WHERE id = ?", id).Scan(&exists)
	if err == sql.ErrNoRows {
		return models.Student{}, utils.ErrorHandler(err, "student not found")
	}
	if err != nil {
		return models.Student{}, utils.ErrorHandler(err, "query failed")
	}
	updatedStudent.ID = id
	_, err = db.Exec("UPDATE students SET first_name = ?, last_name = ?, email = ?, class = ? WHERE id = ?", updatedStudent.FirstName, updatedStudent.LastName, updatedStudent.Email, updatedStudent.Class, updatedStudent.ID)
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
	defer db.Close()
	student, err := GetStudentsById(id)
	if err != nil {
		return models.Student{}, err
	}
	applyUpdate(&student, updates)

	_, err = db.Exec("UPDATE students SET first_name = ?, last_name = ?, email = ?, class = ? WHERE id = ?", student.FirstName, student.LastName, student.Email, student.Class, student.ID)
	if err != nil {
		return models.Student{}, utils.ErrorHandler(err, "update failed")
	}
	return student, nil
}
