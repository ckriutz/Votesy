package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/aidarkhanov/nanoid"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
)

type QuestionEntity struct {
	Id          string `json:"id"`
	Text        string `json:"text"`
	Answer1Id   string `json:"answer1Id"`
	Answer1Text string `json:"answer1Text"`
	Answer2Id   string `json:"answer2Id"`
	Answer2Text string `json:"answer2Text"`
	IsCurrent   bool   `json:"isCurrent"`
	CreatedDate time.Time
}

type VoteEntity struct {
	Id        string `json:"id"`
	VoteCount int32  `json:"voteCount"`
}

func getRedisClient() *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "192.168.0.239:6379",
		Password: "",
		DB:       0,
	})

	return redisClient
}

func getCurrentQuestionFromRedis(redisClient *redis.Client) QuestionEntity {
	fmt.Println("Getting the current question from Redis...")
	var ctx = context.Background()

	// Get the current question.
	question, err := redisClient.Get(ctx, "question:current").Result()
	if err != nil {
		panic(err)
	}

	var currentQuestion QuestionEntity
	err = json.Unmarshal([]byte(question), &currentQuestion)
	if err != nil {
		panic(err)
	}

	return currentQuestion
}

func getAllQuestionsFromStorage(redisClient *redis.Client) []QuestionEntity {
	Quesitons := []QuestionEntity{}

	// Get the number of questions in the questions list in redis.
	var ctx = context.Background()

	qlist, err := redisClient.SMembers(ctx, "questions").Result()
	if err != nil {
		panic(err)
	}

	fmt.Printf("qlist: %v\n", qlist)

	for _, entity := range qlist {
		var myEntity QuestionEntity

		err = json.Unmarshal([]byte(entity), &myEntity)
		if err != nil {
			panic(err)
		}
		Quesitons = append(Quesitons, myEntity)
	}
	//fmt.Printf("qlist: %v\n", qlist)

	//cnt, err := redisClient.LLen(ctx, "questions").Result()
	//if err != nil {
	//panic(err)
	//}

	//questions, _ := redisClient.LRange(ctx, "questions", 0, cnt).Result()
	//for _, entity := range questions {
	//	var myEntity QuestionEntity

	//	err = json.Unmarshal([]byte(entity), &myEntity)
	//	if err != nil {
	//		panic(err)
	//	}
	//	Quesitons = append(Quesitons, myEntity)
	//}

	return Quesitons
}

func getQuestionFromStorage(redisClient *redis.Client, id string) QuestionEntity {
	// Get the question from Redis based on the Id.
	var ctx = context.Background()
	question, err := redisClient.Get(ctx, "question:"+id).Result()
	if err != nil {
		panic(err)
	}

	var questionEnty QuestionEntity
	err = json.Unmarshal([]byte(question), &questionEnty)
	if err != nil {
		panic(err)
	}

	return questionEnty
}

func addQuestionToStorage(redisClient *redis.Client, question QuestionEntity) QuestionEntity {
	// Have to create both the question AND the votes.

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

	question.Id = id
	question.Answer1Id = answerid1
	question.Answer2Id = answerid2
	question.CreatedDate = time.Now()

	marshalledQuestion, err := json.Marshal(question)
	if err != nil {
		panic(err)
	}

	var ctx = context.Background()
	redisClient.SAdd(ctx, "questions", marshalledQuestion)
	redisClient.Set(ctx, "question:"+question.Id, marshalledQuestion, 0)
	redisClient.Set(ctx, "votes:"+question.Answer1Id, 0, 0)
	redisClient.Set(ctx, "votes:"+question.Answer2Id, 0, 0)

	return question
}

func deleteQuestionFromStorage(redisClient *redis.Client, questionId string) bool {
	var ctx = context.Background()

	fmt.Println("Deleting question with id: " + questionId)

	// We need to get the question so we can get the answer id's.
	question := getQuestionFromStorage(redisClient, questionId)

	// Delete the question
	rslt, err := redisClient.Del(ctx, "question:"+questionId).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(rslt)

	// Delete question from the set
	questionString, _ := json.Marshal(question)
	fmt.Println(question)
	rslts, errsrem := redisClient.SRem(ctx, "questions", questionString).Result()
	if errsrem != nil {
		panic(errsrem)
	}
	fmt.Println(rslts)

	// Delete the votes
	rsltsv1, errv1 := redisClient.Del(ctx, "votes:"+question.Answer1Id).Result()
	if errv1 != nil {
		panic(errv1)
	}
	fmt.Println(rsltsv1)

	rsltsv2, errv2 := redisClient.Del(ctx, "votes:"+question.Answer2Id).Result()
	if errv2 != nil {
		panic(errv2)
	}
	fmt.Println(rsltsv2)

	return true
}

func getVotesForQuestionFromStorage(redisClient *redis.Client, questionId string) []VoteEntity {
	Votes := []VoteEntity{}

	// So we need to look up all the votes for the question, based on the quesstionId.
	var ctx = context.Background()

	// Get the question from Redis based on the Id.
	question := getQuestionFromStorage(redisClient, questionId)

	// Now we have the question, we can get the votes.
	vote1, err := redisClient.Get(ctx, "votes:"+question.Answer1Id).Result()
	if err != nil {
		panic(err)
	}

	vote2, err := redisClient.Get(ctx, "votes:"+question.Answer2Id).Result()
	if err != nil {
		panic(err)
	}

	vote1_int64, _ := strconv.ParseInt(vote1, 10, 32)
	vote2_int64, _ := strconv.ParseInt(vote2, 10, 32)

	vote1Entity := VoteEntity{
		Id:        question.Answer1Id,
		VoteCount: int32(vote1_int64),
	}
	vote2Entity := VoteEntity{
		Id:        question.Answer2Id,
		VoteCount: int32(vote2_int64),
	}

	Votes = append(Votes, vote1Entity)
	Votes = append(Votes, vote2Entity)

	return Votes
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "votesy-api in golang!")
	fmt.Println("Endpoint hit: home")
}

func returnAllQuestions(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint hit: returnAllQuestions")

	redisClient := getRedisClient()
	json.NewEncoder(w).Encode(getAllQuestionsFromStorage(redisClient))
}

func returnQuestion(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint hit: returnQuestion")
	vars := mux.Vars(r)
	id := vars["id"]

	redisClient := getRedisClient()
	json.NewEncoder(w).Encode(getQuestionFromStorage(redisClient, id))
}

func createNewQuestion(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint hit: createQuestion")
	// get the body of our POST request
	// return the string response containing the request body
	reqBody, _ := io.ReadAll(r.Body)

	var question QuestionEntity
	json.Unmarshal(reqBody, &question)

	redisClient := getRedisClient()
	newQuestion := addQuestionToStorage(redisClient, question)

	json.NewEncoder(w).Encode(newQuestion)
}

func deleteQuestion(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint hit: deleteQuestion")

	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Println(id)

	// Now that we have the id, we can send that for deletion.
	redisClient := getRedisClient()
	resp := deleteQuestionFromStorage(redisClient, id)
	json.NewEncoder(w).Encode(resp)
}

func getVotesByQuestionId(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint hit: getVotesByQuestionId")
	vars := mux.Vars(r)
	id := vars["questionId"]

	redisClient := getRedisClient()
	json.NewEncoder(w).Encode(getVotesForQuestionFromStorage(redisClient, id))
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

	// Test redis section
	//router.HandleFunc("/redis", checkRedis)

	// First questions.
	router.HandleFunc("/questions", returnAllQuestions)
	router.HandleFunc("/question", createNewQuestion).Methods("POST")
	router.HandleFunc("/question/{id}", deleteQuestion).Methods("DELETE")
	//router.HandleFunc("/question/{id}", updateQuestion).Methods("PUT")
	router.HandleFunc("/question/{id}", returnQuestion)

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

	fmt.Println("Connecting to Redis...")

	rds := getRedisClient()
	pong, err := rds.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}

	fmt.Println(pong)

	var ctx = context.Background()

	// See if there are any questions in storage.
	cnt, err := rds.SCard(ctx, "questions").Result()
	if err != nil {
		panic(err)
	}

	if cnt == 0 {
		// This is a test to see if there are questions. On a new install into an environment, there might
		// not be any questions, and if that's the case, let's add one.
		fmt.Println("It seems there are no questions here, lets add the default one, and the second one.")
		qst1 := QuestionEntity{
			Id:          "0",
			Text:        "~ Bear vs Owl ~",
			Answer1Id:   "0",
			Answer1Text: "Bear",
			Answer2Id:   "1",
			Answer2Text: "Owl",
			IsCurrent:   true,
		}
		u1, _ := json.Marshal(qst1)

		qst2 := QuestionEntity{
			Id:          "1",
			Text:        "Beach Or Mountains?",
			Answer1Id:   "2",
			Answer1Text: "Beach",
			Answer2Id:   "3",
			Answer2Text: "Mountains",
			IsCurrent:   false,
		}
		u2, _ := json.Marshal(qst2)

		rds.SAdd(ctx, "questions", u1)
		rds.SAdd(ctx, "questions", u2)
		rds.Set(ctx, "question:"+qst1.Id, u1, 0)
		rds.Set(ctx, "question:"+qst2.Id, u2, 0)
		rds.Set(ctx, "question:current", u1, 0)
		rds.Set(ctx, "votes:"+qst1.Answer1Id, 0, 0)
		rds.Set(ctx, "votes:"+qst1.Answer2Id, 0, 0)
		rds.Set(ctx, "votes:"+qst2.Answer1Id, 0, 0)
		rds.Set(ctx, "votes:"+qst2.Answer2Id, 0, 0)
	} else {
		fmt.Println("There are", cnt, "questions in storage.")

	}

	// Now lets see if there is a current question.
	currentQuestion := getCurrentQuestionFromRedis(rds)
	fmt.Printf("Current Question: %v\n", currentQuestion)

	handleRequests()
}
