package api

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/RoadTripppin/wazzup/config"
	"github.com/RoadTripppin/wazzup/controllers"
	"github.com/RoadTripppin/wazzup/helpers"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

//var Contex = context.Background()

func StartApi() {
	helpers.LoadEnv()

	config.CreateRedisClient()

	db := config.InitDB()
	defer db.Close()

	router := mux.NewRouter()

	wsServer := controllers.NewWebsocketServer(&helpers.RoomRepository{Db: db}, &helpers.UserRepository{Db: db})
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
	router.HandleFunc("/user", controllers.SearchUser).Methods("GET")
	router.HandleFunc("/register", controllers.Register).Methods("POST")
	router.HandleFunc("/user/update", controllers.UpdateUser).Methods("POST")
	router.HandleFunc("/user/delete", controllers.DeleteUser).Methods("POST")
	router.HandleFunc("/user/interacted", controllers.GetInteractedUsers).Methods("GET")

	port := os.Getenv("SERVER_PORT")
	fmt.Println("App is working on port :" + port)
	log.Fatal(http.ListenAndServe(":"+port, handlers.CORS(headersOk, methodsOk, originsOk)(router)))

}
