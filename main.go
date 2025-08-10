/*
Mini Library Manager - A CLI-based library management system

This application provides a simple interface for managing a personal library collection.
Users can perform standard library operations including adding new books, viewing the
complete catalog, borrowing books, and returning previously borrowed items.

Features:
- SQLite database integration for persistent storage
- GORM ORM for simplified database operations
- Interactive command-line menu system
- Book status tracking (available/borrowed)
- Input validation and error handling

Author: Library Manager Development Team
Version: 2.0 with Database Integration
*/

package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

/*
Book represents a single book entity in our library system.
This struct defines the structure of a book record with all necessary
attributes for library management operations.

Fields:
- ID: Unique identifier for each book (auto-generated primary key)
- Title: The name/title of the book
- Author: The author's name who wrote the book
- Year: Publication year of the book
- isBorrowed: Boolean flag indicating if the book is currently borrowed
- CreatedAt: Timestamp when the book was added to the system
- UpdatedAt: Timestamp when the book record was last modified

GORM Tags:
- primaryKey: Marks ID as the primary key for database operations
- not null: Ensures required fields cannot be empty in the database
- default:false: Sets the default value for isBorrowed to false (available)
- json: Defines JSON field names for potential API integration
*/
type Book struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Title      string    `gorm:"not null" json:"title"`
	Author     string    `gorm:"not null" json:"author"`
	Year       int       `gorm:"not null" json:"year"`
	isBorrowed bool      `gorm:"default:false" json:"is_borrowed"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

/*
Global database connection variable.
This variable holds the database connection instance that will be used
throughout the application for all database operations. Using a global
variable ensures we maintain a single connection pool and can access
the database from any function in our application.
*/
var db *gorm.DB

/*
Legacy global slice for book storage.
NOTE: This slice is maintained for backward compatibility but is not
actively used in the database-integrated version. In the original
in-memory implementation, this slice served as the primary data store.
Future versions should consider removing this to avoid confusion.
*/
var book []Book

/*
initDB initializes the database connection and performs necessary setup operations.

This function handles:
1. Establishing a connection to the SQLite database file
2. Configuring GORM logging settings for development/debugging
3. Performing automatic migration to ensure database schema is up-to-date
4. Error handling for connection and migration failures

The function uses SQLite as the database engine, which creates a local
file-based database perfect for single-user applications. The database
file 'library.db' will be created automatically if it doesn't exist.

Panics:
- If database connection fails
- If schema migration encounters errors

Logs:
- Success messages for connection and migration completion
- SQL queries (when logger.Info is enabled)
*/
func initDB() {
	var err error

	// Establish connection to SQLite database
	// The database file will be created automatically if it doesn't exist
	db, err = gorm.Open(sqlite.Open("library.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Enable SQL query logging for debugging
	})

	// Handle connection errors - critical failure that prevents app from running
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	fmt.Println("✅ Database connected successfully!")

	// Perform automatic migration to ensure database schema matches our Book model
	// This will create tables, add missing columns, and update indexes automatically
	err = db.AutoMigrate(&Book{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	fmt.Println("✅ Database migration completed!")
}

/*
displayMenu renders the main application menu to the console.

This function provides the primary user interface for the application,
presenting all available operations in a clear, numbered format.
The menu serves as the navigation hub for users to access different
features of the library management system.

Menu Options:
1. Add Book - Allows users to add new books to the collection
2. List Books - Displays all books with their current status
3. Borrow Book - Marks available books as borrowed
4. Return Book - Returns previously borrowed books to available status
5. Exit - Safely terminates the application

The menu is displayed repeatedly until the user chooses to exit,
providing a continuous interactive experience.
*/
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

/*
addBook handles the process of adding a new book to the library collection.

This function provides an interactive interface for users to input book details
and stores the new book record in the database. It collects essential book
information including title, author, and publication year.

Process Flow:
1. Prompts user for book title, author, and publication year
2. Creates a new Book instance with the provided information
3. Sets the book's borrowed status to false (available by default)
4. Saves the book record to the database using GORM
5. Provides feedback on success or failure of the operation

Input Validation:
- Basic input collection via fmt.Scanln (reads single words)
- Database-level validation through GORM model constraints

Error Handling:
- Database operation errors are caught and displayed to the user
- Failed insertions won't crash the application

Side Effects:
- Creates a new record in the database
- Displays success message with the generated book ID
- Returns to main menu after completion
*/
func addBook() {
	// Declare variables to store user input for the new book
	// These will be populated through console input prompts
	var (
		title  string // Book title entered by user
		author string // Author name entered by user
		year   int    // Publication year entered by user
	)

	fmt.Println("Add a book")
	
	// Collect book details from user input
	// Note: fmt.Scanln reads until whitespace, so multi-word titles need improvement
	fmt.Print("Enter the Title: ")
	fmt.Scanln(&title)
	fmt.Print("Enter the Author: ")
	fmt.Scanln(&author)
	fmt.Print("Enter the Year: ")
	fmt.Scanln(&year)

	// Create a new Book instance with user-provided data
	// isBorrowed is set to false by default (book is available when added)
	book := Book{
		Title:      title,
		Author:     author,
		Year:       year,
		isBorrowed: false,
	}

	// Attempt to save the new book to the database
	// GORM's Create method handles the SQL INSERT operation
	result := db.Create(&book)
	
	// Check if the database operation encountered any errors
	if result.Error != nil {
		fmt.Printf("❌ Error adding book: %v\n", result.Error)
		return // Exit function early on error
	}

	// Display success confirmation with the auto-generated book ID
	fmt.Println("Book added successfully", book.ID)
}

/*
borrowBook manages the book borrowing process for library users.

This function allows users to borrow available books from the library collection.
It displays available books and processes the borrowing transaction by updating
the book's status in the database.

Current Implementation Note:
This function contains legacy code that checks the global slice instead of
querying the database. This needs to be updated to properly integrate with
the database system for full functionality.

Process Flow (Current - Legacy):
1. Checks if there are books in the global slice
2. Prompts user to select a book by number
3. Validates the selection against available books
4. Updates the book's borrowed status
5. Provides confirmation of the borrowing action

Process Flow (Intended - Database):
1. Query database for available books (isBorrowed = false)
2. Display available books to the user
3. Accept user selection and validate input
4. Update selected book's status in database
5. Confirm successful borrowing operation

TODO: Update this function to query database instead of using global slice
*/
func borrowBook() {
	// Legacy check using global slice - should query database instead
	if len(book) == 0 {
		fmt.Println("No books to borrow..")
		return
	}

	var number int
	fmt.Print("Enter the number of the book you want to borrow: ")
	fmt.Scanln(&number)

	// Convert user input to zero-based array index
	// Users see books numbered 1, 2, 3... but arrays are indexed 0, 1, 2...
	limit := number - 1

	// Validate that the selected number corresponds to an existing book
	if limit < 0 || limit >= len(book) {
		fmt.Println("Invalid number. Please try again.")
		return // Return to main menu due to infinite loop structure
	}

	// Check if the selected book is already borrowed
	// Prevent double-borrowing of the same book
	if book[limit].isBorrowed {
		fmt.Printf("%s has already been borrowed.\n", book[limit].Title)
		return
	}

	// Update the book's status to borrowed
	// TODO: This should update the database record, not just the slice
	book[limit].isBorrowed = true
	fmt.Printf("%s has now been borrowed.\n", book[limit].Title)
}

/*
returnBook processes the return of previously borrowed books to the library.

This function handles the book return workflow, allowing users to return
books they have previously borrowed. It updates the book's availability
status and makes it available for other users to borrow.

Current Implementation Note:
Like borrowBook(), this function uses the legacy global slice approach
instead of properly querying and updating the database records.

Process Flow (Current - Legacy):
1. Displays all books currently in the global slice
2. Checks if there are any books to return
3. Prompts user to select a book by number
4. Validates the selection and checks if book was actually borrowed
5. Updates the book's status to available (not borrowed)
6. Confirms the successful return operation

Process Flow (Intended - Database):
1. Query database for currently borrowed books (isBorrowed = true)
2. Display borrowed books to the user
3. Accept and validate user selection
4. Update selected book's status in database to available
5. Provide confirmation of successful return

Debug Features:
- Displays current books slice for debugging purposes
- This should be removed in production version

TODO: Replace slice operations with database queries and updates
*/
func returnBook() {
	// Debug output showing current state of books slice
	// This should be removed in production as it's not user-friendly
	fmt.Println("Current books: ", book)

	// Check if there are any books in the system to return
	if len(book) == 0 {
		fmt.Println("There are no books to return")
		return
	}

	var number int
	fmt.Print("Enter the number of the book you want to return: ")
	fmt.Scanln(&number)

	// Convert user input to zero-based array index for slice access
	limit := number - 1

	// Validate user selection against available books in the system
	if limit < 0 || limit >= len(book) {
		fmt.Println("Invalid number. Please try again.")
		return
	}

	// Verify that the selected book is actually borrowed
	// Prevent returning books that are already available
	if !book[limit].isBorrowed {
		fmt.Printf("%s has already been returned!\n", book[limit].Title)
		return
	}

	// Update book status to available (not borrowed)
	// TODO: This should update the database record instead of just the slice
	book[limit].isBorrowed = false
	fmt.Printf("%s has now been returned.\n", book[limit].Title)
}

/*
listBooks displays all books in the library collection with their current status.

This function provides a comprehensive view of the entire library catalog,
showing each book's details along with its availability status. This helps
users see what books are available for borrowing and which ones are currently
checked out.

Current Implementation Note:
This function uses the legacy global slice instead of querying the database,
which means it won't display books that were added in the current session
through the database-integrated addBook() function.

Display Format:
- Sequential numbering starting from 1
- Book title, author, and publication year
- Current status (Available/Borrowed)

Status Indicators:
- "nBorrowed" - Indicates book is available (typo: should be "Not Borrowed")
- "Borrowed" - Indicates book is currently checked out

TODO: Replace slice iteration with database query to show all books
TODO: Fix "nBorrowed" typo to display "Available" or "Not Borrowed"
TODO: Improve formatting for better readability
*/
func listBooks() {
	// Check if there are any books to display
	if len(book) == 0 {
		fmt.Println("There are no books to display.")
		return
	}

	fmt.Println("\nList of all books:")
	
	// Iterate through all books in the collection
	for i, books := range book {
		// Determine and format the book's current status
		status := "nBorrowed" // Default status (typo: should be "Available")
		if books.isBorrowed {
			status = "Borrowed" // Book is currently checked out
		}
		
		// Display book information in a formatted, user-readable way
		// Note: "Age" should be "Year" for clarity
		fmt.Printf("%d. Name: %s, Author: %s, Age: %d, Status: %s\n",
			i+1, books.Title, books.Author, books.Year, status)
	}
}

/*
menuSelector implements the main application control loop and user interaction handler.

This function serves as the central coordinator of the application, managing
user input and routing requests to appropriate handler functions. It implements
an infinite loop that continues until the user explicitly chooses to exit,
providing a persistent interactive session.

Application Flow:
1. Displays the main menu options
2. Waits for user input selection
3. Routes the request to the appropriate handler function
4. Returns to menu display after operation completion
5. Continues until user selects exit option

Input Handling:
- Accepts integer input corresponding to menu options
- Validates input against available menu choices
- Provides error feedback for invalid selections

Menu Options Routing:
- Option 1: Add Book - Calls addBook() function
- Option 2: List Books - Calls listBooks() function  
- Option 3: Borrow Book - Calls borrowBook() function
- Option 4: Return Book - Calls returnBook() function
- Option 5: Exit - Terminates application gracefully
- Invalid: Displays error message and returns to menu

Legacy Code:
Contains commented-out test code that was used during development
to verify the Book struct functionality. This should be removed
in production versions.

Error Handling:
- Invalid menu selections are handled gracefully
- User is prompted to try again rather than terminating
- Typo in error message ("ty again" should be "try again")
*/
func menuSelector() {
	var option int // Stores user's menu selection

	/* 
	Legacy testing code - used during development to verify Book struct functionality
	This commented code created a test book instance to validate the struct fields
	and their accessibility. Should be removed in production version.
	
	books := Book{
		Title: "Harry Potter",
		Author: "JK Rolins",
		Year: 2002,
		isBorrowed: true,
	}
	
	// Debug output for testing struct field access
	fmt.Println(books.Title)
	fmt.Println(books.Author)
	fmt.Println(books.Year)
	fmt.Println(books.isBorrowed)
	*/

	// Main application loop - continues until user chooses to exit
	for {
		// Display menu options to user
		displayMenu()
		fmt.Print("Option => ")
		fmt.Scanln(&option) // Read user's menu selection

		// Route user selection to appropriate handler function
		switch option {
		case 1:
			addBook() // Handle book addition
		case 2:
			listBooks() // Display all books
		case 3:
			borrowBook() // Process book borrowing
		case 4:
			returnBook() // Process book return
		case 5:
			// Graceful application termination
			fmt.Println("Thank you for using the program.")
			fmt.Println()
			return // Exit the function and terminate application
		default:
			// Handle invalid menu selections
			fmt.Println("Wrong option. Please try again.") // Typo: "ty" should be "try"
		}
	}
}

/*
main is the application entry point and orchestrates the startup sequence.

This function serves as the primary entry point for the Mini Library Manager
application. Currently, it only initializes the menu system, but it should
be enhanced to include database initialization for proper functionality.

Current Implementation:
- Starts the menu selector loop
- No database initialization (missing critical setup)

Recommended Implementation:
1. Initialize database connection
2. Perform any necessary startup tasks
3. Start the interactive menu system
4. Handle graceful shutdown if needed

TODO: Add initDB() call to properly initialize database connection
TODO: Consider adding application startup messages
TODO: Add error handling for critical startup failures
*/
func main() {
	// TODO: Add database initialization here
	// initDB() should be called before starting the menu system
	
	menuSelector() // Start the main interactive loop
}