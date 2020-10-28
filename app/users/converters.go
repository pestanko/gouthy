package users

func ConvertModelsToList(list []User) []ListUserDTO {
	var listUsers []ListUserDTO

	for _, user := range list {
		item := ListUserDTO{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		}
		listUsers = append(listUsers, item)
	}
	return listUsers
}

func convertModelToUserBase(user *User) *baseUserDTO {
	if user == nil {
		return nil
	}

	return &baseUserDTO{
		Username: user.Username,
		Name:     user.Name,
		Email:    user.Email,
	}
}

func ConvertModelToDTO(user *User) *UserDTO {
	if user == nil {
		return nil
	}
	return user.ToEntity()
}
