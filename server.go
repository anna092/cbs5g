package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"
)

type Alert struct {
	XMLName    xml.Name `xml:"urn:oasis:names:tc:emergency:cap:1.2 alert"`
	Identifier string   `xml:"identifier"`
	Sender     string   `xml:"sender"`
	Sent       string   `xml:"sent"`
	Status     string   `xml:"status"`
	MsgType    string   `xml:"msgType"`
	Scope      string   `xml:"scope"`
	Source     string   `xml:"source"`
	Info       Info     `xml:"info"`
}

type Info struct {
	Language    string `xml:"language"`
	Category    string `xml:"category"`
	Event       string `xml:"event"`
	Urgency     string `xml:"urgency"`
	Severity    string `xml:"severity"`
	Certainty   string `xml:"certainty"`
	Eventcode   string `xml:"eventCode"`
	Expires     string `xml:"expires"`
	SenderName  string `xml:"senderName"`
	Headline    string `xml:"headline"`
	Description string `xml:"description"`
	Instruction string `xml:"instruction"`
	Contact     string `xml:"contact"`
	Area        Area   `xml:"area"`
}

type Area struct {
	AreaDesc string `xml:"areaDesc"`
	Polygon  string `xml:"polygon"`
	GeoCode  string `xml:"geocode"`
}

type CBCFNotification struct {
	IMSI   string `json:"imsi"`
	Status string `json:"status"`
}

type CommandConfig struct {
	Commands []struct {
		IMSI    string `yaml:"IMSI"`
		Command string `yaml:"CMD"`
	} `yaml:"commands"`
}

var (
	db     *sql.DB
	cmdMap map[string]string
)

//var db *sql.DB
//var commandsData CommandConfig

func main() {
	// Initialize database connection
	initDB()
	loadCommandConfig()

	http.HandleFunc("/EmergencyBroadcastRequest", handleEmergencyBroadcastRequest)
	http.HandleFunc("/notify", handleNotify)
	http.HandleFunc("/ue-registration-notify", handleUERegistrationNotify)
	bindAddress := "192.168.56.102:8080"
	fmt.Printf("Listening on http://%s\n", bindAddress)
	log.Fatal(http.ListenAndServe(bindAddress, nil))
}

func loadCommandConfig() {
	// Read the cmdcfg.yaml file
	data, err := ioutil.ReadFile("cmdcfg.yaml")
	if err != nil {
		log.Fatalf("read file error: %vv", err)
	}

	// Unmarshal the YAML content into CommandConfig struct
	var config CommandConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		log.Fatalf("unmarshal YAML error: %v", err)
	}

	// Initialize the command map
	cmdMap = make(map[string]string)
	for _, cmd := range config.Commands {
		cmdMap[cmd.IMSI] = cmd.Command
	}

	fmt.Println("Command configuration loaded successfully")
}

func initDB() {
	dsn := "root:5gpws@tcp(127.0.0.1:3306)/CBSDB"
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	fmt.Println("Database connection successful")
}

func handleUERegistrationNotify(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Test handleUERegistrationNotify Function")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var notification CBCFNotification
	if err := json.NewDecoder(r.Body).Decode(&notification); err != nil {
		http.Error(w, "Failed to decode JSON body", http.StatusBadRequest)
		return
	}

	fmt.Printf("Received notification from AMF: %+v\n", notification)

	// Save IMSI and status to database
	if err := saveIMSI(notification.IMSI, notification.Status); err != nil {
		fmt.Printf("Failed to save IMSI %s with status %s to database: %v\n", notification.IMSI, notification.Status, err)
		http.Error(w, "Failed to save IMSI to database", http.StatusInternalServerError)
		return
	}

	fmt.Println("IMSI saved to database successfully")
	/*
	   if err := saveIMSI(notification.IMSI, notification.Status); err != nil {
	       http.Error(w, "Failed to save IMSI to database", http.StatusInternalServerError)
	       return
	   }
	*/
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Notification received successfully"))
}

func saveIMSI(imsi, status string) error {
	// Prepare statement for inserting data
	stmt, err := db.Prepare("INSERT INTO UERegistrations (imsi, status, registered_at) VALUES (?, ?, NOW())")
	if err != nil {
		return fmt.Errorf("prepare statement error: %v", err)
	}
	defer stmt.Close()

	// Execute the prepared statement with IMSI and status values
	_, err = stmt.Exec(imsi, status)
	if err != nil {
		return fmt.Errorf("execute statement error: %v", err)
	}

	fmt.Printf("IMSI %s with status %s saved to database\n", imsi, status)
	return nil
}

func handleEmergencyBroadcastRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Test handleEmergencyBroadcastRequest Function\n")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	xmlData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error parsing body data", http.StatusBadRequest)
		return
	}

	taiwanTimezone, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		http.Error(w, "Error loading timezone", http.StatusInternalServerError)
		return
	}

	currentTime := time.Now().In(taiwanTimezone)
	formattedTime := currentTime.Format("2006-01-02 15:04:05.000000 MST")
	fmt.Printf("Received Emergency Broadcast Request from CBE at %s\n", formattedTime)

	fmt.Printf("Received XML data\n")

	var alertData Alert
	if err := xml.Unmarshal(xmlData, &alertData); err != nil {
		fmt.Println(err)
		http.Error(w, "Error parsing XML data", http.StatusBadRequest)
		return
	}

	data := make(map[string]string)
	data["serialNumber"] = alertData.Identifier[len(alertData.Identifier)-3:]
	data["messageType"] = alertData.MsgType
	if alertData.Info.Language == "en-US" {
		data["dataCodingScheme"] = "01"
	}
	if alertData.Info.Language == "zh-TW" {
		data["dataCodingScheme"] = "48"
	}
	switch {
	case alertData.Info.Severity == "Extreme" && alertData.Info.Urgency == "Immediate" && alertData.Info.Certainty == "Observed":
		data["messageIdentifier"] = "1113"
	case alertData.Info.Severity == "Extreme" && alertData.Info.Urgency == "Immediate" && alertData.Info.Certainty == "Likely":
		data["messageIdentifier"] = "1114"
	case alertData.Info.Severity == "Extreme" && alertData.Info.Urgency == "Expected" && alertData.Info.Certainty == "Observed":
		data["messageIdentifier"] = "1115"
	case alertData.Info.Severity == "Extreme" && alertData.Info.Urgency == "Expected" && alertData.Info.Certainty == "Likely":
		data["messageIdentifier"] = "1116"
	case alertData.Info.Severity == "Severe" && alertData.Info.Urgency == "Immediate" && alertData.Info.Certainty == "Observed":
		data["messageIdentifier"] = "1117"
	case alertData.Info.Severity == "Severe" && alertData.Info.Urgency == "Immediate" && alertData.Info.Certainty == "Likely":
		data["messageIdentifier"] = "1118"
	case alertData.Info.Severity == "Severe" && alertData.Info.Urgency == "Expected" && alertData.Info.Certainty == "Observed":
		data["messageIdentifier"] = "1119"
	case alertData.Info.Severity == "Severe" && alertData.Info.Urgency == "Expected" && alertData.Info.Certainty == "Likely":
		data["messageIdentifier"] = "111A"
	default:
		data["messageIdentifier"] = "1112"
	}

	timeFormatMsg := "2006-01-02 15:04:05.000 UTC-07:00"
	timeSent, err := time.Parse(timeFormatMsg, alertData.Sent)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error parsing Sent time", http.StatusBadRequest)
		return
	}

	for _, imsi := range getRegisteredIMSIs() {
		if command, found := cmdMap[imsi]; found {
			fmt.Printf("Started distributing the different commands...\n")
			data["warningMessageContents"] = timeSent.Format(timeFormatMsg) + " " + imsi + " " + command
		} else {
			fmt.Printf("Started distributing the instruction...\n")
			data["warningMessageContents"] = timeSent.Format(timeFormatMsg) + " " + alertData.Info.Headline + "\n" + alertData.Info.Instruction
		}
	}

	//data["warningMessageContents"] = timeSent.Format(timeFormatMsg) + " " + alertData.Info.Headline + "\n" + alertData.Info.Instruction

	data["tac"] = alertData.Info.Area.GeoCode
	subscribe()
	transfer(data)

	fmt.Printf("Started distributing the warning message...\n")
	sendEmergencyBroadcastResponse(w)
}

func getRegisteredIMSIs() []string {
	rows, err := db.Query("SELECT imsi FROM UERegistrations")
	if err != nil {
		log.Fatalf("Failed to query IMSIs from database: %v", err)
	}
	defer rows.Close()

	var imsils []string
	for rows.Next() {
		var imsi string
		if err := rows.Scan(&imsi); err != nil {
			log.Fatalf("Failed to scan IMSI from row: %v", err)
		}
		imsils = append(imsils, imsi)
	}
	if err := rows.Err(); err != nil {
		log.Fatalf("Error iterating over rows: %v", err)
	}
	return imsils
}

func sendEmergencyBroadcastResponse(w http.ResponseWriter) {
	fmt.Printf("Test sendEmergencyBroadcastResponse Function\n")
	responseData := `<response><status>Success</status><description>Emergency Broadcast Received</description></response>`
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(responseData))

	taiwanTimezone, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		log.Fatalf("Failed to load timezone: %v", err)
	}
	currentTime := time.Now().In(taiwanTimezone)
	formattedTime := currentTime.Format("2006-01-02 15:04:05.000000 MST")

	fmt.Printf("Sent Emergency Broadcast Response to CBE at %s\n", formattedTime)
}

func handleNotify(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	fmt.Println(string(body))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Received the request body"))
}
