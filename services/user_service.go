package services

import (
	"strings"

	"github.com/ab22/abcd/models"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// Contains all of the logic for the User model.
type userService struct {
	db *gorm.DB
}

// Int type to define statuses for the User model.
type UserStatus int

// Defines all user statuses.
const (
	Enabled UserStatus = iota
	Disabled
)

// SanitizeUsername trims the username string and converts it to a lowercase
// version of it.
//
// In the future, more checks might be added such as not allowing the username
// to start with numbers, not allowing special characters, etc.
func (s *userService) SanitizeUsername(username string) string {
	sanitizedString := strings.Trim(username, " ")
	sanitizedString = strings.ToLower(sanitizedString)

	return sanitizedString
}

// Search user by id.
func (s *userService) FindById(userId int) (*models.User, error) {
	user := &models.User{}

	err := s.db.
		Where("id = ?", userId).
		First(user).Error

	if err != nil {
		if err != gorm.RecordNotFound {
			return nil, err
		}

		return nil, nil
	}

	return user, nil
}

// Search user by Username.
func (s *userService) FindByUsername(username string) (*models.User, error) {
	user := &models.User{}
	username = s.SanitizeUsername(username)

	err := s.db.
		Where("username = ?", username).
		First(user).Error

	if err != nil {
		if err != gorm.RecordNotFound {
			return nil, err
		}

		return nil, nil
	}

	return user, nil
}

// Searches for a User by Email.
// Returns *models.User instance if it finds it, or nil otherwise.
func (s *userService) FindByEmail(email string) (*models.User, error) {
	user := &models.User{}

	err := s.db.
		Where("email = ?", email).
		First(user).Error

	if err != nil {
		if err != gorm.RecordNotFound {
			return nil, err
		}

		return nil, nil
	}

	return user, nil
}

// Encrypts a password with the default password hasher (bcrypt).
// Returns the hashed password []byte and an error.
func (s *userService) EncryptPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	if err != nil {
		return nil, err
	}

	return hashedPassword, nil
}

// Compares if the hashed password equals the plain text password.
func (s *userService) ComparePasswords(hashedPassword []byte, password string) bool {
	err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	return err == nil
}

// Finds all active users in the database.
func (s *userService) FindAll() ([]models.User, error) {
	var users []models.User
	var err error

	err = s.db.
		Order("id asc").
		Find(&users).Error

	if err != nil {
		if err != gorm.RecordNotFound {
			return nil, err
		}

		return users, nil
	}

	return users, nil
}

// Edit user modifies the basic user's models attributes. The function checks
// if the username changed and if it needs to check if the username is already
// taken.
func (s *userService) Edit(newUser *models.User) error {
	user, err := s.FindById(newUser.Id)
	newUser.Username = s.SanitizeUsername(newUser.Username)

	if err != nil {
		return err
	} else if user == nil {
		return nil
	}

	if user.Username != newUser.Username {
		duplicateUser, err := s.FindByUsername(newUser.Username)
		if err != nil {
			return err
		} else if duplicateUser != nil {
			return DuplicateUsernameError
		}
	}

	user.Username = newUser.Username
	user.FirstName = newUser.FirstName
	user.LastName = newUser.LastName
	user.Email = newUser.Email
	user.Status = newUser.Status
	user.IsAdmin = newUser.IsAdmin
	user.IsTeacher = newUser.IsTeacher

	return s.db.Save(&user).Error
}

// ChangePassword finds a user in the database by userId and changes it's
// password.
func (s *userService) ChangePassword(userId int, password string) error {
	encryptedPassword, err := s.EncryptPassword(password)
	if err != nil {
		return err
	}

	err = s.db.
		Table("users").
		Where("id = ?", userId).
		Update("password", string(encryptedPassword)).Error

	if err != nil {
		if err != gorm.RecordNotFound {
			return err
		}
	} else if s.db.RowsAffected == 0 {
		return RecordNotFound
	}

	return nil
}

// Checks if a user with that email already exists in the database. If it does,
// it returns an error, else it hashes the password, saves the new user
// and returns the user.
func (s *userService) Create(user *models.User) error {
	var err error
	user.Username = s.SanitizeUsername(user.Username)

	result, err := s.FindByUsername(user.Username)
	if err != nil {
		return err
	} else if result != nil {
		return DuplicateUsernameError
	}

	hashedPassword, err := s.EncryptPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	user.Status = int(Enabled)

	err = s.db.Create(&user).Error

	if err != nil {
		return err
	}

	return nil
}

// Delete User sets the DeletedAt time for the user which just hides
// the record from any other select query.
func (s *userService) Delete(userId int) error {
	var err error

	err = s.db.
		Table("users").
		Where("id = ?", userId).
		Delete(models.User{}).Error

	return err
}

// Change email for the specified user by user id.
func (s *userService) ChangeEmail(userId int, email string) error {
	email = strings.Trim(email, " ")

	err := s.db.
		Table("users").
		Where("id = ?", userId).
		Update("email", email).Error

	if err != nil {
		if err != gorm.RecordNotFound {
			return err
		}
	} else if s.db.RowsAffected == 0 {
		return RecordNotFound
	}

	return nil
}

// Change the full name of the user.
func (s *userService) ChangeFullName(userId int, firstName, lastName string) error {
	firstName = strings.Trim(firstName, " ")
	lastName = strings.Trim(lastName, " ")

	err := s.db.
		Table("users").
		Where("id = ?", userId).
		Update("first_name", firstName).
		Update("last_name", lastName).Error

	if err != nil {
		if err != gorm.RecordNotFound {
			return err
		}
	} else if s.db.RowsAffected == 0 {
		return RecordNotFound
	}

	return nil
}
