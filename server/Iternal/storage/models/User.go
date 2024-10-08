package models

type User struct {
	userId         int64
	hashedPassword string
	role           string
	surname        string
	name           string
	patronymic     string
	dateOfBirth    string
	phoneNumber    string
	email          string
	gender         string
}

func getAllUsers() (output []User) {
	return
}

func getUserById(userId int64) (output User) {
	return
}

func createNewUser() (isCreated bool) {
	return
}
