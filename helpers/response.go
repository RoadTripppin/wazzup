package helpers

import (
	"encoding/json"
	"fmt"
	"strings"
)

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

func prepareGetInteractedUsersResponse(rooms []string) map[string]interface{} {
	var roomsData []map[string]interface{}
	for _, room := range rooms {
		var tempRoom map[string]interface{}
		fmt.Println(room)
		json.Unmarshal([]byte(room), &tempRoom)
		fmt.Println(tempRoom)

		roomsData = append(roomsData, tempRoom)
	}

	// st := []byte(`"{"id":"room1","name":"Nameofroom1"}:;;:{"id":"room2","name":"Nameofroom2"}"`)

	// // // var temp map[string]interface{}
	// // temp, _ := json.Marshal(st)
	// // fmt.Println(string(temp))
	// temp := map[string]interface{}{
	// 	"id":   "room1",
	// 	"name": "Name of room 1",
	// }
	// tempjson, _ := json.Marshal(temp)

	// fmt.Println(string(tempjson))

	var response = map[string]interface{}{
		"message": "all is fine",
		"rooms":   roomsData,
	}

	return response
}

func prepareGetMessagesResponse(messages string) map[string]interface{} {
	var messageArray []string
	if strings.Contains(messages, ":;;:") {
		messageArray = strings.Split(messages, ":;;:")
	} else {
		messageArray = append(messageArray, messages)
	}

	var messageData []map[string]interface{}
	for _, msg := range messageArray {
		var tempmsg map[string]interface{}

		json.Unmarshal([]byte(msg), &tempmsg)
		// fmt.Println(tempRoom)

		messageData = append(messageData, tempmsg)
	}

	var response = map[string]interface{}{
		"message": "all is fine",
		"rooms":   messageData,
	}

	return response
}
