package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-playground/validator"
)

type HelloHandler struct {
}

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

func (src *Passenger) mergePassenger(dst *Passenger) {
	if dst.Trips != 0 {
		src.Trips = dst.Trips
	}
	if !reflect.DeepEqual(dst.Airlines, ([]Airline{})) {
		src.Airlines = dst.Airlines
	}
	if dst.Name != src.Name {
		src.Name = dst.Name
	}
}

var airlineData []Airline
var passengerData []Passenger

var validate = validator.New()

func loadData(filenames []string) {

	for _, filename := range filenames {
		if filename == "airline.json" {
			rawAirlineData, err := ioutil.ReadFile(filename)
			if err != nil {
				log.Panicf("Error while Reading the airline data. [%+v]", err)
			} else {
				err := json.Unmarshal(rawAirlineData, &airlineData)
				if err != nil {
					fmt.Printf("Error while parsing the airline data. [%+v]", err)
				}
			}
			fmt.Println("Successfully Parsed the airline Data")
		}
		if filename == "passenger.json" {
			rawPassengerData, err := ioutil.ReadFile(filename)
			if err != nil {
				log.Panicf("Error while Reading the Passenger data. [%+v]", err)
			} else {
				err := json.Unmarshal(rawPassengerData, &passengerData)
				if err != nil {
					fmt.Printf("Error while parsing the Passenger data. [%+v]", err)
				}
			}
			fmt.Println("Successfully Parsed the passenger Data")
		}
	}

	fmt.Println("Successfully Parsed the Data")
}

func (h HelloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello %s", r.URL.Path)
}

func GetAirlines(w http.ResponseWriter, r *http.Request) {
	method := r.Method

	p := strings.Split(r.URL.Path, "/")
	if method == http.MethodGet {
		if p[len(p)-1] == "" {
			resArr, _ := json.Marshal(airlineData)
			fmt.Fprintf(w, "%s", string(resArr))
		} else {
			id, _ := strconv.Atoi(p[len(p)-1])
			for _, airline := range airlineData {
				if int(airline.Id) == id {
					resArr, _ := json.Marshal(airline)
					fmt.Fprintf(w, "%s", string(resArr))
					return
				}
			}
			fmt.Fprintf(w, "No Records Found")
		}
	}
	if method == http.MethodPost {
		var newAirline Airline
		rawBody, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()

		json.Unmarshal(rawBody, &newAirline)
		err := validate.Struct(newAirline)

		if err != nil {
			fmt.Fprintf(w, "%+v", err)
		}

		airlineData = append(airlineData, newAirline)

		resArr, _ := json.Marshal(newAirline)
		fmt.Fprintf(w, "%s", string(resArr))
	}
}

func GetPassengers(w http.ResponseWriter, r *http.Request) {
	method := r.Method

	p := strings.Split(r.URL.Path, "/")

	if method == http.MethodGet {
		if p[len(p)-1] == "" {
			pres, _ := json.Marshal(passengerData)
			fmt.Fprintf(w, "%s", string(pres))
		} else {
			id := p[len(p)-1]
			for _, passenger := range passengerData {
				if passenger.Id == id {
					pres, _ := json.Marshal(passenger)
					fmt.Fprintf(w, "%s", string(pres))
					return
				}
			}
			fmt.Fprintf(w, "No Records Found")
		}
	}
	if method == http.MethodPost {
		var newPassenger Passenger
		rawBody, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()

		json.Unmarshal(rawBody, &newPassenger)
		err := validate.Struct(newPassenger)

		if err != nil {
			fmt.Fprintf(w, "Error while reading the Body. \n [%+v]", err)
		}

		passengerData = append(passengerData, newPassenger)

		pres, err := json.Marshal(newPassenger)

		if err != nil {
			fmt.Fprintf(w, "Error while parsing the new Passenger. \n [%+v]", err)
		}
		fmt.Fprintf(w, "%s", string(pres))
		return
	}
	if method == http.MethodDelete {
		var delPassenger Passenger
		if p[len(p)-1] != "" {
			id := p[len(p)-1]
			for index, passenger := range passengerData {
				if passenger.Id == id {
					delPassenger = passenger
					passengerData = append(passengerData[0:index], passengerData[index+1:]...)
					break
				}
			}
		} else {
			fmt.Fprintf(w, "Error : No Id Passed")
		}

		if reflect.DeepEqual(delPassenger, (Passenger{})) {
			fmt.Fprintf(w, "No Record Found")
		} else {
			pres, _ := json.Marshal(delPassenger)
			fmt.Fprintf(w, "%s", string(pres))
		}
	}
	if method == http.MethodPut {
		var newPassenger Passenger

		rawBody, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()

		json.Unmarshal(rawBody, &newPassenger)
		//err := validate.Struct(newPassenger)

		// if err != nil {
		// 	fmt.Fprintf(w, "Error while reading the Body. \n [%+v]", err)
		// 	return
		// }

		if p[len(p)-1] != "" {
			id := p[len(p)-1]
			for index, passenger := range passengerData {
				if passenger.Id == id {
					passenger.mergePassenger(&newPassenger)
					passengerData[index] = passenger
					break
				}
			}
		} else {
			fmt.Fprintf(w, "No Id Passed")
		}

		if reflect.DeepEqual(newPassenger, (Passenger{})) {
			fmt.Fprintf(w, "No Record Found")
		} else {
			pres, _ := json.Marshal(newPassenger)
			fmt.Fprintf(w, "%s", string(pres))
		}
	}
}

func uploadFiles(w http.ResponseWriter, r *http.Request) {
	method := r.Method

	if method == http.MethodPost {
		var filenames []string
		reader, err := r.MultipartReader()
		if err != nil {
			fmt.Fprintf(w, "Something went wrong while getting reader")
			return
		}

		for {
			part, err := reader.NextPart()
			if err == io.EOF {
				fmt.Fprintf(w, "EOF reached")
				break
			}

			fmt.Println("Filename : ", part.FileName())
			data, err := ioutil.ReadAll(part)
			if err != nil {
				fmt.Printf("Error while reading %s Err [%+v]", part.FileName(), err)
			}

			err = ioutil.WriteFile(part.FileName(), data, 0644)
			if err != nil {
				fmt.Printf("Error while writing data to file. Err [%+v]", err)
			} else {
				filenames = append(filenames, part.FileName())
			}
			part.Close()
		}

		if len(filenames) > 0 {
			loadData(filenames)
		}

		fmt.Fprintf(w, "Successfully Uploaded The Files.")

	} else {
		fmt.Fprintf(w, "Invalid method to upload files")
	}
}

func main() {
	a := HelloHandler{}
	srv := &http.Server{
		Addr: ":8080",
	}

	cwd, _ := os.Getwd()
	rootPath := filepath.Dir(filepath.Dir(cwd))
	credPath := path.Join(rootPath, "creds")

	// Routes
	http.HandleFunc("/airlines/", GetAirlines)
	http.HandleFunc("/passenger/", GetPassengers)
	http.HandleFunc("/upload", uploadFiles)
	http.Handle("/hello/", http.StripPrefix("/hello/", a))
	http.Handle("/", http.FileServer(http.Dir("template")))
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	//

	fmt.Println("Listening on Port 8080....")

	srv.ListenAndServeTLS(
		fmt.Sprintf("%s/cert.pem", credPath),
		fmt.Sprintf("%s/key.pem", credPath),
	)
}
