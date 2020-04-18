package services

import (
	"github.com/pestanko/gouthy/app/models"
)

func ConvertModelsToUserList(list []models.User) []ListUser {
	var listUsers []ListUser

	for _, user := range list {
		item := ListUser{
			UserBase: *ConvertModelToUserBase(&user),
		}
		listUsers = append(listUsers, item)
	}
	return listUsers
}

func ConvertModelToUserBase(user *models.User) *UserBase {
	if user == nil {
		return nil
	}

	return &UserBase{
		Username: user.Username,
		Name:     user.Name,
		Email:    user.Email,
		ID:       user.ID,
	}
}

func ConvertModelsToUserDTO(user *models.User) *UserDTO {
	if user == nil {
		return nil
	}
	return &UserDTO{
		UserBase:  *ConvertModelToUserBase(user),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
