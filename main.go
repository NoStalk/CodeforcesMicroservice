package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"

	// "time"

	// utilities "github.com/NoStalk/serviceUtilities"
	platformDatapb "github.com/NoStalk/cfMicroservices/proto"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)



/// This struct is used to parse the JSON data received from the Codeforces(Contests) API.

func UnmarshalCFSubmissionResponse(data []byte) (CFSubmissionResponse, error) {
	var r CFSubmissionResponse
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *CFSubmissionResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type CFSubmissionResponse struct {
	Status string   `json:"status"`
	Submissions []Submissions `json:"result"`
}

type Submissions struct {
	ID                  int64   `json:"id"`                 
	ContestID           int64   `json:"contestId"`          
	CreationTimeSeconds int64   `json:"creationTimeSeconds"`
	RelativeTimeSeconds int64   `json:"relativeTimeSeconds"`
	Problem             Problem `json:"problem"`            
	Author              Author  `json:"author"`             
	ProgrammingLanguage string  `json:"programmingLanguage"`
	Verdict             string  `json:"verdict"`            
	Testset             string  `json:"testset"`            
	PassedTestCount     int64   `json:"passedTestCount"`    
	TimeConsumedMillis  int64   `json:"timeConsumedMillis"` 
	MemoryConsumedBytes int64   `json:"memoryConsumedBytes"`
}

type Author struct {
	ContestID        int64    `json:"contestId"`       
	Members          []Member `json:"members"`         
	ParticipantType  string   `json:"participantType"` 
	Ghost            bool     `json:"ghost"`           
	StartTimeSeconds int64    `json:"startTimeSeconds"`
	Room             *int64   `json:"room,omitempty"`  
}

type Member struct {
	Handle string `json:"handle"`
}

type Problem struct {
	ContestID int64    `json:"contestId"`
	Index     string   `json:"index"`    
	Name      string   `json:"name"`     
	Type      string   `json:"type"`     
	Points    float64    `json:"points"`   
	Rating    int64    `json:"rating"`   
	Tags      []string `json:"tags"`     
}






// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    welcome, err := UnmarshalWelcome(bytes)
//    bytes, err = welcome.Marshal()
// This struct is used to parse the JSON data received from the Codeforces(Contests) API.


func UnmarshalCFContestResponse(data []byte) (Welcome, error) {
	var r Welcome
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Welcome) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Welcome struct {
	Status string   `json:"status"`
	Contests []Contests `json:"result"`
}

type Contests struct {
	ContestID               int64  `json:"contestId"`              
	ContestName             string `json:"contestName"`            
	Handle                  Handle `json:"handle"`                 
	Rank                    int64  `json:"rank"`                   
	RatingUpdateTimeSeconds int64  `json:"ratingUpdateTimeSeconds"`
	OldRating               int64  `json:"oldRating"`              
	NewRating               int64  `json:"newRating"`              
}

type Handle string;









type server struct{
	platformDatapb.UnimplementedFetchPlatformDataServer;
}



/**
* @brief This is unary grpc function that is invoked when a client calls the GetUserSubmissions function
* @param ctx The context of the grpc call.
* @param req The request of the type codeforcesMSpb.Request
* @return The response of the type codeforcesMSpb.Response
**/


func (*server) GetUserSubmissions(ctx context.Context, req *platformDatapb.Request) (*platformDatapb.SubmissionResponse, error){
	fmt.Println("GetUserSubmissions function invoked with user handle");
	cfHandle := req.GetUserHandle();
	submissions := codeforcesSubmissionsRequestHandler(cfHandle);
	response := &platformDatapb.SubmissionResponse{
		Submissions: submissions,
	}

	return response, nil;
}


func (*server) GetUserContests(ctx context.Context, req *platformDatapb.Request) (*platformDatapb.ContestResponse, error){
	fmt.Println("GetUserContests function invoked with user handle");
	cfHandle := req.GetUserHandle();
	contests := codeforcesContestRequestHandler(cfHandle);
	response := &platformDatapb.ContestResponse{
		Contests: contests,
	}
	return response, nil;
}



/**
* @brief Queries the codeforces API with a user handle and returns the user's recent submissions
* @param handle The user handle to query
* @return An array user's recent submissions and an array of user's contests
**/


func codeforcesSubmissionsRequestHandler(cfHandle string) ([]*platformDatapb.Submission) { 
	queryString := "https://codeforces.com/api/user.status?handle=" + cfHandle;
	response, err := http.Get(queryString);

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body);
	if err != nil {
		log.Fatal(err)
	}
	responseCFSubmissions, err := UnmarshalCFSubmissionResponse(responseData);
	if err!=nil {
		fmt.Printf("Couldnt unmarshal the byte slice: %v", err);
	}
	fmt.Println(len(responseCFSubmissions.Submissions));
	submissions := []*platformDatapb.Submission{
		{
			Date: "24oct18",
			Language: "C++",
			ProblemStatus: "Accepted",
			ProblemTitle: "Brackets",
			ProblemLink: "www.codeforces.com/problem/1020/",
			CodeLink: "www.codeforces.com/submission/1020/",
		},
	};

	return submissions;
}

func codeforcesContestRequestHandler(cfHandle string) ([]*platformDatapb.Contest){
	queryString := "https://codeforces.com/api/user.rating?handle=" + cfHandle;
	response, err := http.Get(queryString);

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body);
	if err != nil {
		log.Fatal(err)
	}
	responseCFContests, err := UnmarshalCFContestResponse(responseData);
	if err!=nil {
		fmt.Printf("Couldnt unmarshal the byte slice: %v", err);
	}
	fmt.Println(len(responseCFContests.Contests));
	contests := []*platformDatapb.Contest{
		{
			ContestName: "Codeforces Round #664 (Div. 2)",
			Rank: 12345,
			OldRating: 1234,
			NewRating: 1268,
			RatingUpdateTimeSeconds: 1539098200,
			ContestId: 664,
		},
	};
	return contests;
}




func main(){
	err := godotenv.Load();
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	lis, err := net.Listen("tcp", "0.0.0.0:5003");
	if err != nil {
		log.Fatalf("Failed to listen: %v", err);
	}
	s := grpc.NewServer();
	platformDatapb.RegisterFetchPlatformDataServer(s, &server{});
	fmt.Println("Server started on port 5003");
	if err = s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err);
	} 
	submissionsArray := codeforcesSubmissionsRequestHandler("zeus_codes");
	contestsArray := codeforcesContestRequestHandler("zeus_codes");
	fmt.Println(len(submissionsArray), len(contestsArray));

}