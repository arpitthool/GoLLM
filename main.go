package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func getHello(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, "hello!")
}

func getHomePage(context *gin.Context) {
	context.HTML(http.StatusOK, "index.html", nil)
}

func postHomePageForm(context *gin.Context) {
	userInput := context.PostForm("userInput")
    context.HTML(http.StatusOK, "result.html", gin.H{"userInput": userInput})
}

func main() {
	router := gin.Default()
	
	router.GET("/hello", getHello)

	// Define a route for the home page
	router.GET("/", getHomePage)

    // Define a route to handle the form submission
    router.POST("/submit", postHomePageForm)

    // Serve static files (e.g., CSS, JavaScript) from the "static" directory
    router.Static("/static", "./static")

    // Load HTML templates
    router.LoadHTMLFiles("templates/index.html", "templates/result.html")

	router.Run("localhost:8080")
}
