//Author: Piyush Sharma
//Date: April 2, 2019
//This is the http server that contains endpoints that return json responses based on the request

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

//struct for to put the query data into
type buildingInfo struct {
	Base_BBL   string `json: "base_bbl"`
	Bin        string `json: "bin"`
	Cnstrct_yr string `json: "cnstrct_yr"`
	Doitt_id   string `json: "doitt_id"`
	Feat_code  string `json: "feat_code"`
	Geomsource string `json: "geomsource"`
	Lstmoddate string `json: "lstmoddate"`
	Lststatype string `json: lststatype`
	GroundElev string `json: "groundelev"`
	Heightroof string `json: "heightroof"`
	Mpluto_bbl string `json: "mpluto_bbl"`
	Shape_area string `json: "shape_area"`
	Shape_len  string `json: "shape_len"`
}

type pageData struct {
	PageTitle string
}

//for area average endpoint
type areaAverage struct {
	AreaAverage string `json: area_average`
}

func main() {

	//initialized a router and enforcing url redirect in case of a missing slash
	router := mux.NewRouter().StrictSlash(true)

	//index page
	router.HandleFunc("/", Index)

	//finds the row with that specific base bbl and prints the json response
	router.HandleFunc("/BuildingInfo/base_bbl/{base_bbl}", pullDataByBBL)

	//returns the latest average of the shape areas for all buildings in the database

	router.HandleFunc("/BuildingInfo/getLatestAreaAvearge/", pullShapeAreaAverage)

	//Starts the server and listens on port 8080. Port number can be changed.
	log.Fatal(http.ListenAndServe(":8080", router))
}

//Shows the user the basic information about the Server and how to use the API
func Index(w http.ResponseWriter, r *http.Request) {

	data := pageData{
		PageTitle: "BuildingInfo API by Piyush Sharma",
	}

	tmpl, err := template.ParseFiles("index.html")
	check(err)
	tmpl.Execute(w, data)
}

func pullShapeAreaAverage(w http.ResponseWriter, r *http.Request) {

	db, err := sql.Open("sqlite3", "building_data.db")
	check(err)
	defer db.Close()

	q, err := db.Query("SELECT shape_area FROM BuildingInfo")
	check(err)

	var area_arr []float64

	var sum float64

	//will be parsed into float upon retrieval
	var Shape_area string
	for q.Next() {

		err = q.Scan(&Shape_area)
		check(err)
		conv_shape_area, err := strconv.ParseFloat(Shape_area, 32)
		area_arr = append(area_arr, conv_shape_area)
		check(err)

	}

	for i := range area_arr {
		sum += area_arr[i]
	}

	//conversions
	var float_average = sum / float64(len(area_arr))
	var string_average = strconv.FormatFloat(float_average, 'f', '6', 64)

	//pack final response into the struct
	final_response := areaAverage{
		AreaAverage: string_average,
	}

	//displays response on the webpage
	json.NewEncoder(w).Encode(final_response)

}

func pullDataByBBL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	base_bbl_r := vars["base_bbl"]
	// var all_responses []buildingInfo

	db, err := sql.Open("sqlite3", "building_data.db")
	check(err)
	defer db.Close()

	//all vars for db query usage
	var (
		Base_BBL   string
		Bin        string
		Cnstrct_yr string
		Doitt_id   string
		Feat_code  string
		Geomsource string
		Lstmoddate string
		Lststatype string
		GroundElev string
		Heightroof string
		Mpluto_bbl string
		Shape_area string
		Shape_len  string
	)

	//Queries the database for a base_bbl that matches the request
	row := db.QueryRow("SELECT * FROM BuildingInfo WHERE base_bbl = ?", base_bbl_r)
	err = row.Scan(&Base_BBL, &Bin, &Cnstrct_yr, &Doitt_id, &Feat_code, &Geomsource, &Lstmoddate, &Lststatype, &GroundElev, &Heightroof, &Mpluto_bbl, &Shape_area, &Shape_len)

	resp := buildingInfo{

		Base_BBL:   Base_BBL,
		Bin:        Bin,
		Cnstrct_yr: Cnstrct_yr,
		Doitt_id:   Doitt_id,
		Feat_code:  Feat_code,
		Geomsource: Geomsource,
		Lstmoddate: Lstmoddate,
		Lststatype: Lststatype,
		GroundElev: GroundElev,
		Heightroof: Heightroof,
		Mpluto_bbl: Mpluto_bbl,
		Shape_area: Shape_area,
		Shape_len:  Shape_len,
	}

	//preparing json response
	json_resp, err := json.Marshal(resp)
	check(err)

	//Finally outputs the json data to the site
	fmt.Fprintf(w, string(json_resp))

}

//Basic error checking
func check(err error) {
	if err != nil {
		panic(err)
	}
}
