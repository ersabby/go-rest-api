/*
	Title:		"RESTful JSON API with GOlang"
	Description:	This code contains certain API's allowing users to create 
			and view events.
*/

package main
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"github.com/gorilla/mux"
)

//Creating a database structure
type event struct{
	ID		string `json:"ID"`
	Title		string `json:"Title"`
	Description	string `json:"Description"`
}

type allEvents []event

//Creating an initial pre-defined entry
var events = allEvents{
	{
		ID: "1",
		Title: "Sabby",
		Description: "Learning Golang",
	},
   }

func homeLink(w http.ResponseWriter, r *http.Request){
fmt.Fprintf(w, "welcome home..!");
}


/*
create event function is used to handle data coming from user's end
The data incoming through http.Request is not in human readable form,
so the data is sliced using 'ioutil' and then unmarshal is used to store sliced
data into defined struct.
(Unmarshal is used to decode JSON array elements into corresponding Go array
elements. Also, it stores the JSON objects into the map)
*/

func createEvent(w http.ResponseWriter, r *http.Request){
	var newEvent event

	//Convert r.Body into a readable format
	reqBody, err:= ioutil.ReadAll(r.Body)
	if err!=nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description")
	}

	json.Unmarshal(reqBody, &newEvent)

	//Adding newly created event to array of events
	events = append(events, newEvent)

	//Return 201 created status code
	w.WriteHeader(http.StatusCreated)

	//Return newly created url
	json.NewEncoder(w).Encode(newEvent)
}

func getOneEvent(w http.ResponseWriter, r *http.Request){
	//get id from url received
	eventID := mux.Vars(r)["id"]

	for _,singleEvent := range events {
		if singleEvent.ID == eventID {
		json.NewEncoder(w).Encode(singleEvent)
		}
   }
}

func getAllEvents(w http.ResponseWriter, r *http.Request){
	json.NewEncoder(w).Encode(events)
}

func updateEvent(w http.ResponseWriter, r* http.Request){

	//get ID by url
	eventID := mux.Vars(r)["id"]
	var updatedEvent event
	//Converting r.Body into readable format
	reqBody, err:=ioutil.ReadAll(r.Body)
	if err!=nil {
		fmt.Fprintf(w, "Kindly enter the data with the event title and description")
	}
	json.Unmarshal(reqBody, &updatedEvent)

	for i, singleEvent := range events{
		if singleEvent.ID == eventID {
			singleEvent.Title = updatedEvent.Title
			singleEvent.Description = updatedEvent.Description
			events = append(events[:i], singleEvent)
			json.NewEncoder(w).Encode(singleEvent)
		}
	}
}

//Deleting an event
func deleteEvent(w http.ResponseWriter, r *http.Request){
	//Get id from url
	eventID := mux.Vars(r)["id"]

	//Get Details from an existing event
	for i, singleEvent := range events{
		if singleEvent.ID == eventID{
			events = append(events[:i], events[i+1:]...)
			fmt.Fprintf(w, "The event with ID %v has been successfully deleted", eventID)
		}
	}
}




func main(){
	router:=mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/event", createEvent).Methods("POST")
	router.HandleFunc("/events", getAllEvents).Methods("GET")
	router.HandleFunc("/events/{id}", getOneEvent).Methods("GET")
	router.HandleFunc("/events/{id}", updateEvent).Methods("PATCH")
	router.HandleFunc("/events/{id}", deleteEvent).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}
