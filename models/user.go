package models

type Users interface {
	GetId() string
	GetName() string
}

type UserRepository interface {
	AddUser(user Users)
	RemoveUser(user Users)
	FindUserById(ID string) Users
	GetAllUsers() []Users
}
