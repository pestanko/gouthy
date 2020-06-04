package applications

func ConvertModelsToList(list []ApplicationModel) (result []ListApplicationDTO) {
	for _, user := range list {
		item := ListApplicationDTO{
			baseApplicationDTO: convertModelToBase(&user),
		}
		result = append(result, item)
	}
	return result
}

func convertModelToBase(app *ApplicationModel) baseApplicationDTO {
	if app == nil {
		return baseApplicationDTO{}
	}

	return baseApplicationDTO{
		Codename: app.Codename,
		Name:     app.Name,
		ID:       app.ID,
	}
}

func ConvertModelToDTO(app *ApplicationModel) *Application {
	if app == nil {
		return nil
	}
	return &Application{
		baseApplicationDTO: convertModelToBase(app),
		Description:        app.Description,
		ClientId:           app.ClientId,
		State:              app.State,
		Type:               app.Type,
		CreatedAt:          app.CreatedAt,
		UpdatedAt:          app.UpdatedAt,
		DeletedAt:          app.DeletedAt,
	}
}
