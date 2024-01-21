package main

import (
	"errors"
	"net/http"
	"strings"
	"github.com/gin-gonic/gin"
	"strconv"
)

type book struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

var books = []book{
	{ID: 1, Title: "In Search of Lost Time", Author: "Marcel Proust", Quantity: 2},
	{ID: 2, Title: "Fundamentos da Matemática Elementar - 3", Author: "Iezzy", Quantity: 5},
	{ID: 3, Title: "War and Peace", Author: "Leo Tolstoy", Quantity: 6},
}

func getBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books)
}

func createBook(c *gin.Context) {
	var newBook book
	if err := c.BindJSON(&newBook); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := isValidNewBook(&newBook); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := isBookAlreadyExists(newBook.Title); err {
		c.IndentedJSON(http.StatusConflict, gin.H{"error": "Book already exist"})
		return 
	}
	
	// auto increment ID
	newBook.ID = int(books[len(books)-1].ID) + 1
	books = append(books, newBook)

	c.IndentedJSON(http.StatusCreated, newBook)
}

func isValidNewBook(newBook *book) error {

	if newBook.ID != 0 {
		return errors.New("the ID field is autoincremented")
	}
	if strings.TrimSpace(newBook.Title) == "" {
		return errors.New("missing title")
	}
	if strings.TrimSpace(newBook.Author) == "" {
		return errors.New("missing author name")
	}
	if newBook.Quantity == 0 {
		return errors.New("missing quantity")
	}
	return nil
}

func isBookAlreadyExists(title string) bool {
	for _, existingBook := range books {
			if strings.EqualFold(existingBook.Title, title) {
					return true 
			}
	}
	return false 
}
func bookById(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid Id"})
		return
	}
	book, err, _ := getBookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound,gin.H{"error": err.Error()} )
		return
	}

	c.IndentedJSON(http.StatusFound, book)
}
func getBookById(ID int) (*book, error, int) {
	for i, bo := range books {
			if bo.ID == ID {
					return &books[i], nil, i
			}
	}

	return nil, errors.New("book not found"), -1
}


func updateBookById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid Id"})
		return
	}

	var newBook book
	if err := c.BindJSON(&newBook); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if newBook.ID != 0 {
		e := errors.New("the ID field is autoincremented")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": e.Error()})
		return
	}
	book, err, _ := getBookById(id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if strings.TrimSpace(newBook.Title) != "" {
		book.Title = newBook.Title
	}else	if strings.TrimSpace(newBook.Author) != "" {
		book.Author = newBook.Title
	}else if newBook.Quantity >= 0 {
	book.Quantity = newBook.Quantity
	}

	c.IndentedJSON(http.StatusOK, book)
}

func removeBook(index int) {
	 books = append(books[:index], books[index+1:]...)
}
func deleteBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid Id"})
		return
	}

	book, err, index := getBookById(id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	deletedBook := *book
	removeBook(index)

	c.IndentedJSON(http.StatusOK, gin.H{"deleted": true, "book": deletedBook})
}

func main() {
	router := gin.Default()
	defer router.Run("localhost:8080")
	router.GET("/books", getBooks)
	router.GET("/books/:id", bookById)
	router.POST("/books", createBook)
	router.PUT("/books/:id", updateBookById)
	router.DELETE("/books/:id", deleteBook)
}
