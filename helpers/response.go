package helpers

func prepareAuthResponse(user *User) map[string]interface{} {

	userData := map[string]interface{}{
		"id":         user.Id,
		"name":       user.Name,
		"email":      user.Email,
		"password":   user.Password,
		"profilepic": user.ProfilePic,
	}

	var token = prepareToken(user.Id)
	var response = map[string]interface{}{"message": "all is fine"}
	response["token"] = token
	response["user"] = userData

	return response
}

func prepareSearchUserResponse(users []User) map[string]interface{} {
	// for _, user = range users {

	// }
	// userData := map[string]interface{}{
	// 	"id":    user.Id,
	// 	"name":  user.Name,
	// 	"email": user.Email,
	// }

	// userData := map[string]interface{}{
	// 	"users": users,
	// }

	var usersData []map[string]interface{}

	for _, user := range users {
		tempuser := map[string]interface{}{
			"id":    user.Id,
			"name":  user.Name,
			"email": user.Email,
		}

		usersData = append(usersData, tempuser)
	}

	var response = map[string]interface{}{
		"message": "all is fine",
		"users":   usersData,
	}

	return response
}
