package main

import "fmt"

// Model
type Student struct {
	Name  string
	RollNo string
}

func (s *Student) GetName() string {
	return s.Name
}

func (s *Student) SetName(name string) {
	s.Name = name
}

func (s *Student) GetRollNo() string {
	return s.RollNo
}

func (s *Student) SetRollNo(rollNo string) {
	s.RollNo = rollNo
}

// View
type StudentView struct{}

func (sv *StudentView) PrintStudentDetails(studentName, studentRollNo string) {
	fmt.Println("Student:")
	fmt.Println("Name:", studentName)
	fmt.Println("Roll No:", studentRollNo)
}

// Controller
type StudentController struct {
	model *Student
	view  *StudentView
}

func NewStudentController(model *Student, view *StudentView) *StudentController {
	return &StudentController{
		model: model,
		view:  view,
	}
}

func (sc *StudentController) SetStudentName(name string) {
	sc.model.SetName(name)
}

func (sc *StudentController) GetStudentName() string {
	return sc.model.GetName()
}

func (sc *StudentController) SetStudentRollNo(rollNo string) {
	sc.model.SetRollNo(rollNo)
}

func (sc *StudentController) GetStudentRollNo() string {
	return sc.model.GetRollNo()
}

func (sc *StudentController) UpdateView() {
	sc.view.PrintStudentDetails(sc.model.GetName(), sc.model.GetRollNo())
}

func main() {
	// Fetch student record based on roll no from the database
	model := &Student{}
	model.SetName("Robert")
	model.SetRollNo("10")

	// Create a view to write student details on console
	view := &StudentView{}

	controller := NewStudentController(model, view)
	controller.UpdateView()

	// Update the model data
	controller.SetStudentName("John")
	controller.UpdateView()
}
