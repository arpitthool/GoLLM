package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"bytes"
	// "encoding/json"
	"regexp"
    // "fmt"
)

// global variable
var history string = ""

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

func parseJSONResponseFromLLM(JSONResponse *bytes.Buffer) string {
	r, _ := regexp.Compile("response\":\"([a-zA-Z .!,'?]*)")
	// stringArray := r.FindAllStringSubmatch(JSONResponse.String(), -1)
	stringArray := r.FindAllString(JSONResponse.String(), -1)
	finalString := ""
	for i := range stringArray {
		finalString += stringArray[i][11:]
	}
	return finalString 
	// return finalString + "      " + JSONResponse.String()
	// return JSONResponse.String()
}

// // decalre a struct
// type LLMResponse struct {
// 	// defining struct variables
// 	model 		string
// 	created_at 	string
// 	response	string
// 	done		bool
// }

// // function to parse JSON response from LLM model 
// func parseJSONResponseFromLLM(JSONResponse *bytes.Buffer) string {
// 	// defining a struct instance
// 	var LLMResponseArray []LLMResponse

// 	// decoding JSON array
// 	err := json.Unmarshal(JSONResponse, &LLMResponseArray)

// 	if err != nil {
// 		return "ERROR!!!!!!"
// 	}

// 	finalString := ""
// 	for i := range LLMResponseArray {
// 		finalString += LLMResponseArray[i].response
// 	}

// 	// return JSONResponse.String()
// 	return finalString
// }

func postToLLMModel(context *gin.Context) {
	modelRestriction := "[Give your answer in one word]"
	userInput := context.PostForm("userInput")

	// url
	url := "http://localhost:11434/api/generate"

	x := `{
		"model": "llama2-uncensored:cpu",
		"prompt": "`+userInput+modelRestriction+`"
	}`

	// JSON payload for LLM Model
	payload := []byte(x)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
        if err != nil {
            context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        // Set the Content-Type header to specify that you're sending JSON data
        req.Header.Set("Content-Type", "application/json")

        // Create an HTTP client
        client := &http.Client{}

        // Send the POST request
        resp, err := client.Do(req)
        if err != nil {
            context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        defer resp.Body.Close()

        // Read and print the response body
        responseBody := new(bytes.Buffer)
        _, err = responseBody.ReadFrom(resp.Body)
        if err != nil {
            context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

		// parse the json
		JSONParsedResponseString := parseJSONResponseFromLLM(responseBody)

		// update history
		history += userInput + " -> " + JSONParsedResponseString + "\n\n"

        // Respond with the response status and body as JSON
        // context.JSON(http.StatusOK, gin.H{
        //     "status": resp.Status,
        //     "body":   responseBody.String(),
        // })

    	// context.HTML(http.StatusOK, "index.html", gin.H{"userInput": userInput, "history": history})
    	context.HTML(http.StatusOK, "result.html", gin.H{"userInput": userInput, "history": history})
}

func main() {
	router := gin.Default()

	router.GET("/hello", getHello)

	// Define a route for the home page
	router.GET("/", getHomePage)

    // Define a route to handle the form submission
    router.POST("/submit", postToLLMModel)
    router.POST("/", postToLLMModel)

    // Serve static files (e.g., CSS, JavaScript) from the "static" directory
    router.Static("/static", "./static")

    // Load HTML templates
    router.LoadHTMLFiles("templates/index.html", "templates/result.html")

	router.Run("localhost:8080")
}
