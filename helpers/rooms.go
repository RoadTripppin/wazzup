package helpers

import (
	"database/sql"
	"log"
	"strings"

	"github.com/RoadTripppin/wazzup/config"
)

func GetMessages(token string, roomid string) map[string]interface{} {
	usr := decodeToken(token)

	if strings.Contains(usr, "Error") {
		return map[string]interface{}{
			"message": usr,
		}
	}

	db := config.InitDB()
	var messages string

	row := db.QueryRow("SELECT messages FROM room WHERE id = ?", roomid)

	if err := row.Scan(&messages); err != nil {
		if err == sql.ErrNoRows {
			log.Println("No room found")
			return map[string]interface{}{"message": "Room not found"}
		}
		//panic(err)
	}

	if messages == "" {
		return map[string]interface{}{"message": "No messages found"}
	}

	defer db.Close()

	var response = prepareGetMessagesResponse(messages)
	return response

}
