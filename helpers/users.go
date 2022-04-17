package helpers

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/RoadTripppin/wazzup/config"
)

func SearchUser(token string, query string) map[string]interface{} {
	usr := decodeToken(token)

	if strings.Contains(usr, "Error") {
		return map[string]interface{}{
			"message": usr,
		}
	}

	db := config.InitDB()

	fmt.Println(query)
	rows, err := db.Query("SELECT id, name, email FROM user WHERE email LIKE ?", "%"+query+"%@%")

	var users []User
	for rows.Next() {
		var user User

		err = rows.Scan(&user.Id, &user.Name, &user.Email)
		if err != nil {
			fmt.Println(err)
			return map[string]interface{}{"message": "User not found"}
		}

		users = append(users, user)
	}

	defer db.Close()

	if users == nil {
		return map[string]interface{}{"message": "No user found"}
	}

	fmt.Println(users)
	var response = prepareSearchUserResponse(users)
	return response

}

func GetInteractedUsers(token string) map[string]interface{} {
	usr := decodeToken(token)

	if strings.Contains(usr, "Error") {
		return map[string]interface{}{
			"message": usr,
		}
	}

	db := config.InitDB()
	user := &User{}
	fmt.Println(usr)
	row := db.QueryRow("SELECT email, rooms FROM user WHERE id LIKE ?", usr)

	if err := row.Scan(&user.Email, &user.Rooms); err != nil {
		if err == sql.ErrNoRows {
			return map[string]interface{}{"message": "No User found"}
		}
		//panic(err)
	}

	if user.Rooms == "" {
		return map[string]interface{}{"message": "No chats found"}
	}

	rooms := strings.Split(user.Rooms, ":;;:")

	var response = prepareGetInteractedUsersResponse(rooms)
	return response
}
