package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/free5gc/openapi/Namf_Communication"
	"github.com/free5gc/openapi/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func transfer(data map[string]string) {
	var err error
	taiwanTimezone, err := time.LoadLocation("Asia/Taipei")
	reqData := models.N2InformationTransferReqData{}
	message := models.NonUeN2MessageTransferRequest{}
	jsonString := []byte(`{
		"taiList": [
		  {
			"tac": "",
			"plmnId": {
			  "mnc": "",
			  "mcc": ""
			}
		  }
		],
		"ratSelector": "PWS",
		"ecgiList": [
		  {
			"eutraCellId": "",
			"plmnId": {
			  "mnc": "",
			  "mcc": ""
			}
		  }
		],
		"ncgiList": [
		  {
			"nrCellId": "",
			"plmnId": {
			  "mnc": "",
			  "mcc": ""
			}
		  }
		],
		"globalRanNodeList": [
		  {
			"gNbId": {
			  "bitLength": 24,
			  "gNBValue": ""
			},
			"plmnId": {
			  "mnc": "",
			  "mcc": ""
			},
			"n3IwfId": "",
			"ngeNbId": ""
		  }
		],
		"n2Information": {
		  "n2InformationClass": "PWS",
		  "smInfo": {
			"subjectToHo": false,
			"pduSessionId": 29,
			"n2InfoContent": {
			  "ngapData": {
				"contentId": "contentId"
			  },
			  "ngapIeType": "",
			  "ngapMessageType": 32
			},
			"sNssai": {
			  "sd": "sd",
			  "sst": 32
			}
		  },
		  "ranInfo": {
			"n2InfoContent": {
			  "ngapData": {
				"contentId": "contentId"
			  },
			  "ngapIeType": "",
			  "ngapMessageType": 32
			}
		  },
		  "nrppaInfo": {
			"nfId": "nfId",
			"nrppaPdu": {
			  "n2InfoContent": {
				"ngapData": {
				  "contentId": "contentId"
				},
				"ngapIeType": "",
				"ngapMessageType": 32
			  }
			}
		  },
		  "pwsInfo": {
			"messageIdentifier": 0,
			"serialNumber": 0,
			"pwsContainer": {
			  "ngapData": {
				"contentId": "n2msg"
			  },
			  "ngapIeType": "",
			  "ngapMessageType": 51
			},
			"sendRanResponse": true,
			"omcId": true
		  }
		},
		"supportedFeatures": ""
	  }
	  `)
	BinaryDataN2informationString := `"messageType": "",
		"messageIdentifier": "",
		"serialNumber": "",
		"warningAreaList": "",
		"repetitionPeriod": "",
		"numberOfBroadcast": "",
		"warningType": "",
		"warningSecurityInformation": "",
		"dataCodingScheme": "",
		"warningMessageContents" : "",
		"concurrentWarningMessageIndicator": "",
		"warningAreaCoordinates": ""
		`
	BinaryDataN2InformationKeyValue := make(map[string]interface{})
	json.Unmarshal([]byte(BinaryDataN2informationString), &BinaryDataN2InformationKeyValue)
	BinaryDataN2InformationKeyValue["messageIdentifier"] = data["messageIdentifier"]
	BinaryDataN2InformationKeyValue["repetitionPeriod"] = "10"
	BinaryDataN2InformationKeyValue["numberOfBroadcastsRequested"] = "3"
	BinaryDataN2InformationKeyValue["dataCodingScheme"] = data["dataCodingScheme"]
	BinaryDataN2InformationKeyValue["warningMessageContents"] = data["warningMessageContents"]
	json.Unmarshal(jsonString, &reqData)
	message.JsonData = &reqData
	if data["ratSelector"] == "NR" {
		message.JsonData.RatSelector = models.RatSelector_NR
	}
	if data["ratSelector"] == "E-UTRA" {
		message.JsonData.RatSelector = models.RatSelector_E_UTRA
	}
	(*message.JsonData.TaiList)[0].PlmnId.Mcc = data["mcc"]
	(*message.JsonData.TaiList)[0].PlmnId.Mnc = data["mnc"]
	(*message.JsonData.TaiList)[0].Tac = data["tac"]
	(*message.JsonData.EcgiList)[0].PlmnId.Mcc = data["mcc"]
	(*message.JsonData.EcgiList)[0].PlmnId.Mnc = data["mnc"]
	(*message.JsonData.NcgiList)[0].PlmnId.Mcc = data["mcc"]
	(*message.JsonData.NcgiList)[0].PlmnId.Mnc = data["mnc"]
	(*message.JsonData.GlobalRanNodeList)[0].PlmnId.Mcc = data["mcc"]
	(*message.JsonData.GlobalRanNodeList)[0].PlmnId.Mnc = data["mnc"]

	countNumber := countMessageFromDatabase(data["messageIdentifier"], data["serialNumber"])
	var serialNumber int64
	var serialNumberInteger int
	if countNumber >= 0 {
		serialNumberInteger, err = strconv.Atoi(data["serialNumber"])
		if err != nil {
			fmt.Println(err)
		}
		serialNumberBits := "01" + "00" + fmt.Sprintf("%08b", serialNumberInteger) + fmt.Sprintf("%04b", countNumber)
		serialNumber, err = strconv.ParseInt(serialNumberBits, 2, 64)
	}
	message.JsonData.N2Information.PwsInfo.SerialNumber = int32(serialNumberInteger)
	messageIdentifier, err := strconv.Atoi(data["messageIdentifier"])
	message.JsonData.N2Information.PwsInfo.MessageIdentifier = int32(messageIdentifier)
	BinaryDataN2InformationKeyValue["serialNumber"] = fmt.Sprintf("%x", serialNumber)
	(*&message.BinaryDataN2Information), err = json.Marshal(BinaryDataN2InformationKeyValue)
	jsonString, err = json.Marshal(message)
	message.JsonData.N2Information.PwsInfo.SerialNumber = int32(serialNumber)
	namfConfiguration := Namf_Communication.NewConfiguration()
	namfConfiguration.SetBasePath("http://127.0.0.18:8000")
	apiClient := Namf_Communication.NewAPIClient(namfConfiguration)

	_, _, err = apiClient.NonUEN2MessagesCollectionDocumentApi.NonUeN2MessageTransfer(context.TODO(), message)
	currentTime := time.Now().In(taiwanTimezone)
	formattedTime := currentTime.Format("2006-01-02 15:04:05.000000 MST")
	fmt.Printf("Sent Write Replace Rquest to AMF at %s\n", formattedTime)
	if err != nil {
		log.Fatalf("Failed to send HTTP request: %v", err)
	}
}

/*
	func insertToDatabase(message models.NonUeN2MessageTransferRequest, messageTimestamp MessageTimestamp) {
		clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
		client, err := mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			log.Fatal(err)
		}
		err = client.Ping(context.TODO(), nil)
		if err != nil {
			log.Fatal(err)
		}
		collection := client.Database("local").Collection("cbcf")
		collectionTimestamp := client.Database("local").Collection("cbcfTimestamp")
		_, err = collectionTimestamp.InsertOne(context.TODO(), messageTimestamp)
		if err != nil {
			log.Fatal(err)
		}
		_, err = collection.InsertOne(context.Background(), message)
		if err != nil {
			log.Fatal(err)
		}

		err = client.Disconnect(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
	}
*/
func countMessageFromDatabase(messageIdentifier string, serialNumber string) int64 {
	messageIdentifierint, err := strconv.Atoi(messageIdentifier)
	serialNumberint, err := strconv.Atoi(serialNumber)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	collection := client.Database("local").Collection("cbcf")
	var result int64
	condition := bson.M{"jsondata.n2information.pwsinfo.messageidentifier": messageIdentifierint, "jsondata.n2information.pwsinfo.serialnumber": serialNumberint}
	result, err = collection.CountDocuments(context.Background(), condition)
	if err != nil {
		log.Fatal(err)
	}
	return result
}
