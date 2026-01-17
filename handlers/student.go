package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"http-api.com/models"
	"http-api.com/respository/sqlconnect"
	"http-api.com/utils"
)

var validates = validator.New()

func GetStudentsHandler(w http.ResponseWriter, r *http.Request) {
	var students []models.Student
	students, err := sqlconnect.GetStudentsDbHandler(students, r)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteJSONList(w, http.StatusOK, students, len(students))
}

func GetOneStudentHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "Invalid ID")
		return
	}
	student, err := sqlconnect.GetStudentByID(id)
	if err != nil {
		utils.WriteJSONError(w, http.StatusNotFound, err.Error())
		return
	}
	utils.WriteJSONSuccess(w, http.StatusOK, student)
}

func AddStudentHandler(w http.ResponseWriter, r *http.Request) {
	var newStudents []models.Student

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&newStudents); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if len(newStudents) == 0 {
		utils.WriteJSONError(w, http.StatusBadRequest, "At least one student is required")
		return
	}

	for i, student := range newStudents {
		if err := validates.Struct(student); err != nil {
			if validationErrors, ok := err.(validator.ValidationErrors); ok {
				errorMsg := utils.FormatValidationErrors(validationErrors)
				utils.WriteJSONError(w, http.StatusBadRequest, "Student #"+strconv.Itoa(i+1)+": "+errorMsg)
				return
			}
			utils.WriteJSONError(w, http.StatusBadRequest, "Validation failed for student #"+strconv.Itoa(i+1))
			return
		}
	}

	addedStudents, err := sqlconnect.AddStudentsDBHandler(newStudents)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteJSONList(w, http.StatusCreated, addedStudents, len(addedStudents))
}

func UpdateStudentHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	var updatedStudent models.Student
	if err := json.NewDecoder(r.Body).Decode(&updatedStudent); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := validates.Struct(updatedStudent); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errorMsg := utils.FormatValidationErrors(validationErrors)
			utils.WriteJSONError(w, http.StatusBadRequest, errorMsg)
			return
		}
		utils.WriteJSONError(w, http.StatusBadRequest, "Validation failed")
		return
	}

	updatedStudentFromDB, err := sqlconnect.UpdateStudent(id, updatedStudent)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteJSONSuccess(w, http.StatusOK, updatedStudentFromDB)
}

func PatchOneStudentHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if len(updates) == 0 {
		utils.WriteJSONError(w, http.StatusBadRequest, "No fields to update")
		return
	}

	patchedStudent, err := sqlconnect.PatchOneStudent(id, updates)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteJSONSuccess(w, http.StatusOK, patchedStudent)
}

func PatchStudentsHandler(w http.ResponseWriter, r *http.Request) {
	var updates []map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if len(updates) == 0 {
		utils.WriteJSONError(w, http.StatusBadRequest, "At least one student update is required")
		return
	}
	patchedStudents, err := sqlconnect.PatchStudents(updates)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteJSONList(w, http.StatusOK, patchedStudents, len(patchedStudents))
}

func DeleteOneStudentHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "Invalid ID")
		return
	}
	if err := sqlconnect.DeleteOneStudent(id); err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteJSONSuccess(w, http.StatusOK, id)
}

func DeleteStudentsHandler(w http.ResponseWriter, r *http.Request) {
	var ids []int
	if err := json.NewDecoder(r.Body).Decode(&ids); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if len(ids) == 0 {
		utils.WriteJSONError(w, http.StatusBadRequest, "At least one student ID is required")
		return
	}
	deleteIds, err := sqlconnect.DeleteStudents(ids)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteJSONList(w, http.StatusOK, deleteIds, len(deleteIds))
}
