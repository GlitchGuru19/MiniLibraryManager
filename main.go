// Mini Library Manager
/* User can add, list, borrow and return books and they can exit
 */

package main

import (
	"fmt"
)

// A simple struct with exported fields
type Book struct{
	Title string 
	Author string 
	Year int
	// isBorrowed bool, comented out to test and run the rest..
}

// Function to display the menu
func displayMenu(){
	fmt.Println("\nWelcome to the Mini Library Manager")
	fmt.Println()
	fmt.Println("please select an option")
	fmt.Println("1. Add Book")
	fmt.Println("2. List Books")
	fmt.Println("3. Borrow Book")
	fmt.Println("4. Return Book")
	fmt.Println("5. Exit")
}

// Function to add a book
func addBook(){
	// Slice of type Book
	var book []Book
	// Here we variables we will pass into the slice 
	var (
		title string
		author string
		year int
		// isBorrowed bool
	)

	fmt.Println("Add a book")
	fmt.Print("Enter the Title: ")
	fmt.Scanln(&title)
	fmt.Print("Enter the Author: ")
	fmt.Scanln(&author)
	fmt.Print("Enter the Year: ")
	fmt.Scanln(&year)

	books := Book{Title: title, Author: author, Year: year}
	book = append(book, books)
	fmt.Println("Book added successfully")
}

// Function to add a book
func borrowBook(){
	// function to borrow a book
}

// Function to Return a borrowed book
func returnBook(){
	// Function to return a book
}

// Function to display/list all the books
func listBooks(book []Book){
	fmt.Println("Function to display all the books")

}


func main() {
	var option int

	/* this is just for testing purposes, as we want to know if it is working or not
	books := Book{
		Title: "Harry Potter",
		Author: "JK Rolins",
		Year: 2002,
		isBorrowed: true,
	} */

	// fmt.Println(books.Title)
	// fmt.Println(books.Author)
	// fmt.Println(books.Year)
	// fmt.Println(books.isBorrowed)


	for {

		displayMenu()
		fmt.Print("Option => ")
		fmt.Scanln(&option)

		switch option {
		case 1:
			addBook()
		case 2:
			//listBooks()
		case 3:
			//borrowBook()
			fmt.Println("Borrow book")
		case 4:
			//returnBook()
			fmt.Println("Return book?")
		case 5:
			fmt.Println("Thank you for using the program.")
			fmt.Println()
			return
		default:
			fmt.Println("Wrong option. Please ty again.")

		}
	}
}
