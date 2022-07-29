package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// type responseObject struct{
//     Status string
//     Result []result
// }

// type result struct{
//     Id int
//     ContestId int
//     CreationTimeSeconds string
//     RelativeTimeSeconds int
//     Problem problem
//     Author author
//     ProgrammingLanguage string
//     Verdict string
//     Testset string
//     PassedTestCount int
//     TimeConsumedMillis int
//     MemoryConsumedBytes int64
// }

// type problem struct{
//     ContestId int
//     Index string
//     Name string
//     Type string
//     Points int
//     Ratings int
//     Tags []string
// }

// type author struct{
//     ContestId int
//     Members []member
//     ParticipantType string
//     Ghost bool
//     StartTimeSeconds int
// }

// type member struct{
//     Handle string
// }


func UnmarshalCFResponse(data []byte) (CFResponse, error) {
	var r CFResponse
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *CFResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type CFResponse struct {
	Status string   `json:"status"`
	Result []Result `json:"result"`
}

type Result struct {
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



func main(){
	response, err := http.Get("https://codeforces.com/api/user.status?handle=zeus_codes")

    if err != nil {
        fmt.Print(err.Error())
        os.Exit(1)
    }

    responseData, err := ioutil.ReadAll(response.Body)
    if err != nil {
        log.Fatal(err)
    }
    responseCF, err := UnmarshalCFResponse(responseData);
    if err!=nil {
        fmt.Printf("Couldnt unmarshal the byte slice: %v", err);
    }
    fmt.Println(len(responseCF.Result));

}