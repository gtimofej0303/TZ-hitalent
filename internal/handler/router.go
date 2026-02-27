package handler

import (
	"net/http"

	"github.com/gtimofej0303/org-structure-api/internal/repository/mygorm"
	"github.com/gtimofej0303/org-structure-api/internal/service"

	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB) http.Handler {
	deptRepo := mygorm.NewDepartmentRepository(db)
	empRepo := mygorm.NewEmployeeRepository(db)

	deptService := service.NewDepartmentService(deptRepo, empRepo)
	empService := service.NewEmployeeService(empRepo, deptRepo)

	deptHandler := NewDepartmentHandler(deptService)
	empHandler := NewEmployeeHandler(empService)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /departments/", deptHandler.CreateDepartment)
	mux.HandleFunc("GET /departments/{id}", deptHandler.GetDepartmentTree)
	mux.HandleFunc("PATCH /departments/{id}", deptHandler.UpdateDepartment)
	mux.HandleFunc("DELETE /departments/{id}", deptHandler.DeleteDepartment)

	mux.HandleFunc("POST /departments/{id}/employees/", empHandler.CreateEmployee)
	mux.HandleFunc("GET /departments/{id}/employees/", empHandler.GetEmployeesByDepartment)

	return Logger(Recoverer(mux))
}
