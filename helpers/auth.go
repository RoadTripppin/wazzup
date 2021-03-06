package helpers

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/RoadTripppin/wazzup/config"
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
	valid, field := Validation(
		[]models.Validation{
			{Value: email, Valid: "email"},
			{Value: pass, Valid: "password"},
		})
	if valid {
		// Connect DB
		db := config.InitDB()
		//user := &models.User{}
		user := &User{}
		//db.Prepare()
		//stmt, err := repo.Db.Prepare("INSERT INTO user(id, name) values(?,?)")
		row := db.QueryRow("SELECT * FROM user WHERE email = ?", email)

		if err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.ProfilePic, &user.Rooms); err != nil {
			if err == sql.ErrNoRows {
				return map[string]interface{}{"message": "User not found"}
			}
			//panic(err)
		}
		//if db.Where("email = ? ", email).First(&user).RecordNotFound() {
		//stmt, err := db.Prepare("INSERT INTO user(id, name) values(?,?)")
		//CheckErr(err)

		//_, err = stmt.Exec(user.GetId(), user.GetName())
		//CheckErr(err)

		//}
		// Verify password
		//fmt.Println(email)
		passErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass))

		if passErr == bcrypt.ErrMismatchedHashAndPassword && passErr != nil {
			return map[string]interface{}{"message": "Wrong password"}
		}

		defer db.Close()

		// var response = prepareResponse(user)
		var response = prepareAuthResponse(user)

		return response
	} else {
		return map[string]interface{}{"message": "Invalid " + field + " value"}
	}
}

func Register(name string, email string, pass string, pic string) map[string]interface{} {
	// Add validation to registration
	valid, field := Validation(
		[]models.Validation{
			{Value: name, Valid: "name"},
			{Value: email, Valid: "email"},
			{Value: pass, Valid: "password"},
		})

	if valid {
		// Create registration logic
		// Connect DB
		db := config.InitDB()
		//db := ConnectDB()
		// User := &models.User{}
		// db.AutoMigrate(&User)
		generatedPassword := HashAndSalt([]byte(pass))
		user := &User{Id: uuid.New().String(), Name: name, Email: email, Password: generatedPassword, ProfilePic: pic}

		stmt, err := db.Prepare("INSERT INTO user(id, name, email, password, profilepic) values(?,?,?,?,?)")
		CheckErr(err)

		_, err = stmt.Exec(user.Id, user.Name, user.Email, user.Password, user.ProfilePic)
		CheckErr(err)
		//db.Create(&user)

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
		return map[string]interface{}{"message": "Invalid " + field + " value"}
	}
}

func UpdateUser(token string, body *models.Register) map[string]interface{} {
	usr := decodeToken(token)

	if strings.Contains(usr, "Error") {
		return map[string]interface{}{
			"message": usr,
		}
	}

	db := config.InitDB()

	// valid, field := Validation(
	// 	[]models.Validation{
	// 		{Value: body.Name, Valid: "name"},
	// 		{Value: body.Email, Valid: "email"},
	// 		{Value: body.Password, Valid: "password"},
	// 	})

	// if valid {
	user := &User{}

	row := db.QueryRow("SELECT * FROM user WHERE id = ?", usr)
	if err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.ProfilePic, &user.Rooms); err != nil {
		if err == sql.ErrNoRows {
			return map[string]interface{}{"message": "User not found"}
		}
		//panic(err)
	}

	// Check for empty fields
	if body.Name == "" {
		body.Name = user.Name
	}

	if body.ProfilePic == "" {
		body.ProfilePic = user.ProfilePic
	}

	if body.Email == "" {
		body.Email = user.Email
	}

	var passw string
	if body.Password == "" {
		passw = user.Password
		//body.Password = user.Password
	} else {
		passErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

		if passErr == bcrypt.ErrMismatchedHashAndPassword && passErr != nil {
			fmt.Println("Updating Password")
			passw = HashAndSalt([]byte(body.Password))
		} else {
			passw = user.Password
		}

	}

	// Update Query
	stmt, err := db.Prepare("UPDATE user set name = ?, password = ?, profilepic = ? where id = ?")
	CheckErr(err)
	defer stmt.Close()

	res, err := stmt.Exec(body.Name, passw, body.ProfilePic, usr)
	CheckErr(err)

	affected, err := res.RowsAffected()
	CheckErr(err)

	if affected != 1 {
		return map[string]interface{}{"message": "Error while updating user"}
	}

	updatedUser := map[string]interface{}{
		"name":       body.Name,
		"email":      body.Email,
		"password":   passw,
		"profilepic": body.ProfilePic,
	}

	return map[string]interface{}{
		"message": "all is fine",
		"user":    updatedUser,
	}

	// user1 := &models.User{}
	// if db.Where("id = ? ", usr).First(&user1).RecordNotFound() {
	// 	return map[string]interface{}{"message": "User not found"}
	// }

	// user2 := &models.User{}
	// db.Where("email = ? ", body.Email).First(&user2)

	// if user1.ID != user2.ID {
	// 	return map[string]interface{}{"message": "Email already in use. Use different email"}
	// }

	// user1.Name = body.Name
	// user1.Email = body.Email
	// if user1.Password != body.Password {
	// 	user1.Password = HashAndSalt([]byte(body.Password))
	// }
	// user1.ProfilePic = body.ProfilePic

	// if dbc := db.Model(&user1).Updates(&user1); dbc.Error != nil {
	// 	fmt.Printf(dbc.Error.Error())
	// 	return map[string]interface{}{"message": "Error while updating user"}
	// 	// response["userID"] = usr
	// }

	// updatedUser := map[string]interface{}{
	// 	"name":       user1.Name,
	// 	"email":      user1.Email,
	// 	"password":   user1.Password,
	// 	"profilepic": user1.ProfilePic,
	// }

	// return map[string]interface{}{
	// 	"message": "all is fine",
	// 	"user":    updatedUser,
	// }

	// } else {
	// 	return map[string]interface{}{"message": "Invalid " + field + " value"}
	// }
}

func DeleteUser(token string, body *models.Register) map[string]interface{} {
	usr := decodeToken(token)

	if strings.Contains(usr, "Error") {
		return map[string]interface{}{
			"message": usr,
		}
	}

	db := ConnectDB()
	user1 := &models.User{}
	if db.Where("id = ? ", usr).First(&user1).RecordNotFound() {
		return map[string]interface{}{"message": "User not found"}
	}

	if dbc := db.Delete(&user1); dbc.Error != nil {
		return map[string]interface{}{
			"message": "Error while deleting user",
		}
	}

	return map[string]interface{}{
		"message": "all is fine",
	}
}
