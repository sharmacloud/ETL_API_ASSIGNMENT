//Author: Piyush Sharma
//Date: April 2, 2019
//This file gets the json data from the url and inserts it into the empty database
package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

//the json_response is unpacked into this structure
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

func main() {

	//make a get request to the url
	resp, err := http.Get("https://data.cityofnewyork.us/resource/5kmz-p3b7.json")

	//basic error checking
	check(err)

	//reads the response body
	dat, err := ioutil.ReadAll(resp.Body)
	check(err)

	//Used to store multiple building's data
	var result []buildingInfo

	//Unmarshals the read json
	stuff := json.Unmarshal(dat, &result)
	check(stuff)

	//inserts the result array into the database
	insertData(result)

}

//Simple error checker
func check(e error) {
	if e != nil {
		panic(e)
	}
}

//Takes an array of building's infos and populates the database with the information
func insertData(data []buildingInfo) {

	//requires database type and the filename/location if file is located elsewhere
	db, err := sql.Open("sqlite3", "building_data.db")
	check(err)

	defer db.Close()

	//CONVERT THE FIELDS INTO INTS,FLOATS ETC TO INSERT INTO THE DATABASE

	// CONVERT STRUCT BIN TO INT64
	for i := range data {
		_, err := strconv.ParseInt(data[i].Bin, 10, 0)
		check(err)
	}

	//CONVERT THE CONSTRCT_YR IN STRUCT TO INT
	for i := range data {
		_, err := strconv.ParseInt(data[i].Cnstrct_yr, 10, 0)
		check(err)
	}

	//CONVERT FEAT CODE IN STRUCT TO INT
	for i := range data {
		_, err := strconv.ParseInt(data[i].Feat_code, 10, 0)
		check(err)
	}

	//CONVERT Ground_elev IN STRUCT TO float
	for i := range data {
		_, err := strconv.ParseFloat(data[i].GroundElev, 32)
		check(err)
	}

	//CONVERT height_roof IN STRUCT TO float
	for i := range data {
		_, err := strconv.ParseFloat(data[i].Heightroof, 32)
		check(err)
	}

	//CONVERT shape_area IN STRUCT TO float
	for i := range data {
		_, err := strconv.ParseFloat(data[i].GroundElev, 32)
		check(err)
	}

	//CONVERT shape_len IN STRUCT TO float
	for i := range data {
		_, err := strconv.ParseFloat(data[i].Shape_len, 32)
		check(err)
	}

	//loops through the data and performs insertions into the BuildingInfo table
	for i := range data {
		_, err := db.Exec(
			"INSERT INTO BuildingInfo (base_bbl,bin,cnstrct_yr,doitt_id,feat_code,geom_source,ground_elev,height_roof,lst_mod_date,lst_statype,mpluto_bbl,shape_area,shape_len) VALUES(?, ?, ?, ? ,? ,? ,? ,? ,? ,? ,? ,? ,? )",
			data[i].Base_BBL,
			data[i].Bin,
			data[i].Cnstrct_yr,
			data[i].Doitt_id,
			data[i].Feat_code,
			data[i].Geomsource,
			data[i].GroundElev,
			data[i].Heightroof,
			data[i].Lstmoddate,
			data[i].Lststatype,
			data[i].Mpluto_bbl,
			data[i].Shape_area,
			data[i].Shape_len)

		//checking for any errors during the insertions
		check(err)
	}

}
