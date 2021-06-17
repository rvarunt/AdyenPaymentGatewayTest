package main

import (
	. "varunAdyenTest/src/api"

	"fmt"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func indexHandler(c *gin.Context) {
	fmt.Print("Loading Summary Page")
	renderPagewithData(c, gin.H{
		"page":      "summary",
	})
}

func paymentpageHandler(c *gin.Context) {
	fmt.Print("Loading Payment Page")
	renderPagewithData(c, gin.H{
		"page":      "payment",
		"type": c.Param("type"),
		"clientKey": os.Getenv("CLIENT_KEY"),
	})
}

func endpageHandler(c *gin.Context) {
	log.Println("Loading preview page")
	renderPagewithData(c, gin.H{
		"page": "end",
		"status": c.Param("status"),
	})
}

func renderPagewithData(c *gin.Context,d gin.H) {
	// Call the HTML method of the Context to render a template
	data := d
	c.HTML(
		// Set the HTTP status to 200 (OK)
		http.StatusOK,
		// Use the index.html template
		"index.html",
		// Pass the data that the page uses
		data,
	)
}


func main() {

	err := godotenv.Load()
	if err != nil {
			log.Println("Error loading .env file")
	}

	fmt.Println(os.Getenv("API_KEY"))

	// Set the router as the default one shipped with Gin
	router := gin.Default()
	// Serve HTML templates
	router.LoadHTMLGlob("./templates/*")
	// pre-fetch static files
	router.Use(static.Serve("/static", static.LocalFile("./static", true)))

	//http.HandleFunc("/",indexHandler)
	//http.HandleFunc( "/payment",paymentHandler)

	//log.Fatal(http.ListenAndServe(":8080", nil))

	// Setup route group for the API
	api := router.Group("/api")
	{
		api.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H {
				"message": "pong",
			})
		})
	}

	api.POST("/getPaymentMethods", PaymentMethodHandler)

	api.POST("/initiatePayment", PaymentsHandler)

	api.GET("/handleUIRedirect", RedirectUIHandler)

	api.POST("/handleUIRedirect", RedirectUIHandler)

	//api.POST("/handlePaymentAddtionaldetails", PaymentAdditionalDetailsHandler)

	router.GET("/",indexHandler)

	router.POST("/payment/:type",paymentpageHandler)

	router.GET("/end/:status", endpageHandler)

	router.Run(":8080")
}
