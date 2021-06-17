package main


import (
"log"
"net/http"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Hello World!</h1>"))
}

func main() {
	/*
		err := godotenv.Load()
		if err != nil {
			log.Println("Error loading .env file")
		}
	*/
	http.HandleFunc("/",indexHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
