package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
	"github.com/aidarkhanov/nanoid"
	"github.com/gorilla/mux"
)

type QuestionEntity struct {
	aztables.Entity
	Id          string `json:"id"`
	Text        string `json:"text"`
	Answer1Id   string `json:"answer1Id"`
	Answer1Text string `json:"answer1Text"`
	Answer2Id   string `json:"answer2Id"`
	Answer2Text string `json:"answer2Text"`
	Answer3Id   string `json:"answer3Id"`
	Answer3Text string `json:"answer3Text"`
	Answer4Id   string `json:"answer4Id"`
	Answer4Text string `json:"answer4Text"`
	IsCurrent   bool   `json:"isCurrent"`
	IsUsed      bool   `json:"isUsed"`
	CreatedDate time.Time
}

type VoteEntity struct {
	aztables.Entity
	Id        string `json:"id"`
	VoteCount int32  `json:"voteCount"`
}

func getTableServiceClient() *aztables.ServiceClient {
	// Get the connection string from the environment.
	storageConnectionString := os.Getenv("STORAGE_CONNECTION_STRING")
	serviceClient, err := aztables.NewServiceClientFromConnectionString(storageConnectionString, nil)
	if err != nil {
		panic(err)
	}

	return serviceClient
}

func getCurrentQuestionFromStorage(serviceClient *aztables.ServiceClient) QuestionEntity {
	fmt.Println("Getting the current question from Table Storage...")

	// Get the current question from storage.
	filter := "isCurrent eq true"
	options := &aztables.ListEntitiesOptions{
		Filter: &filter,
	}

	// Now check to see if the count is 0, if it is, we need to add some questions.
	pager := serviceClient.NewClient("questions").NewListEntitiesPager(options)
	entity := QuestionEntity{}

	for {
		if !pager.More() {
			break
		}
		resp, _ := pager.NextPage(context.Background())
		if resp.Entities == nil {
			fmt.Println("No questions found.")
			break
		}
		err := json.Unmarshal([]byte(resp.Entities[0]), &entity)
		if err != nil {
			panic(err)
		}
		break
	}

	return entity
}

func getAllQuestionsFromStorage(serviceClient *aztables.ServiceClient, partitionKey string) []QuestionEntity {
	questions := []QuestionEntity{}
	ctx := context.Background()

	filter := fmt.Sprintf("PartitionKey eq '%s'", partitionKey)
	options := &aztables.ListEntitiesOptions{
		Filter: &filter,
	}
	// Create a pager to list entities
	pager := serviceClient.NewClient("questions").NewListEntitiesPager(options)

	for pager.More() {
		resp, err := pager.NextPage(ctx)
		if err != nil {
			panic(err)
		}
		for _, entity := range resp.Entities {
			var questionEntity QuestionEntity
			err = json.Unmarshal(entity, &questionEntity)
			if err != nil {
				panic(err)
			}
			questions = append(questions, questionEntity)
		}
	}

	return questions
}

func getQuestionFromStorage(serviceClient *aztables.ServiceClient, partitionKey, rowKey string) (QuestionEntity, error) {
	var ctx = context.Background()
	tableClient := serviceClient.NewClient("questions")

	// Get the entity from Azure Table Storage
	resp, err := tableClient.GetEntity(ctx, partitionKey, rowKey, nil)
	if err != nil {
		return QuestionEntity{}, err
	}

	// Unmarshal the entity into a QuestionEntity struct
	var questionEntity QuestionEntity
	err = json.Unmarshal(resp.Value, &questionEntity)
	if err != nil {
		return QuestionEntity{}, err
	}

	return questionEntity, nil
}

func addQuestionToStorage(serviceClient *aztables.ServiceClient, question QuestionEntity) QuestionEntity {
	fmt.Println("Adding Entity to Azure Table Storage...")

	// Have to create both the question AND the votes.
	ctx := context.Background()
	alphabet := "abcdefghijklmnopqrstuvwxyz0123456789"

	// Row Id
	id, err := nanoid.Generate(alphabet, 5)
	if err != nil {
		panic(err)
	}

	// Answer 1 Id
	answer1Id, err := nanoid.Generate(alphabet, 5)
	if err != nil {
		panic(err)
	}

	// Answer 2 Id
	answer2Id, err := nanoid.Generate(alphabet, 5)
	if err != nil {
		panic(err)
	}

	// Create the votes for the answers.
	votes1 := VoteEntity{
		Entity: aztables.Entity{
			PartitionKey: "votes",
			RowKey:       answer1Id,
		},
		Id:        answer1Id,
		VoteCount: 0,
	}
	votes1Bytes, err := json.Marshal(votes1)
	if err != nil {
		panic(err)
	}
	serviceClient.NewClient("votes").AddEntity(ctx, votes1Bytes, nil)

	votes2 := VoteEntity{
		Entity: aztables.Entity{
			PartitionKey: "votes",
			RowKey:       answer2Id,
		},
		Id:        answer2Id,
		VoteCount: 0,
	}
	votes2Bytes, err := json.Marshal(votes2)
	if err != nil {
		panic(err)
	}
	serviceClient.NewClient("votes").AddEntity(ctx, votes2Bytes, nil)

	// Sometimes we have 4 answers.
	if question.Answer3Text != "" {
		// Answer 3 Id
		answer3Id, err := nanoid.Generate(alphabet, 5)
		if err != nil {
			panic(err)
		}
		question.Answer3Id = answer3Id
		votes3 := VoteEntity{
			Entity: aztables.Entity{
				PartitionKey: "votes",
				RowKey:       answer3Id,
			},
			Id:        answer3Id,
			VoteCount: 0,
		}
		votes3Bytes, err := json.Marshal(votes3)
		if err != nil {
			panic(err)
		}
		serviceClient.NewClient("votes").AddEntity(ctx, votes3Bytes, nil)
	}

	// Sometimes we have 4 answers.
	if question.Answer4Text != "" {
		// Answer 3 Id
		answer4Id, err := nanoid.Generate(alphabet, 5)
		if err != nil {
			panic(err)
		}
		question.Answer4Id = answer4Id
		votes4 := VoteEntity{
			Entity: aztables.Entity{
				PartitionKey: "votes",
				RowKey:       answer4Id,
			},
			Id:        answer4Id,
			VoteCount: 0,
		}
		votes4Bytes, err := json.Marshal(votes4)
		if err != nil {
			panic(err)
		}
		serviceClient.NewClient("votes").AddEntity(ctx, votes4Bytes, nil)
	}

	question.Entity.PartitionKey = "questions"
	question.Entity.RowKey = id
	question.Id = id
	question.Answer1Id = answer1Id
	question.Answer2Id = answer2Id
	question.CreatedDate = time.Now()
	question.IsUsed = false

	questionBytes, err := json.Marshal(question)
	if err != nil {
		panic(err)
	}

	serviceClient.NewClient("questions").AddEntity(ctx, questionBytes, nil)

	fmt.Println("Done adding Entity to Azure Table Storage...")
	return question

}

func deleteQuestionFromStorage(serviceClient *aztables.ServiceClient, partitionKey, rowKey string) bool {
	var ctx = context.Background()

	fmt.Println("Deleting question with id: " + rowKey)

	// We need to get the question so we can get the answer id's.
	question, err := getQuestionFromStorage(serviceClient, partitionKey, rowKey)
	if err != nil {
		panic(err)
	}

	// Delete the question
	serviceClient.NewClient("questions").DeleteEntity(ctx, partitionKey, rowKey, nil)

	// Then, delete the votes.
	if question.Answer1Id != "" {
		serviceClient.NewClient("votes").DeleteEntity(ctx, "votes", question.Answer1Id, nil)
	}

	if question.Answer2Id != "" {
		serviceClient.NewClient("votes").DeleteEntity(ctx, "votes", question.Answer2Id, nil)
	}

	if question.Answer3Id != "" {
		serviceClient.NewClient("votes").DeleteEntity(ctx, "votes", question.Answer3Id, nil)
	}

	if question.Answer4Id != "" {
		serviceClient.NewClient("votes").DeleteEntity(ctx, "votes", question.Answer4Id, nil)
	}

	return true
}

func getVotesForQuestionFromStorage(serviceClient *aztables.ServiceClient, partitionKey, rowKey string) []VoteEntity {
	Votes := []VoteEntity{}

	// 	// So we need to look up all the votes for the question, based on the quesstionId.
	var ctx = context.Background()

	// Get the question from Redis based on the Id.
	question, err := getQuestionFromStorage(serviceClient, partitionKey, rowKey)
	if err != nil {
		panic(err)
	}

	// Get all the votes for the question.
	tableClient := serviceClient.NewClient("votes")

	// Get the entity from Azure Table Storage
	vote1, err := tableClient.GetEntity(ctx, "votes", question.Answer1Id, nil)
	if err != nil {
		panic(err)
	}
	var vote1Entity VoteEntity
	err = json.Unmarshal(vote1.Value, &vote1Entity)
	if err != nil {
		panic(err)
	}
	Votes = append(Votes, vote1Entity)

	vote2, err := tableClient.GetEntity(ctx, "votes", question.Answer2Id, nil)
	if err != nil {
		panic(err)
	}
	var vote2Entity VoteEntity
	err = json.Unmarshal(vote2.Value, &vote2Entity)
	if err != nil {
		panic(err)
	}
	Votes = append(Votes, vote2Entity)

	// If there are more...
	if question.Answer3Id != "" {
		vote3, err := tableClient.GetEntity(ctx, "votes", question.Answer3Id, nil)
		if err != nil {
			panic(err)
		}
		var vote3Entity VoteEntity
		err = json.Unmarshal(vote3.Value, &vote3Entity)
		if err != nil {
			panic(err)
		}
		Votes = append(Votes, vote3Entity)
	}

	if question.Answer4Id != "" {
		vote4, err := tableClient.GetEntity(ctx, "votes", question.Answer4Id, nil)
		if err != nil {
			panic(err)
		}
		var vote4Entity VoteEntity
		err = json.Unmarshal(vote4.Value, &vote4Entity)
		if err != nil {
			panic(err)
		}
		Votes = append(Votes, vote4Entity)
	}

	return Votes
}

func updateQuestionInStorage(serviceClient *aztables.ServiceClient, partitionKey, rowKey string, question QuestionEntity) QuestionEntity {
	fmt.Println("Updating Entity in Azure Table Storage...")

	// Get the question from storage.
	questionEntity, err := getQuestionFromStorage(serviceClient, partitionKey, rowKey)
	if err != nil {
		panic(err)
	}

	// Update the question.
	questionEntity.Text = question.Text
	questionEntity.Answer1Text = question.Answer1Text
	questionEntity.Answer2Text = question.Answer2Text
	questionEntity.Answer3Text = question.Answer3Text
	questionEntity.Answer4Text = question.Answer4Text
	questionEntity.IsCurrent = question.IsCurrent
	questionEntity.IsUsed = question.IsUsed

	// Update the question in storage.
	questionBytes, err := json.Marshal(questionEntity)
	if err != nil {
		panic(err)
	}

	serviceClient.NewClient("questions").UpdateEntity(context.Background(), questionBytes, nil)

	fmt.Println("Done updating Entity in Azure Table Storage...")
	return questionEntity
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "votesy-api in golang!")
	fmt.Println("Endpoint hit: home")
}

func returnAllQuestions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Println("Endpoint hit: returnAllQuestions")

	client := getTableServiceClient()
	json.NewEncoder(w).Encode(getAllQuestionsFromStorage(client, "questions"))
}

func returnQuestion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Println("Endpoint hit: returnQuestion")
	vars := mux.Vars(r)
	partitionKey := vars["partitionKey"]
	rowKey := vars["rowKey"]

	client := getTableServiceClient()
	question, err := getQuestionFromStorage(client, partitionKey, rowKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(question)
}

func returnCurrentQuestion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Println("Endpoint hit: returnCurrentQuestion")

	client := getTableServiceClient()
	json.NewEncoder(w).Encode(getCurrentQuestionFromStorage(client))
}

func createNewQuestion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Println("Endpoint hit: createQuestion")
	// get the body of our POST request
	// return the string response containing the request body
	reqBody, _ := io.ReadAll(r.Body)

	var question QuestionEntity
	json.Unmarshal(reqBody, &question)

	client := getTableServiceClient()
	newQuestion := addQuestionToStorage(client, question)

	json.NewEncoder(w).Encode(newQuestion)
}

func deleteQuestion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Println("Endpoint hit: deleteQuestion")

	vars := mux.Vars(r)
	partitionKey := vars["partitionKey"]
	rowKey := vars["rowKey"]

	// Now that we have the id, we can send that for deletion.
	client := getTableServiceClient()
	resp := deleteQuestionFromStorage(client, partitionKey, rowKey)
	json.NewEncoder(w).Encode(resp)
}

func getVotesByQuestionId(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Println("Endpoint hit: getVotesByQuestionId")
	vars := mux.Vars(r)
	partitionKey := vars["partitionKey"]
	rowKey := vars["rowKey"]

	client := getTableServiceClient()
	json.NewEncoder(w).Encode(getVotesForQuestionFromStorage(client, partitionKey, rowKey))
}

func updateQuestion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Println("Endpoint hit: updateQuestion")
	vars := mux.Vars(r)
	partitionKey := vars["partitionKey"]
	rowKey := vars["rowKey"]

	// get the body of our POST request
	// return the string response containing the request body
	reqBody, _ := io.ReadAll(r.Body)

	var question QuestionEntity
	json.Unmarshal(reqBody, &question)

	client := getTableServiceClient()
	updatedQuestion := updateQuestionInStorage(client, partitionKey, rowKey, question)

	json.NewEncoder(w).Encode(updatedQuestion)
}

// Readyness functions
func ready(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Second)
	w.WriteHeader(http.StatusOK)
}

func liveness(w http.ResponseWriter, r *http.Request) {
	time.Sleep(250 * time.Millisecond)
	w.WriteHeader(http.StatusOK)
}

func startup(w http.ResponseWriter, r *http.Request) {
	time.Sleep(400 * time.Millisecond)
	w.WriteHeader(http.StatusOK)
}

func handleRequests() {
	// Lets use the Mux Router, since everyone else does.
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", home)

	// First questions.
	router.HandleFunc("/questions", returnAllQuestions).Methods("GET")
	router.HandleFunc("/question", createNewQuestion).Methods("POST")
	router.HandleFunc("/question/{partitionKey}/{rowKey}", deleteQuestion).Methods("DELETE")
	router.HandleFunc("/question/{partitionKey}/{rowKey}", updateQuestion).Methods("PUT")
	router.HandleFunc("/question/{partitionKey}/{rowKey}", returnQuestion).Methods("GET")
	router.HandleFunc("/questions/current", returnCurrentQuestion).Methods("GET")

	// Get Votes by question id
	router.HandleFunc("/votes/{questionId}", getVotesByQuestionId)

	// Let us add readyness probes!
	router.HandleFunc("/health/readiness", ready)
	router.HandleFunc("/health/liveness", liveness)
	router.HandleFunc("/health/startup", startup)

	router.Use(mux.CORSMethodMiddleware(router))

	log.Fatal(http.ListenAndServe(":10000", router))
}

func main() {
	fmt.Println("Hello Votesy!")

	fmt.Println("Connecting to Azure Table Storage...")

	ctx := context.Background()
	serviceClient := getTableServiceClient()

	// See if there are any questions in storage.
	filter := "PartitionKey eq 'questions'"
	options := &aztables.ListEntitiesOptions{
		Filter: &filter,
		Select: to.Ptr("RowKey"),
		// Only get the first entity
	}

	// Now check to see if the count is 0, if it is, we need to add some questions.
	pager := serviceClient.NewClient("questions").NewListEntitiesPager(options)

	cnt := 0
	for {
		if !pager.More() {
			break
		}
		resp, _ := pager.NextPage(ctx)
		cnt = len(resp.Entities)
		if resp.Entities == nil {
			fmt.Println("No questions found.")
			break
		}
		fmt.Println("Entities: ", cnt)
		break
	}

	if cnt == 0 {
		fmt.Println("Adding question...")
		question := QuestionEntity{
			Entity: aztables.Entity{
				PartitionKey: "questions",
			},
			Text:        "Bear or Owl?",
			Answer1Text: "Bear",
			Answer2Text: "Owl",
			IsCurrent:   true,
			CreatedDate: time.Now(),
		}

		addQuestionToStorage(serviceClient, question)
	}

	// Now lets see if there is a current question.
	currentQuestion := getCurrentQuestionFromStorage(serviceClient)
	fmt.Printf("Current Question: %v\n", currentQuestion.Text)

	// Now we can start the server.
	handleRequests()
}
