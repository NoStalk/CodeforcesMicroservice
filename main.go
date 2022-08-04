package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	// "time"
	platformDatapb "github.com/NoStalk/cfMicroservices/proto"
	utilities "github.com/NoStalk/serviceUtilities"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

/// This struct is used to parse the JSON data received from the Codeforces(Contests) API.

func UnmarshalCFSubmissionResponse(data []byte) (CFSubmissionResponse, error) {
	var r CFSubmissionResponse;
	err := json.Unmarshal(data, &r);
	return r, err;
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


func UnmarshalCFContestResponse(data []byte) (CFContestResponse, error) {
	var r CFContestResponse
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *CFContestResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type CFContestResponse struct {
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
	userEmail := req.GetEmail();
	rawSubmissionsArray := codeforcesSubmissionsRequestHandler(cfHandle);
	grpcSubmissionsResponseArray := submissionDataConverterforGrpcResponse(rawSubmissionsArray);
	response := &platformDatapb.SubmissionResponse{
		Submissions: grpcSubmissionsResponseArray,
	}
	mongoURI := os.Getenv("DB_URI");
	start := time.Now();
	var dbResources utilities.DBResources;
	var err error;
	dbResources, err = utilities.OpenDatabaseConnection(mongoURI);
	if err != nil {
		log.Printf("Couldnt connect to Database: %v", err);
	}
	writeUserSubmissionsToDB(dbResources, userEmail, grpcSubmissionsResponseArray);
	utilities.CloseDatabaseConnection(dbResources);
	fmt.Printf("Time taken to cast and write to DB: %f\n",time.Since(start).Seconds());
	return response, nil;
}




/**
* @brief This is unary grpc function that is invoked when a client calls the GetUserContests function
* @param ctx The context of the grpc call and the request object.
* @return The response of the type platformDatapb.Response
**/


func (*server) GetUserContests(ctx context.Context, req *platformDatapb.Request) (*platformDatapb.ContestResponse, error){
	fmt.Println("GetUserContests function invoked with user handle");
	cfHandle := req.GetUserHandle();
	userEmail := req.GetEmail();
	rawContestsArray := codeforcesContestRequestHandler(cfHandle);
	grpcContestsResponseArray := contestDataConverterforGrpcresponse(rawContestsArray);
	response := &platformDatapb.ContestResponse{
		Contests: grpcContestsResponseArray,
	}
	mongoURI := os.Getenv("DB_URI");
	start := time.Now();
	var dbResources utilities.DBResources;
	var err error;
	dbResources, err = utilities.OpenDatabaseConnection(mongoURI);
	if err != nil {
		log.Printf("Couldnt connect to Database: %v", err);
	}
	writeUserContestsToDB(dbResources, userEmail, grpcContestsResponseArray);
	utilities.CloseDatabaseConnection(dbResources);
	fmt.Printf("Time taken to cast and write to DB: %f\n",time.Since(start).Seconds());

	return response, nil;
}



func(*server) GetAllUserData(stream platformDatapb.FetchPlatformData_GetAllUserDataServer) error{
	fmt.Println("Bi-Directional Streamin function invoked");
	mongoURI := os.Getenv("DB_URI");
	start := time.Now();
	var dbResources utilities.DBResources;
	var err error;
	dbResources, err = utilities.OpenDatabaseConnection(mongoURI);
	if err != nil {
		log.Printf("Couldnt connect to Database: %v", err);
	}

	for{
		req, err := stream.Recv();

		if err == io.EOF {
			utilities.CloseDatabaseConnection(dbResources);
			fmt.Printf("Time taken to write all the data and send confirmation: %f\n", time.Since(start).Seconds());
			return nil;
		}

		if err != nil {
			fmt.Printf("Error while reading the stream: %v", err);
		}
		cfHandle := req.GetUserHandle();
		userEmail := req.GetEmail();
		
		//The conversion is necessary for writing to the database, maybe will optimise later
		SubmissionsData := submissionDataConverterforGrpcResponse(codeforcesSubmissionsRequestHandler(cfHandle));
		ContestsData := contestDataConverterforGrpcresponse(codeforcesContestRequestHandler(cfHandle));

		errWritingSubmissionToDB := writeUserSubmissionsToDB(dbResources,userEmail,SubmissionsData);
		errWritingContestToDB := writeUserContestsToDB(dbResources,userEmail,ContestsData);
		var userStatus bool = true;
		if errWritingContestToDB != nil || errWritingSubmissionToDB != nil{
			userStatus = false;
		}

		sendErr := stream.Send(&platformDatapb.OperationStatus{
			Status: userStatus,
			UserHandle: cfHandle,
		})		
		
		if sendErr != nil {
			fmt.Printf("Error while sending data to client: %v", sendErr);
			utilities.CloseDatabaseConnection(dbResources);
			return sendErr;
		}


	}
}





/**
* The function that is used to type convert the Codeforces API response to the format that is expected by the DB.
* @param The database resoources, the email id of the user, submissions array i.e. the response from the grpc request.
* @return None.
* Room for improvement, maybe reduce number of type conversions in the future!
**/


func writeUserSubmissionsToDB(dbResources utilities.DBResources, userEmail string, submissions []*platformDatapb.Submission) error {
	var newSubmissionsArray []utilities.SubmissionData;
	for _, submission := range submissions{
		submissonObject := utilities.SubmissionData{
			ProblemUrl: submission.ProblemLink,
			ProblemName: submission.ProblemTitle,
			SubmissionDate: submission.Date,
			SubmissionLanguage: submission.Language,
			SubmissionStatus: submission.ProblemStatus,
			CodeUrl: submission.CodeLink,
		}
		newSubmissionsArray = append(newSubmissionsArray, submissonObject);
	}
	err := utilities.AppendSubmissionData(dbResources,userEmail,"codeforces",newSubmissionsArray);
	if err != nil {
		fmt.Printf("Couldnt write Submissions Data to DB: %v", err);
		return err;
	}	
	fmt.Println("successfully wrote data to DB");
	return nil;
}




/**
* The function that is used to type convert the Codeforces API response to the format that is expected by the DB.
* @param The database resoources, th email id of the user, contest array i.e. the response from the grpc request.
* @return None.
* Room for improvement, maybe reduce number of type conversions in the future!
**/


func writeUserContestsToDB(dbResources utilities.DBResources, userEmail string, contests []*platformDatapb.Contest) error {
	var newContestsArray []utilities.ContestData;
	for _, contest := range contests{
		contestObject := utilities.ContestData{
		   ContestName: contest.ContestName,
		   Rank: contest.Rank,
		   OldRating: contest.OldRating,
		   NewRating: contest.NewRating,
		   ContestID: contest.ContestId,		
		}
		newContestsArray = append(newContestsArray, contestObject);
	}
	err := utilities.AppendContestData(dbResources,userEmail,"codeforces",newContestsArray);

	if err != nil {
		fmt.Printf("Couldnt write Contest Data to DB: %v", err);
		return err;
	}
	
	fmt.Println("successfully wrote data to DB");
	return err;
}



/**
* @brief Queries the codeforces API with a user handle and returns the user's recent submissions
* @param handle The user handle to query
* @return An array user's recent submissions data
**/


func codeforcesSubmissionsRequestHandler(cfHandle string) ([]Submissions) { 
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

	return responseCFSubmissions.Submissions;
}


/**
* @brief Queries the codeforces API with a user handle and returns the user's recent contests
* @param handle The user handle to query
* @return An array user's recent contests data
**/


func codeforcesContestRequestHandler(cfHandle string) ([]Contests){
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

	return responseCFContests.Contests;
}



/**
* @brief Helper function to typecast the submissions response array recieved from querying the codeforces service
* @param The submissions array that is modelled like the json response from the api call
* @return The submissions array typecasted(type: []*platformDatapb.Submission) so as to be sent as a response for the unary grpc call
**/

func submissionDataConverterforGrpcResponse(submissionArray []Submissions) ([]*platformDatapb.Submission){
	submissionsArrayforGRPC := []*platformDatapb.Submission{};

	for _, submission := range submissionArray{
		submissionResponseObject := platformDatapb.Submission{
			Date: strconv.FormatInt(submission.CreationTimeSeconds,10),
			Language: submission.ProgrammingLanguage,
			ProblemStatus: submission.Verdict,
			ProblemTitle: submission.Problem.Name,
			ProblemLink: "https://codeforces.com/contest/"+strconv.FormatInt(submission.ContestID, 10)+"/problem/"+submission.Problem.Index,
			CodeLink: "https://codeforces.com/contest/"+strconv.FormatInt(submission.ContestID, 10)+"/submission/"+strconv.FormatInt(submission.ID, 10),
		}
		submissionsArrayforGRPC = append(submissionsArrayforGRPC, &submissionResponseObject);
	}
		return submissionsArrayforGRPC;
}


/**
* @brief Helper function to typecast the contests response array recieved from querying the codeforces service
* @param The contests array that is modelled like the json response from the api call
* @return The contests array typecasted(type: []*platformDatapb.Contest) so as to be sent as a response for the unary grpc call
**/


func contestDataConverterforGrpcresponse(contestArray []Contests) ([]*platformDatapb.Contest){
		contestsResponseforGrpc := []*platformDatapb.Contest{};

	for _, contest := range contestArray{
		contestResponseObject := platformDatapb.Contest{
			ContestName: contest.ContestName,
			Rank: contest.Rank,
			OldRating: contest.OldRating,
			NewRating: contest.NewRating,
			ContestId: contest.ContestID,
		}
		contestsResponseforGrpc = append(contestsResponseforGrpc, &contestResponseObject);
	}

	return contestsResponseforGrpc;
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