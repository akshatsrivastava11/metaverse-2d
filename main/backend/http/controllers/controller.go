package controllers

import (
	"backendMetaverse/http/client"
	"backendMetaverse/http/utils"
	"backendMetaverse/prisma/db"
	"strconv"
	"strings"

	//  "backendMetaverse/utils"
	"encoding/json"
	"fmt"

	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type SignupRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UpdateMetaDataRequest struct {
	Token    string `json:"token"`
	AvatarId string `json:"avatarId"`
}

type SpaceCreationRequst struct {
	Name      string `json:"name"`
	Dimension string `json:"dimensions"`
	MapId     string `json:"mapId"`
}

type AddingElementRequest struct {
	ElementId string `json:"elementId"`
	SpaceId   string `json:"spaceId"`
	X         string `json:"x"`
	Y         string `json:"y"`
}

// var client *db.PrismaClient=client.Client()
// working
func SignupControllers(w http.ResponseWriter, r *http.Request) {

	prismaClient := client.GetClient()
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	var user SignupRequest
	json.NewDecoder(r.Body).Decode(&user)
	// client := db.NewClient()
	data, err := prismaClient.User.CreateOne(db.User.Username.Set(user.Username), db.User.Password.Set(utils.HashedPassword(user.Password)), db.User.Role.Set("Admin")).Exec(r.Context())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Data is ", data.ID)
	json.NewEncoder(w).Encode(struct {
		UserID string `json:"userId"`
	}{
		UserID: data.ID,
	})
}

// working
func SignInControllers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	prismaClient := client.GetClient()
	var user SignupRequest
	json.NewDecoder(r.Body).Decode(&user)
	fmt.Println(user.Username)
	data, err := prismaClient.User.FindFirst(db.User.Username.Equals(user.Username)).Exec(r.Context())
	if err != nil {
		log.Fatal("Error in the signin Controller", err)
	}
	// json.NewEncoder(w).Encode(data)
	if utils.CheckHashedPassword(user.Password, data.Password) {
		token := utils.GetToken(data.Username)
		w.Header().Set("Authorization", "Bearer "+token) // <-- Set header first!
		json.NewEncoder(w).Encode(struct {
			Token string `json:"token"`
		}{
			Token: token,
		})
	} else {
		json.NewEncoder(w).Encode("Nahh bro ,,u cheated")
	}

}

// working
func ShowAllUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	prismaClient := client.GetClient()
	// client := db.NewClient()
	resp, _ := prismaClient.User.FindMany().Exec(r.Context())
	json.NewEncoder(w).Encode(resp)
}

// working
func UpdateMetaData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")

	var data UpdateMetaDataRequest
	json.NewDecoder(r.Body).Decode(&data)
	// fmt.Println(decode)
	token := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	fmt.Println("header is", r.Header.Get("Authorization"))
	fmt.Println("token is ", token)
	// json.NewDecoder(r.Body).Decode(&data)
	username, err := utils.VerifyToken(token[1])
	// fmt.Println("username is ", username)
	if err != nil {
		log.Fatal("Error in authentication")
	}
	prismaClient := client.GetClient()
	updated, err := prismaClient.User.FindMany(db.User.Username.Equals(username)).Update(db.User.AvatarID.Set(data.AvatarId)).Exec(r.Context())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("updated is ", updated)
	if err := json.NewEncoder(w).Encode(struct {
		Msg string `json:"message"`
	}{
		Msg: "Updated",
	}); err != nil {
		log.Println("JSON encoding error:", err)
	}

}

// working
func GetAvailableAvatars(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	prismaClient := client.GetClient()
	rsp, _ := prismaClient.Avatar.FindMany().Exec(r.Context())
	json.NewEncoder(w).Encode(rsp)

}

// kaam chl rha janab...
func CreateSpace(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	prismaClient := client.GetClient()

	var data SpaceCreationRequst

	token := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	fmt.Println("header is", r.Header.Get("Authorization"))
	fmt.Println("token is ", token)
	// json.NewDecoder(r.Body).Decode(&data)
	username, err := utils.VerifyToken(token[1])
	if err != nil {
		log.Fatal("error in validating token", err)
	}
	resp, err := prismaClient.User.FindFirst(db.User.Username.Equals(username)).Exec(r.Context())
	if err != nil {
		log.Fatal("error for user finding", err)
	}
	json.NewDecoder(r.Body).Decode(&data)
	fmt.Println("the data from the request is ", data)
	fmt.Println("The dimensioons are ", data.Dimension)
	widthStr := strings.Split(data.Dimension, "x")[0]
	width, _ := strconv.Atoi(widthStr)
	heightStr := strings.Split(data.Dimension, "x")[1]
	height, _ := strconv.Atoi(heightStr)
	// cretorID, _ := strconv.Atoi(resp.ID)
	// db.Space.CreatorID.Set(strconv.Itoa(cretorID))
	if data.MapId == "" {
		fmt.Println("Map id is given")
		createdSpace, _ := prismaClient.Space.CreateOne(
			db.Space.Name.Set(data.Name),
			db.Space.Width.Set(width),
			db.Space.Height.Set(height),
			db.Space.Creator.Link(
				db.User.ID.Equals(resp.ID),
			),
		).Exec(r.Context())
		fmt.Println(createdSpace)
		return
	}
	fmt.Println("mapid", data.MapId)
	mapResp, _ := prismaClient.Map.FindFirst(db.Map.ID.Equals(data.MapId)).Exec(r.Context())
	fmt.Println("Mapresp is", mapResp)
	if mapResp != nil {
		// meaning the user selects a pre existing map
		createdSpace, _ := prismaClient.Space.CreateOne(
			db.Space.Name.Equals(data.Name),
			db.Space.Width.Equals(mapResp.Width),
			db.Space.Height.Equals(mapResp.Height),
			db.Space.Creator.Link(
				db.User.ID.Equals(resp.ID),
			),
		).Exec(r.Context())
		fmt.Println("Created spave is ", createdSpace)
		if mapResp != nil {
			fmt.Println("In the map response function")
		}
		// create all elements in spaccem
		mapElems := mapResp.MapElements
		fmt.Println("map elems is ", mapElems)
		// for i := range mapElems() {
		// 	prismaClient.SpaceElements.CreateOne(
		// 		db.SpaceElements.ElementID.Equals(i.mapelemetns)
		// 	)
		// }

	}

}

// working
func DeleteSpace(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	prismaClient := client.GetClient()
	deleted_elems, _ := prismaClient.Space.FindUnique(db.Space.ID.Equals(id)).Delete().Exec(r.Context())
	fmt.Println("elem deleted successfully", deleted_elems)
	json.NewEncoder(w).Encode("Deleted successfully seerr")
}

// working
func GetAllExistingSpaces(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	token := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	fmt.Println("header is", r.Header.Get("Authorization"))
	fmt.Println("token is ", token)
	username, err := utils.VerifyToken(token[1])
	if err != nil {
		fmt.Println("Error in verifying token ")
		log.Fatal(err)
	}

	prismaClient := client.GetClient()
	data, err := prismaClient.User.FindFirst(db.User.Username.Equals(username)).Exec(r.Context())

	resp, err := prismaClient.Space.FindMany(db.Space.CreatorID.Equals(data.ID)).Exec(r.Context())
	fmt.Println(resp)
	json.NewEncoder(w).Encode(resp)

}

// kaam chl rha inpe bh
func GetaSpace(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var vars = mux.Vars(r)
	spaceId := vars["spaceId"]
	fmt.Println(spaceId)
	prismaClient := client.GetClient()
	resp, err := prismaClient.Space.FindFirst(db.Space.ID.Equals(spaceId)).Exec(r.Context())
	if err != nil {
		fmt.Println(err)
	}
	json.NewEncoder(w).Encode(resp)

}

// working
func AddElementinSpace(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var resp AddingElementRequest

	json.NewDecoder(r.Body).Decode(&resp)
	fmt.Println(resp)

	fmt.Println(resp.X)

	prismaClient := client.GetClient()

	data, err := prismaClient.Space.FindFirst(db.Space.ID.Equals(resp.SpaceId)).Exec(r.Context())

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(data)

	givenX, _ := strconv.Atoi(resp.X)
	givenY, _ := strconv.Atoi(resp.Y)

	if givenX < 0 || givenY < 0 || givenX < data.Width || givenY < data.Height {
		json.NewEncoder(w).Encode("Point is outside the boundary")
		return
	}
	response, _ := prismaClient.SpaceElements.CreateOne(db.SpaceElements.X.Set(givenX), db.SpaceElements.Y.Set(givenY), db.SpaceElements.Space.Link(
		db.Space.ID.Equals(resp.SpaceId),
	), db.SpaceElements.Element.Link(
		db.Element.ID.Equals(resp.ElementId),
	)).Exec(r.Context())
	// json.Decoder(w).DE(response)
	json.NewEncoder(w).Encode(response)
}

type DeleteRequestBody struct {
	Id string `json:"id"`
}

// not done
func DeleteElement(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var resp DeleteRequestBody
	json.NewDecoder(r.Body).Decode(&resp)
	fmt.Println(resp)
	prismaClient := client.GetClient()
	deletedElem, err := prismaClient.Element.FindUnique(db.Element.ID.Equals(resp.Id)).Delete().Exec(r.Context())
	if err != nil {
		fmt.Println(err)
	}
	json.NewEncoder(w).Encode(deletedElem)
}

// not done
func SeeAllElems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// token := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	// fmt.Println("header is", r.Header.Get("Authorization"))
	// fmt.Println("token is ", token)
	// // json.NewDecoder(r.Body).Decode(&data)
	// username, err := utils.VerifyToken(token[1])
	// if err != nil {
	// 	log.Fatal(err)
	// }
	prismaClient := client.GetClient()
	// user, _ := prismaClient.User.FindUnique(db.User.Username.Equals(username)).Exec(r.Context())
	data, _ := prismaClient.Element.FindMany().Exec(r.Context())
	json.NewEncoder(w).Encode(data)
}

type CreateElementRequestBody struct {
	ImageURL string `json:"imageUrl"`
	Width    string `json:"width"`
	Height   string `json:"height"`
	Static   string `json:"static"`
}

// working
func CreateElement(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var resp CreateElementRequestBody
	json.NewDecoder(r.Body).Decode(&resp)
	fmt.Println(resp)
	fmt.Println(r.Body)
	prismaClient := client.GetClient()
	givenW, _ := strconv.Atoi(resp.Width)
	givenH, _ := strconv.Atoi(resp.Height)
	static, _ := strconv.ParseBool(resp.Static)
	createdElem, err := prismaClient.Element.CreateOne(db.Element.Width.Set(givenW), db.Element.Height.Set(givenH), db.Element.Static.Set(static), db.Element.ImageURL.Set(resp.ImageURL)).Exec(r.Context())
	if err != nil {
		fmt.Println(err)
	}
	json.NewEncoder(w).Encode(createdElem)

}

type updateElemBody struct {
	ImageURL string `json:"imageUrl"`
}

// working
func UpdateElem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	elemId := vars["elementId"]
	primsaClient := client.GetClient()
	var resp updateElemBody
	json.NewDecoder(r.Body).Decode(&resp)
	updatedElem, _ := primsaClient.Element.FindUnique(db.Element.ID.Equals(elemId)).Update(db.Element.ImageURL.Set(resp.ImageURL)).Exec(r.Context())

	json.NewEncoder(w).Encode(updatedElem)
}

type CreateAvaterResponse struct {
	ImageURl string `json:"imageUrl"`
	Name     string `json:"name"`
}

// working
func CreateAvatar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	prismaClient := client.GetClient()
	var resp CreateAvaterResponse
	json.NewDecoder(r.Body).Decode(&resp)
	token := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	fmt.Println("header is", r.Header.Get("Authorization"))
	fmt.Println("token is ", token)
	// json.NewDecoder(r.Body).Decode(&data)
	username, err := utils.VerifyToken(token[1])
	if err != nil {
		log.Fatal("error in validating token", err)
	}
	createdAvatar, _ := prismaClient.Avatar.CreateOne(db.Avatar.ImageURL.Set(resp.ImageURl), db.Avatar.Name.Set(resp.Name), db.Avatar.Users.Link(db.User.Username.Equals(username))).Exec(r.Context())
	json.NewEncoder(w).Encode(createdAvatar)
}

type ArrayTypeCreateMap struct {
	ElementId string `json:"elementId"`
	X         string `json:"x"`
	Y         string `json:"y"`
}
type CreateMapResponse struct {
	Thumbnail       string               `json:"thumbnail"`
	Dimensions      string               `json:"dimensions"`
	Name            string               `json:"name"`
	DefaultElements []ArrayTypeCreateMap `json:"defaultElements"`
}

func CreateMap(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	prismaClient := client.GetClient()
	var resp CreateMapResponse
	json.NewDecoder(r.Body).Decode(&resp)
	givenX, _ := strconv.Atoi(strings.Split(resp.Dimensions, "x")[0])
	givenY, _ := strconv.Atoi(strings.Split(resp.Dimensions, "x")[1])

	cretedMapWithoutElemsAdded, _ := prismaClient.Map.CreateOne(db.Map.Width.Set(givenX), db.Map.Height.Set(givenY), db.Map.Name.Set(resp.Name), db.Map.Thumbnail.Set(resp.Thumbnail)).Exec(r.Context())
	// var elementLinks []db.MapElementsCursorParam
	for _, val := range resp.DefaultElements {

		givenXElem, _ := strconv.Atoi(val.X)
		givenYElem, _ := strconv.Atoi(val.Y)

		prismaClient.MapElements.CreateOne(
			db.MapElements.X.Equals(givenXElem),
			db.MapElements.Y.Set(givenYElem),
			db.MapElements.Map.Link(db.Map.ID.Equals(cretedMapWithoutElemsAdded.ID)),
			db.MapElements.Element.Link(db.Element.ID.Equals(val.ElementId)),
		)

	}
	json.NewEncoder(w).Encode(cretedMapWithoutElemsAdded.ID)

}
