package users

func ConvertModelsToList(list []UserModel) []ListUserDTO {
	var listUsers []ListUserDTO

	for _, user := range list {
		item := ListUserDTO{
			baseUserDTO: *convertModelToUserBase(&user),
		}
		listUsers = append(listUsers, item)
	}
	return listUsers
}

func convertModelToUserBase(user *UserModel) *baseUserDTO {
	if user == nil {
		return nil
	}

	return &baseUserDTO{
		Username: user.Username,
		Name:     user.Name,
		Email:    user.Email,
		ID:       user.ID,
	}
}

func ConvertModelToDTO(user *UserModel) *User {
	if user == nil {
		return nil
	}
	return user.ToEntity()
}
