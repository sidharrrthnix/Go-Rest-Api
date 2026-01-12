package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"http-api.com/models"
	"http-api.com/respository/sqlconnect"
)

// Helper: write JSON error response
func writeJSONError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"status":  "error",
		"message": msg,
	})
}

// GET /teachers - Get all teachers with optional filters
func GetTeachersHandler(w http.ResponseWriter, r *http.Request) {
	var teachers []models.Teacher
	teachers, err := sqlconnect.GetTeachersDbHandler(teachers, r)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Teacher `json:"data"`
	}{
		Status: "success",
		Count:  len(teachers),
		Data:   teachers,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GET /teachers/{id} - Get single teacher by ID
func GetOneTeacherHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	teacher, err := sqlconnect.GetTeacherByID(id)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, err.Error())
		return
	}

	response := struct {
		Status string         `json:"status"`
		Data   models.Teacher `json:"data"`
	}{
		Status: "success",
		Data:   teacher,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// POST /teachers - Add new teachers
func AddTeacherHandler(w http.ResponseWriter, r *http.Request) {
	var newTeachers []models.Teacher
	if err := json.NewDecoder(r.Body).Decode(&newTeachers); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if len(newTeachers) == 0 {
		writeJSONError(w, http.StatusBadRequest, "No teachers to add")
		return
	}

	addedTeachers, err := sqlconnect.AddTeachersDBHandler(newTeachers)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Teacher `json:"data"`
	}{
		Status: "success",
		Count:  len(addedTeachers),
		Data:   addedTeachers,
	}
	json.NewEncoder(w).Encode(response)
}

// PUT /teachers/{id} - Update existing teacher (full update)
func UpdateTeacherHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	var updatedTeacher models.Teacher
	if err := json.NewDecoder(r.Body).Decode(&updatedTeacher); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate required fields
	if updatedTeacher.FirstName == "" || updatedTeacher.LastName == "" || updatedTeacher.Email == "" {
		writeJSONError(w, http.StatusBadRequest, "Missing required fields: firstName, lastName, email")
		return
	}

	updatedTeacherFromDB, err := sqlconnect.UpdateTeacher(id, updatedTeacher)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := struct {
		Status  string         `json:"status"`
		Data    models.Teacher `json:"data"`
		Message string         `json:"message,omitempty"`
	}{
		Status:  "success",
		Data:    updatedTeacherFromDB,
		Message: "Teacher updated successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// PATCH /teachers/{id} - Partial update single teacher
func PatchOneTeacherHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if len(updates) == 0 {
		writeJSONError(w, http.StatusBadRequest, "No fields to update")
		return
	}

	updatedTeacher, err := sqlconnect.PatchOneTeacher(id, updates)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := struct {
		Status  string         `json:"status"`
		Data    models.Teacher `json:"data"`
		Message string         `json:"message,omitempty"`
	}{
		Status:  "success",
		Data:    updatedTeacher,
		Message: "Teacher updated successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// PATCH /teachers - Batch partial update multiple teachers
func PatchTeachersHandler(w http.ResponseWriter, r *http.Request) {
	var updates []map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if len(updates) == 0 {
		writeJSONError(w, http.StatusBadRequest, "No teachers to update")
		return
	}

	updatedTeachers, err := sqlconnect.PatchTeachers(updates)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := struct {
		Status  string           `json:"status"`
		Count   int              `json:"count"`
		Data    []models.Teacher `json:"data"`
		Message string           `json:"message,omitempty"`
	}{
		Status:  "success",
		Count:   len(updatedTeachers),
		Data:    updatedTeachers,
		Message: "Teachers updated successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// DELETE /teachers/{id} - Delete single teacher
func DeleteOneTeacherHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	err = sqlconnect.DeleteOneTeacher(id)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		ID      int    `json:"id"`
	}{
		Status:  "success",
		Message: "Teacher successfully deleted",
		ID:      id,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// DELETE /teachers - Batch delete multiple teachers
func DeleteTeachersHandler(w http.ResponseWriter, r *http.Request) {
	var ids []int
	if err := json.NewDecoder(r.Body).Decode(&ids); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if len(ids) == 0 {
		writeJSONError(w, http.StatusBadRequest, "No teacher IDs provided")
		return
	}

	deletedIds, err := sqlconnect.DeleteTeachers(ids)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := struct {
		Status     string `json:"status"`
		Message    string `json:"message"`
		DeletedIDs []int  `json:"deleted_ids"`
	}{
		Status:     "success",
		Message:    "Teachers successfully deleted",
		DeletedIDs: deletedIds,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
