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

// Global validator instance (reuse across requests for performance)
var validate = validator.New()

// ============================================================================
// HANDLERS
// ============================================================================

// GET /teachers - Get all teachers with optional filters and sorting
func GetTeachersHandler(w http.ResponseWriter, r *http.Request) {
	var teachers []models.Teacher
	teachers, err := sqlconnect.GetTeachersDbHandler(teachers, r)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteJSONList(w, http.StatusOK, teachers, len(teachers))
}

// GET /teachers/{id} - Get single teacher by ID
func GetOneTeacherHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	teacher, err := sqlconnect.GetTeacherByID(id)
	if err != nil {
		utils.WriteJSONError(w, http.StatusNotFound, err.Error())
		return
	}
	utils.WriteJSONSuccess(w, http.StatusOK, teacher)
}

// POST /teachers - Add new teachers
func AddTeacherHandler(w http.ResponseWriter, r *http.Request) {
	var newTeachers []models.Teacher

	// Strict decoder blocks unknown fields
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&newTeachers); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if len(newTeachers) == 0 {
		utils.WriteJSONError(w, http.StatusBadRequest, "At least one teacher is required")
		return
	}

	// Validate each teacher
	for i, teacher := range newTeachers {
		if err := validate.Struct(teacher); err != nil {
			if validationErrors, ok := err.(validator.ValidationErrors); ok {
				errorMsg := utils.FormatValidationErrors(validationErrors)
				utils.WriteJSONError(w, http.StatusBadRequest, "Teacher #"+strconv.Itoa(i+1)+": "+errorMsg)
				return
			}
			utils.WriteJSONError(w, http.StatusBadRequest, "Validation failed for teacher #"+strconv.Itoa(i+1))
			return
		}
	}

	addedTeachers, err := sqlconnect.AddTeachersDBHandler(newTeachers)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteJSONList(w, http.StatusCreated, addedTeachers, len(addedTeachers))
}

// PUT /teachers/{id} - Full update of teacher
func UpdateTeacherHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	var updatedTeacher models.Teacher
	if err := json.NewDecoder(r.Body).Decode(&updatedTeacher); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate the teacher data
	if err := validate.Struct(updatedTeacher); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errorMsg := utils.FormatValidationErrors(validationErrors)
			utils.WriteJSONError(w, http.StatusBadRequest, errorMsg)
			return
		}
		utils.WriteJSONError(w, http.StatusBadRequest, "Validation failed")
		return
	}

	updatedTeacherFromDB, err := sqlconnect.UpdateTeacher(id, updatedTeacher)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteJSONSuccess(w, http.StatusOK, updatedTeacherFromDB)
}

// PATCH /teachers/{id} - Partial update of single teacher
func PatchOneTeacherHandler(w http.ResponseWriter, r *http.Request) {
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

	updatedTeacher, err := sqlconnect.PatchOneTeacher(id, updates)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteJSONSuccess(w, http.StatusOK, updatedTeacher)
}

// PATCH /teachers - Batch partial update of multiple teachers
func PatchTeachersHandler(w http.ResponseWriter, r *http.Request) {
	var updates []map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if len(updates) == 0 {
		utils.WriteJSONError(w, http.StatusBadRequest, "At least one teacher update is required")
		return
	}

	updatedTeachers, err := sqlconnect.PatchTeachers(updates)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteJSONList(w, http.StatusOK, updatedTeachers, len(updatedTeachers))
}

// DELETE /teachers/{id} - Delete single teacher
func DeleteOneTeacherHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	if err := sqlconnect.DeleteOneTeacher(id); err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteJSONSuccess(w, http.StatusOK, id)
}

// DELETE /teachers - Batch delete multiple teachers
func DeleteTeachersHandler(w http.ResponseWriter, r *http.Request) {
	var ids []int
	if err := json.NewDecoder(r.Body).Decode(&ids); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if len(ids) == 0 {
		utils.WriteJSONError(w, http.StatusBadRequest, "At least one teacher ID is required")
		return
	}

	deletedIds, err := sqlconnect.DeleteTeachers(ids)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteJSONList(w, http.StatusOK, deletedIds, len(deletedIds))
}

func GetStudentsByTeacherId(w http.ResponseWriter, r *http.Request) {
	teacherId := r.PathValue("id")
	var students []models.Student

	student, err := sqlconnect.GetStudentsByTeacherIdFromDb(teacherId, students)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteJSONList(w, http.StatusOK, student, len(student))
}
