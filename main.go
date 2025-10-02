
/*
Mini Library Manager - A CLI-based library management system

This application provides a simple interface for managing a personal library collection.
Users can perform standard library operations including adding new books, viewing the
complete catalog, borrowing books, and returning previously borrowed items.

Features:

- Interactive command-line menu system
- Book status tracking (available/borrowed)
- Input validation and error handling

Author: Glitch Guru 19
Version: 1.0
*/

package main

import "fmt"

type Book struct {
	Title    string
	Author   string
	Year     int
	IsBorrowed bool
}

func displayMenu() {
	fmt.Println("\nWelcome to the Mini Library Manager")
	fmt.Println()
	fmt.Println("Please select an option:")
	fmt.Println("1. Add Book")
	fmt.Println("2. List Books")
	fmt.Println("3. Borrow Book")
	fmt.Println("4. Return Book")
	fmt.Println("5. Exit")
} 

func main() {
	displayMenu()
}