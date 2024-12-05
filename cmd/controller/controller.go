package schoolsystem

import (
	"sort"

	dbutils "github.com/AndreDrummer/school-system-api/cmd/db/utils"
	"github.com/AndreDrummer/school-system-api/cmd/domain"
)

func getClassRoomInstance() *domain.ClassRoom {
	return domain.NewClassRoom()
}

var classRoomInstance = getClassRoomInstance()

func getNextAvailableID() (int, error) {
	studentIDs := make([]int, 0)
	students, err := AllStudents()

	if err != nil {
		return 0, err
	}

	for _, student := range students {
		studentID := student.ID
		studentIDs = append(studentIDs, studentID)
	}

	sort.Ints(studentIDs)
	startID := 1
	for _, ID := range studentIDs {
		if ID-startID == 0 {
			startID++
			continue
		} else {
			return startID, nil
		}
	}

	return len(studentIDs) + 1, nil
}

func AllStudents() ([]*domain.Student, error) {
	return classRoomInstance.AllStudents()
}

func AddStudent(student *domain.Student) (bool, error) {
	studentID, err := getNextAvailableID()

	if err != nil {
		return false, err
	}

	student.ID = studentID
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
