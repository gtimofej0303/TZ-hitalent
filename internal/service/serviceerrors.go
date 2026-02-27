package service

import "errors"

var (
    //Department Errors
    ErrDepartmentNotFound       = errors.New("department not found")
    ErrParentNotFound           = errors.New("parent department not found")
    ErrDuplicateName            = errors.New("department name already exists in parent")
    ErrCycleDetected            = errors.New("parent change creates cycle")
    ErrInvalidMode              = errors.New("invalid delete mode")
    ErrInvalidReassignTarget    = errors.New("invalid reassign target department")
    ErrInvalidName              = errors.New("invalid department name")
    ErrSelfParent               = errors.New("element becomes selfparent")

    //Employee errors
    ErrInvalidFullName          = errors.New("full_name must be between 1 and 200 characters")
    ErrInvalidPosition          = errors.New("position must be between 1 and 200 characters")
)
