package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gtimofej0303/TZ-hitalent/internal/dto"
	"github.com/gtimofej0303/TZ-hitalent/internal/service"
)

type DepartmentHandler struct {
	service service.DepartmentService
}

func NewDepartmentHandler(s service.DepartmentService) *DepartmentHandler {
	return &DepartmentHandler{service: s}
}

func (h *DepartmentHandler) CreateDepartment(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateDepartmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	dept, err := h.service.Create(r.Context(), req.Name, req.ParentID)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, dept)
}

func (h *DepartmentHandler) GetDepartmentTree(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	depth := 1
	if d := r.URL.Query().Get("depth"); d != "" {
		if parsed, err := strconv.Atoi(d); (err == nil) && (parsed >= 1) && (parsed <= 5) {
			depth = parsed
		} else {
			writeError(w, http.StatusBadRequest, "invalid depth")
		}
	}

	includeEmployees := true
	if ie := r.URL.Query().Get("include_employees"); ie == "false" {
		includeEmployees = false
	}

	children, employees, err := h.service.GetTree(r.Context(), id, depth, includeEmployees)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	resp := dto.BuildDepartmentTree(id, children, employees)
	writeJSON(w, http.StatusOK, resp)
}

func (h *DepartmentHandler) UpdateDepartment(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req dto.UpdateDepartmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	dept, err := h.service.Update(r.Context(), id, req.Name, req.ParentID)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, dept)
}

func (h *DepartmentHandler) DeleteDepartment(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	mode := r.URL.Query().Get("mode")
	if mode == "" {
		writeError(w, http.StatusBadRequest, "mode is required (cascade or reassign)")
		return
	}

	var reassignTo *uint
	if mode == "reassign" {
		rt := r.URL.Query().Get("reassign_to")

		if rt == "" {
			writeError(w, http.StatusBadRequest, "invalid reassign_to")
			return
		}

		parsed, err := strconv.ParseUint(rt, 10, 64)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid reassign_to")
			return
		}
		val := uint(parsed)
		reassignTo = &val
	}

	if err := h.service.Delete(r.Context(), id, mode, reassignTo); err != nil {
		writeServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func parseID(r *http.Request) (uint, error) {
	raw := r.PathValue("id")
	parsed, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(parsed), nil
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrDepartmentNotFound):
		writeError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, service.ErrParentNotFound):
		writeError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, service.ErrDuplicateName):
		writeError(w, http.StatusConflict, err.Error())
	case errors.Is(err, service.ErrCycleDetected):
		writeError(w, http.StatusConflict, err.Error())
	case errors.Is(err, service.ErrSelfParent):
		writeError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, service.ErrInvalidMode):
		writeError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, service.ErrInvalidReassignTarget):
		writeError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, service.ErrInvalidName):
		writeError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, service.ErrInvalidFullName):
   		writeError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, service.ErrInvalidPosition):
    	writeError(w, http.StatusBadRequest, err.Error())
	default:
		log.Printf("internal error: %v", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}
