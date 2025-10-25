package main

import (
	"fmt"
)

func gradingStudents(grades []int) []int {
	for i, grade := range grades {
		if grade >= 38 {
			nextMultiple := ((grade / 5) + 1) * 5
			if nextMultiple-grade < 3 {
				grades[i] = nextMultiple
			}
		}
	}
	return grades
}

func main() {
	var n int
	fmt.Print("Enter number of students: ")
	fmt.Scan(&n)

	grades := make([]int, n)
	fmt.Println("Enter grades:")
	for i := 0; i < n; i++ {
		fmt.Scan(&grades[i])
	}

	finalGrades := gradingStudents(grades)

	fmt.Println("\nRounded Grades:")
	for _, g := range finalGrades {
		fmt.Println(g)
	}
}