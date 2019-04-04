ETL Process and API by Piyush Sharma


This folder contains all the files required to run and test the ETL process and the API


To get started:

1. Type go run main.go in the terminal to pull the JSON data from the url and insert into the database.
  
2. The database is now populated and ready to receive API requests.

3. To test the API, simply type go run api.go in the terminal and this will start the HTTP Server. You can access the site at 127.0.0.1:8080

4. The index page contains the instructions on which urls to use to get JSON responses.

Thank you for taking the time to view this project and I hope you liked it :)

Sincerely,

Piyush Sharma


*If the output is an error that says the buildingInfo cant be used, that means that the API url has changed and it can be fixed by visiting https://data.cityofnewyork.us/Housing-Development/Building-Footprints/nqwf-w8eh and selecting export and then SODA API, that will provide another link which can be pasted in the resp variable in main.go*
i.e. resp, err := http.Get("NEW_API_URL")

