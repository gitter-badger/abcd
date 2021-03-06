package services

import (
	"github.com/ab22/abcd/models"
	"github.com/jinzhu/gorm"
)

// Contains all logic to handle students in the database.
type studentService struct {
	db *gorm.DB
}

// Finds all active students in the database.
func (s *studentService) FindAll() ([]models.Student, error) {
	var (
		students []models.Student
		err      error
	)

	err = s.db.
		Order("id asc").
		Find(&students).Error

	if err != nil {
		if err != gorm.RecordNotFound {
			return nil, err
		}

		return students, nil
	}

	return students, nil
}

// Search student by id
func (s *studentService) FindById(studentId int) (*models.Student, error) {
	student := &models.Student{}

	err := s.db.
		Where("id = ?", studentId).
		First(student).Error

	if err != nil {
		if err != gorm.RecordNotFound {
			return nil, err
		}

		return nil, nil
	}

	return student, nil
}

// Creates a new student.
func (s *studentService) Create(student *models.Student) error {
	return s.db.Create(student).Error
}

// Edit an existing student
func (s *studentService) Edit(newStudent *models.Student) error {
	student, err := s.FindById(newStudent.Id)

	if err != nil {
		return err
	} else if student == nil {
		return nil
	}

	student.FirstName = newStudent.FirstName
	student.LastName = newStudent.LastName
	student.Status = newStudent.Status

	return s.db.Save(&student).Error
}
