package helpers

import (
	"github.com/RoadTripppin/wazzup/models"
	"golang.org/x/crypto/bcrypt"
	"github.com/dgrijalva/jwt-go"
	"time"
	"os"
)

func Login(email string, pass string) map[string]interface{} {
	// Connect DB
	valid := Validation(
		[]models.Validation{
			{Value: email, Valid: "email"},
			{Value: pass, Valid: "password"},
		})
	if valid {
		// Connect DB
		db := ConnectDB()
		user := &models.User{}
		if db.Where("email = ? ", email).First(&user).RecordNotFound() {
			return map[string]interface{}{"message": "User not found"}
		}
		// Verify password
		passErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass))

		if passErr == bcrypt.ErrMismatchedHashAndPassword && passErr != nil {
			return map[string]interface{}{"message": "Wrong password"}
		}

		defer db.Close()

		var response = prepareResponse(user)

		return response
	} else {
		return map[string]interface{}{"message": "not valid values"}
	}
}

func Register(name string, email string, pass string) map[string]interface{} {
	// Add validation to registration
	valid := Validation(
		[]models.Validation{
			{Value: name, Valid: "name"},
			{Value: email, Valid: "email"},
			{Value: pass, Valid: "password"},
		})
	if valid {
		// Create registration logic
		// Connect DB
		db := ConnectDB()
		generatedPassword := HashAndSalt([]byte(pass))
		user := &models.User{Name: name, Email: email, Password: generatedPassword}
		db.Create(&user)

		defer db.Close()
		var response = prepareResponse(user)

		return response
	} else {
		return map[string]interface{}{"message": "not valid values"}
	}

}

func prepareToken(user *models.User) string {
	tokenContent := jwt.MapClaims{
		"user_id": user.ID,
		"expiry":  time.Now().Add(time.Minute * 60).Unix(),
	}
	jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tokenContent)
	token, err := jwtToken.SignedString([]byte(os.Getenv("SECRET_KEY")))
	HandleErr(err)

	return token
}

func prepareResponse(user *models.User) map[string]interface{} {
	responseUser := map[string]interface{}{
		"ID":    user.ID,
		"Name":  user.Name,
		"Email": user.Email,
	}

	var token = prepareToken(user)
	var response = map[string]interface{}{"message": "all is fine"}
	response["jwt"] = token
	response["data"] = responseUser

	return response
}