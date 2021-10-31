package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"golang.org/x/net/http2"
)

var url = "https://localhost:8080"

type Airline struct {
	Id           int64  `json:"id" validate:"required,max=9999"`
	Name         string `json:"name"`
	Country      string `json:"country"`
	Logo         string `json:"logo" validate:"omitempty,url"`
	Slogan       string `json:"slogan"`
	Head_Quaters string `json:"head_quaters"`
	Website      string `json:"website" validate:"omitempty,url"`
	Established  string `json:"established" validate:"omitempty"`
}

type Passenger struct {
	Id       string    `json:"_id" validate:"required"`
	Name     string    `json:"name"`
	Trips    int       `json:"trips" validate:"lte=1000"`
	Airlines []Airline `json:"airline"`
	Version  int       `json:"__v" validate:"required"`
}

func main() {
	client := GetClient()

	for {
		fmt.Println("Please Select the option")

		fmt.Println("Press 1 - Get All Passengers")
		fmt.Println("Press 2 - Get Passenger by Id")
		fmt.Println("Press 3 - Post Passenger")
		fmt.Println("Press 4 - Put Passenger")
		fmt.Println("Press 5 - Delete Passenger")
		fmt.Println("Press 6 - Get All Airlines")
		fmt.Println("Press 7 - Get Airlines by ID")
		fmt.Println("Press 8 - Post Airlines")

		var option string
		fmt.Scanln(&option)

		switch option {
		case "1":
			GetPassengerList(client)

		case "2":
			GetPassengerById(client)

		case "3":
			PostPassenger(client)

		case "4":
			PutPassenger(client)

		case "5":
			DeletePassenger(client)

		case "6":
			GetAirlines(client)

		case "7":
			GetAirlineById(client)

		case "8":
			PostAirline(client)

		default:
			fmt.Println("Error - Please Select the correct option")
		}
	}
}

func DeletePassenger(client http.Client) {
	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/passenger/92be20cc395e11ec8d3d0242ac130003", url), nil)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error while fetching the request. Err [%+v] ", err)
	}
	resArr, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error while reading the response. Err [%+v]", err)
		return
	}
	var passenger Passenger
	err = json.Unmarshal(resArr, &passenger)
	if err != nil {
		fmt.Printf("Error while parsing the data. Err [%+v]", err)
		return
	}
	fmt.Printf("%+v", passenger)
}

func PutPassenger(client http.Client) {
	updatePassenger := Passenger{
		Id:       "92be20cc395e11ec8d3d0242ac130003",
		Name:     "Donald Buck",
		Trips:    300,
		Airlines: []Airline{},
		Version:  1,
	}

	reqArr, err := json.Marshal(updatePassenger)
	if err != nil {
		fmt.Printf("Error while parsing the request body. Err [%+v]", err)
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/passenger/92be20cc395e11ec8d3d0242ac130003", url), bytes.NewBuffer(reqArr))

	if err != nil {
		fmt.Printf("Error while creating request. Err [%+v]", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("Error on getting response. Err [%+v]", err)
	}

	var resPassenger Passenger
	resArr, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Printf("Error while reading the response body. Error : [%+v]", err)
	}

	err = json.Unmarshal(resArr, &resPassenger)

	if err != nil {
		fmt.Printf("Error while parsing the response body. Error : [%+v]", err)
	}

	fmt.Printf("Passenger has been added : %+v", resPassenger)
}

func PostPassenger(client http.Client) {
	newPassenger := Passenger{
		Id:       "92be20cc395e11ec8d3d0242ac130003",
		Name:     "Donald Duck",
		Trips:    500,
		Airlines: []Airline{},
		Version:  1,
	}

	reqArr, err := json.Marshal(newPassenger)
	if err != nil {
		fmt.Printf("Error while parsing the request body. Err [%+v]", err)
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/passenger/", url), bytes.NewBuffer(reqArr))

	if err != nil {
		fmt.Printf("Error while creating request. Err [%+v]", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("Error on getting response. Err [%+v]", err)
	}

	var resPassenger Passenger
	resArr, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Printf("Error while reading the response body. Error : [%+v]", err)
	}

	err = json.Unmarshal(resArr, &resPassenger)

	if err != nil {
		fmt.Printf("Error while parsing the response body. Error : [%+v]", err)
	}

	fmt.Printf("Passenger has been added : %+v", resPassenger)
}

func GetPassengerList(client http.Client) {
	resp, err := client.Get(fmt.Sprintf("%s/passenger", url))
	if err != nil {
		fmt.Printf("Error while fetching the request. Err [%+v] ", err)
	}
	resArr, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error while reading the response. Err [%+v]", err)
		return
	}
	var passengers []Passenger
	err = json.Unmarshal(resArr, &passengers)
	if err != nil {
		fmt.Printf("Error while parsing the data. Err [%+v]", err)
		return
	}
	fmt.Printf("%+v", passengers)
}

func GetPassengerById(client http.Client) {
	resp, err := client.Get(fmt.Sprintf("%s/passenger/92be20cc395e11ec8d3d0242ac130003", url))
	if err != nil {
		fmt.Printf("Error while fetching the request. Err [%+v] ", err)
	}
	resArr, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error while reading the response. Err [%+v]", err)
		return
	}
	var passenger Passenger
	err = json.Unmarshal(resArr, &passenger)
	if err != nil {
		fmt.Printf("Error while parsing the data. Err [%+v]", err)
		return
	}
	fmt.Printf("%+v", passenger)
}

func GetAirlines(client http.Client) {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/airlines/", url), nil)
	resp, _ := client.Do(req)

	var airlines []Airline
	respArr, _ := ioutil.ReadAll(resp.Body)
	err := json.Unmarshal(respArr, &airlines)
	if err != nil {
		fmt.Printf("Something went wrong while parsing data. Err - [%+v]", err)
	}

	fmt.Println(airlines)
}

func GetAirlineById(client http.Client) {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/airlines/1990", url), nil)
	resp, _ := client.Do(req)

	var airline Airline
	respArr, _ := ioutil.ReadAll(resp.Body)
	err := json.Unmarshal(respArr, &airline)
	if err != nil {
		fmt.Printf("Something went wrong while parsing data. Err - [%+v]", err)
	}

	fmt.Println(airline)
}

func PostAirline(client http.Client) {
	newAirline := Airline{
		Id:           6970,
		Name:         "Air India",
		Country:      "India",
		Logo:         "https://images.hindustantimes.com/img/2021/10/08/1600x900/maharaja_1633694839415_1633694844518.jpg",
		Slogan:       "Unbeatable Service",
		Head_Quaters: "New Delhi",
		Website:      "https://www.google.com/url?sa=t&rct=j&q=&esrc=s&source=web&cd=&cad=rja&uact=8&ved=2ahUKEwj09YWSq_LzAhUPVysKHQqEAoIQFnoECAcQAQ&url=https%3A%2F%2Fwww.airindia.in%2F&usg=AOvVaw1-CB6i0XUVT2n-P34RfRQo",
		Established:  "1930",
	}

	reqArr, _ := json.Marshal(newAirline)

	req, _ := http.NewRequest(http.MethodPost,
		fmt.Sprintf("%s/airlines/", url),
		bytes.NewBuffer(reqArr))

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Something went wromg while making the request. Err [%+v]", err)
	}

	respArr, _ := ioutil.ReadAll(resp.Body)

	var resAirline Airline
	err = json.Unmarshal(respArr, &resAirline)
	if err != nil {
		fmt.Printf("Error while parsing the response. Err [%+v]", err)
	}

	fmt.Printf("New Airline has been added. Airline [%+v]", resAirline)
}

func GetClient() http.Client {

	cwd, _ := os.Getwd()
	rootPath := filepath.Dir(filepath.Dir(cwd))
	credPath := path.Join(rootPath, "creds")

	cacert, err := ioutil.ReadFile(fmt.Sprintf("%s/cert.pem", credPath))
	if err != nil {
		log.Fatalf("Reading server certificate: %s", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(cacert)

	// dialer := net.Dialer{
	// 	Timeout:   5 * time.Second,
	// 	KeepAlive: 1 * time.Minute,
	// }
	tlsConfig := tls.Config{
		RootCAs: caCertPool,
	}
	transport := http2.Transport{
		TLSClientConfig: &tlsConfig,
	}
	client := http.Client{
		Transport: &transport,
	}

	return client
}
