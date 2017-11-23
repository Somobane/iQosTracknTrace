package main

import (

	"errors"
	"fmt"
	"time"
	"strconv"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	
)



// TnT is a high level smart contract that collaborate together business artifact based smart contracts

type TnT struct {

}


//==============================================================================================================================

//	 Participant types - Each participant type is mapped to an integer which we use to compare to the value stored in a

//						 user's eCert and Specific Assembly Statuses

//==============================================================================================================================

const   ASSEMBLYLINE_ROLE  		=	"assemblyline_role"
const   PACKAGELINE_ROLE   		=	"packageline_role"
const   ASSEMBLYSTATUS_RFP   	=	"6" //Ready For Packaging"
const  	ASSEMBLYSTATUS_PKG 		=	"7" //Packaged" 
const  	ASSEMBLYSTATUS_CAN 		=	"8" //Cancelled"
const  	ASSEMBLYSTATUS_QAF 		=	"2" //QA Failed"
const   FIL_BATCH  				=	"FilamentBatchId"
const   LED_BATCH  				=	"LedBatchId"
const   CIR_BATCH  				=	"CircuitBoardBatchId"
const   WRE_BATCH  				=	"WireBatchId"
const   CAS_BATCH  				=	"CasingBatchId"
const   ADP_BATCH  				=	"AdaptorBatchId"
const   STK_BATCH  				=	"StickPodBatchId"

// Assembly Line Structure

type AssemblyLine struct{	

	AssemblyId string `json:"assemblyId"`
	DeviceSerialNo string `json:"deviceSerialNo"`
	DeviceType string `json:"deviceType"`
	FilamentBatchId string `json:"filamentBatchId"`
	LedBatchId string `json:"ledBatchId"`
	CircuitBoardBatchId string `json:"circuitBoardBatchId"`
	WireBatchId string `json:"wireBatchId"`
	CasingBatchId string `json:"casingBatchId"`
	AdaptorBatchId string `json:"adaptorBatchId"`
	StickPodBatchId string `json:"stickPodBatchId"`
	ManufacturingPlant string `json:"manufacturingPlant"`
	AssemblyStatus string `json:"assemblyStatus"`
	AssemblyDate string `json:"assemblyDate"` // New
	AssemblyCreationDate string `json:"assemblyCreationDate"`
	AssemblyLastUpdatedOn string `json:"assemblyLastUpdateOn"`
	AssemblyCreatedBy string `json:"assemblyCreatedBy"`
	AssemblyLastUpdatedBy string `json:"assemblyLastUpdatedBy"`
	AssemblyPackage string `json:"assemblyPackage"`
	AssemblyInfo1 string `json:"assemblyInfo1"`
	AssemblyInfo2 string `json:"assemblyInfo2"`

	//_assemblyPackage,_assemblyInfo1,_assemblyInfo2

	}


//AssemblyID Holder

type AssemblyID_Holder struct {
	AssemblyIDs 	[]string `json:"assemblyIDs"`
}


//AssemblyLine Holder

type AssemblyLine_Holder struct {

	AssemblyLines 	[]AssemblyLine `json:"assemblyLines"`

}

// Package Line Structure

type PackageLine struct{	

	CaseId string `json:"caseId"`
	HolderAssemblyId string `json:"holderAssemblyId"`
	ChargerAssemblyId string `json:"chargerAssemblyId"`
	PackageStatus string `json:"packageStatus"`
	PackagingDate string `json:"packagingDate"`
	ShippingToAddress string `json:"shippingToAddress"`
	PackageCreationDate string `json:"packageCreationDate"`
	PackageLastUpdatedOn string `json:"packageLastUpdateOn"`
	PackageCreatedBy string `json:"packageCreatedBy"`
	PackageLastUpdatedBy string `json:"packageLastUpdatedBy"`

	}

type PackageCaseID_Holder struct {
		PackageCaseIDs 	[]string `json:"packageCaseIDs"`
	
	}
	
	//PackageLine Holder
type PackageLine_Holder struct {
	
		PackageLines 	[]PackageLine `json:"packageLines"`
	
	}
/* Assembly Section */


//API to create an assembly

//"args": [ "ASM0101","DEV0101","HOLDER","FIL0002","LED0002","CIR0002","WIR0002","CAS0002","ADA0002","STK0002","MAN0002","1","20170608","aluser1"]

//_assemblyId,_deviceSerialNo,_deviceType,_filamentBatchId,_ledBatchId,_circuitBoardBatchId,_wireBatchId,_casingBatchId,_adaptorBatchId,_stickPodBatchId,_manufacturingPlant,_assemblyStatus _assemblyDate,_assemblyPackage,_assemblyInfo1,_assemblyInfo2 ,user_name

func (t *TnT) createAssembly(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 17 {

			return nil, fmt.Errorf("Incorrect number of arguments. Expecting 17. Got: %d.", len(args))
		}

	/* Access check -------------------------------------------- Starts*/
	user_name := args[16]
	if len(user_name) == 0 { return nil, errors.New("User name supplied as empty")} 
			
	if len(user_name) > 0 {

		ecert_role, err := t.get_ecert(stub, user_name)
		if err != nil {return nil, errors.New("userrole couldn't be retrieved")}
				if ecert_role == nil {return nil, errors.New("username not defined")}
	
		user_role := string(ecert_role)
		if user_role != ASSEMBLYLINE_ROLE {
			return nil, errors.New("Permission denied not AssemblyLine Role")
			
		}

	}

	/* Access check -------------------------------------------- Ends*/

		_assemblyId := args[0]
		_deviceSerialNo:= args[1]
		_deviceType:=args[2]
		_filamentBatchId:=args[3]
		_ledBatchId:=args[4]
		_circuitBoardBatchId:=args[5]
		_wireBatchId:=args[6]
		_casingBatchId:=args[7]
		_adaptorBatchId:=args[8]
		_stickPodBatchId:=args[9]
		_manufacturingPlant:=args[10]
		_assemblyStatus:= args[11]
		_assemblyDate:= args[12]
		_assemblyPackage:= args[13]
		_assemblyInfo1:= args[14]
		_assemblyInfo2:= args[15]
		_time:= time.Now().Local()
		_assemblyCreationDate := _time.Format("20060102150405")
		_assemblyLastUpdatedOn := _time.Format("20060102150405")
		_assemblyCreatedBy := user_name
		_assemblyLastUpdatedBy := user_name

	//Check Date
	if len(_assemblyDate) != 14 {return nil, errors.New("AssemblyDate must be 14 digit datetime field.")}	
		
	//Checking if the Assembly already exists
		assemblyAsBytes, err := stub.GetState(_assemblyId)
		if err != nil { return nil, errors.New("Failed to get assembly Id") }
		if assemblyAsBytes != nil { return nil, errors.New("Assembly already exists") }
		
		/* AssemblyLine history -----------------Starts */

		var assemLine_HolderInit AssemblyLine_Holder
		assemLine_HolderKey := _assemblyId + "H" // Indicates history key
		bytesAssemblyLinesInit, err := json.Marshal(assemLine_HolderInit)
		if err != nil { return nil, errors.New("Error creating assemID_Holder record") }
		err = stub.PutState(assemLine_HolderKey, bytesAssemblyLinesInit)

		/* AssemblyLine history -----------------Ends */
		//setting the AssemblyLine to create
		assem := AssemblyLine{}
		assem.AssemblyId = _assemblyId
		assem.DeviceSerialNo = _deviceSerialNo
		assem.DeviceType = _deviceType
		assem.FilamentBatchId = _filamentBatchId
		assem.LedBatchId = _ledBatchId
		assem.CircuitBoardBatchId = _circuitBoardBatchId
		assem.WireBatchId = _wireBatchId
		assem.CasingBatchId = _casingBatchId
		assem.AdaptorBatchId = _adaptorBatchId
		assem.StickPodBatchId = _stickPodBatchId
		assem.ManufacturingPlant = _manufacturingPlant
		assem.AssemblyStatus = _assemblyStatus
		assem.AssemblyDate = _assemblyDate
		assem.AssemblyCreationDate = _assemblyCreationDate
		assem.AssemblyLastUpdatedOn = _assemblyLastUpdatedOn
		assem.AssemblyCreatedBy = _assemblyCreatedBy
		assem.AssemblyLastUpdatedBy = _assemblyLastUpdatedBy
		assem.AssemblyPackage = _assemblyPackage
		assem.AssemblyInfo1 = _assemblyInfo1
		assem.AssemblyInfo2 = _assemblyInfo2

		

		bytes, err := json.Marshal(assem)
		if err != nil { fmt.Printf("SAVE_CHANGES: Error converting Assembly record: %s", err); 
			return nil, errors.New("Error converting Assembly record")
			}

		err = stub.PutState(_assemblyId, bytes)
		if err != nil { fmt.Printf("SAVE_CHANGES: Error storing Assembly record: %s", err); 
			return nil, errors.New("Error storing Assembly record") 
		}

		/* GetAll changes-------------------------starts--------------------------*/
		// Holding the AssemblyIDs in State separately
		bytesAssemHolder, err := stub.GetState("Assemblies")
		if err != nil { return nil, errors.New("Unable to get Assemblies") }
		var assemID_Holder AssemblyID_Holder
		err = json.Unmarshal(bytesAssemHolder, &assemID_Holder)
		if err != nil {	return nil, errors.New("Corrupt Assemblies record") }
		assemID_Holder.AssemblyIDs = append(assemID_Holder.AssemblyIDs, _assemblyId)


		bytesAssemHolder, err = json.Marshal(assemID_Holder)
		if err != nil { return nil, errors.New("Error creating Assembly_Holder record") }

		err = stub.PutState("Assemblies", bytesAssemHolder)
		if err != nil { return nil, errors.New("Unable to put the state") }
		
		/* GetAll changes---------------------------ends------------------------ */

		/* AssemblyLine history ------------------------------------------Starts */

		bytesAssemblyLines, err := stub.GetState(assemLine_HolderKey)
		if err != nil { return nil, errors.New("Unable to get Assemblies") }
		
		var assemLine_Holder AssemblyLine_Holder
		err = json.Unmarshal(bytesAssemblyLines, &assemLine_Holder)
		if err != nil {	return nil, errors.New("Corrupt AssemblyLines record") }
		assemLine_Holder.AssemblyLines = append(assemLine_Holder.AssemblyLines, assem) //appending the newly created AssemblyLine
		
		bytesAssemblyLines, err = json.Marshal(assemLine_Holder)
		if err != nil { return nil, errors.New("Error creating AssemblyLine_Holder record") }
		

		err = stub.PutState(assemLine_HolderKey, bytesAssemblyLines)
		if err != nil { return nil, errors.New("Unable to put the state") }
		

		/* AssemblyLine history ------------------------------------------Ends */

		//fmt.Println("Created Assembly successfully")
		fmt.Printf("Created Assembly successfully")

		
		return nil, nil
		

}



//Update Assembly based on Id - All except AssemblyId, DeviceSerialNo,DeviceType and AssemblyCreationDate and AssemblyCreatedBy
//"args": [ "ASM0101","DEV0101","HOLDER","FIL0002","LED0002","CIR0002","WIR0002","CAS0002","ADA0002","STK0002","MAN0002","1","20170608","CASE0001","INFO1","INFO2"aluser1"]
//_assemblyId,_deviceSerialNo,_deviceType,_filamentBatchId,_ledBatchId,_circuitBoardBatchId,_wireBatchId,_casingBatchId,_adaptorBatchId,_stickPodBatchId,_manufacturingPlant,_assemblyStatus _assemblyDate,_assemblyPackage,_assemblyInfo1,_assemblyInfo2 ,user_name

func (t *TnT) updateAssemblyByID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 17 {
		return nil, errors.New("Incorrect number of arguments. Expecting 17.")
		
	} 

	/* Access check -------------------------------------------- Starts*/
	user_name := args[16]
	if len(user_name) == 0 { return nil, errors.New("User name supplied as empty") }
	
	if len(user_name) > 0 {
		ecert_role, err := t.get_ecert(stub, user_name)
		if err != nil {return nil, errors.New("userrole couldn't be retrieved")}
		if ecert_role == nil {return nil, errors.New("username not defined")}	
		
	}

	/* Access check -------------------------------------------- Ends*/

		_assemblyId := args[0]
		_deviceSerialNo:= args[1]
		//_deviceType:=args[2] - No Change
		_filamentBatchId:=args[3]
		_ledBatchId:=args[4]
		_circuitBoardBatchId:=args[5]
		_wireBatchId:=args[6]
		_casingBatchId:=args[7]
		_adaptorBatchId:=args[8]
		_stickPodBatchId:=args[9]
		_manufacturingPlant:=args[10]
		_assemblyStatus:= args[11]
		_assemblyDate:= args[12] 
		_assemblyPackage:= args[13]
		_assemblyInfo1:= args[14]
		_assemblyInfo2:= args[15]
		_time:= time.Now().Local()
		//_assemblyCreationDate - No change
		_assemblyLastUpdatedOn := _time.Format("20060102150405")
		//_assemblyCreatedBy - No change
		_assemblyLastUpdatedBy := user_name

		//Check Date
		if len(_assemblyDate) != 14 {return nil, errors.New("AssemblyDate must be 14 digit datetime field.")}	
		
		//get the Assembly

		assemblyAsBytes, err := stub.GetState(_assemblyId)
		if err != nil {	return nil, errors.New("Failed to get assembly Id")	}
		if assemblyAsBytes == nil { return nil, errors.New("Assembly doesn't exists") }
		

		assem := AssemblyLine{}
		json.Unmarshal(assemblyAsBytes, &assem)

		//update the AssemblyLine 
		//assem.AssemblyId = _assemblyId
		assem.DeviceSerialNo = _deviceSerialNo
		//assem.DeviceType = _deviceType
		assem.FilamentBatchId = _filamentBatchId
		assem.LedBatchId = _ledBatchId
		assem.CircuitBoardBatchId = _circuitBoardBatchId
		assem.WireBatchId = _wireBatchId
		assem.CasingBatchId = _casingBatchId
		assem.AdaptorBatchId = _adaptorBatchId
		assem.StickPodBatchId = _stickPodBatchId
		assem.ManufacturingPlant = _manufacturingPlant
		assem.AssemblyStatus = _assemblyStatus
		assem.AssemblyDate = _assemblyDate
		//assem.AssemblyCreationDate = _assemblyCreationDate
		assem.AssemblyLastUpdatedOn = _assemblyLastUpdatedOn
		//assem.AssemblyCreatedBy = _assemblyCreatedBy
		assem.AssemblyLastUpdatedBy = _assemblyLastUpdatedBy
		assem.AssemblyPackage = _assemblyPackage
		assem.AssemblyInfo1 = _assemblyInfo1
		assem.AssemblyInfo2 = _assemblyInfo2


		bytes, err := json.Marshal(assem)
		if err != nil { fmt.Printf("SAVE_CHANGES: Error converting Assembly record: %s", err); 
			return nil, errors.New("Error converting Assembly record") }
			

		err = stub.PutState(_assemblyId, bytes)
		if err != nil { fmt.Printf("SAVE_CHANGES: Error storing Assembly record: %s", err); 
			return nil, errors.New("Error storing Assembly record") }
			 

		/* AssemblyLine history ------------------------------------------Starts */

		// assemLine_HolderKey := _assemblyId + "H" // Indicates history key
		assemLine_HolderKey := _assemblyId + "H" // Indicates History Key for Assembly with ID = _assemblyId
		bytesAssemblyLines, err := stub.GetState(assemLine_HolderKey)
		if err != nil { return nil, errors.New("Unable to get Assemblies")}
		

		var assemLine_Holder AssemblyLine_Holder
		err = json.Unmarshal(bytesAssemblyLines, &assemLine_Holder)
		if err != nil {	return nil, errors.New("Corrupt AssemblyLines record") }
		

		assemLine_Holder.AssemblyLines = append(assemLine_Holder.AssemblyLines, assem) //appending the updated AssemblyLine
		bytesAssemblyLines, err = json.Marshal(assemLine_Holder)

		if err != nil { return nil, errors.New("Error creating AssemblyLine_Holder record") }
		err = stub.PutState(assemLine_HolderKey, bytesAssemblyLines)

		if err != nil { return nil, errors.New("Unable to put the state") }
		
		/* AssemblyLine history ------------------------------------------Ends */

		return nil, nil	
	

}


//Update Assembly based on Id - AssemblyStatus

func (t *TnT) updateAssemblyStatusByID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2.")
		
	} 

	/* Access check -------------------------------------------- Starts*/
	user_name := args[2]
	if len(user_name) == 0 { return nil, errors.New("User name supplied as empty") }
	
	if len(user_name) > 0 {
		ecert_role, err := t.get_ecert(stub, user_name)
		if err != nil {return nil, errors.New("userrole couldn't be retrieved")}
		if ecert_role == nil {return nil, errors.New("username not defined")}	
			
	}

	/* Access check -------------------------------------------- Ends*/

		_assemblyId := args[0]
		_assemblyStatus:= args[1]
		_time:= time.Now().Local()
		_assemblyLastUpdatedOn := _time.Format("20060102150405")
		_assemblyLastUpdatedBy := ""

		//get the Assembly
		assemblyAsBytes, err := stub.GetState(_assemblyId)
		if err != nil {	return nil, errors.New("Failed to get assembly Id")	}
		if assemblyAsBytes == nil { return nil, errors.New("Assembly doesn't exists") }
		

		assem := AssemblyLine{}
		json.Unmarshal(assemblyAsBytes, &assem)
		//update the AssemblyLine status

		assem.AssemblyStatus = _assemblyStatus
		assem.AssemblyLastUpdatedOn = _assemblyLastUpdatedOn
		assem.AssemblyLastUpdatedBy = _assemblyLastUpdatedBy
		

		bytes, err := json.Marshal(assem)

		if err != nil { fmt.Printf("SAVE_CHANGES: Error converting Assembly record: %s", err); 
			return nil, errors.New("Error converting Assembly record") }
			

		err = stub.PutState(_assemblyId, bytes)
		if err != nil { fmt.Printf("SAVE_CHANGES: Error storing Assembly record: %s", err); 
			return nil, errors.New("Error storing Assembly record") }
			

		/* AssemblyLine history ------------------------------------------Starts */
		// assemLine_HolderKey := _assemblyId + "H" // Indicates history key
		assemLine_HolderKey := _assemblyId + "H" // Indicates History Key for Assembly with ID = _assemblyId
		bytesAssemblyLines, err := stub.GetState(assemLine_HolderKey)

		if err != nil { return nil, errors.New("Unable to get Assemblies") }

		var assemLine_Holder AssemblyLine_Holder
		err = json.Unmarshal(bytesAssemblyLines, &assemLine_Holder)

		if err != nil {	return nil, errors.New("Corrupt AssemblyLines record") }
		assemLine_Holder.AssemblyLines = append(assemLine_Holder.AssemblyLines, assem) //appending the updated AssemblyLine

		bytesAssemblyLines, err = json.Marshal(assemLine_Holder)
		if err != nil { return nil, errors.New("Error creating AssemblyLine_Holder record") }
		
		err = stub.PutState(assemLine_HolderKey, bytesAssemblyLines)
		if err != nil { return nil, errors.New("Unable to put the state") }
		/* AssemblyLine history ------------------------------------------Ends */

		return nil, nil
		
		
}


//get the Assembly against ID
func (t *TnT) getAssemblyByID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 1 {

		return nil, errors.New("Incorrect number of arguments. Expecting AssemblyID to query")
		
	}

	_assemblyId := args[0]
	//get the var from chaincode state

	valAsbytes, err := stub.GetState(_assemblyId)									
	if err != nil {

		jsonResp := "{\"Error\":\"Failed to get state for " +  _assemblyId  + "\"}"
		return nil, errors.New(jsonResp)
		
	}

	return valAsbytes, nil
	
}


//get all Assemblies

func (t *TnT) getAllAssemblies(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	/* Access check -------------------------------------------- Starts*/
	if len(args) != 1 {

			return nil, errors.New("Incorrect number of arguments. Expecting 1.")
		}

	user_name := args[0]
	if len(user_name) == 0 { return nil, errors.New("User name supplied as empty") }
	
	if len(user_name) > 0 {

		ecert_role, err := t.get_ecert(stub, user_name)
		if err != nil {return nil, errors.New("userrole couldn't be retrieved")}
		if ecert_role == nil {return nil, errors.New("username not defined")}
		
		user_role := string(ecert_role)
		if user_role != ASSEMBLYLINE_ROLE {

			return nil, errors.New("Permission denied not AssemblyLine Role")
			
		}

	}

	/* Access check -------------------------------------------- Ends*/

	bytes, err := stub.GetState("Assemblies")
	if err != nil { return nil, errors.New("Unable to get Assemblies") }
	var assemID_Holder AssemblyID_Holder

	err = json.Unmarshal(bytes, &assemID_Holder)
	if err != nil {	return nil, errors.New("Corrupt Assemblies") }
	
	res2E:= []*AssemblyLine{}	

	for _, assemblyId := range assemID_Holder.AssemblyIDs {

		//Get the existing AssemblyLine

		assemblyAsBytes, err := stub.GetState(assemblyId)
		if err != nil { return nil, errors.New("Failed to get Assembly")}
		
		if assemblyAsBytes != nil 	{ 
				res := new(AssemblyLine)
				json.Unmarshal(assemblyAsBytes, &res)

				// Append Assembly to Assembly Array
				res2E=append(res2E,res)

			} // If ends
		} // For ends

    mapB, _ := json.Marshal(res2E)

    //fmt.Println(string(mapB))
	return mapB, nil
	
}

//get all Assemblies based on Type & BatchNo

func (t *TnT) getAssembliesByBatchNumber(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	
	/* Access check -------------------------------------------- Starts*/
	if len(args) != 3 {

			return nil, errors.New("Incorrect number of arguments. Expecting 3.")
		}

	user_name := args[2]
	if len(user_name) == 0 { return nil, errors.New("User name supplied as empty") }
	
	if len(user_name) > 0 {

		ecert_role, err := t.get_ecert(stub, user_name)
		if err != nil {return nil, errors.New("userrole couldn't be retrieved")}
		if ecert_role == nil {return nil, errors.New("username not defined")}
		
		user_role := string(ecert_role)
		if user_role != ASSEMBLYLINE_ROLE {
			return nil, errors.New("Permission denied not AssemblyLine Role")
		}

	}

	/* Access check -------------------------------------------- Ends*/

	_batchType:= args[0]
	_batchNumber:= args[1]
	_assemblyFlag:= 0

	bytes, err := stub.GetState("Assemblies")
	if err != nil { return nil, errors.New("Unable to get Assemblies") }
	var assemID_Holder AssemblyID_Holder
	err = json.Unmarshal(bytes, &assemID_Holder)

	if err != nil {	return nil, errors.New("Corrupt Assemblies") }
	
	res2E:= []*AssemblyLine{}	
	for _, assemblyId := range assemID_Holder.AssemblyIDs {

		//Get the existing AssemblyLine
		assemblyAsBytes, err := stub.GetState(assemblyId)
		if err != nil { return nil, errors.New("Failed to get Assembly")}
		if assemblyAsBytes != nil { 
			res := new(AssemblyLine)
			json.Unmarshal(assemblyAsBytes, &res)

			//Check the filter condition

			if 	_batchType == FIL_BATCH &&
						res.FilamentBatchId == _batchNumber		{ 
						_assemblyFlag = 1
			} else if  _batchType == LED_BATCH	&&
						res.LedBatchId == _batchNumber			{ 
						_assemblyFlag = 1
			} else if  _batchType == CIR_BATCH	&&
						res.CircuitBoardBatchId == _batchNumber	{ 
						_assemblyFlag = 1
			} else if  _batchType == WRE_BATCH &&
						res.WireBatchId == _batchNumber			{ 
						_assemblyFlag = 1
			} else if  _batchType == CAS_BATCH &&
						res.CasingBatchId == _batchNumber		{ 
						_assemblyFlag = 1
			} else if  _batchType == ADP_BATCH	&&
						res.AdaptorBatchId == _batchNumber		{ 
						_assemblyFlag = 1
			} else if  _batchType == STK_BATCH	&&
						res.StickPodBatchId == _batchNumber		{ 
						_assemblyFlag = 1

			}

			// Append Assembly to Assembly Array if the flag is 1 (indicates valid for filter criteria)

			if _assemblyFlag == 1 {
				res2E=append(res2E,res)

			}

		} // If ends

		//re-setting the flag to 0
		_assemblyFlag = 0

	} // For ends


    mapB, _ := json.Marshal(res2E)

    //fmt.Println(string(mapB))
	return mapB, nil
    
}

//get all Assemblies based on FromDate & ToDate

func (t *TnT) getAssembliesByDate(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	/* Access check -------------------------------------------- Starts*/
	if len(args) != 3 {

			return nil, errors.New("Incorrect number of arguments. Expecting 3.")
			
		}

	user_name := args[2]
	if len(user_name) == 0 { return nil, errors.New("User name supplied as empty") }
	if len(user_name) > 0 {

		ecert_role, err := t.get_ecert(stub, user_name)
		if err != nil {return nil, errors.New("userrole couldn't be retrieved")}
		if ecert_role == nil {return nil, errors.New("username not defined")}
		

		user_role := string(ecert_role)
		if user_role != ASSEMBLYLINE_ROLE {

			return nil, errors.New("Permission denied not AssemblyLine Role")
		}

	}

	/* Access check -------------------------------------------- Ends*/

	// YYYYMMDDHHMMSS (e.g. 20170612235959) handled as Int64
	//var _fromDate int64
	//var _toDate int64


	_fromDate, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil { return nil, errors.New ("Error in converting FromDate to int64")}
	
	_toDate, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil { return nil, errors.New ("Error in converting ToDate to int64")}
	

	_assemblyFlag:= 0
	bytes, err := stub.GetState("Assemblies")
	if err != nil { return nil, errors.New("Unable to get Assemblies") }
	
	var assemID_Holder AssemblyID_Holder
	var _assemblyDateInt64 int64

	err = json.Unmarshal(bytes, &assemID_Holder)
	if err != nil {	return nil, errors.New("Corrupt Assemblies") }
	
	res2E:= []*AssemblyLine{}	

	for _, assemblyId := range assemID_Holder.AssemblyIDs {

		//Get the existing AssemblyLine History
		assemblyAsBytes, err := stub.GetState(assemblyId)
		if err != nil { return nil, errors.New("Failed to get Assembly")}
		if assemblyAsBytes == nil { return nil, errors.New("Failed to get Assembly")}
		
		res := new(AssemblyLine)
		json.Unmarshal(assemblyAsBytes, &res)


		//fmt.Printf("%T, %v\n", _fromDate, _fromDate)
		//fmt.Printf("%T, %v\n", _toDate, _toDate)
		//if _fromDate == _toDate { return nil, errors.New("Failed to get Assembly")}
		
		//Check the filter condition YYYYMMDDHHMMSS

		if len(res.AssemblyDate) != 14 {return nil, errors.New("AssemblyDate must be 14 digit datetime field.")}
		
		if _assemblyDateInt64, err = strconv.ParseInt(res.AssemblyDate, 10, 64); err != nil { errors.New ("Error in converting AssemblyDate to int64")}

		if	_assemblyDateInt64 >= _fromDate	&&
			_assemblyDateInt64 <= _toDate		{ 
			_assemblyFlag = 1

		} 

		// Append Assembly to Assembly Array if the flag is 1 (indicates valid for filter criteria)

		if _assemblyFlag == 1 {
			res2E=append(res2E,res)
		}

	//re-setting the flag and AssemblyDate

		_assemblyFlag = 0
		_assemblyDateInt64 = 0

	} // For ends

    mapB, _ := json.Marshal(res2E)
    //fmt.Println(string(mapB))
	return mapB, nil
	
}



//get all Assemblies History based on FromDate & ToDate

func (t *TnT) getAssembliesHistoryByDate(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	

	/* Access check -------------------------------------------- Starts*/
	if len(args) != 3 {
			return nil, errors.New("Incorrect number of arguments. Expecting 3.")
		}

	user_name := args[2]
	if len(user_name) == 0 { return nil, errors.New("User name supplied as empty") }
	
	if len(user_name) > 0 {
		ecert_role, err := t.get_ecert(stub, user_name)
		if err != nil {return nil, errors.New("userrole couldn't be retrieved")}
		if ecert_role == nil {return nil, errors.New("username not defined")}
		
		user_role := string(ecert_role)
		if user_role != ASSEMBLYLINE_ROLE {

			return nil, errors.New("Permission denied not AssemblyLine Role")
		}

	}

	/* Access check -------------------------------------------- Ends*/

	// YYYYMMDDHHMMSS (e.g. 20170612235959) handled as Int64
	//var _fromDate int64
	//var _toDate int64


	_fromDate, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil { return nil, errors.New ("Error in converting FromDate to int64")}
		
	_toDate, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil { return nil, errors.New ("Error in converting ToDate to int64")}
	
	_assemblyFlag:= 0

	bytes, err := stub.GetState("Assemblies")
	if err != nil { return nil, errors.New("Unable to get Assemblies") }
	

	var assemID_Holder AssemblyID_Holder
	var _assemblyDateInt64 int64

	
	err = json.Unmarshal(bytes, &assemID_Holder)
	if err != nil {	return nil, errors.New("Corrupt Assemblies") }
	
	// Array of filtered Assemblies
	res2E:= []AssemblyLine{}	

	// Filtered Assembly
	//res := new(AssemblyLine)

	//Looping through the array of assemblyids
	for _, assemblyId := range assemID_Holder.AssemblyIDs {

		//Get the AssemblyLine History for each AssemblyID
		assemLine_HolderKey := assemblyId + "H" // Indicates History Key for Assembly with ID = _assemblyId
		bytesAssemblyLinesHistoryByID, err := stub.GetState(assemLine_HolderKey)
		if err != nil { return nil, errors.New("Unable to get bytesAssemblyLinesHistoryByID") }
		
		var assemLineHistory_Holder AssemblyLine_Holder
		err = json.Unmarshal(bytesAssemblyLinesHistoryByID, &assemLineHistory_Holder)

		if err != nil {	return nil, errors.New("Corrupt assemLineHistory_Holder record") }
		
		//Looping through the array of assemblies
		for _, res := range assemLineHistory_Holder.AssemblyLines {

			//Check the filter condition YYYYMMDDHHMMSS
			/*

			if len(res.AssemblyDate) != 14 {return nil, errors.New("AssemblyDate must be 14 digit datetime field.")}
			if _assemblyDateInt64, err = strconv.ParseInt(res.AssemblyDate, 10, 64); err != nil { errors.New ("Error in converting AssemblyDate to int64")}

			*/
			//Skip if not a valid date YYYYMMDDHHMMSS

			if len(res.AssemblyDate) == 14 {
				if _assemblyDateInt64, err = strconv.ParseInt(res.AssemblyDate, 10, 64); err == nil { 
					if	_assemblyDateInt64 >= _fromDate	&&
						_assemblyDateInt64 <= _toDate { 
						_assemblyFlag = 1
					} 

				}

			}

			// Append Assembly to Assembly Array if the flag is 1 (indicates valid for filter criteria)

			if _assemblyFlag == 1 {
				res2E=append(res2E,res)

			}
		

			//re-setting the flag and AssemblyDate

				_assemblyFlag = 0
				_assemblyDateInt64 = 0

		} // For assemLineHistory_Holder.AssemblyLines ends

	} // For assemID_Holder.AssemblyIDs ends


    mapB, _ := json.Marshal(res2E)
    //fmt.Println(string(mapB))

	return mapB, nil
	
}



//get all Assemblies based on Type & BatchNo & From & To Date

func (t *TnT) getAssembliesByBatchNumberAndByDate(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {


	/* Access check -------------------------------------------- Starts*/

	if len(args) != 5 {
			return nil, errors.New("Incorrect number of arguments. Expecting 5.")
		}

	user_name := args[4]
	if len(user_name) == 0 { return nil, errors.New("User name supplied as empty") }
	
	if len(user_name) > 0 {
		ecert_role, err := t.get_ecert(stub, user_name)

		if err != nil {return nil, errors.New("userrole couldn't be retrieved")}
		if ecert_role == nil {return nil, errors.New("username not defined")}

		user_role := string(ecert_role)
		if user_role != ASSEMBLYLINE_ROLE {
			return nil, errors.New("Permission denied not AssemblyLine Role")
		}

	}

	/* Access check -------------------------------------------- Ends*/

	_batchType:= args[0]
	_batchNumber:= args[1]
	_assemblyFlag:= 0

	_fromDate, err := strconv.ParseInt(args[2], 10, 64)
	if err != nil { return nil, errors.New ("Error in converting FromDate to int64")}
	
	_toDate, err := strconv.ParseInt(args[3], 10, 64)
	if err != nil { return nil, errors.New ("Error in converting ToDate to int64")}
	
	bytes, err := stub.GetState("Assemblies")
	if err != nil { return nil, errors.New("Unable to get Assemblies") }
	
	var assemID_Holder AssemblyID_Holder
	var _assemblyDateInt64 int64

	err = json.Unmarshal(bytes, &assemID_Holder)
	if err != nil {	return nil, errors.New("Corrupt Assemblies") }
	
	res2E:= []*AssemblyLine{}	

	for _, assemblyId := range assemID_Holder.AssemblyIDs {

		//Get the existing AssemblyLine
		assemblyAsBytes, err := stub.GetState(assemblyId)
		if err != nil { return nil, errors.New("Failed to get Assembly")}
		if assemblyAsBytes != nil { 
			res := new(AssemblyLine)
			json.Unmarshal(assemblyAsBytes, &res)

			//Check the filter condition
			if len(res.AssemblyDate) == 14 {
				if _assemblyDateInt64, err = strconv.ParseInt(res.AssemblyDate, 10, 64); err == nil { 
					if	_assemblyDateInt64 >= _fromDate	&&
						_assemblyDateInt64 <= _toDate	{
							if  _batchType == FIL_BATCH &&
										res.FilamentBatchId == _batchNumber	{ 
										_assemblyFlag = 1
							} else if  _batchType == LED_BATCH &&
										res.LedBatchId == _batchNumber	{ 
										_assemblyFlag = 1
							} else if  _batchType == CIR_BATCH &&
										res.CircuitBoardBatchId == _batchNumber	{ 
										_assemblyFlag = 1
							} else if  _batchType == WRE_BATCH &&
										res.WireBatchId == _batchNumber			{ 
										_assemblyFlag = 1
							} else if  _batchType == CAS_BATCH &&
										res.CasingBatchId == _batchNumber		{ 
										_assemblyFlag = 1
							} else if  _batchType == ADP_BATCH &&
										res.AdaptorBatchId == _batchNumber		{ 
										_assemblyFlag = 1

							} else if  _batchType == STK_BATCH &&
										res.StickPodBatchId == _batchNumber		{ 
										_assemblyFlag = 1

							}

						}// from date and to date check

				}// if date parse

			}// if date lenght

			// Append Assembly to Assembly Array if the flag is 1 (indicates valid for filter criteria)
			if _assemblyFlag == 1 {

				res2E=append(res2E,res)

			}

		} // If ends

		//re-setting the flag to 0

		_assemblyFlag = 0

		_assemblyDateInt64 = 0

	} // For ends

    mapB, _ := json.Marshal(res2E)

    //fmt.Println(string(mapB))
	return mapB, nil
}

/* Package section*/

//API to create an Package

// Assemblies related to the package is updated with status = PACKAGED

func (t *TnT) createPackage(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 8 {
			return nil, fmt.Errorf("Incorrect number of arguments. Expecting 8. Got: %d.", len(args))
			
		}

	/* Access check -------------------------------------------- Starts*/

	user_name := args[7]
	if len(user_name) == 0 { return nil, errors.New("User name supplied as empty") }
	
	if len(user_name) > 0 {

		ecert_role, err := t.get_ecert(stub, user_name)

		if err != nil {return nil, errors.New("userrole couldn't be retrieved")}
		if ecert_role == nil {return nil, errors.New("username not defined")}

		user_role := string(ecert_role)

		if user_role != PACKAGELINE_ROLE {
			return nil, errors.New("Permission denied not PackageLine Role")
		}

	}

	/* Access check -------------------------------------------- Ends*/
	

		_caseId := args[0]
		_holderAssemblyId := args[1]
		_chargerAssemblyId := args[2]
		_packageStatus := args[3]
		_packagingDate := args[4]
		_shippingToAddress := args[5]
		// Status of associated Assemblies	
		_assemblyStatus:= args[6]
		_time:= time.Now().Local()
		_packageCreationDate := _time.Format("20060102150405")
		_packageLastUpdatedOn := _time.Format("20060102150405")
		_packageCreatedBy := user_name
		_packageLastUpdatedBy := user_name

	//Checking if the Package already exists

		packageAsBytes, err := stub.GetState(_caseId)

		if err != nil { return nil, errors.New("Failed to get Package") }
		if packageAsBytes != nil { return nil, errors.New("Package already exists") }

	

		//setting the Package to create

		pack := PackageLine{}
		pack.CaseId = _caseId
		pack.HolderAssemblyId = _holderAssemblyId
		pack.ChargerAssemblyId = _chargerAssemblyId
		pack.PackageStatus = _packageStatus
		pack.PackagingDate = _packagingDate
		pack.ShippingToAddress = _shippingToAddress
		pack.PackageCreationDate = _packageCreationDate
		pack.PackageLastUpdatedOn = _packageLastUpdatedOn
		pack.PackageCreatedBy = _packageCreatedBy
		pack.PackageLastUpdatedBy = _packageLastUpdatedBy

		bytes, err := json.Marshal(pack)

		if err != nil { fmt.Printf("SAVE_CHANGES: Error converting Package record: %s", err); return nil, errors.New("Error converting Package record") }
		
		err = stub.PutState(_caseId, bytes)

		if err != nil { fmt.Printf("SAVE_CHANGES: Error storing Package record: %s", err); return nil, errors.New("Error storing Package record") }


		/* PackageLine history -----------------Starts */

		// Initialises the PackageLine_Holder

		var packLine_HolderInit PackageLine_Holder
		packLine_HolderKey := _caseId + "H" // Indicates history key
		bytesPackLinesInit, err := json.Marshal(packLine_HolderInit)

		if err != nil { return nil, errors.New("Error creating packLine_HolderInit record") }
		
		err = stub.PutState(packLine_HolderKey, bytesPackLinesInit)

		/* PackageLine history -----------------Ends */
		/* AssemblyLine history ------------------------------------------Starts */

		//packLine_HolderKey := _caseId + "H" // Indicates history key

		bytesPackageLines, err := stub.GetState(packLine_HolderKey)
		if err != nil { return nil, errors.New("Unable to get bytesPackageLines") }
		
		var packLine_Holder PackageLine_Holder
		err = json.Unmarshal(bytesPackageLines, &packLine_Holder)

		if err != nil {	return nil, errors.New("Corrupt bytesPackageLines record") }
		
		packLine_Holder.PackageLines = append(packLine_Holder.PackageLines, pack) //appending the newly created pack
		bytesPackageLines, err = json.Marshal(packLine_Holder)

		if err != nil { return nil, errors.New("Error creating AssemblyLine_Holder record") }
		
		err = stub.PutState(packLine_HolderKey, bytesPackageLines)
		if err != nil { return nil, errors.New("Unable to put the state") }
		
		/* PackageLine history ------------------------------------------Ends */
		//fmt.Println("Created Package successfully")
		fmt.Printf("Created Package successfully")
		//Update Holder Assemblies to Packaged status

		if 	len(_holderAssemblyId) > 0	{

			//_assemblyStatus:= "PACKAGED"
			_time:= time.Now().Local()
			_assemblyLastUpdatedOn := _time.Format("20060102150405")
			_assemblyLastUpdatedBy := _packageCreatedBy
			_assemblyPackage:= _caseId // Keeping reference

			//get the Assembly
			assemblyHolderAsBytes, err := stub.GetState(_holderAssemblyId)
			if err != nil {	return nil, errors.New("Failed to get assembly Id")	}
			
			if assemblyHolderAsBytes == nil { return nil, errors.New("Assembly doesn't exists") }
			
			assemHolder := AssemblyLine{}
			json.Unmarshal(assemblyHolderAsBytes, &assemHolder)

			//update the AssemblyLine status
			assemHolder.AssemblyStatus = _assemblyStatus
			assemHolder.AssemblyLastUpdatedOn = _assemblyLastUpdatedOn
			assemHolder.AssemblyLastUpdatedBy = _assemblyLastUpdatedBy
			assemHolder.AssemblyPackage = _assemblyPackage
			bytesHolder, err := json.Marshal(assemHolder)

			if err != nil { fmt.Printf("SAVE_CHANGES: Error converting Assembly record: %s", err); return nil, errors.New("Error converting Assembly record") }
			err = stub.PutState(_holderAssemblyId, bytesHolder)

			if err != nil { fmt.Printf("SAVE_CHANGES: Error storing Assembly record: %s", err); return nil, errors.New("Error storing Assembly record") }

			

			/* AssemblyLine history ------------------------------------------Starts */

			holderAssemLine_HolderKey := _holderAssemblyId + "H" // Indicates History Key for Assembly with ID = _assemblyId
			bytesHolderAssemblyLines, err := stub.GetState(holderAssemLine_HolderKey)

			if err != nil { return nil, errors.New("Unable to get Assemblies") }
			
			var holderAssemLine_Holder AssemblyLine_Holder
			err = json.Unmarshal(bytesHolderAssemblyLines, &holderAssemLine_Holder)

			if err != nil {	return nil, errors.New("Corrupt AssemblyLines record") }
			
			holderAssemLine_Holder.AssemblyLines = append(holderAssemLine_Holder.AssemblyLines, assemHolder) //appending the updated AssemblyLine
			bytesHolderAssemblyLines, err = json.Marshal(holderAssemLine_Holder)

			if err != nil { return nil, errors.New("Error creating AssemblyLine_Holder record") }
			
			err = stub.PutState(holderAssemLine_HolderKey, bytesHolderAssemblyLines)
			if err != nil { return nil, errors.New("Unable to put the state") }
			
			/* AssemblyLine history ------------------------------------------Ends */
			}



		//Update Charger Assemblies to Packaged status

		if 	len(_chargerAssemblyId) > 0		{
			//_assemblyStatus:= "PACKAGED"
			_time:= time.Now().Local()
			_assemblyLastUpdatedOn := _time.Format("20060102150405")
			_assemblyLastUpdatedBy := _packageCreatedBy
			_assemblyPackage:= _caseId // Keeping reference

			//get the Assembly
			assemblyChargerAsBytes, err := stub.GetState(_chargerAssemblyId)

			if err != nil {	return nil, errors.New("Failed to get assembly Id")	}
			
			if assemblyChargerAsBytes == nil { return nil, errors.New("Assembly doesn't exists") }
			
			assemCharger := AssemblyLine{}
			json.Unmarshal(assemblyChargerAsBytes, &assemCharger)

			//update the AssemblyLine status
			assemCharger.AssemblyStatus = _assemblyStatus
			assemCharger.AssemblyLastUpdatedOn = _assemblyLastUpdatedOn
			assemCharger.AssemblyLastUpdatedBy = _assemblyLastUpdatedBy
			assemCharger.AssemblyPackage = _assemblyPackage

			bytesCharger, err := json.Marshal(assemCharger)

			if err != nil { fmt.Printf("SAVE_CHANGES: Error converting Assembly record: %s", err); return nil, errors.New("Error converting Assembly record") }

			err = stub.PutState(_chargerAssemblyId, bytesCharger)

			if err != nil { fmt.Printf("SAVE_CHANGES: Error storing Assembly record: %s", err); return nil, errors.New("Error storing Assembly record") }
			
			/* AssemblyLine history ------------------------------------------Starts */

			chargerAssemLine_HolderKey := _chargerAssemblyId + "H" // Indicates History Key for Assembly with ID = _assemblyId
			bytesChargerAssemblyLines, err := stub.GetState(chargerAssemLine_HolderKey)

			if err != nil { return nil, errors.New("Unable to get Assemblies") }
			
			var chargerAssemLine_Holder AssemblyLine_Holder
			err = json.Unmarshal(bytesChargerAssemblyLines, &chargerAssemLine_Holder)

			if err != nil {	return nil, errors.New("Corrupt AssemblyLines record") }
			
			chargerAssemLine_Holder.AssemblyLines = append(chargerAssemLine_Holder.AssemblyLines, assemCharger) //appending the updated AssemblyLine
			bytesChargerAssemblyLines, err = json.Marshal(chargerAssemLine_Holder)

			if err != nil { return nil, errors.New("Error creating AssemblyLine_Holder record") }
			
			err = stub.PutState(chargerAssemLine_HolderKey, bytesChargerAssemblyLines)

			if err != nil { return nil, errors.New("Unable to put the state") }
			
			/* AssemblyLine history ------------------------------------------Ends */
		}

	/* GetAll changes-------------------------starts--------------------------*/

		// Holding the PackageCaseIDs in State separately
		bytesPackageCaseHolder, err := stub.GetState("Packages")
		if err != nil { return nil, errors.New("Unable to get Packages") }
		
		var packageCaseID_Holder PackageCaseID_Holder
		err = json.Unmarshal(bytesPackageCaseHolder, &packageCaseID_Holder)

		if err != nil {	return nil, errors.New("Corrupt Packages record") }
		
		packageCaseID_Holder.PackageCaseIDs = append(packageCaseID_Holder.PackageCaseIDs, _caseId)
		bytesPackageCaseHolder, err = json.Marshal(packageCaseID_Holder)

		if err != nil { return nil, errors.New("Error creating PackageCaseID_Holder record") }
		err = stub.PutState("Packages", bytesPackageCaseHolder)
		if err != nil { return nil, errors.New("Unable to put the state") }
		

	/* GetAll changes---------------------------ends------------------------ */
		return nil, nil
		
}

//API to update an Package

// Assemblies related to the package is updated with status sent as parameter

func (t *TnT) updatePackage(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 8 {

			return nil, fmt.Errorf("Incorrect number of arguments. Expecting 8. Got: %d.", len(args))
		}

	/* Access check -------------------------------------------- Starts*/
	user_name := args[7]
	if len(user_name) == 0 { return nil, errors.New("User name supplied as empty") }

	if len(user_name) > 0 {

		ecert_role, err := t.get_ecert(stub, user_name)
		if err != nil {return nil, errors.New("userrole couldn't be retrieved")}
		if ecert_role == nil {return nil, errors.New("username not defined")}


		user_role := string(ecert_role)
		if user_role != PACKAGELINE_ROLE {

			return nil, errors.New("Permission denied not PackageLine Role")

		}

	}

	/* Access check -------------------------------------------- Ends*/

		_caseId := args[0]
		//_holderAssemblyId := args[1]
		//_chargerAssemblyId := args[2]
		_packageStatus := args[3]
		//_packagingDate := args[4]
		_shippingToAddress := args[5]
		// Status of associated Assemblies	
		_assemblyStatus := args[6]
		_time:= time.Now().Local()
		//_packageCreationDate := _time.Format("2006-01-02")
		_packageLastUpdatedOn := _time.Format("20060102150405")
		//_packageCreatedBy := ""
		_packageLastUpdatedBy := user_name

	//Checking if the Package already exists

		packageAsBytes, err := stub.GetState(_caseId)

		if err != nil { return nil, errors.New("Failed to get Package") }
		if packageAsBytes == nil { return nil, errors.New("Package doesn't exists") }
		
		//setting the Package to update
		pack := PackageLine{}
		json.Unmarshal(packageAsBytes, &pack)

		//pack.CaseId = _caseId
		//pack.HolderAssemblyId = _holderAssemblyId
		//pack.ChargerAssemblyId = _chargerAssemblyId
		pack.PackageStatus = _packageStatus
		//pack.PackagingDate = _packagingDate
		pack.ShippingToAddress = _shippingToAddress
		//pack.PackageCreationDate = _packageCreationDate
		pack.PackageLastUpdatedOn = _packageLastUpdatedOn
		//pack.PackageCreatedBy = _packageCreatedBy
		pack.PackageLastUpdatedBy = _packageLastUpdatedBy

		// Getting associate Assembly IDs
		_holderAssemblyId := pack.HolderAssemblyId
		_chargerAssemblyId := pack.ChargerAssemblyId


		bytes, err := json.Marshal(pack)

	if err != nil { fmt.Printf("SAVE_CHANGES: Error converting Package record: %s", err); return nil, errors.New("Error converting Package record") }
		
	err = stub.PutState(_caseId, bytes)

	if err != nil { fmt.Printf("SAVE_CHANGES: Error storing Package record: %s", err); return nil, errors.New("Error storing Package record") }
	

		/* PackageLine history ------------------------------------------Starts */

		packLine_HolderKey := _caseId + "H" // Indicates history key
		bytesPackageLines, err := stub.GetState(packLine_HolderKey)
		if err != nil { return nil, errors.New("Unable to get bytesPackageLines") }
		
		var packLine_Holder PackageLine_Holder
		err = json.Unmarshal(bytesPackageLines, &packLine_Holder)

		if err != nil {	return nil, errors.New("Corrupt bytesPackageLines record") }

		packLine_Holder.PackageLines = append(packLine_Holder.PackageLines, pack) //appending the newly created pack
		bytesPackageLines, err = json.Marshal(packLine_Holder)

		if err != nil { return nil, errors.New("Error creating AssemblyLine_Holder record") }
		
		err = stub.PutState(packLine_HolderKey, bytesPackageLines)

		if err != nil { return nil, errors.New("Unable to put the state") }
		
		/* PackageLine history ------------------------------------------Ends */

		//fmt.Println("Created Package successfully")
		fmt.Printf("Created Package successfully")

		//Update Holder Assemblies status

		if 	len(_holderAssemblyId) > 0	{
			//_assemblyStatus:= "PACKAGED"
			_time:= time.Now().Local()
			_assemblyLastUpdatedOn := _time.Format("20060102150405")
			_assemblyLastUpdatedBy := _packageLastUpdatedBy
			_assemblyPackage:= _caseId // Keeping reference

			//get the Assembly

			assemblyHolderAsBytes, err := stub.GetState(_holderAssemblyId)
			if err != nil {	return nil, errors.New("Failed to get assembly Id")	}
			if assemblyHolderAsBytes == nil { return nil, errors.New("Assembly doesn't exists") }

			assemHolder := AssemblyLine{}
			json.Unmarshal(assemblyHolderAsBytes, &assemHolder)

			// Don't update assembly if there is no chnage in status
			// Update only when status moves say from Packaged -> Cancelled	

			if assemHolder.AssemblyStatus != _assemblyStatus {

				//update the AssemblyLine status
				assemHolder.AssemblyStatus = _assemblyStatus
				assemHolder.AssemblyLastUpdatedOn = _assemblyLastUpdatedOn
				assemHolder.AssemblyLastUpdatedBy = _assemblyLastUpdatedBy
				assemHolder.AssemblyPackage = _assemblyPackage

				bytesHolder, err := json.Marshal(assemHolder)

				if err != nil { fmt.Printf("SAVE_CHANGES: Error converting Assembly record: %s", err); return nil, errors.New("Error converting Assembly record") }
				
				err = stub.PutState(_holderAssemblyId, bytesHolder)

				if err != nil { fmt.Printf("SAVE_CHANGES: Error storing Assembly record: %s", err); return nil, errors.New("Error storing Assembly record") }
			

				/* AssemblyLine history ------------------------------------------Starts */

				holderAssemLine_HolderKey := _holderAssemblyId + "H" // Indicates History Key for Assembly with ID = _assemblyId
				bytesHolderAssemblyLines, err := stub.GetState(holderAssemLine_HolderKey)

				if err != nil { return nil, errors.New("Unable to get Assemblies") }
				var holderAssemLine_Holder AssemblyLine_Holder
				err = json.Unmarshal(bytesHolderAssemblyLines, &holderAssemLine_Holder)

				if err != nil {	return nil, errors.New("Corrupt AssemblyLines record") }
				holderAssemLine_Holder.AssemblyLines = append(holderAssemLine_Holder.AssemblyLines, assemHolder) //appending the updated AssemblyLine
		

				bytesHolderAssemblyLines, err = json.Marshal(holderAssemLine_Holder)

				if err != nil { return nil, errors.New("Error creating AssemblyLine_Holder record") }
				
				err = stub.PutState(holderAssemLine_HolderKey, bytesHolderAssemblyLines)

				if err != nil { return nil, errors.New("Unable to put the state") }
				
				/* AssemblyLine history ------------------------------------------Ends */

			}// Change of Status ends	

		}

		//Update Charger Assemblies status

		if 	len(_chargerAssemblyId) > 0		{

			//_assemblyStatus:= "PACKAGED"
			_time:= time.Now().Local()
			_assemblyLastUpdatedOn := _time.Format("20060102150405")
			_assemblyLastUpdatedBy := _packageLastUpdatedBy
			_assemblyPackage:= _caseId // Keeping reference

			//get the Assembly
			assemblyChargerAsBytes, err := stub.GetState(_chargerAssemblyId)
			if err != nil {	return nil, errors.New("Failed to get assembly Id")	}
			
			if assemblyChargerAsBytes == nil { return nil, errors.New("Assembly doesn't exists") }
			
			assemCharger := AssemblyLine{}
			json.Unmarshal(assemblyChargerAsBytes, &assemCharger)

			// Don't update assembly if there is no chnage in status
			// Update only when status moves say from Packaged -> Cancelled	

			if assemCharger.AssemblyStatus != _assemblyStatus {

				//update the AssemblyLine status
				assemCharger.AssemblyStatus = _assemblyStatus
				assemCharger.AssemblyLastUpdatedOn = _assemblyLastUpdatedOn
				assemCharger.AssemblyLastUpdatedBy = _assemblyLastUpdatedBy
				assemCharger.AssemblyPackage = _assemblyPackage

				bytesCharger, err := json.Marshal(assemCharger)

				if err != nil { fmt.Printf("SAVE_CHANGES: Error converting Assembly record: %s", err); return nil, errors.New("Error converting Assembly record") }
			
				err = stub.PutState(_chargerAssemblyId, bytesCharger)

				if err != nil { fmt.Printf("SAVE_CHANGES: Error storing Assembly record: %s", err); return nil, errors.New("Error storing Assembly record") }
				

				/* AssemblyLine history ------------------------------------------Starts */

				chargerAssemLine_HolderKey := _chargerAssemblyId + "H" // Indicates History Key for Assembly with ID = _assemblyId
				bytesChargerAssemblyLines, err := stub.GetState(chargerAssemLine_HolderKey)

				if err != nil { return nil, errors.New("Unable to get Assemblies") }
				
				var chargerAssemLine_Holder AssemblyLine_Holder
				err = json.Unmarshal(bytesChargerAssemblyLines, &chargerAssemLine_Holder)

				if err != nil {	return nil, errors.New("Corrupt AssemblyLines record") }
				
				chargerAssemLine_Holder.AssemblyLines = append(chargerAssemLine_Holder.AssemblyLines, assemCharger) //appending the updated AssemblyLine
	
				bytesChargerAssemblyLines, err = json.Marshal(chargerAssemLine_Holder)

				if err != nil { return nil, errors.New("Error creating AssemblyLine_Holder record") }
				
				err = stub.PutState(chargerAssemLine_HolderKey, bytesChargerAssemblyLines)

				if err != nil { return nil, errors.New("Unable to put the state") }
				
				/* AssemblyLine history ------------------------------------------Ends */

			}// Check if status changes

		}

		return nil, nil
		
}


//get the Package against ID

func (t *TnT) getPackageByID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting CaseId to query")
		
	}

	_caseId := args[0]
	//get the var from chaincode state

	valAsbytes, err := stub.GetState(_caseId)									
	if err != nil {

		jsonResp := "{\"Error\":\"Failed to get state for " +  _caseId  + "\"}"
		return  nil, errors.New(jsonResp)
		
	}
	return valAsbytes, nil	
	
}


//get all Packages

func (t *TnT) getAllPackages(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {


	/* Access check -------------------------------------------- Starts*/

	if len(args) != 1 {

			return nil, errors.New("Incorrect number of arguments. Expecting 1.")
		}

	user_name := args[0]
	if len(user_name) == 0 { return nil, errors.New("User name supplied as empty") }
	if len(user_name) > 0 {

		ecert_role, err := t.get_ecert(stub, user_name)
		if err != nil {return nil, errors.New("userrole couldn't be retrieved")}
		if ecert_role == nil {return nil, errors.New("username not defined")}


		user_role := string(ecert_role)

		if user_role != PACKAGELINE_ROLE {
			return nil, errors.New("Permission denied not PackageLine Role")
			
		}

	}

	/* Access check -------------------------------------------- Ends*/

	bytesPackageCaseHolder, err := stub.GetState("Packages")

	if err != nil { return nil, errors.New("Unable to get Packages") }
	
	var packageCaseID_Holder PackageCaseID_Holder
	err = json.Unmarshal(bytesPackageCaseHolder, &packageCaseID_Holder)

	if err != nil {	return nil, errors.New("Corrupt Assemblies") }
	res2E:= []*PackageLine{}	

	for _, caseId := range packageCaseID_Holder.PackageCaseIDs {

		//Get the existing AssemblyLine
		packageAsBytes, err := stub.GetState(caseId)
		if err != nil { return nil, errors.New("Failed to get Assembly")}
		
		if packageAsBytes != nil { 
		res := new(PackageLine)
		json.Unmarshal(packageAsBytes, &res)
		// Append Assembly to Assembly Array
		res2E=append(res2E,res)

		} // If ends

		} // For ends

    mapB, _ := json.Marshal(res2E)

    //fmt.Println(string(mapB))
	return mapB, nil
	
}

//All Validators to be called before Invoke
// Validator before createAssembly invoke call

func (t *TnT) validateCreateAssembly(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 17 {

			return nil, fmt.Errorf("Incorrect number of arguments. Expecting 17. Got: %d.", len(args))
		}

	/* Access check -------------------------------------------- Starts*/

	user_name := args[16]
	if len(user_name) == 0 { return nil, errors.New("User name supplied as empty") }
	
	if len(user_name) > 0 {

		ecert_role, err := t.get_ecert(stub, user_name)

		if err != nil {return nil, errors.New("userrole couldn't be retrieved")}
		if ecert_role == nil {return nil, errors.New("username not defined")}

		user_role := string(ecert_role)
		if user_role != ASSEMBLYLINE_ROLE {

			return nil, errors.New("Permission denied, not an AssemblyLine Role")
		}

	}

	/* Access check -------------------------------------------- Ends*/

	//Checking if the Assembly already exists

	_assemblyId := args[0]
	assemblyAsBytes, err := stub.GetState(_assemblyId)

	if err != nil { return nil, errors.New("Failed to get assembly Id") }
	if assemblyAsBytes != nil { return nil, errors.New("Assembly already exists") }

	//Check Date

	_assemblyDate:= args[12]
	if len(_assemblyDate) != 14 {return nil, errors.New("AssemblyDate must be 14 digit datetime field.")}	
			
	//No validation error proceed to call Invoke command
	return nil, nil
	
}

// Validator before updateAssembly invoke call

func (t *TnT) validateUpdateAssembly(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {


	if len(args) != 17 {
			return nil, errors.New("Incorrect number of arguments. Expecting 17.")
			
		} 

	_assemblyId := args[0]
	_assemblyStatus:= args[11]

	//get the Assembly

	assemblyAsBytes, err := stub.GetState(_assemblyId)
	if err != nil {	return nil, errors.New("Failed to get assembly Id")	}
	if assemblyAsBytes == nil { return nil, errors.New("Assembly doesn't exists") }

	//Check Date
	_assemblyDate:= args[12]
	if len(_assemblyDate) != 14 {return nil, errors.New("AssemblyDate must be 14 digit datetime field.")}	
		
	assem := AssemblyLine{}
	json.Unmarshal(assemblyAsBytes, &assem)

	/* Access check -------------------------------------------- Starts*/

	user_name := args[16]

	if len(user_name) == 0 { return nil, errors.New("User name supplied as empty") }
	
	if len(user_name) > 0 {

		ecert_role, err := t.get_ecert(stub, user_name)

		if err != nil {return nil, errors.New("userrole couldn't be retrieved")}
		if ecert_role == nil {return nil, errors.New("username not defined")}
		user_role := string(ecert_role)
		if user_role != ASSEMBLYLINE_ROLE {

			return nil, errors.New("Permission denied, not an AssemblyLine Role")
		}

		// AssemblyLine can't edit an Assembly in certain statuses

		if (user_role == ASSEMBLYLINE_ROLE 	&&
		assem.AssemblyStatus == ASSEMBLYSTATUS_RFP) {

			return nil, errors.New("Permission denied for AssemblyLine Role to update Assembly if status = 'Ready For Packaging'")
			
		}
	

		// AssemblyLine user can't move an AssemblyLine from QA Failed to Ready For packaging status
		if (user_role 			== ASSEMBLYLINE_ROLE &&
		assem.AssemblyStatus 	== ASSEMBLYSTATUS_QAF &&
		_assemblyStatus		 	== ASSEMBLYSTATUS_RFP) {

			return nil, errors.New("Permission denied for updating AssemblyLine with status = 'QA Failed' to 'Ready For Packaging' status")
		}

	}

	/* Access check -------------------------------------------- Ends*/	

	//No validation error proceed to call Invoke command
	return nil, nil
	
}

// Validator before createPackage invoke call

func (t *TnT) validateCreatePackage(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

		if len(args) != 8 {

			return nil, fmt.Errorf("Incorrect number of arguments. Expecting 8. Got: %d.", len(args))
			
		}

	/* Access check -------------------------------------------- Starts*/

	user_name := args[7]
	if len(user_name) == 0 { return nil, errors.New("User name supplied as empty") }
	
	if len(user_name) > 0 {

		ecert_role, err := t.get_ecert(stub, user_name)
		if err != nil {return nil, errors.New("userrole couldn't be retrieved")}
		if ecert_role == nil {return nil, errors.New("username not defined")}

		user_role := string(ecert_role)
		if user_role != PACKAGELINE_ROLE {

			return nil, errors.New("Permission denied not PackageLine Role")
			
		}

	}

	/* Access check -------------------------------------------- Ends*/

	//Checking if the Package already exists

		_caseId := args[0]
		packageAsBytes, err := stub.GetState(_caseId)

		if err != nil { return nil, errors.New("Failed to get Package") }
		if packageAsBytes != nil { return nil, errors.New("Package already exists") }

	//No validation error proceed to call Invoke command
	return nil, nil
	
}

// Validator before updatePackagey invoke call

func (t *TnT) validateUpdatePackage(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 8 {
				return nil, fmt.Errorf("Incorrect number of arguments. Expecting 8. Got: %d.", len(args))
			}

		/* Access check -------------------------------------------- Starts*/
		user_name := args[7]

		if len(user_name) == 0 { return nil, errors.New("User name supplied as empty") }
		
		if len(user_name) > 0 {
			ecert_role, err := t.get_ecert(stub, user_name)

			if err != nil {return nil, errors.New("userrole couldn't be retrieved")}
			if ecert_role == nil {return nil, errors.New("username not defined")}

			
			user_role := string(ecert_role)

			if user_role != PACKAGELINE_ROLE {
				return nil, errors.New("Permission denied not PackageLine Role")
			}

		}

		/* Access check -------------------------------------------- Ends*/

		//Checking if the Package already exists
		_caseId := args[0]
		packageAsBytes, err := stub.GetState(_caseId)

		if err != nil { return nil, errors.New("Failed to get Package") }
		if packageAsBytes == nil { return nil, errors.New("Package doesn't exists") }
		

	//No validation error proceed to call Invoke command
	return nil, nil
	
}

//AllAssemblyIDS

//get the all Assembly IDs from AssemblyID_Holder - To Test only

func (t *TnT) getAllAssemblyIDs(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 0 {

		return nil, errors.New("Incorrect number of arguments. Expecting zero argument to query")
	}
	bytesAssemHolder, err := stub.GetState("Assemblies")
	if err != nil { return nil, errors.New("Unable to get Assemblies") }
	return bytesAssemHolder, nil	

}



//AllPackageCaseIDs

//get the all Package CaseIDs from PackageCaseID_Holder - To Test only

func (t *TnT) getAllPackageCaseIDs(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 0 {
		return nil, errors.New("Incorrect number of arguments. Expecting zero argument to query")
		
	}

	bytesPackageCaseHolder, err := stub.GetState("Packages")
	if err != nil { return nil, errors.New("Unable to get Packages") }
	return bytesPackageCaseHolder, nil	
	
}


// All AssemblyLine history

func (t *TnT) getAssemblyLineHistoryByID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 2 {

		return nil, errors.New("Incorrect number of arguments. Expecting 2 arguments to query")
		
	}

	_assemblyId := args[0]
	//_userName:= args[1]	
	assemLine_HolderKey := _assemblyId + "H" // Indicates history key
	bytesAssemLineHolder, err := stub.GetState(assemLine_HolderKey)
	
	if err != nil { return nil, errors.New("Unable to get Assemblies") }

	return bytesAssemLineHolder, nil	

}

// All PackageLine history

func (t *TnT) getPackageLineHistoryByID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 2 {

		return nil, errors.New("Incorrect number of arguments. Expecting 2 arguments to query")
	}

	_caseId := args[0]
	//_userName:= args[1]	

	packLine_HolderKey := _caseId + "H" // Indicates history key
	bytesPackLineHolder, err := stub.GetState(packLine_HolderKey)
		if err != nil { return nil, errors.New("Unable to get PackageLine history") }
		
	return bytesPackLineHolder, nil
	
}

//Security & Access



//==============================================================================================================================

//	 General Functions

//==============================================================================================================================

//	 get_ecert - Takes the name passed and calls out to the REST API for HyperLedger to retrieve the ecert

//				 for that user. Returns the ecert as retrived including html encoding.

//==============================================================================================================================

func (t *TnT) get_ecert(stub shim.ChaincodeStubInterface, name string) ([]byte, error) {
	ecert, err := stub.GetState(name)

	if err != nil { return nil, errors.New("Couldn't retrieve ecert for user " + name) }
	return ecert, nil

}

//==============================================================================================================================

//	 add_ecert - Adds a new ecert and user pair to the table of ecerts

//==============================================================================================================================


func (t *TnT) add_ecert(stub shim.ChaincodeStubInterface, name string, ecert string) ([]byte, error) {


	err := stub.PutState(name, []byte(ecert))
	if err == nil {

		return nil, errors.New("Error storing eCert for user " + name + " identity: " + ecert)
	}

	return nil, nil
	
}



/*Standard Calls*/



// Init initializes the smart contracts

func (t *TnT) Init(stub shim.ChaincodeStubInterface) pb.Response {

	// Get the args from the transaction proposal
	_,args := stub.GetFunctionAndParameters()
	/* GetAll changes-------------------------starts--------------------------*/
	var assemID_Holder AssemblyID_Holder
	bytesAssembly, err := json.Marshal(assemID_Holder)
    if err != nil {return shim.Error(fmt.Sprintf("Error creating assemID_Holder record"))}
	err = stub.PutState("Assemblies", bytesAssembly)



	var packageCaseID_Holder PackageCaseID_Holder
	bytesPackage, err := json.Marshal(packageCaseID_Holder)
	if err != nil {return shim.Error(fmt.Sprintf("Error creating packageCaseID_Holder record"))}
   	err = stub.PutState("Packages", bytesPackage)

	
	/* GetAll changes---------------------------ends------------------------ */
	// creating minimum default user and roles
	//"AssemblyLine_User1","assemblyline_role","PackageLine_User1", "packageline_role"

	for i:=0; i < len(args); i=i+2 {

		t.add_ecert(stub, args[i], args[i+1])

	}
	return shim.Success(nil)

}


// Invoke callback representing the invocation of a chaincode

func (t *TnT) Invoke(stub shim.ChaincodeStubInterface)pb.Response {

	fmt.Println("Invoke called, determining function")
	 // Extract the function and args from the transaction proposal
	 function, args := stub.GetFunctionAndParameters()
	
	// Handle different functions
	/*if function == "init" {
		fmt.Println("Function is init")
		return t.Init(stub)}*/
	 var result []byte
	 var err error
	 	// Handle different functions
		if function == "init" {
		fmt.Println("Function is init")
		//result = t.Init(stub)
		//return shim.Success(result)
		return t.Init(stub)
		}else if function == "createAssembly" {
		fmt.Println("Function is createAssembly")
		//return t.createAssembly(stub, args)
		result, err = t.createAssembly(stub, args)
		return shim.Success(result)
		} else if function == "updateAssemblyByID" {
		fmt.Println("Function is updateAssemblyByID")
		//return t.updateAssemblyByID(stub, args)
		result, err = t.updateAssemblyByID(stub, args)
		return shim.Success(result)
		}  else if function == "createPackage" {
		fmt.Println("Function is createPackage")
		//return t.createPackage(stub, args)
		result, err = t.createPackage(stub, args)
		return shim.Success(result)
		} else if function == "updatePackage" {
		fmt.Println("Function is updatePackage")
		//return t.updatePackage(stub, args)
		result, err = t.updatePackage(stub, args)
		return shim.Success(result)
	}else  if function == "getAssemblyByID" { 
		t := TnT{}
		//return t.getAssemblyByID(stub, args)
		result, err = t.getAssemblyByID(stub, args)
		return shim.Success(result)

	} else if function == "getPackageByID" { 
		t := TnT{}
		//return t.getPackageByID(stub, args)
		result, err = t.getPackageByID(stub, args)
		return shim.Success(result)
	} else if function == "getAllAssemblies" { 
		t := TnT{}
		//return t.getAllAssemblies(stub, args)
		result, err = t.getAllAssemblies(stub, args)
		return shim.Success(result)
	} else if function == "getAllPackages" { 
		t := TnT{}
		//return t.getAllPackages(stub, args)
		result, err = t.getAllPackages(stub, args)
		return shim.Success(result)
		} else if function == "getAllAssemblyIDs" { 
		t := TnT{}
		//return t.getAllAssemblyIDs(stub, args)
		result, err = t.getAllAssemblyIDs(stub, args)
		return shim.Success(result)
		} else if function == "getAllPackageCaseIDs" { 
		t := TnT{}
		//return t.getAllPackageCaseIDs(stub, args)
		result, err = t.getAllPackageCaseIDs(stub, args)
		return shim.Success(result)
		} else if function == "get_ecert" {
		t := TnT{}
		//return t.get_ecert(stub, args[0])
		result, err = t.get_ecert(stub, args[0])
		return shim.Success(result)
		} else if function == "validateCreateAssembly" {
		t := TnT{}
		//return t.validateCreateAssembly(stub, args)
		result, err = t.validateCreateAssembly(stub, args)
		return shim.Success(result)
		} else if function == "validateUpdateAssembly" {
		t := TnT{}
		//return t.validateUpdateAssembly(stub, args)
		result, err = t.validateUpdateAssembly(stub, args)
		return shim.Success(result)
		} else if function == "validateCreatePackage" {
		t := TnT{}
		//return t.validateCreatePackage(stub, args)
		result, err = t.validateCreatePackage(stub, args)
		return shim.Success(result)
		} else if function == "validateUpdatePackage" {
		t := TnT{}
		//return t.validateUpdatePackage(stub, args)
		result, err = t.validateUpdatePackage(stub, args)
		return shim.Success(result)
		} else if function == "getAssemblyLineHistoryByID" {
		t := TnT{}
		//return t.getAssemblyLineHistoryByID(stub, args)
		result, err = t.getAssemblyLineHistoryByID(stub, args)
		return shim.Success(result)
		} else if function == "getPackageLineHistoryByID" {
		t := TnT{}
		//return t.getPackageLineHistoryByID(stub, args)
		result, err = t.getPackageLineHistoryByID(stub, args)
		return shim.Success(result)
		} else if function == "getAssembliesByBatchNumber" {
		t := TnT{}
		//return t.getAssembliesByBatchNumber(stub, args)
		result, err = t.getAssembliesByBatchNumber(stub, args)
		return shim.Success(result)
		} else if function == "getAssembliesByDate" {
		t := TnT{}
		//return t.getAssembliesByDate(stub, args)
		result, err = t.getAssembliesByDate(stub, args)
		return shim.Success(result)
		} else if function == "getAssembliesHistoryByDate" {
		t := TnT{}
		//return t.getAssembliesHistoryByDate(stub, args)
		result, err = t.getAssembliesHistoryByDate(stub, args)
		return shim.Success(result)
		} else if function == "getAssembliesByBatchNumberAndByDate" {
		t := TnT{}
		//return t.getAssembliesByBatchNumberAndByDate(stub, args)
		result, err = t.getAssembliesByBatchNumberAndByDate(stub, args)
		return shim.Success(result)
	}  

	if err != nil {
		return shim.Error(err.Error())
	}


	//return nil, errors.New("Received unknown function invocation")
	return shim.Error("Received unknown function invocation")
}


// query queries the chaincode
/*func (t *TnT) Query(stub shim.ChaincodeStubInterface) pb.Response {

	fmt.Println("Query called, determining function")
 	// Extract the function and args from the transaction proposal
 	function, args := stub.GetFunctionAndParameters()

	 var result []byte
	 var err error
	 if function == "getAssemblyByID" { 
		t := TnT{}
		//return t.getAssemblyByID(stub, args)
		result, err = t.getAssemblyByID(stub, args)
		return shim.Success(result)

	} else if function == "getPackageByID" { 
		t := TnT{}
		//return t.getPackageByID(stub, args)
		result, err = t.getPackageByID(stub, args)
		return shim.Success(result)

	} else if function == "getAllAssemblies" { 

		t := TnT{}
		//return t.getAllAssemblies(stub, args)
		result, err = t.getAllAssemblies(stub, args)
		return shim.Success(result)

	} else if function == "getAllPackages" { 

		t := TnT{}
		//return t.getAllPackages(stub, args)
		result, err = t.getAllPackages(stub, args)
		return shim.Success(result)

	} else if function == "getAllAssemblyIDs" { 

		t := TnT{}
		//return t.getAllAssemblyIDs(stub, args)
		result, err = t.getAllAssemblyIDs(stub, args)
		return shim.Success(result)

	} else if function == "getAllPackageCaseIDs" { 

		t := TnT{}
		//return t.getAllPackageCaseIDs(stub, args)
		result, err = t.getAllPackageCaseIDs(stub, args)
		return shim.Success(result)

	} else if function == "get_ecert" {

		t := TnT{}
		//return t.get_ecert(stub, args[0])
		result, err = t.get_ecert(stub, args[0])
		return shim.Success(result)

	} else if function == "validateCreateAssembly" {

		t := TnT{}
		//return t.validateCreateAssembly(stub, args)
		result, err = t.validateCreateAssembly(stub, args)
		return shim.Success(result)

	} else if function == "validateUpdateAssembly" {

		t := TnT{}
		//return t.validateUpdateAssembly(stub, args)
		result, err = t.validateUpdateAssembly(stub, args)
		return shim.Success(result)

	} else if function == "validateCreatePackage" {

		t := TnT{}
		//return t.validateCreatePackage(stub, args)
		result, err = t.validateCreatePackage(stub, args)
		return shim.Success(result)

	} else if function == "validateUpdatePackage" {

		t := TnT{}
		//return t.validateUpdatePackage(stub, args)
		result, err = t.validateUpdatePackage(stub, args)
		return shim.Success(result)

	} else if function == "getAssemblyLineHistoryByID" {

		t := TnT{}
		//return t.getAssemblyLineHistoryByID(stub, args)
		result, err = t.getAssemblyLineHistoryByID(stub, args)
		return shim.Success(result)

	} else if function == "getPackageLineHistoryByID" {

		t := TnT{}
		//return t.getPackageLineHistoryByID(stub, args)
		result, err = t.getPackageLineHistoryByID(stub, args)
		return shim.Success(result)

	} else if function == "getAssembliesByBatchNumber" {

		t := TnT{}
		//return t.getAssembliesByBatchNumber(stub, args)
		result, err = t.getAssembliesByBatchNumber(stub, args)
		return shim.Success(result)

	} else if function == "getAssembliesByDate" {

		t := TnT{}
		//return t.getAssembliesByDate(stub, args)
		result, err = t.getAssembliesByDate(stub, args)
		return shim.Success(result)

	} else if function == "getAssembliesHistoryByDate" {

		t := TnT{}
		//return t.getAssembliesHistoryByDate(stub, args)
		result, err = t.getAssembliesHistoryByDate(stub, args)
		return shim.Success(result)

	} else if function == "getAssembliesByBatchNumberAndByDate" {

		t := TnT{}
		//return t.getAssembliesByBatchNumberAndByDate(stub, args)
		result, err = t.getAssembliesByBatchNumberAndByDate(stub, args)
		return shim.Success(result)

	} 

	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Error("Received unknown function query")
	
}*/


//main function
// main function starts up the chaincode in the container during instantiate
func main() {
    if err := shim.Start(new(TnT)); err != nil {
            fmt.Printf("Error starting SimpleAsset chaincode: %s", err)
	}
}


