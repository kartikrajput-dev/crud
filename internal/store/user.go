package store

import (
	"context"
	stderrors "errors"
	"strings"
	"time"

	"github.com/KARTIKrocks/apikit/dbx"
	"github.com/KARTIKrocks/apikit/errors"
	"github.com/KARTIKrocks/apikit/request"
	"github.com/KARTIKrocks/apikit/sqlbuilder"
)

type User struct {
	ID        int       `db:"id"         json:"id"`
	Name      string    `db:"name"       json:"name"`
	Email     string    `db:"email"      json:"email"`
	Role      string    `db:"role"       json:"role"`
	Active    bool      `db:"active"     json:"active"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type CreateUserInput struct {
	Name  string `json:"name"  validate:"required,min=2,max=100"`
	Email string `json:"email" validate:"required,email"`
	Role  string `json:"role"  validate:"oneof=admin user mod"`
}

type UpdateUserInput struct {
	Name   *string `json:"name"   validate:"omitempty,min=2,max=100"`
	Email  *string `json:"email"  validate:"omitempty,email"`
	Role   *string `json:"role"   validate:"omitempty,oneof=admin user mod"`
	Active *bool   `json:"active"`
}

var allowedSortCols = map[string]string{
	"id":         "id",
	"name":       "name",
	"email":      "email",
	"created_at": "created_at",
}

var allowedFilterCols = map[string]string{
	"role":   "role",
	"active": "active",
}

func ListUsers(ctx context.Context, pg request.Pagination, sorts []request.SortField, filters []request.Filter) ([]User, int, error) {
	type countRow struct {
		Count int `db:"count"`
	}

	countSQL, countArgs := sqlbuilder.Select("COUNT(*) as count").
		From("users").
		ApplyFilters(filters, allowedFilterCols).
		Build()

	row, err := dbx.QueryOne[countRow](ctx, countSQL, countArgs...)
	if err != nil {
		return nil, 0, errors.Internal("failed to count users").Wrap(err)
	}

	dataSQL, dataArgs := sqlbuilder.Select("id", "name", "email", "role", "active", "created_at", "updated_at").
		From("users").
		ApplyFilters(filters, allowedFilterCols).
		ApplySort(sorts, allowedSortCols).
		ApplyPagination(pg).
		Build()

	users, err := dbx.QueryAll[User](ctx, dataSQL, dataArgs...)
	if err != nil {
		return nil, 0, errors.Internal("failed to list users").Wrap(err)
	}

	return users, row.Count, nil
}

func GetUser(ctx context.Context, id int) (*User, error) {
	sql, args := sqlbuilder.Select("id", "name", "email", "role", "active", "created_at", "updated_at").
		From("users").
		WhereEq("id", id).
		Build()

	user, err := dbx.QueryOne[User](ctx, sql, args...)
	if err != nil {
		if stderrors.Is(err, errors.ErrNotFound) {
			return nil, errors.NotFound("User")
		}
		return nil, errors.Internal("failed to get user").Wrap(err)
	}
	return &user, nil
}

func CreateUser(ctx context.Context, in CreateUserInput) (*User, error) {
	sql, args := sqlbuilder.Insert("users").
		Columns("name", "email", "role").
		Values(in.Name, in.Email, in.Role).
		Returning("id", "name", "email", "role", "active", "created_at", "updated_at").
		Build()

	user, err := dbx.QueryOne[User](ctx, sql, args...)
	if err != nil {
		if isUniqueViolation(err) {
			return nil, errors.Conflict("Email already in use").WithField("email", "already exists")
		}
		return nil, errors.Internal("failed to create user").Wrap(err)
	}
	return &user, nil
}

func UpdateUser(ctx context.Context, id int, in UpdateUserInput) (*User, error) {
	b := sqlbuilder.Update("users").WhereEq("id", id)

	if in.Name != nil {
		b.Set("name", *in.Name)
	}
	if in.Email != nil {
		b.Set("email", *in.Email)
	}
	if in.Role != nil {
		b.Set("role", *in.Role)
	}
	if in.Active != nil {
		b.Set("active", *in.Active)
	}
	b.SetExpr("updated_at", sqlbuilder.Raw("NOW()"))

	sql, args := b.Returning("id", "name", "email", "role", "active", "created_at", "updated_at").Build()

	user, err := dbx.QueryOne[User](ctx, sql, args...)
	if err != nil {
		if stderrors.Is(err, errors.ErrNotFound) {
			return nil, errors.NotFound("User")
		}
		if isUniqueViolation(err) {
			return nil, errors.Conflict("Email already in use").WithField("email", "already exists")
		}
		return nil, errors.Internal("failed to update user").Wrap(err)
	}
	return &user, nil
}

func DeleteUser(ctx context.Context, id int) error {
	type deletedRow struct {
		ID int `db:"id"`
	}

	sql, args := sqlbuilder.Delete("users").
		WhereEq("id", id).
		Returning("id").
		Build()

	_, err := dbx.QueryOne[deletedRow](ctx, sql, args...)
	if err != nil {
		if stderrors.Is(err, errors.ErrNotFound) {
			return errors.NotFound("User")
		}
		return errors.Internal("failed to delete user").Wrap(err)
	}
	return nil
}

func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "unique") || strings.Contains(msg, "duplicate") || strings.Contains(msg, "23505")
}
