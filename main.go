package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"bytes"
	"regexp"
)

var historyArray []string
var llmModel = "llama2-uncensored:cpu"
var ollamaAPIURL = "http://localhost:11434/api/generate"

func getHomePage(context *gin.Context) {
	context.HTML(http.StatusOK, "index.html", gin.H{"historyArray": historyArray})

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
}

func postToLLMModel(context *gin.Context) {
	modelRestriction := "[Give your answer in less than 250 character]"
	userInput := context.PostForm("userInput")

	// url
	url := ollamaAPIURL

	x := `{
		"model": `+llmModel+`,
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

		// historyArray = append(historyArray, userInput + " -> " + JSONParsedResponseString)
		historyArray = append(historyArray, userInput )
		historyArray = append(historyArray, JSONParsedResponseString)

    	context.HTML(http.StatusOK, "index.html", gin.H{"userInput": userInput, "historyArray": historyArray})
}

func main() {
	router := gin.Default()

	// Define a route for the home page
	router.GET("/GoLLM", getHomePage)

    // Define a route to handle the form submission
    router.POST("/GoLLM", postToLLMModel)

    // Serve static files (e.g., CSS, JavaScript) from the "static" directory
    router.Static("/static", "./static")

    // Load HTML templates
    router.LoadHTMLFiles("templates/index.html")

	router.Run("localhost:8080")
}
