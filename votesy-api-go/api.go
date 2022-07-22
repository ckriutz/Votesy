package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
	"github.com/aidarkhanov/nanoid"
	"github.com/gorilla/mux"
)

// Used this for reference.
// https://docs.microsoft.com/en-us/azure/cosmos-db/table/how-to-use-go?tabs=bash

type QuestionEntity struct {
	aztables.Entity
	Text        string `json:"text"`
	Answer1Id   string `json:"answer1Id"`
	Answer1Text string `json:"answer1Text"`
	Answer2Id   string `json:"answer2Id"`
	Answer2Text string `json:"answer2Text"`
	IsCurrent   bool   `json:"isCurrent"`
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

func addQuestionToTableStorage(serviceClient *aztables.ServiceClient, question QuestionEntity) QuestionEntity {
	// Have to create both the question AND the votes.
	questionClient := serviceClient.NewClient("questions")
	voteClient := serviceClient.NewClient("votes")

	// So, first thing we need to do is update the id's of the Question.
	alphabet := nanoid.DefaultAlphabet

	id, err := nanoid.Generate(alphabet, 5)
	if err != nil {
		panic(err)
	}

	answerid1, err := nanoid.Generate(alphabet, 5) //> "i25_4"
	if err != nil {
		panic(err)
	}

	answerid2, err := nanoid.Generate(alphabet, 5) //> "i25_4"
	if err != nil {
		panic(err)
	}

	question.PartitionKey = "Questions"
	question.RowKey = id

	question.Answer1Id = answerid1
	question.Answer2Id = answerid2
	question.CreatedDate = time.Now()

	marshalledQuestion, err := json.Marshal(question)
	if err != nil {
		panic(err)
	}

	respQ, err := questionClient.AddEntity(context.TODO(), marshalledQuestion, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(respQ)

	// At this point, we also need to add the enteries for the votes.
	var a VoteEntity
	a.PartitionKey = question.RowKey
	a.RowKey = question.Answer1Id
	a.VoteCount = 0

	marshalledVote1, err := json.Marshal(a)
	if err != nil {
		panic(err)
	}

	respA, err := voteClient.AddEntity(context.TODO(), marshalledVote1, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(respA)

	var b VoteEntity
	b.PartitionKey = question.RowKey
	b.RowKey = question.Answer2Id
	b.VoteCount = 0

	marshalledVote2, err := json.Marshal(b)
	if err != nil {
		panic(err)
	}

	respB, err := voteClient.AddEntity(context.TODO(), marshalledVote2, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(respB)

	return question
}

func deleteQuestionFromTableStorage(serviceClient *aztables.ServiceClient, question QuestionEntity) {
	// Have to delete both the question AND the votes.
	questionClient := serviceClient.NewClient("questions")
	voteClient := serviceClient.NewClient("votes")

	// First, lets get rid of the votes.
	voteClient.DeleteEntity(context.TODO(), question.RowKey, question.Answer1Id, nil)
	voteClient.DeleteEntity(context.TODO(), question.RowKey, question.Answer2Id, nil)

	// Then the question.
	questionClient.DeleteEntity(context.TODO(), question.PartitionKey, question.RowKey, nil)

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
	fmt.Println("Endpoint hit: returnQuestion")
	vars := mux.Vars(r)
	id := vars["id"]

	serviceClient := getServiceClient()
	json.NewEncoder(w).Encode(getQuestionFromTableStorage(serviceClient, id))
}

func createNewQuestion(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint hit: createQuestion")
	// get the body of our POST request
	// return the string response containing the request body
	reqBody, _ := ioutil.ReadAll(r.Body)

	var question QuestionEntity
	json.Unmarshal(reqBody, &question)

	serviceClient := getServiceClient()
	newQuestion := addQuestionToTableStorage(serviceClient, question)

	json.NewEncoder(w).Encode(newQuestion)
}

func deleteQuestion(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint hit: deleteQuestion")
	// get the body of our DELETE request
	// return the string response containing the request body
	reqBody, _ := ioutil.ReadAll(r.Body)

	var question QuestionEntity
	json.Unmarshal(reqBody, &question)

	serviceClient := getServiceClient()
	deleteQuestionFromTableStorage(serviceClient, question)
	w.WriteHeader(200)
}

func handleRequests() {
	// Lets use the Mux Router, since everyone else does.
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", home)
	router.HandleFunc("/questions", returnAllQuestions)
	router.HandleFunc("/question", createNewQuestion).Methods("POST")
	router.HandleFunc("/question/{id}", deleteQuestion).Methods("DELETE")
	//router.HandleFunc("/question/{id}", returnQuestion)

	router.Use(mux.CORSMethodMiddleware(router))

	log.Fatal(http.ListenAndServe(":10000", router))
}

func main() {
	fmt.Println("Hello Votesy!")

	fmt.Println("Authenicating...")
	serviceClient := getServiceClient()
	getAllQuestionsFromTableStorage(serviceClient)

	handleRequests()
}
