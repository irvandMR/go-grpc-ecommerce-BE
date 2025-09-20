package entity

import "time"

const (
	CUSTOMER_ROLE = "customer"
	ADMIN_ROLE    = "admin"
)

type UserRole struct {
	Id        string
	Name	  string
	Code      string
	CreatedAt time.Time
	CreatedBy *string
	UpdatedAt time.Time
	UpdatedBy *string
	DeletedAt time.Time
	DeletedBy *string
	IsDeleted bool
}

type User struct {
	Id        string
	Fullname  string
	Email     string
	Password  string
	RoleCode  string
	CreatedAt time.Time
	CreatedBy *string
	UpdatedAt time.Time
	UpdatedBy *string
	DeletedAt time.Time
	DeletedBy *string
	IsDeleted bool
}