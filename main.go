package main

import (
	"errors"
	"net/http" // standard

	"github.com/gin-gonic/gin"
)

/*
	use gin.Context.JSON() rather than gin.Context.IndentedJSON()
	JSON() sets the content-type of our resonse to "application/JSON"
	IndentedJSON() is generally used for dev purposes - it uses up more CPU to format the JSON
*/


// start field names with caps - this makes them get exported to JSON later
// make it serialisable-to-json with the backticks
type book struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

type author struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// create a slice of books
var test_books = []book{
	{ID: "1", Title: "Man's search for meaning", Author: "Viktor Frankl", Quantity: 1},
	{ID: "2", Title: "Another book", Author: "Viktor Frankl", Quantity: 3},
	{ID: "3", Title: "Archetypes", Author: "Jung", Quantity: 5},
}

var authors = []author{
	{ID: "1", Name: "Viktor Frankl"},
}

// GETTING
func getBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, test_books) // books slice will be serialised automatically
}

func getAuthors(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, authors)
}

func getBooksByAuthor(c *gin.Context) {
	var result []book
	name := c.Param("name")
	for _, b := range test_books {
		if b.Author == name {
			result = append(result, b)
		}
	}

	if len(result) == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "No books for author " + name})
		return
	}

	c.IndentedJSON(http.StatusOK, result)
}

// the function for our router
func bookById(c *gin.Context) {
	id := c.Param("id")
	book, err := getBookById(id)

	if err != nil {
		// gin.H lets us write our own JSON back to the page.  shorthand.
		// we needed to add this so there would be a status code.  return alone doesn't do that
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found!"})
		return
	}
	c.IndentedJSON(http.StatusOK, book)
}

// return a pointer to a book, and an error
func getBookById(id string) (*book, error) {
	// iterate over all books
	for i, b := range test_books {
		if b.ID == id {
			return &test_books[i], nil
		}
	}
	return nil, errors.New("book not found")
}

// POSTING
func createBook(c *gin.Context) {
	var newBook book

	//try to bind our new book to json
	// BindJSON will return a BAD REQUEST for us if it fails
	// (rather than the return statement)
	if err := c.BindJSON(&newBook); err != nil {
		return
	}

	test_books = append(test_books, newBook)
	c.IndentedJSON(http.StatusCreated, newBook)
}

// PATCHING
func checkoutBook(c *gin.Context) {
	id, ok := c.GetQuery("id") // get from querystring
	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing ID parameter"})
		return
	}

	book, err := getBookById(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}

	if book.Quantity <= 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Book not available"})
		return
	}

	book.Quantity -= 1
	c.IndentedJSON(http.StatusOK, book)
}

func checkinBook(c *gin.Context) {
	id, ok := c.GetQuery("id")
	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing ID Parameter"})
		return
	}

	book, err := getBookById(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "we didn't find the book"})
		return
	}

	book.Quantity += 1
	c.IndentedJSON(http.StatusOK, book)
}

func main() {
	// create our router which will handle our requests and point to functions
	router := gin.Default()
	router.GET("/books", getBooks)
	router.GET("/books/:id", bookById)
	router.GET("/authors", getAuthors)
	router.GET("/search/:name", getBooksByAuthor)

	router.POST("/books", createBook)

	router.PATCH("/checkout", checkoutBook)
	router.PATCH("/checkin", checkinBook)

	router.Run("127.0.0.1:8080")
}
