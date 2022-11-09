package main

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"

	// "time"
	platformDatapb "github.com/NoStalk/protoDefinitions"
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
//    CFContestResponse, err := UnmarshalWelcome(bytes)
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
	Handle                  string `json:"handle"`                 
	Rank                    int64  `json:"rank"`                   
	RatingUpdateTimeSeconds int64  `json:"ratingUpdateTimeSeconds"`
	OldRating               int64  `json:"oldRating"`              
	NewRating               int64  `json:"newRating"`              
}








func UnmarshalCFContestListResponse(data []byte) (ContestList, error) {
	var r ContestList
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *ContestList) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type ContestList struct {
	Status string   `json:"status"`
	AllContests []ContestResponse `json:"result"`
}

type ContestResponse struct {
	ID                  int64  `json:"id"`                 
	Name                string `json:"name"`               
	Type                string `json:"type"`               
	Phase               string `json:"phase"`              
	Frozen              bool   `json:"frozen"`             
	DurationSeconds     int64  `json:"durationSeconds"`    
	StartTimeSeconds    int64  `json:"startTimeSeconds"`   
	RelativeTimeSeconds int64  `json:"relativeTimeSeconds"`
}












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
	log.Println("GetUserSubmissions function invoked with user handle");
	cfHandle := req.GetUserHandle();
	userEmail := req.GetEmail();
	rawSubmissionsArray := codeforcesSubmissionsRequestHandler(cfHandle);
	
	
	submissionsArrayforDB := submissionDataConverterforDB(rawSubmissionsArray);
	mongoURI := os.Getenv("DB_URI");
	var dbResources utilities.DBResources;
	var err error;
	dbResources, err = utilities.OpenDatabaseConnection(mongoURI);
	if err != nil {
		log.Printf("Couldnt connect to Database: %v", err);
	}
	utilities.AppendSubmissionData(dbResources,userEmail,"Codeforces",submissionsArrayforDB);
	grpcSubmissionsResponseArray := utilities.CreateGRPCSubmissionResponseFromSubmissionSchema(submissionsArrayforDB);
	utilities.CloseDatabaseConnection(dbResources);

	response := &platformDatapb.SubmissionResponse{
		Submissions: grpcSubmissionsResponseArray,
	}
	
	
	
	return response, nil;
}




/**
* @brief This is unary grpc function that is invoked when a client calls the GetUserContests function
* @param ctx The context of the grpc call and the request object.
* @return The response of the type platformDatapb.Response
**/


func (*server) GetUserContests(ctx context.Context, req *platformDatapb.Request) (*platformDatapb.ContestResponse, error){
	log.Println("GetUserContests function invoked with user handle");
	cfHandle := req.GetUserHandle();
	userEmail := req.GetEmail();
	
	rawContestsArray := codeforcesContestRequestHandler(cfHandle);
	contestArrayforDB := contestDataConverterforDB(rawContestsArray);
	mongoURI := os.Getenv("DB_URI");
	var dbResources utilities.DBResources;
	var err error;

	dbResources, err = utilities.OpenDatabaseConnection(mongoURI);
	if err != nil {
		log.Printf("Couldnt connect to Database: %v", err);
	}
	utilities.AppendContestData(dbResources,userEmail,"Codeforces",contestArrayforDB);
	grpcContestsResponseArray := utilities.CreateGRPCContestResponseFromContestSchema(contestArrayforDB);
	utilities.CloseDatabaseConnection(dbResources);
	response := &platformDatapb.ContestResponse{
		Contests: grpcContestsResponseArray,
	}
	return response, nil;
}


/**
* @brief This is a Bi-Directional streaming grpc function that is invoked when a client calls the GetAllUserData function.
* @param ctx The stream interface that is has methods for sending and recieving data.
* @return The response is an error object that returns if there is an error during the streamin of data between the user and the client.
**/



func(*server) GetAllUserData(stream platformDatapb.FetchPlatformData_GetAllUserDataServer) error{
	log.Println("Bi-Directional Streamin function invoked");
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
			log.Printf("Time taken to write all the data and send confirmation: %f\n", time.Since(start).Seconds());
			return nil;
		}

		if err != nil {
			log.Printf("Error while reading the stream: %v", err);
		}
		cfHandle := req.GetUserHandle();
		userEmail := req.GetEmail();
		
		rawSubmissionsData := codeforcesSubmissionsRequestHandler(cfHandle);
		rawContestsData := codeforcesContestRequestHandler(cfHandle);

		submissionArrayforDB := submissionDataConverterforDB(rawSubmissionsData);
		contestDataforDB := contestDataConverterforDB(rawContestsData);

		errWritingSubmissionToDB :=utilities.AppendSubmissionData(dbResources,userEmail,"Codeforces",submissionArrayforDB);
		errWritingContestToDB :=utilities.AppendContestData(dbResources,userEmail,"Codeforces",contestDataforDB);

		var userStatus bool = true;
		if errWritingContestToDB != nil || errWritingSubmissionToDB != nil{
			userStatus = false;
		}

		sendErr := stream.Send(&platformDatapb.OperationStatus{
			Status: userStatus,
			UserHandle: cfHandle,
		})		
		
		if sendErr != nil {
			log.Printf("Error while sending data to client: %v", sendErr);
			utilities.CloseDatabaseConnection(dbResources);
			return sendErr;
		}

	}
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
			log.Print(err.Error())
			os.Exit(1)
		}
	
		responseData, err := ioutil.ReadAll(response.Body);
		if err != nil {
			log.Fatal(err)
		}
		responseCFSubmissions, err := UnmarshalCFSubmissionResponse(responseData);
		if err!=nil {
			log.Printf("Couldnt unmarshal the byte slice: %v", err);
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
		log.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body);
	if err != nil {
		log.Fatal(err)
	}
	responseCFContests, err := UnmarshalCFContestResponse(responseData);
	if err!=nil {
		log.Printf("Couldnt unmarshal the byte slice: %v", err);
	}

	return responseCFContests.Contests;
}



/**
* @brief Helper function to typecast the submissions response array recieved from querying the codeforces service
* @param The submissions array that is modelled like the json response from the api call
* @return The submissions array typecasted(type: []*platformDatapb.Submission) so as to be sent as a response for the unary grpc call
**/

func submissionDataConverterforDB(submissionArray []Submissions) ([]utilities.SubmissionData){
	 submissionsArrayforDB := []utilities.SubmissionData{};

	for _, submission := range submissionArray{
		submissionResponseObject := utilities.SubmissionData{
			ProblemUrl: "https://codeforces.com/contest/"+strconv.FormatInt(submission.ContestID, 10)+"/problem/"+submission.Problem.Index,
			ProblemName: submission.Problem.Name,
			SubmissionDate: strconv.FormatInt(submission.CreationTimeSeconds,10),
			SubmissionLanguage: submission.ProgrammingLanguage,
			SubmissionStatus: submission.Verdict,
			CodeUrl: "https://codeforces.com/contest/"+strconv.FormatInt(submission.ContestID, 10)+"/submission/"+strconv.FormatInt(submission.ID, 10),
		}
		submissionsArrayforDB = append(submissionsArrayforDB, submissionResponseObject);
	}
		return submissionsArrayforDB;
}


/**
* @brief Helper function to typecast the contests response array recieved from querying the codeforces service
* @param The contests array that is modelled like the json response from the api call
* @return The contests array typecasted(type: []*platformDatapb.Contest) so as to be sent as a response for the unary grpc call
**/


func contestDataConverterforDB(contestArray []Contests) ([]utilities.ContestData){
		contestsResponseforDB := []utilities.ContestData{};

		contestList := fetchAllContests();

	for _, contest := range contestArray{
		contestResponseObject := utilities.ContestData{
			ContestName: contest.ContestName,
			ContestDate: findContestAndReturnDate(contestList,contest.ContestID),
			Rank: float64(contest.Rank),
			Rating: float64(contest.NewRating),
			ContestID: strconv.FormatInt(contest.ContestID, 10),
		}
		contestsResponseforDB = append(contestsResponseforDB, contestResponseObject);
	}

	return contestsResponseforDB;
}

/**
* @brief Helper function to find the list of all contests
* @param None
* @return The list of all contests
**/

func fetchAllContests() []ContestResponse {
	queryString := "https://codeforces.com/api/contest.list";
	response, err := http.Get(queryString);

	if err != nil {
		log.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body);
	if err != nil {
		log.Fatal(err)
	}
	responseCFContestList, err := UnmarshalCFContestListResponse(responseData);
	if err!=nil {
		log.Printf("Couldnt unmarshal the byte slice: %v", err);
	}
	allContests := responseCFContestList.AllContests;

	sort.Slice(allContests,func(i,j int) bool{
		return allContests[i].StartTimeSeconds < allContests[j].StartTimeSeconds;
	})

	return allContests;

}


/**
* @brief Helper function to find the date of a contest
* @param The list of all contests and the contest id of the contest to find the date of
* @return The date of the contest
**/

func findContestAndReturnDate(AllContests []ContestResponse, contestID int64) string {
	//Binary Search, the function only returns lower bounds or upper bounds, to check for the existence of an element, use == separately.
	indexOfContest := sort.Search(len(AllContests),func(index int) bool {return AllContests[index].ID>=contestID});
	
	//Linear Search incase the Binary Search doesnt work
	// for i := 0; i<len(AllContests);i++{
	// 	if(contestID==AllContests[i].ID){
	// 		indexOfContest = i;
	// 		break;
	// 	}
			
	// }

	if(indexOfContest==len(AllContests)){
		log.Fatalln("Couldnt find the contest");
	}
	return strconv.FormatInt(AllContests[indexOfContest].StartTimeSeconds, 10);
	
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
	fetchAllContests();
	platformDatapb.RegisterFetchPlatformDataServer(s, &server{});
	log.Println("Server started on port 5003");
	if err = s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err);
	} 

}