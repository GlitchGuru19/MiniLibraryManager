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
Version: 1.0 with Database Integration
*/

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
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
- IsBorrowed: Boolean flag indicating if the book is currently borrowed (exported field)
- CreatedAt: Timestamp when the book was added to the system
- UpdatedAt: Timestamp when the book record was last modified

GORM Tags:
- primaryKey: Marks ID as the primary key for database operations
- not null: Ensures required fields cannot be empty in the database
- default:false: Sets the default value for IsBorrowed to false (available)
- json: Defines JSON field names for potential API integration
*/
type Book struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Title      string    `gorm:"not null" json:"title"`
	Author     string    `gorm:"not null" json:"author"`
	Year       int       `gorm:"not null" json:"year"`
	IsBorrowed bool      `gorm:"default:false;column:is_borrowed" json:"is_borrowed"`
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
	fmt.Println("Please select an option:")
	fmt.Println("1. Add Book")
	fmt.Println("2. List Books")
	fmt.Println("3. Borrow Book")
	fmt.Println("4. Return Book")
	fmt.Println("5. Exit")
}

/*
readInput reads a line of input from the user, handling multi-word strings properly.
*/
func readInput(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

/*
readInt reads an integer from user input with basic validation.
*/
func readInt(prompt string) int {
	for {
		input := readInput(prompt)
		if value, err := strconv.Atoi(input); err == nil {
			return value
		}
		fmt.Println("❌ Please enter a valid number.")
	}
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
- Handles multi-word titles and authors properly
- Basic validation for year input

Error Handling:
- Database operation errors are caught and displayed to the user
- Failed insertions won't crash the application

Side Effects:
- Creates a new record in the database
- Displays success message with the generated book ID
- Returns to main menu after completion
*/
func addBook() {
	fmt.Println("\n📖 Add a New Book")
	fmt.Println("-----------------")

	// Collect book details from user input using improved input handling
	title := readInput("Enter the Title: ")
	if title == "" {
		fmt.Println("❌ Title cannot be empty.")
		return
	}

	author := readInput("Enter the Author: ")
	if author == "" {
		fmt.Println("❌ Author cannot be empty.")
		return
	}

	year := readInt("Enter the Year: ")
	if year < 1000 || year > 2030 {
		fmt.Println("❌ Please enter a valid year between 1000 and 2030.")
		return
	}

	// Create a new Book instance with user-provided data
	// IsBorrowed is set to false by default (book is available when added)
	newBook := Book{
		Title:      title,
		Author:     author,
		Year:       year,
		IsBorrowed: false,
	}

	// Attempt to save the new book to the database
	// GORM's Create method handles the SQL INSERT operation
	result := db.Create(&newBook)

	// Check if the database operation encountered any errors
	if result.Error != nil {
		fmt.Printf("❌ Error adding book: %v\n", result.Error)
		return // Exit function early on error
	}

	// Display success confirmation with the auto-generated book ID
	fmt.Printf("✅ Book added successfully! (ID: %d)\n", newBook.ID)
}

/*
borrowBook manages the book borrowing process for library users.

This function allows users to borrow available books from the library collection.
It queries the database for available books, displays them to the user, and
processes the borrowing transaction by updating the book's status in the database.

Process Flow:
1. Query database for available books (IsBorrowed = false)
2. Display available books to the user with clear numbering
3. Accept user selection and validate input range
4. Update selected book's status in database to borrowed
5. Provide confirmation of the borrowing operation

Database Operations:
- SELECT query to fetch available books
- UPDATE query to mark book as borrowed

Error Handling:
- Database query failures are caught and displayed
- Invalid user selections are handled gracefully
- Prevents borrowing when no books are available

Input Validation:
- Ensures selection is within valid range of available books
- Provides clear feedback for invalid selections
*/
func borrowBook() {
	var availableBooks []Book

	// Query database for books that are currently available (not borrowed)
	result := db.Where("is_borrowed = ?", false).Order("title").Find(&availableBooks)
	if result.Error != nil {
		fmt.Printf("❌ Error fetching available books: %v\n", result.Error)
		return
	}

	// Check if there are any books available to borrow
	if len(availableBooks) == 0 {
		fmt.Println("📚 No books available to borrow.")
		return
	}

	fmt.Println("\n📖 Available Books:")
	fmt.Println("-------------------")

	// Display all available books with clear numbering for user selection
	for i, book := range availableBooks {
		fmt.Printf("%d. \"%s\" by %s (%d)\n",
			i+1, book.Title, book.Author, book.Year)
	}

	number := readInt("\nEnter the number of the book you want to borrow: ")

	// Convert user input to zero-based array index and validate selection
	if number < 1 || number > len(availableBooks) {
		fmt.Println("❌ Invalid selection. Please try again.")
		return
	}

	// Get the selected book from the available books list
	selectedBook := availableBooks[number-1]

	// Update the book's borrowed status in the database
	updateResult := db.Model(&selectedBook).Update("is_borrowed", true)
	if updateResult.Error != nil {
		fmt.Printf("❌ Error borrowing book: %v\n", updateResult.Error)
		return
	}

	// Confirm successful borrowing operation
	fmt.Printf("✅ Successfully borrowed \"%s\" by %s!\n", selectedBook.Title, selectedBook.Author)
}

/*
returnBook processes the return of previously borrowed books to the library.

This function handles the book return workflow, allowing users to return
books they have previously borrowed. It queries the database for borrowed books,
displays them to the user, and updates the book's availability status in the database.

Process Flow:
1. Query database for currently borrowed books (IsBorrowed = true)
2. Display borrowed books to the user with clear numbering
3. Accept and validate user selection
4. Update selected book's status in database to available (IsBorrowed = false)
5. Provide confirmation of successful return operation

Database Operations:
- SELECT query to fetch borrowed books
- UPDATE query to mark book as available

Error Handling:
- Database query failures are caught and reported
- Invalid user selections are handled gracefully
- Prevents return operations when no books are borrowed

User Experience:
- Shows only borrowed books for selection
- Provides clear confirmation of return operations
- Handles edge cases (no borrowed books) gracefully
*/
func returnBook() {
	var borrowedBooks []Book

	// Query database for books that are currently borrowed
	result := db.Where("is_borrowed = ?", true).Order("title").Find(&borrowedBooks)
	if result.Error != nil {
		fmt.Printf("❌ Error fetching borrowed books: %v\n", result.Error)
		return
	}

	// Check if there are any books currently borrowed
	if len(borrowedBooks) == 0 {
		fmt.Println("📚 No books currently borrowed.")
		return
	}

	fmt.Println("\n📤 Borrowed Books:")
	fmt.Println("------------------")

	// Display all borrowed books with clear numbering for user selection
	for i, book := range borrowedBooks {
		fmt.Printf("%d. \"%s\" by %s (%d)\n",
			i+1, book.Title, book.Author, book.Year)
	}

	number := readInt("\nEnter the number of the book you want to return: ")

	// Convert user input to zero-based array index and validate selection
	if number < 1 || number > len(borrowedBooks) {
		fmt.Println("❌ Invalid selection. Please try again.")
		return
	}

	// Get the selected book from the borrowed books list
	selectedBook := borrowedBooks[number-1]

	// Update the book's borrowed status in the database to available
	updateResult := db.Model(&selectedBook).Update("is_borrowed", false)
	if updateResult.Error != nil {
		fmt.Printf("❌ Error returning book: %v\n", updateResult.Error)
		return
	}

	// Confirm successful return operation
	fmt.Printf("✅ Successfully returned \"%s\" by %s!\n", selectedBook.Title, selectedBook.Author)
}

/*
listBooks displays all books in the library collection with their current status.

This function provides a comprehensive view of the entire library catalog by
querying the database for all books and displaying each book's details along
with its current availability status. This helps users see what books are
available for borrowing and which ones are currently checked out.

Process Flow:
1. Query database for all books in the collection
2. Check if any books exist in the database
3. Display each book with formatted information including status
4. Show total count of books in the collection

Display Format:
- Sequential numbering starting from 1 (user-friendly)
- Book title, author, and publication year
- Current status with visual indicators
- Total book count summary

Status Indicators:
- "Available ✅" - Book can be borrowed
- "Borrowed 📤" - Book is currently checked out

Database Operations:
- SELECT query to fetch all books ordered by creation date (newest first)

Error Handling:
- Database query failures are caught and reported
- Empty collection is handled gracefully with informative message
*/
func listBooks() {
	var allBooks []Book

	// Query database for all books, ordered by creation date (newest first)
	result := db.Order("created_at desc").Find(&allBooks)
	if result.Error != nil {
		fmt.Printf("❌ Error fetching books from database: %v\n", result.Error)
		return
	}

	// Check if there are any books in the database
	if len(allBooks) == 0 {
		fmt.Println("📚 No books in the library collection.")
		return
	}

	fmt.Println("\n📖 Library Collection:")
	fmt.Println("======================")

	// Iterate through all books and display their information
	for i, book := range allBooks {
		// Determine the book's current status with visual indicators
		status := "Available ✅" // Book is available for borrowing
		if book.IsBorrowed {
			status = "Borrowed 📤" // Book is currently checked out
		}

		// Display book information in a formatted, user-readable way
		fmt.Printf("%d. \"%s\" by %s (%d) - %s\n",
			i+1, book.Title, book.Author, book.Year, status)
	}

	// Display summary information
	fmt.Printf("\nTotal books in collection: %d\n", len(allBooks))

	// Show breakdown of available vs borrowed books
	var availableCount, borrowedCount int
	for _, book := range allBooks {
		if book.IsBorrowed {
			borrowedCount++
		} else {
			availableCount++
		}
	}

	fmt.Printf("Available: %d | Borrowed: %d\n", availableCount, borrowedCount)
}

/*
performStartupTasks handles additional initialization operations after database setup.

This function performs various startup tasks that prepare the application
for user interaction. It can include data validation, system checks,
welcome data setup, and other initialization routines.

Startup Tasks:
1. Verify database connectivity and integrity
2. Display system information and statistics
3. Perform any necessary data cleanup or validation
4. Set up initial data if needed (sample books for first-time users)
5. Display welcome information to the user

This function is called after successful database initialization but before
the main menu system starts, ensuring the application is in a proper state
for user operations.
*/
func performStartupTasks() {
	fmt.Println("⚙️  Performing startup tasks...")

	// Task 1: Verify database connectivity by running a simple query
	var bookCount int64
	result := db.Model(&Book{}).Count(&bookCount)
	if result.Error != nil {
		log.Fatal("❌ Database connectivity check failed:", result.Error)
	}

	// Task 2: Display current library statistics
	var availableCount, borrowedCount int64
	db.Model(&Book{}).Where("is_borrowed = ?", false).Count(&availableCount)
	db.Model(&Book{}).Where("is_borrowed = ?", true).Count(&borrowedCount)

	fmt.Printf("📊 Current Library Status:\n")
	fmt.Printf("   📚 Total Books: %d\n", bookCount)
	fmt.Printf("   ✅ Available: %d\n", availableCount)
	fmt.Printf("   📤 Borrowed: %d\n", borrowedCount)

	// Task 3: Add sample data for new users (first-time setup)
	if bookCount == 0 {
		fmt.Println("🔧 First-time setup detected - adding sample books...")
		addSampleBooks()

		// Recount after adding sample data
		db.Model(&Book{}).Count(&bookCount)
		fmt.Printf("✅ Added sample books! New total: %d books\n", bookCount)
	}

	// Task 4: Perform basic data integrity checks
	fmt.Println("🔍 Performing data integrity checks...")

	// Check for any books with invalid years (basic validation)
	var invalidBooks int64
	db.Model(&Book{}).Where("year < ? OR year > ?", 1000, 2030).Count(&invalidBooks)
	if invalidBooks > 0 {
		fmt.Printf("⚠️  Warning: Found %d books with potentially invalid years\n", invalidBooks)
	} else {
		fmt.Println("✅ Data integrity check passed")
	}

	fmt.Println("✅ Startup tasks completed successfully!")
	fmt.Println()
}

/*
addSampleBooks adds a set of sample books to the database for new users.

This function is called during first-time setup when the database is empty.
It populates the library with a diverse collection of programming and
technical books to demonstrate the application's functionality and provide
users with immediate data to work with.

Sample Books Include:
- Programming language books (Go, Python, etc.)
- Software engineering principles
- Classic computer science texts
- Modern development practices

All sample books are added as "available" status by default.
*/
func addSampleBooks() {
	sampleBooks := []Book{
		{Title: "The Go Programming Language", Author: "Alan Donovan", Year: 2015, IsBorrowed: false},
		{Title: "Clean Code", Author: "Robert Martin", Year: 2008, IsBorrowed: false},
		{Title: "Design Patterns", Author: "Gang of Four", Year: 1994, IsBorrowed: true}, // One borrowed for demo
		{Title: "The Pragmatic Programmer", Author: "Andy Hunt", Year: 1999, IsBorrowed: false},
		{Title: "Effective Go", Author: "Google Team", Year: 2020, IsBorrowed: false},
		{Title: "You Don't Know JS", Author: "Kyle Simpson", Year: 2014, IsBorrowed: false},
		{Title: "Python Crash Course", Author: "Eric Matthes", Year: 2019, IsBorrowed: true}, // Another borrowed for demo
		{Title: "The Art of Computer Programming", Author: "Donald Knuth", Year: 1968, IsBorrowed: false},
	}

	// Insert each sample book into the database
	for _, book := range sampleBooks {
		result := db.Create(&book)
		if result.Error != nil {
			fmt.Printf("⚠️  Warning: Could not add sample book '%s': %v\n", book.Title, result.Error)
		}
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

Error Handling:
- Invalid menu selections are handled gracefully
- User is prompted to try again rather than terminating
*/
func menuSelector() {
	// Main application loop - continues until user chooses to exit
	for {
		// Display menu options to user
		displayMenu()
		option := readInt("Option => ")

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
			fmt.Println("✅ Thank you for using Mini Library Manager!")
			fmt.Println("💾 All your data has been safely saved to the database.")
			fmt.Println("👋 Goodbye!")
			return // Exit the function and terminate application
		default:
			// Handle invalid menu selections
			fmt.Println("❌ Invalid option. Please try again.")
		}

		// Add a pause after each operation for better user experience
		fmt.Print("\n📖 Press Enter to continue...")
		readInput("")
	}
}

/*
main is the application entry point and orchestrates the complete startup sequence.

This function serves as the primary entry point for the Mini Library Manager
application. It handles the complete application lifecycle from initialization
to shutdown, ensuring all components are properly set up before user interaction begins.

Complete Startup Sequence:
1. Display application welcome banner and startup information
2. Initialize database connection and perform schema migrations
3. Perform necessary startup tasks (data validation, sample data, etc.)
4. Start the interactive menu system for user operations
5. Handle graceful application termination with confirmation

Critical Operations:
- Database initialization (creates library.db file and tables if needed)
- Error handling for startup failures (will terminate if critical errors occur)
- User interface initialization with clear status reporting
- Startup task execution for optimal user experience

The function ensures that all necessary components are properly initialized
and validated before allowing user interaction, preventing runtime errors
related to uninitialized database connections or corrupted data.

Application Architecture:
This follows a proper application startup pattern with clear separation
of concerns: initialization → validation → user interface → shutdown.
*/
func main() {
	// Step 1: Display welcome banner and startup information
	fmt.Println("🚀 Starting Mini Library Manager with Database Integration...")
	fmt.Println("=" + strings.Repeat("=", 59)) // Create a nice banner line
	fmt.Println("📚 A CLI-based library management system with SQLite & GORM")
	fmt.Println("=" + strings.Repeat("=", 59))
	fmt.Println()

	// Step 2: Initialize database connection and perform migrations
	fmt.Println("🔧 Initializing database connection...")
	initDB()
	fmt.Println()

	// Step 3: Perform necessary startup tasks
	performStartupTasks()

	// Step 4: Start the interactive menu system for user operations
	fmt.Println("🎯 Library Manager is ready for use!")
	fmt.Println("📋 Use the menu below to manage your book collection.")

	menuSelector()

	// Step 5: Handle graceful application termination
	fmt.Println()
	fmt.Println("=" + strings.Repeat("=", 40))
	fmt.Println("📚 Mini Library Manager - Session Complete")
	fmt.Println("💾 Database: library.db (all data preserved)")
	fmt.Println("🔧 Status: All operations completed successfully")
	fmt.Println("=" + strings.Repeat("=", 40))

	fmt.Println("Thank you for using the application")
}