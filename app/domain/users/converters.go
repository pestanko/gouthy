package users

func ConvertModelsToUserList(list []User) []ListUserDTO {
	var listUsers []ListUserDTO

	for _, user := range list {
		item := ListUserDTO{
			UserBaseDTO: *convertModelToUserBase(&user),
		}
		listUsers = append(listUsers, item)
	}
	return listUsers
}

func convertModelToUserBase(user *User) *UserBaseDTO {
	if user == nil {
		return nil
	}

	return &UserBaseDTO{
		Username: user.Username,
		Name:     user.Name,
		Email:    user.Email,
		ID:       user.ID,
	}
}

func ConvertModelsToUserDTO(user *User) *UserDTO {
	if user == nil {
		return nil
	}
	return &UserDTO{
		UserBaseDTO: *convertModelToUserBase(user),
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}
