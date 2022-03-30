package api

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/RoadTripppin/wazzup/controllers"
	"github.com/RoadTripppin/wazzup/helpers"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func StartApi() {
	helpers.LoadEnv()

	router := mux.NewRouter()

	wsServer := controllers.NewWebsocketServer()
	go wsServer.Run()

	// CORS Handler
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{os.Getenv("ORIGIN_ALLOWED")})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"})

	//routes
	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		controllers.ServeWs(wsServer, w, r)
	})
	router.HandleFunc("/login", controllers.Login).Methods("POST")
	router.HandleFunc("/user", controllers.GetUser).Methods("GET")
	router.HandleFunc("/register", controllers.Register).Methods("POST")
	router.HandleFunc("/user/update", controllers.UpdateUser).Methods("POST")
	router.HandleFunc("/user/delete", controllers.DeleteUser).Methods("POST")

	port := os.Getenv("SERVER_PORT")
	fmt.Println("App is working on port :" + port)
	log.Fatal(http.ListenAndServe(":"+port, handlers.CORS(headersOk, methodsOk, originsOk)(router)))
}
