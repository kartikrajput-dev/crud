package handler

import (
	"net/http"

	"github.com/KARTIKrocks/apikit/errors"
	"github.com/KARTIKrocks/apikit/request"
	"github.com/KARTIKrocks/apikit/response"
	"github.com/kartikrajput-dev/crud/internal/store"
)

func ListUsers(w http.ResponseWriter, r *http.Request) error {
	pg, err := request.Paginate(r)
	if err != nil {
		return err
	}

	sorts, err := request.ParseSort(r, request.SortConfig{
		AllowedFields: []string{"id", "name", "email", "created_at"},
		Default:       []request.SortField{{Field: "id", Direction: request.SortAsc}},
	})
	if err != nil {
		return err
	}

	filters, err := request.ParseFilters(r, request.FilterConfig{
		AllowedFields: []string{"role", "active"},
	})
	if err != nil {
		return err
	}

	users, total, err := store.ListUsers(r.Context(), pg, sorts, filters)
	if err != nil {
		return err
	}

	response.New().
		Message("Users retrieved").
		Data(users).
		Pagination(pg.Page, pg.PerPage, total).
		Send(w)
	return nil
}

func GetUser(w http.ResponseWriter, r *http.Request) error {
	id, err := request.PathParamInt(r, "id")
	if err != nil {
		return errors.BadRequest("invalid user id")
	}

	user, err := store.GetUser(r.Context(), id)
	if err != nil {
		return err
	}

	response.OK(w, "User retrieved", user)
	return nil
}

func CreateUser(w http.ResponseWriter, r *http.Request) error {
	in, err := request.Bind[store.CreateUserInput](r)
	if err != nil {
		return err
	}

	user, err := store.CreateUser(r.Context(), in)
	if err != nil {
		return err
	}

	response.Created(w, "User created", user)
	return nil
}

func UpdateUser(w http.ResponseWriter, r *http.Request) error {
	id, err := request.PathParamInt(r, "id")
	if err != nil {
		return errors.BadRequest("invalid user id")
	}

	in, err := request.Bind[store.UpdateUserInput](r)
	if err != nil {
		return err
	}

	user, err := store.UpdateUser(r.Context(), id, in)
	if err != nil {
		return err
	}

	response.OK(w, "User updated", user)
	return nil
}

func DeleteUser(w http.ResponseWriter, r *http.Request) error {
	id, err := request.PathParamInt(r, "id")
	if err != nil {
		return errors.BadRequest("invalid user id")
	}

	if err := store.DeleteUser(r.Context(), id); err != nil {
		return err
	}

	response.OK(w, "User deleted", nil)
	return nil
}
