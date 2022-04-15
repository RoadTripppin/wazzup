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

func prepareGetUserResponse(user *User) map[string]interface{} {
	userData := map[string]interface{}{
		"id":         user.Id,
		"name":       user.Name,
		"email":      user.Email,
		"profilepic": user.ProfilePic,
	}

	var response = map[string]interface{}{
		"message": "all is fine",
		"user":    userData,
	}

	return response
}
