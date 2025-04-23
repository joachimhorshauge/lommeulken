package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"lommeulken/gen/dbstore"
	"net/http"
)

type userDisplay struct {
	Id    string `json:"id"`
	Value string `json:"value"`
	Label string `json:"label"`
	Email string `json:"email"`
	Image string `json:"img"`
}

func (h *Handler) HandleListUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	params := dbstore.ListUsersParams{
		Limit:  1000,
		Offset: 0,
	}

	users, err := h.queries.ListUsers(context.Background(), params)
	if err != nil {
		slog.Error("Failed to retrieve users", "Error", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	userList := []userDisplay{}

	for _, user := range users {
		userList = append(userList, userDisplay{
			Id:    user.ID.String(),
			Value: fmt.Sprintf("%s %s", user.FirstName.String, user.LastName.String),
			Label: fmt.Sprintf("%s %s", user.FirstName.String, user.LastName.String),
			Email: user.Email.String,
			Image: "https://cdn.pixabay.com/photo/2023/02/18/11/00/icon-7797704_640.png",
		})
	}

	userListJson, err := json.Marshal(userList)
	if err != nil {
		slog.Error("Failed to marshal list of users", "Error", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Write(userListJson)
	return
}
