// Mini Library Manager
/* User can add, list, borrow and return books and they can exit 
*/

package main

import "fmt"

func main() {
	var option int

	fmt.Println("\nWelcome to the Mini Library Manager")
	fmt.Println()
	fmt.Println("please choose an option")
	fmt.Println("1. Add Book")
	fmt.Println("2. List Books")
	fmt.Println("3. Borrow Book")
	fmt.Println("4. Return Book")
	fmt.Println("5. Exit")
	fmt.Print("Choice => ")
	fmt.Scanln(&option)

	fmt.Println("Option: ", option, "has been chosen")

	fmt.Println()
}
