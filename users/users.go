package users

import (
	"time"

	"github.com/RoadTripppin/wazzup/helpers"
	"github.com/RoadTripppin/wazzup/interfaces"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func prepareToken(user *interfaces.User) string {
	tokenContent := jwt.MapClaims{
		"user_id": user.ID,
		"expiry": time.Now().Add(time.Minute * 60).Unix(),
	}
	jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tokenContent)
	token, err := jwtToken.SignedString([]byte("TokenPassword"))
	helpers.HandleErr(err)

	return token
}

func prepareResponse(user *interfaces.User) map[string]interface{} {
	responseUser := &interfaces.ResponseUser{
		ID: user.ID,
		Name: user.Name,
		Email: user.Email,
	}

	var token = prepareToken(user);
	var response = map[string]interface{}{"message": "all is fine"}
	response["jwt"] = token
	response["data"] = responseUser

	return response
}

func Login(email string, pass string) map[string]interface{} {
	// Connect DB
	valid := helpers.Validation(
		[]interfaces.Validation{
			{Value: email, Valid: "email"},
			{Value: pass, Valid: "password"},
		})
	if valid {
		// Connect DB
		db := helpers.ConnectDB()
		user := &interfaces.User{}
		if db.Where("email = ? ", email).First(&user).RecordNotFound() {
			return map[string]interface{}{"message": "User not found"}
		}
		// Verify password
		passErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass))

		if passErr == bcrypt.ErrMismatchedHashAndPassword && passErr != nil {
			return map[string]interface{}{"message": "Wrong password"}
		}

		defer db.Close()

		var response = prepareResponse(user);

		return response
	} else {
		return map[string]interface{}{"message": "not valid values"}
	}
}

func Register(name string, email string, pass string) map[string]interface{} {
	// Add validation to registration
	valid := helpers.Validation(
		[]interfaces.Validation{
			{Value: name, Valid: "name"},
			{Value: email, Valid: "email"},
			{Value: pass, Valid: "password"},
		})
	if valid {
		// Create registration logic
		// Connect DB
		db := helpers.ConnectDB()
		generatedPassword := helpers.HashAndSalt([]byte(pass))
		user := &interfaces.User{Name: name, Email: email, Password: generatedPassword}
		db.Create(&user)

		defer db.Close()
		var response = prepareResponse(user)

		return response
	} else {
		return map[string]interface{}{"message": "not valid values"}
	}
	
}