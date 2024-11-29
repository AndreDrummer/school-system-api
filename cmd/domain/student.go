package domain

type Student struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Grades []int  `json:"grades"`
}

func (s *Student) AddGrade(grade int) {
	s.Grades = append(s.Grades, grade)
}

func (s *Student) GetAverage() int {
	studentGrades := s.Grades
	if len(studentGrades) == 0 {
		return 0
	}

	var sumUpGrades int
	for _, grade := range studentGrades {
		sumUpGrades = sumUpGrades + grade
	}
	average := sumUpGrades / len(studentGrades)
	return average
}
