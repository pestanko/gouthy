package services

import (
	"github.com/pestanko/gouthy/app/models"
)

func ConvertModelsToUserList(list []models.User) []ListUser {
	var listUsers []ListUser

	for _, user := range list {
		item := ListUser{
			Username: user.Username,
			Email:    user.Email,
			ID:       user.ID,
		}
		listUsers = append(listUsers, item)
	}
	return listUsers
}

func ConvertModelsToUserDTO(user *models.User) UserDTO {
	return UserDTO{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Username:  user.Username,
		Name:      user.Name,
		Email:     user.Email,
	}
}
