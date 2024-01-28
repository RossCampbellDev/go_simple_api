package main

import (
	"net/http" // standard

	"github.com/gin-gonic/gin"
)

// start field names with caps - this makes them get exproted to JSON later
// make it serialisable-to-json
type book struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

// create a slice of books
var test_books = []book{
	{ID: "1", Title: "Man's search for meaning", Author: "Viktor Frankl", Quantity: 1},
	{ID: "2", Title: "Archetypes", Author: "Jung", Quantity: 5},
}

func getBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, test_books) // books slice will be serialised
}

func main() {
	// create our router which will handle our routes
	router := gin.Default()
	router.GET("/books", getBooks)
	router.Run("127.0.0.1:8080")
}
