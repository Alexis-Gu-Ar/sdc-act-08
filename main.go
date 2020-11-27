/*
- Hello:
    - Recibe el nombre del usuario en un *`string`.*
    - Retorna la palabra "Hello" mas el *`string`* recibido.
- Negate:
    - Recibe un número en un `int64`.
    - Retorna el negativo del `int64` recibido.
- Min:
    - Recibe un slice de tipo `int64`.
    - Retorna el mínimo de entre los elementos del slice recibido.
*/

package main

import (
	"errors"
	"fmt"
	"net"
	"net/rpc"
)

func average(arr []float64) (float64, error) {
	if len(arr) == 0 {
		return 0, errors.New("No hay nada que promediar para la peticion")
	}
	totalSum := 0.0
	for _, num := range arr {
		totalSum += num
	}
	return totalSum / float64(len(arr)), nil
}

type Grade struct {
	StudentName string
	SubjectName string
	Total       float64
}

type Student struct {
	grades map[string]float64
}

func (student *Student) getAverageGrade() (float64, error) {
	var grades []float64
	for _, grade := range student.grades {
		grades = append(grades, grade)
	}
	return average(grades)
}

type StudentsGradesBook struct {
	students map[string]Student
}

func (studentsGradesBook *StudentsGradesBook) PostGrade(grade Grade, reply *string) error {
	student, ok := studentsGradesBook.students[grade.StudentName]
	if !ok {
		student = Student{
			grades: map[string]float64{},
		}
		studentsGradesBook.students[grade.StudentName] = student
	}
	_, ok = student.grades[grade.SubjectName]
	if ok {
		return errors.New("El estudiante ya esta calificado para la materia" + grade.SubjectName)
	}
	student.grades[grade.SubjectName] = grade.Total
	*reply = "La calificación se publico con exito"
	return nil
}

func (studentGradesbook *StudentsGradesBook) GetAverageGradeOfStudent(name string, reply *float64) error {
	student, ok := studentGradesbook.students[name]

	if ok {
		var err error
		*reply, err = student.getAverageGrade()
		return err
	}

	return errors.New("Student doesn't exists")
}

func (studentGradesbook *StudentsGradesBook) GetAverageGrade(args string, reply *float64) error {
	var averages []float64
	for _, student := range studentGradesbook.students {
		average, err := student.getAverageGrade()
		if err == nil {
			averages = append(averages, average)
		}
	}
	var err error
	*reply, err = average(averages)
	return err
}

func (StudentsGradesBook *StudentsGradesBook) GetAverageGradeOfSubject(name string, reply *float64) error {
	var grades []float64
	for _, student := range StudentsGradesBook.students {
		grade, ok := student.grades[name]
		if ok {
			grades = append(grades, grade)
		}
	}
	var err error
	*reply, err = average(grades)
	return err
}

/*
func (s *Server) Hello(name string, reply *string) error {
	*reply = "Hello " + name
	return nil
}

func (s *Server) Negate(num int64, reply *int64) error {
	*reply = -num
	return nil
}

func (s *Server) Min(slice []int64, reply *int64) error {
	if len(slice) > 0 {
		min := slice[0]
		for _, myInt := range slice {
			if myInt < min {
				min = myInt
			}
		}
		*reply = min
		return nil
	} else {
		return errors.New("Empty Slice")
	}
}
*/

func server() {
	book := StudentsGradesBook{
		students: map[string]Student{},
	}
	rpc.Register(&book)
	ln, err := net.Listen("tcp", ":4242")
	if err != nil {
		panic(err)
	}

	for {
		c, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		go rpc.ServeConn(c)
	}

}

func main() {
	go server()
	var input string
	fmt.Scanln(&input)
}
