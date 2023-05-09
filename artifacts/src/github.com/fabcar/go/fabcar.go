package main

import (
	//"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	//"time"

	//"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric/common/flogging"
)

type SmartContract struct {
	contractapi.Contract
}

var logger = flogging.MustGetLogger("fabcar_cc")


type Claim struct {
	FHIRID     string    `json:"fhir_id"`
	HospitalID string    `json:"hospital_id"`
	PatientID  string    `json:"patient_id"`
	InsurerID  string    `json:"insurer_id"`
	Status     bool   `json:"status"`
}




func (s *SmartContract) CreateClaim(ctx contractapi.TransactionContextInterface, claimData string) (string, error) {

	if len(claimData) == 0 {
		return "", fmt.Errorf("Please pass the correct car data")
	}

	var claim Claim
	err := json.Unmarshal([]byte(claimData), &claim)
	if err != nil {
		return "", fmt.Errorf("Failed while unmarshling car. %s", err.Error())
	}

	claimBytes, err := json.Marshal(claim)
	if err != nil {
		return "", fmt.Errorf("Failed while marshling car. %s", err.Error())
	}

	ctx.GetStub().SetEvent("CreateAsset", claimBytes)
	id := "claim/"+ claim.FHIRID
	return ctx.GetStub().GetTxID(), ctx.GetStub().PutState(id, claimBytes)
}

func (s *SmartContract) GetclaimDataById(ctx contractapi.TransactionContextInterface, claimId string) (*Claim, error) {
	if len(claimId) == 0 {
		return nil, fmt.Errorf("Please provide correct contract Id")
		// return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	claimBytes, err := ctx.GetStub().GetState(claimId)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if claimBytes == nil {
		return nil, fmt.Errorf("%s does not exist", claimId)
	}

	claim := new(Claim)
	_ = json.Unmarshal(claimBytes, claim)

	return claim, nil

}


func (s *SmartContract) UpdateClaimStatus(ctx contractapi.TransactionContextInterface, claimID string, status string) (string, error) {

	if len(claimID) == 0 {
		return "", fmt.Errorf("Please pass the correct car id")
	}

	claimBytes, err := ctx.GetStub().GetState(claimID)

	if err != nil {
		return "", fmt.Errorf("Failed to get car data. %s", err.Error())
	}

	if claimBytes == nil {
		return "", fmt.Errorf("%s does not exist", claimID)
	}

	claim := new(Claim)
	_ = json.Unmarshal(claimBytes, claim)
	bool1, _ := strconv.ParseBool(status)


	claim.Status = bool1

	claimBytes, err = json.Marshal(claim)
	if err != nil {
		return "", fmt.Errorf("Failed while marshling car. %s", err.Error())
	}
	id := "claim/"+claim.FHIRID

	//  txId := ctx.GetStub().GetTxID()

	return ctx.GetStub().GetTxID(), ctx.GetStub().PutState(id, claimBytes)

}







func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))
	if err != nil {
		fmt.Printf("Error create fabcar chaincode: %s", err.Error())
		return
	}
	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting chaincodes: %s", err.Error())
	}

}
