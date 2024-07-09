package main

import (
	"context"
	"encoding/json"

	"github.com/free5gc/openapi/Namf_Communication"
	"github.com/free5gc/openapi/models"
)

func subscribe() {
	subscribe := models.NonUeN2InfoSubscriptionCreateData{}
	// Specify the URL you want to send the request to
	// Create the request body
	jsonString := []byte(`{
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
		"anTypeList":[

		],
		"n2InformationClass": "PWS",
		"n2NotifyCallbackUri": "192.168.56.102:8000/notify",
		"nfId": "",
		"supportedFeatures": ""
	  }`)
	json.Unmarshal(jsonString, &subscribe)
	namfConfiguration := Namf_Communication.NewConfiguration()
	namfConfiguration.SetBasePath("http://127.0.0.18:8000")
	apiClient := Namf_Communication.NewAPIClient(namfConfiguration)
	_, _, _ = apiClient.NonUEN2MessagesSubscriptionsCollectionDocumentApi.NonUeN2InfoSubscribe(context.TODO(), subscribe)
}
