package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
	"github.com/gorilla/mux"
)

// Used this for reference.
// https://docs.microsoft.com/en-us/azure/cosmos-db/table/how-to-use-go?tabs=bash

type QuestionEntity struct {
	aztables.Entity
	Text        string `json:"text"`
	Answer1Id   string
	Answer1Text string
	Answer2Id   string
	Answer2Text string
	IsCurrent   bool
	CreatedDate time.Time
}

type VoteEntity struct {
	aztables.Entity
	VoteCount int32
}

func getServiceClient() *aztables.ServiceClient {

	connectionString, ok := os.LookupEnv("AZURE_CONNECTION_STRING")
	if !ok {
		panic("AZURE_CONNECTION_STRING environment variable not found.")
	}

	serviceClient, err := aztables.NewServiceClientFromConnectionString(connectionString, nil)
	if err != nil {
		panic(err)
	}

	return serviceClient
}

func getAllQuestionsFromTableStorage(serviceClient *aztables.ServiceClient) []QuestionEntity {
	Quesitons := []QuestionEntity{}

	client := serviceClient.NewClient("questions")
	listPager := client.NewListEntitiesPager(nil)
	for listPager.More() {
		response, err := listPager.NextPage(context.TODO())
		if err != nil {
			panic(err)
		}

		for _, entity := range response.Entities {
			var myEntity QuestionEntity
			err = json.Unmarshal(entity, &myEntity)
			if err != nil {
				panic(err)
			}

			Quesitons = append(Quesitons, myEntity)
		}
	}

	return Quesitons

}

func getQuestionFromTableStorage(serviceClient *aztables.ServiceClient, id string) QuestionEntity {
	client := serviceClient.NewClient("questions")
	options := &aztables.GetEntityOptions{}

	var question QuestionEntity
	entity, err := client.GetEntity(context.TODO(), "Questions", id, options)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(entity.Value, &question)
	if err != nil {
		panic(err)
	}

	return question

}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "votesy-api in golang!")
	fmt.Println("Endpoint hit: home")
}

func returnAllQuestions(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint hit: returnAllQuestions")
	serviceClient := getServiceClient()
	json.NewEncoder(w).Encode(getAllQuestionsFromTableStorage(serviceClient))
}

func returnQuestion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	serviceClient := getServiceClient()
	json.NewEncoder(w).Encode(getQuestionFromTableStorage(serviceClient, id))
}

func handleRequests() {
	// Lets use the Mux Router, since everyone else does.
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", home)
	router.HandleFunc("/questions", returnAllQuestions)
	router.HandleFunc("/questions/{id}", returnQuestion)
	log.Fatal(http.ListenAndServe(":10000", router))
}

func main() {
	fmt.Println("Hello Votesy!")

	fmt.Println("Authenicating...")
	serviceClient := getServiceClient()
	getAllQuestionsFromTableStorage(serviceClient)

	handleRequests()
}
