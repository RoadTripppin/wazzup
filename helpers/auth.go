package helpers

import (
	"fmt"
	"time"

	"github.com/RoadTripppin/wazzup/models"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func prepareToken(userID string) string {
	tokenContent := jwt.MapClaims{
		"user_id": userID,
		"expiry":  time.Now().Add(time.Hour * 60).Unix(),
	}
	jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tokenContent)
	token, err := jwtToken.SignedString([]byte("TokenPassword"))
	HandleErr(err)

	return token
}

func decodeToken(token string) string {
	tkn, err := jwt.Parse(token, func(tkn *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := tkn.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Error: unexpected signing method: %v", tkn.Header["alg"])
		}

		return []byte("TokenPassword"), nil
	})

	if err != nil {
		return err.Error()
	} else if !tkn.Valid {
		return "Error: Invalid token"
	}

	var user string
	if claims, ok := tkn.Claims.(jwt.MapClaims); ok && tkn.Valid {
		user = claims["user_id"].(string)
	} else {
		fmt.Println(err)
	}

	return user
}

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

		// var response = prepareResponse(user)
		var response = prepareAuthResponse(user)

		return response
	} else {
		return map[string]interface{}{"message": "Invalid value"}
	}
}

func Register(name string, email string, pass string, pic string) map[string]interface{} {
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
		User := &models.User{}
		db.AutoMigrate(&User)
		generatedPassword := HashAndSalt([]byte(pass))
		user := &models.User{Name: name, Email: email, Password: generatedPassword, ProfilePic: pic}
		db.Create(&user)

		// Error handling for creation ---------
		// var errMessage = createdUser.Error
		// if createdUser.Error != nil {
		// 	fmt.Println(errMessage)
		// }
		// -----------

		defer db.Close()
		var response = prepareAuthResponse(user)

		return response
	} else {
		return map[string]interface{}{"message": "Invalid value"}
	}
}
