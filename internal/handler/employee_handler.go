package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gtimofej0303/TZ-hitalent/internal/dto"
	"github.com/gtimofej0303/TZ-hitalent/internal/service"
)

type EmployeeHandler struct {
	service service.EmployeeService
}

func NewEmployeeHandler(s service.EmployeeService) *EmployeeHandler {
    return &EmployeeHandler{service: s}
}

func (h *EmployeeHandler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
    deptID, err := parseID(r)
    if err != nil {
        writeError(w, http.StatusBadRequest, "invalid id")
        return
    }

	var req dto.CreateEmployeeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, http.StatusBadRequest, "invalid request body")
        return
    }

	emp, err := h.service.Create(r.Context(), deptID, req.Fullname, req.Position, req.HiredAt)
    if err != nil {
        writeServiceError(w, err)
        return
    }

	writeJSON(w, http.StatusCreated, emp)
}

func (h *EmployeeHandler) GetEmployeesByDepartment(w http.ResponseWriter, r *http.Request) {
	deptID, err := parseID(r)
    if err != nil {
        writeError(w, http.StatusBadRequest, "invalid id")
        return
    }

	employees, err := h.service.GetByDepartmentID(r.Context(), deptID)
    if err != nil {
        writeServiceError(w, err)
        return
    }

	writeJSON(w, http.StatusOK, employees)
}