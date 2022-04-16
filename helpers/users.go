package helpers

import "github.com/RoadTripppin/wazzup/models"

func GetUser(email string) map[string]interface{} {
	valid, field := Validation(
		[]models.Validation{
			{Value: email, Valid: "email"},
		})

	if valid {
		db := ConnectDB()
		user := &User{}
		if db.Where("email = ? ", email).First(&user).RecordNotFound() {
			return map[string]interface{}{"message": "User not found"}
		}

		defer db.Close()

		var response = prepareGetUserResponse(user)
		return response
	} else {
		return map[string]interface{}{"message": "Invalid " + field + " value"}
	}
}
