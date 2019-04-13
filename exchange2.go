package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type SimpleChaincode struct {
}

type exchangeWR struct {
	ObjectType    string  `json:"doctype"`
	User          string  `json:"user"`
	OriginalWR    int     `json:"originalwr"`
	VariableWR    int     `json:"variablewr"`
	FlagOfBS      int     `json:"flagofbs"`
	TransVolume   int     `json:"transvolume"`
	TransPrice    float64 `json:"transprice"`
	TransDeadLine int     `json:"transdeadline"`
	TransType     int     `json:"transtype"`
}

// =============================
// Main
// =============================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init initializes chaincode
// ===========================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "listingTrans()" {
		return t.listingTrans()(stub, args)
	} else if function == "queryTrans" {
		return t.queryTrans(stub, args)
	}

	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}

func (t *SimpleChaincode) listingTrans(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error

	//   0       1       2      3     4         5           6     7     
	// "asdf",  "1" ,   "2" ,  "3",  "4",  "0.3212310",    "6",  "7" 
	if len(args) != 8 {
		return shim.Error("Incorrect number of arguments. Expecting 8")
	}
	// ==== Input  ====
	fmt.Println("- start init ")
	if len(args[0]) <= 0 {
		return shim.Error("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return shim.Error("2nd argument must be int")
	}
	if len(args[2]) <= 0 {
		return shim.Error("3rd argument must be int")
	}
	if len(args[3]) <= 0 {
		return shim.Error("4th argument must be int")
	}
	if len(args[4]) <= 0 {
		return shim.Error("5rd argument must be int")
	}
	if len(args[5]) <= 0 {
		return shim.Error("6th argument must be float64")
	}
	if len(args[6]) <= 0 {
		return shim.Error("7nd argument must be int")
	}
	if len(args[7]) <= 0 {
		return shim.Error("8rd argument must be int")
	}

	user := args[0]
	originalwr, err := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("2rd argument must be a numeric string")
	}
	variablewr, err := strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("3rd argument must be a numeric string")
	}
	flagofbs, err := strconv.Atoi(args[3])
	if err != nil {
		return shim.Error("4rd argument must be a numeric string")
	}
	transvolume, err := strconv.Atoi(args[4])
	if err != nil {
		return shim.Error("5rd argument must be a numeric string")
	}
	transprice, err := strconv.ParseFloat(args[5], 64)
	if err != nil {
		return shim.Error("6rd argument must be a numeric(float64) string")
	}
	transdeadline, err := strconv.Atoi(args[6])
	if err != nil {
		return shim.Error("7rd argument must be a numeric string")
	}
	transtype, err := strconv.Atoi(args[7])
	if err != nil {
		return shim.Error("8rd argument must be a numeric string")
	}

	objectType := "exchangeWR"
	exchangeWR := &exchangeWR{objectType, user, originalwr, variablewr, flagofbs, transvolume, transprice, transdeadline, transtype}
	exchangeWRJSONasBytes, err := json.Marshal(exchangeWR)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(user, exchangeWRJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Return success ====
	fmt.Println("- end init ")
	return shim.Success(nil)
}

func (t *SimpleChaincode) queryTrans(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var username, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name query")
	}

	username = args[0]
	valAsbytes, err := stub.GetState(username)

	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + username + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Exchange does not exist: " + username + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)
}


