package domain

import (
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/AndreDrummer/school-system-api/cmd/db"
	dbutils "github.com/AndreDrummer/school-system-api/cmd/db/utils"
)

type ClassRoom struct {
	Students            map[int]*Student
	database            *db.DB
	StudentsQty         int
	MinimumPassingGrade int
}

func NewClassRoom() *ClassRoom {
	return &ClassRoom{
		Students:            make(map[int]*Student),
		database:            db.GetDB(),
		StudentsQty:         0,
		MinimumPassingGrade: 60,
	}
}

func (c *ClassRoom) GetStudentByID(studentID int) (string, error) {
	return c.database.GetByID(studentID)
}

func (c *ClassRoom) AllStudents() ([]*Student, error) {
	content, err := c.database.GetAll()

	if err != nil {
		return []*Student{}, err
	}

	list := make([]*Student, len(content))

	if len(content) > 0 {
		for i, v := range content {
			studentIDString := strings.Split(v, " ")[0]
			studentID, err := strconv.Atoi(studentIDString)
			if v == "" {
				continue
			}
			studentName, grades := dbutils.GetStudentNameAndGrades(v)

			if err != nil {
				log.Fatal(err)
			}

			newStudent := &Student{
				ID:     studentID,
				Grades: dbutils.ConvertGradesToIntSlice(grades),
				Name:   studentName,
			}

			list[i] = newStudent
		}
	}

	return list, nil
}

func (c *ClassRoom) AddStudent(newStudent *Student) (bool, error) {
	c.Students[newStudent.ID] = newStudent
	c.StudentsQty++

	return c.database.Insert(*newStudent)
}

func (c *ClassRoom) AddGrade(studentID, grade int) (bool, error) {
	student := c.Students[studentID]
	student.AddGrade(grade)

	return c.database.Update(studentID, *student)
}

func (c *ClassRoom) UpdateStudent(student *Student) (bool, error) {
	return c.database.Update(student.ID, *student)
}

func (c *ClassRoom) RemoveStudent(studentID int) (bool, error) {
	delete(c.Students, studentID)

	return c.database.Delete(studentID)
}

func (c *ClassRoom) CalculateAverage(studentID int) (int, error) {
	student, ok := c.Students[studentID]

	if !ok {
		return 0, errors.New("not found... searching on DB")
	}

	return student.GetAverage(), nil
}

func (c *ClassRoom) CheckPassOrFail(studentID int) bool {
	student, ok := c.Students[studentID]

	if ok {
		return student.GetAverage() > c.MinimumPassingGrade
	}

	return false
}

func (c *ClassRoom) ClearAll() (bool, error) {
	return c.database.Clear()
}
