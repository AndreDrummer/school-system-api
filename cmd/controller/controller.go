package schoolsystem

import (
	dbutils "github.com/AndreDrummer/school-system-api/cmd/db/utils"
	"github.com/AndreDrummer/school-system-api/cmd/domain"
)

func getInstance() *domain.ClassRoom {
	return domain.NewClassRoom()
}

var instance = getInstance()

func AllStudents() ([]*domain.Student, error) {
	return instance.AllStudents()
}

func AddStudent(student *domain.Student) (bool, error) {
	return instance.AddStudent(student)
}

func AddGrade(studentID, grade int) (bool, error) {
	return instance.AddGrade(studentID, grade)
}

func UpdateStudent(student *domain.Student) (bool, error) {
	return instance.UpdateStudent(student)
}

func RemoveStudent(studentID int) (bool, error) {
	return instance.RemoveStudent(studentID)
}

func CalculateAverage(studentID int) (int, error) {
	avg, err := instance.CalculateAverage(studentID)

	if err != nil {
		student, ok := GetStudentByID(studentID)

		if !ok {
			return 0, err
		}

		avg = student.GetAverage()
	}

	return avg, nil
}

func GetStudentByID(studentID int) (*domain.Student, bool) {
	student, ok := instance.Students[studentID]

	if ok {
		return student, ok
	} else {
		studentInfo, err := instance.GetStudentByID(studentID)
		if err != nil {
			return nil, false
		}
		studentName, grades := dbutils.GetStudentNameAndGrades(studentInfo)
		student := &domain.Student{
			ID:     studentID,
			Grades: dbutils.ConvertGradesToIntSlice(grades),
			Name:   studentName,
		}
		return student, true
	}
}

func CheckPassOrFail(studentID int) bool {
	return instance.CheckPassOrFail(studentID)
}

func ClearAll() (bool, error) {
	return instance.ClearAll()
}
