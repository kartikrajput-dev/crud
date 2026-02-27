package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/KARTIKrocks/apikit"
)

func CreateUser(ctx *apikit.Context) {
	var user User
	if err := json.NewDecoder(ctx.Request.Body).Decode(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, apikit.H{"error": "Invalid request"})
		return
	}

	err := DB.QueryRow(
		"INSERT INTO users(name,email) VALUES($1,$2) RETURNING id",
		user.Name, user.Email,
	).Scan(&user.ID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, apikit.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

func GetUsers(ctx *apikit.Context) {
	rows, err := DB.Query("SELECT id,name,email FROM users")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, apikit.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var u User
		rows.Scan(&u.ID, &u.Name, &u.Email)
		users = append(users, u)
	}

	ctx.JSON(http.StatusOK, users)
}

func GetUser(ctx *apikit.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))

	var user User
	err := DB.QueryRow(
		"SELECT id,name,email FROM users WHERE id=$1",
		id,
	).Scan(&user.ID, &user.Name, &user.Email)

	if err != nil {
		ctx.JSON(http.StatusNotFound, apikit.H{"error": "User not found"})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func UpdateUser(ctx *apikit.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))

	var user User
	json.NewDecoder(ctx.Request.Body).Decode(&user)

	_, err := DB.Exec(
		"UPDATE users SET name=$1,email=$2 WHERE id=$3",
		user.Name, user.Email, id,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, apikit.H{"error": err.Error()})
		return
	}

	user.ID = id
	ctx.JSON(http.StatusOK, user)
}

func DeleteUser(ctx *apikit.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))

	_, err := DB.Exec("DELETE FROM users WHERE id=$1", id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, apikit.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, apikit.H{"message": "User deleted"})
}