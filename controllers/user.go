package controllers

import "transaction_project/models"

type User struct {
	userService *models.UserService
}

func (uC *User) NewModel(email, password, last string, middle, first, phone interface{}) (*models.User, error) {
	var newMiddle string
	if middle != nil {
		newMiddle = middle.(string)
	}

	var newFirst string
	if first != nil {
		newFirst = first.(string)
	}

	var newPhone string
	if phone != nil {
		newPhone = phone.(string)
	}

	newUser := &models.User{
		Email:    email,
		Password: password,
		Last:     last,
		Middle:   newMiddle,
		First:    newFirst,
		Phone:    newPhone,
	}
	return newUser, uC.userService.Create(newUser)
}

func (uC *User) UpdateModel(id uint, email, password, last, middle, first, phone interface{}) (*models.User, error) {
	// Get the exist model
	user, err := uC.userService.ReadByID(id)
	if err != nil {
		return nil, models.ErrNotFound
	}

	// Update exist model email
	updateEmail := user.Email
	if email != nil {
		updateEmail = email.(string)
	}
	user.Email = updateEmail

	// Update exist model password
	updatePassword := user.Password
	if password != nil {
		updatePassword = password.(string)
	}
	user.Password = updatePassword

	// Update exist model last
	updateLast := user.Last
	if last != nil {
		updateLast = last.(string)
	}
	user.Last = updateLast

	// Update exist model middle
	updateMiddle := user.Middle
	if middle != nil {
		updateMiddle = middle.(string)
	}
	user.Middle = updateMiddle

	// Update exist model first
	updateFirst := user.First
	if first != nil {
		updateFirst = first.(string)
	}
	user.First = updateFirst

	// Update exist model phone
	updatePhone := user.Phone
	if phone != nil {
		updatePhone = phone.(string)
	}
	user.Phone = updatePhone

	return user, uC.userService.Update(user)
}

func NewUserController(userService *models.UserService) *User {
	return &User{
		userService: userService,
	}
}
