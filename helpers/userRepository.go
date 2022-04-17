package helpers

import (
	"database/sql"
	"log"

	"github.com/RoadTripppin/wazzup/models"
)

type User struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	ProfilePic string `json:"profile"`
	Rooms      string `json:"rooms"`
}

func (user *User) GetId() string {
	return user.Id
}

func (user *User) GetName() string {
	return user.Name
}

type UserRepository struct {
	Db *sql.DB
}

func (repo *UserRepository) AddUser(user models.Users) {
	stmt, err := repo.Db.Prepare("INSERT INTO user(id, name) values(?,?)")
	CheckErr(err)

	_, err = stmt.Exec(user.GetId(), user.GetName())
	CheckErr(err)
}

func (repo *UserRepository) RemoveUser(user models.Users) {
	stmt, err := repo.Db.Prepare("DELETE FROM user WHERE id = ?")
	CheckErr(err)

	_, err = stmt.Exec(user.GetId())
	CheckErr(err)
}

func (repo *UserRepository) FindUserById(ID string) models.Users {

	row := repo.Db.QueryRow("SELECT id, name FROM user where id = ? LIMIT 1", ID)

	var user User

	if err := row.Scan(&user.Id, &user.Name); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		panic(err)
	}

	return &user

}

func (repo *UserRepository) GetAllUsers() []models.Users {

	rows, err := repo.Db.Query("SELECT id, name FROM user")

	if err != nil {
		log.Fatal(err)
	}
	var users []models.Users
	defer rows.Close()
	for rows.Next() {
		var user User
		rows.Scan(&user.Id, &user.Name)
		users = append(users, &user)
	}

	return users
}
