package router

import (
	"backendMetaverse/controllers"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/signup", controllers.SignupControllers).Methods("POST")
	r.HandleFunc("/allusers", controllers.ShowAllUser).Methods("GET")
	r.HandleFunc("/signin", controllers.SignInControllers).Methods("POST")
	r.HandleFunc("/update", controllers.UpdateMetaData).Methods("PUT")
	r.HandleFunc("/createSpace", controllers.CreateSpace).Methods("POST")
	r.HandleFunc("/space/{id}", controllers.DeleteSpace).Methods("DELETE")
	r.HandleFunc("/space/{spaceId}", controllers.GetaSpace).Methods("GET")

	r.HandleFunc("/addelem", controllers.AddElementinSpace).Methods("POST")
	r.HandleFunc("/avatar", controllers.CreateAvatar).Methods("POST")
	r.HandleFunc("/createMap", controllers.CreateMap).Methods("POST")
	r.HandleFunc("/allavatars", controllers.GetAvailableAvatars).Methods("POST")
	r.HandleFunc("/createElem", controllers.CreateElement).Methods("POST")
	// r.HandleFunc("/deleteSpace/:id", controllers.DeleteElement).Methods("DELETE")
	r.HandleFunc("/getAllExistingSpace", controllers.GetAllExistingSpaces).Methods("GET")
	r.HandleFunc("/element/{elementId}", controllers.UpdateElem).Methods("PUT")
	r.HandleFunc("/element/{elementId}", controllers.DeleteElement).Methods("DELETE")

	return r
}
