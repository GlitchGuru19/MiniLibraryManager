// Mini Library Manager
/* User can add, list, borrow and return books and they can exit
 */

package main

import (
	"fmt"
)

// A simple struct with exported fields
type Book struct {
	Title      string
	Author     string
	Year       int
	isBorrowed bool //, comented out to test and run the rest..
}

// Slice of type Book declared globally
var book []Book

// Function to display the menu
func displayMenu() {
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
func addBook() {

	// Here we variables we will pass into the slice
	var (
		title  string
		author string
		year   int
	)

	fmt.Println("Add a book")
	fmt.Print("Enter the Title: ")
	fmt.Scanln(&title)
	fmt.Print("Enter the Author: ")
	fmt.Scanln(&author)
	fmt.Print("Enter the Year: ")
	fmt.Scanln(&year)

	books := Book{Title: title, Author: author, Year: year, isBorrowed: false}
	book = append(book, books)
	fmt.Println("Book added successfully")
}

// Function to add a book
func borrowBook() {

	if len(book) == 0 {
		fmt.Println("No books to borrow..")
	}

	var number int
	fmt.Print("Enter the number of the book you want to borrow: ")
	fmt.Scanln(&number)

	// converting to zero index as we want the counting to start from 0
	limit := number - 1

	// checking if the number is valid
	if limit < 0 || limit >= len(book) {
		fmt.Println("Invalid number. PLease try again.")
		return // We return because the main function is in an infinite for loop
	}

	//Check if already borrowed
	if book[limit].isBorrowed {
		fmt.Printf("%s has already been borrowed.\n", book[limit].Title)
		return
	}

	//Mark as Borrowed
	book[limit].isBorrowed = true
	fmt.Printf("%s has now been borrowed.\n", book[limit].Title)

}

// Function to Return a borrowed book
func returnBook() {

	fmt.Println("Current books: ", book)

	if len(book) == 0 {
		fmt.Println("There are no books to return")
		return
	}

	var number int
	fmt.Print("Enter the number of the book you want to return: ")
	fmt.Scanln(&number)

	// converting to zero index as we want the counting to start from 0
	limit := number - 1

	if limit < 0 || limit >= len(book) {
		fmt.Println("Invalid number. Please try again.")
		return
	}

	// Check if already active
	if !book[limit].isBorrowed {
		fmt.Printf("%s has already been returned!\n", book[limit].Title)
		return
	}

	book[limit].isBorrowed = false
	fmt.Printf("%s has now been returned.\n", book[limit].Title)

}

// Function to display/list all the books
func listBooks() {

	if len(book) == 0 {
		fmt.Println("There are no books to display.")
		return
	}

	fmt.Println("\nList of all books:")
	for i, books := range book {
		status := "nBorrowed"
		if books.isBorrowed {
			status = "Borrowed"
		}
		fmt.Printf("%d. Name: %s, Role: %s, Age: %d, Status: %s\n",
			i+1, books.Title, books.Author, books.Year, status)
	}
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
			listBooks()
		case 3:
			borrowBook()
		case 4:
			returnBook()
		case 5:
			fmt.Println("Thank you for using the program.")
			fmt.Println()
			return
		default:
			fmt.Println("Wrong option. Please ty again.")
		}
	}
}
