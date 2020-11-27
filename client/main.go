package main

import (
	"fmt"
	"net/rpc"
)

type Grade struct {
	StudentName string
	SubjectName string
	Total       float64
}

func printMenu() {
	fmt.Println("0) Salir")
	fmt.Println("1) Agregar calificación de una materia")
	fmt.Println("2) Mostrar el promedio de un Alumno")
	fmt.Println("3) Mostrar el promedio general")
	fmt.Println("4) Mostrar el promedio de una materia")
}

func startClient() {
	c, err := rpc.Dial("tcp", "127.0.0.1:4242")
	if err != nil {
		panic(err)
	}
	var opc int64
	for {
		printMenu()
		fmt.Scanln(&opc)

		switch opc {
		case 1:
			var studentName string
			var subjectName string
			var grade float64
			fmt.Print("Nombre de estudiante: ")
			fmt.Scanln(&studentName)
			fmt.Print("Nombre de materia: ")
			fmt.Scanln(&subjectName)
			fmt.Print("calificación: ")
			fmt.Scanln(&grade)

			gradeObj := Grade{
				StudentName: studentName,
				SubjectName: subjectName,
				Total:       grade,
			}

			var result string
			err = c.Call("StudentsGradesBook.PostGrade", gradeObj, &result)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(result)
			}
		case 2:
			var name string
			fmt.Print("Nombre de alumno: ")
			fmt.Scanln(&name)

			var result float64
			err = c.Call("StudentsGradesBook.GetAverageGradeOfStudent", name, &result)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("El promedio de ", name, " es: ", result)
			}
		case 3:
			var result float64
			err = c.Call("StudentsGradesBook.GetAverageGrade", "", &result)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("la calificación promedio es: ", result)
			}
		case 4:
			var name string
			fmt.Print("Nombre de materia: ")
			fmt.Scanln(&name)

			var result float64
			err = c.Call("StudentsGradesBook.GetAverageGradeOfSubject", name, &result)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("El promedio de la materia ", name, " es ", result)
			}
		case 0:
			return
		default:
			fmt.Println("Opcion no valida")
		}

	}
}

func main() {
	startClient()
}
