package schoolsystem

import (
	dbutils "github.com/AndreDrummer/school-system-api/cmd/db/utils"
	"github.com/AndreDrummer/school-system-api/cmd/domain"
)

func getClassRoomInstance() *domain.ClassRoom {
	return domain.NewClassRoom()
}

var classRoomInstance = getClassRoomInstance()

func AllStudents() ([]*domain.Student, error) {
	return classRoomInstance.AllStudents()
}

func AddStudent(student *domain.Student) (bool, error) {
	return classRoomInstance.AddStudent(student)
}

func AddGrade(studentID, grade int) (bool, error) {
	return classRoomInstance.AddGrade(studentID, grade)
}

func UpdateStudent(student *domain.Student) (bool, error) {
	return classRoomInstance.UpdateStudent(student)
}

func RemoveStudent(studentID int) (bool, error) {
	return classRoomInstance.RemoveStudent(studentID)
}

func CalculateAverage(studentID int) (int, error) {
	avg, err := classRoomInstance.CalculateAverage(studentID)

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
	student, ok := classRoomInstance.Students[studentID]

	if ok {
		return student, ok
	} else {
		studentInfo, err := classRoomInstance.GetStudentByID(studentID)

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
	return classRoomInstance.CheckPassOrFail(studentID)
}

func ClearAll() (bool, error) {
	return classRoomInstance.ClearAll()
}
