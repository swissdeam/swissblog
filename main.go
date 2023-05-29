package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Magazine struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Tel   string `json:"tel"`
	Price int    `json:"price`
	P     int
	P2    int
}

type Request struct {
	Id             int    `json:"id"`
	Name           string `json:"name`
	Content        string `json:"content"`
	Timestamp      string `json:"timestamp"`
	Theme_id       int
	Advertisers_id int
}

type Contract_table struct {
	Id           int    `json:"id"`
	Timestamp    string `json:"timestamp"`
	Duration     int    `json:"duration"`
	Request_id   int
	Magazines_id int
	Amount       int
	Editorial    int
}

type Contract_page struct {
	Id                int    `json:"id"`
	Timestamp         string `json:"timestamp"`
	Duration          int    `json:"duration"`
	Request_Name      string
	Request_Content   string
	Request_timestamp string
	Theme_name        string
	Advertisers_id    int
	Advertisers_name  string
	Advertisers_Tel   string
	Advertisers_email string
	Magazines_name    string
	Magazine_tel      string
	Magazine_price    int
	Amount            int
	Editorial         int
}

type Advertiser struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Tel   string `json:"tel"`
	Email string `json:"email"`
}

type Themes struct {
	Id   int
	Name string
}

type Magazines_has_themes struct {
	magazines_id int
	themes_id    int
}

func home_page(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("src/home_page.html")
	tmpl.Execute(w, "")
}

func advertisers_page(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:bitard671K-On@tcp(127.0.0.1:3306)/mydb")
	if err != nil {
		panic(err)
	}

	defer db.Close()
	res, err := db.Query("SELECT * FROM `advertisers`")
	if err != nil {
		panic(err)
	}
	advertisers := []Advertiser{}
	for res.Next() {
		var advertiser Advertiser
		err = res.Scan(&advertiser.Id, &advertiser.Name, &advertiser.Tel, &advertiser.Email)
		if err != nil {
			panic(err)
		}
		advertisers = append(advertisers, advertiser)
	}
	tmpl, _ := template.ParseFiles("src/advertisers_page.html")
	tmpl.Execute(w, advertisers)
}

func magazines_page(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:bitard671K-On@tcp(127.0.0.1:3306)/mydb")
	if err != nil {
		panic(err)
	}

	defer db.Close()
	res, err := db.Query("SELECT * FROM `magazines`")
	if err != nil {
		panic(err)
	}
	magazines := []Magazine{}
	for res.Next() {
		var magazine Magazine
		err = res.Scan(&magazine.Id, &magazine.Name, &magazine.Email, &magazine.Tel, &magazine.Price)
		if err != nil {
			panic(err)
		}
		magazines = append(magazines, magazine)
	}
	tmpl, _ := template.ParseFiles("src/magazines_page.html")
	tmpl.Execute(w, magazines)
}

func requests_page(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:bitard671K-On@tcp(127.0.0.1:3306)/mydb")
	if err != nil {
		panic(err)
	}

	defer db.Close()
	res, err := db.Query("SELECT * FROM `requests`")
	if err != nil {
		panic(err)
	}
	feed := []Request{}
	for res.Next() {
		var request Request
		err = res.Scan(&request.Id, &request.Advertisers_id, &request.Theme_id, &request.Name, &request.Content, &request.Timestamp)
		if err != nil {
			panic(err)
		}
		feed = append(feed, request)
	}

	tmpl, _ := template.ParseFiles("src/requests_page.html")
	tmpl.Execute(w, feed)
}

func contracts_page(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:bitard671K-On@tcp(127.0.0.1:3306)/mydb")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	show := []Contract_page{}
	var list Contract_page

	res, err := db.Query("SELECT * FROM `contracts`")
	if err != nil {
		panic(err)
	}
	storage := []Contract_table{}
	for res.Next() {
		var contract Contract_table
		err = res.Scan(&contract.Id, &contract.Request_id, &contract.Magazines_id, &contract.Amount, &contract.Duration, &contract.Timestamp, &contract.Editorial)
		if err != nil {
			panic(err)
		}
		storage = append(storage, contract)
	}

	for _, contract := range storage {

		list.Duration = contract.Duration
		list.Id = contract.Id
		list.Timestamp = contract.Timestamp
		list.Amount = contract.Amount
		list.Editorial = contract.Editorial

		db, err := sql.Open("mysql", "root:bitard671K-On@tcp(127.0.0.1:3306)/mydb")
		if err != nil {
			panic(err)
		}

		defer db.Close()

		req_res, err := db.Query(fmt.Sprintf("SELECT * FROM `requests` WHERE id = %d", contract.Request_id))
		if err != nil {
			panic(err)
		}
		for req_res.Next() {
			var request Request
			err = req_res.Scan(&request.Id, &request.Advertisers_id, &request.Theme_id, &request.Name, &request.Content, &request.Timestamp)
			if err != nil {
				panic(err)
			}
			list.Request_Name = request.Name
			list.Request_Content = request.Content
			list.Request_timestamp = request.Timestamp
			list.Advertisers_id = request.Advertisers_id

			ad_res, err := db.Query(fmt.Sprintf("SELECT * FROM `advertisers` WHERE id = %d", list.Advertisers_id))
			if err != nil {
				panic(err)
			}

			for ad_res.Next() {
				var advertiser Advertiser

				err = ad_res.Scan(&advertiser.Id, &advertiser.Name, &advertiser.Tel, &advertiser.Email)
				if err != nil {
					panic(err)
				}
				list.Advertisers_Tel = advertiser.Tel
				list.Advertisers_name = advertiser.Name
				list.Advertisers_email = advertiser.Email

			}

			th_res, err := db.Query(fmt.Sprintf("SELECT * FROM `themes` WHERE id = %d", request.Theme_id))
			if err != nil {
				panic(err)
			}

			for th_res.Next() {
				var theme Themes

				err = th_res.Scan(&theme.Id, &theme.Name)
				if err != nil {
					panic(err)
				}
				list.Theme_name = theme.Name

			}
		}

		// ad_res, err := db.Query(fmt.Sprintf("SELECT * FROM `advertisers` WHERE id = %d", contract.Advertisers_id))
		// log.Println("Advertisers id from list", list.Advertisers_id)
		// log.Println("Advertisers id from contract", contract.Advertisers_id)

		mag_res, err := db.Query(fmt.Sprintf("SELECT * FROM `magazines` WHERE id = %d", contract.Magazines_id))
		if err != nil {
			panic(err)
		}
		for mag_res.Next() {

			var magazine Magazine
			err = mag_res.Scan(&magazine.Id, &magazine.Name, &magazine.Tel, &magazine.Email, &magazine.Price)
			if err != nil {
				panic(err)
			}
			list.Magazines_name = magazine.Name
			list.Magazine_tel = magazine.Tel
			list.Magazine_price = magazine.Price
		}

		show = append(show, list)
	}

	tmpl, _ := template.ParseFiles("src/contracts_page.html")
	tmpl.Execute(w, show)
}

func new_requests_page(w http.ResponseWriter, r *http.Request) {
	var advertiser_id int
	db, err := sql.Open("mysql", "root:bitard671K-On@tcp(127.0.0.1:3306)/mydb")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	res, err := db.Query("SELECT `id` FROM `advertisers` ORDER BY id DESC LIMIT 1")
	if err != nil {
		panic(err)
	}

	for res.Next() {
		err = res.Scan(&advertiser_id)
		if err != nil {
			panic(err)
		}
	}

	log.Println("advertiser_id from table", advertiser_id)

	res1, err := db.Query("SELECT * FROM `themes`")
	if err != nil {
		panic(err)
	}
	themes_names := []Themes{}
	for res1.Next() {
		var theme Themes
		err = res1.Scan(&theme.Id, &theme.Name)
		if err != nil {
			panic(err)
		}
		themes_names = append(themes_names, theme)
	}
	// log.Println("list of themes from table", themes_names)
	log.Println("request", r)
	if r.Method == "POST" {
		log.Println("после пост")
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
		}
		var new_request Request
		new_request.Name = r.FormValue("name")
		log.Println(new_request.Name)
		new_request.Content = r.FormValue("message")
		log.Println(new_request.Content)
		new_request.Advertisers_id = advertiser_id
		log.Println(new_request.Advertisers_id)
		temp := time.Now()
		new_request.Timestamp = temp.Format("2006-01-02 15:04:05")
		for i, abc := range themes_names {
			log.Println("в цикле", i, abc, new_request.Theme_id)
			if abc.Name == r.FormValue("theme") {
				new_request.Theme_id = i + 1
			}
			log.Println("после условия цикле", i, abc, new_request.Theme_id)
		}
		log.Println("прокинутые данные", new_request)
		db, err := sql.Open("mysql", "root:bitard671K-On@tcp(127.0.0.1:3306)/mydb")
		if err != nil {
			panic(err)
		}
		defer db.Close()

		insert, err := db.Query(fmt.Sprintf("INSERT INTO `requests` (`name`, `content`, `timestamp`,`themes_id`, `advertisers_id`)  VALUES( '%s','%s','%s','%d','%d')", new_request.Name, new_request.Content, new_request.Timestamp, new_request.Theme_id, new_request.Advertisers_id))
		if err != nil {
			panic(err)
		}
		defer insert.Close()
		defer http.Redirect(w, r, "/requests/", 301)

	} else {
		log.Println("элс")
		tmpl, _ := template.ParseFiles("src/new_request_page.html")
		log.Println("после элс")
		tmpl.Execute(w, themes_names)
		log.Println("после после элс")
	}
}

func new_advertiser_page(w http.ResponseWriter, r *http.Request) {
	log.Println(" nachalo ")
	if r.Method == "POST" {
		log.Println("mux")
		err := r.ParseForm()
		if err != nil {
			panic(err)
		}
		var new_advertiser Advertiser
		new_advertiser.Name = r.FormValue("name")
		new_advertiser.Tel = r.FormValue("phone-number")
		new_advertiser.Email = r.FormValue("email")
		log.Println(new_advertiser, new_advertiser.Name, r.FormValue("name"))

		db, err := sql.Open("mysql", "root:bitard671K-On@tcp(127.0.0.1:3306)/mydb")
		if err != nil {
			panic(err)
		}
		log.Println("do bazi")
		defer db.Close()
		log.Println("posle bazi")
		insert, err := db.Query(fmt.Sprintf("INSERT INTO `advertisers` (`name`, `tel`, `email`) VALUES( '%s','%s','%s')", new_advertiser.Name, new_advertiser.Tel, new_advertiser.Email))
		if err != nil {
			panic(err)
		}
		log.Println("do insert")
		insert.Close()
		log.Println("testing log")
		// tmpl, err := template.ParseFiles("src/new_request_page.html")
		// if err != nil {
		// 	panic(err)
		// }
		// tmpl.Execute(w, "")
		defer http.Redirect(w, r, "/newrequest/", 301)

	} else {
		tmpl, _ := template.ParseFiles("src/new_advertiser_page.html")
		tmpl.Execute(w, "")
	}
}

func choosen_request_page(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["request_id"]
	db, err := sql.Open("mysql", "root:bitard671K-On@tcp(127.0.0.1:3306)/mydb")
	if err != nil {
		panic(err)
	}

	//convert to integer
	int_id, _ := strconv.Atoi(id)

	defer db.Close()
	res, err := db.Query(fmt.Sprintf("SELECT * FROM `requests` WHERE id = %d", int_id))
	if err != nil {
		panic(err)
	}

	var request Request
	for res.Next() {

		err = res.Scan(&request.Id, &request.Advertisers_id, &request.Theme_id, &request.Name, &request.Content, &request.Timestamp)
		if err != nil {
			panic(err)
		}

	}
	log.Println(request)
	tmpl, _ := template.ParseFiles("src/choosen_request_page.html")
	tmpl.Execute(w, request)
}

func relevant_page(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["theme_id"]
	int_request, _ := strconv.Atoi(vars["request_id"])
	db, err := sql.Open("mysql", "root:bitard671K-On@tcp(127.0.0.1:3306)/mydb")
	if err != nil {
		panic(err)
	}

	//convert to integer
	int_id, _ := strconv.Atoi(id)

	defer db.Close()

	res1, err := db.Query(fmt.Sprintf("SELECT * FROM `magazines_has_themes` WHERE themes_id = %d", int_id))
	if err != nil {
		panic(err)
	}

	indexes := []Magazines_has_themes{}
	for res1.Next() {
		var index Magazines_has_themes
		err = res1.Scan(&index.magazines_id, &index.themes_id)
		if err != nil {
			panic(err)
		}
		indexes = append(indexes, index)
	}
	magazines := []Magazine{}
	for _, index := range indexes {
		res, err := db.Query(fmt.Sprintf("SELECT * FROM `magazines` WHERE id = %d", index.magazines_id))
		if err != nil {
			panic(err)
		}
		for res.Next() {
			var magazine Magazine
			magazine.P = int_request
			magazine.P2 = int_id
			err = res.Scan(&magazine.Id, &magazine.Name, &magazine.Tel, &magazine.Email, &magazine.Price)
			if err != nil {
				panic(err)
			}
			magazines = append(magazines, magazine)
		}
	}

	tmpl, _ := template.ParseFiles("src/relevant.html")
	log.Println("данные", magazines)
	tmpl.Execute(w, magazines)
}

func new_contract_page(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	request_id, _ := strconv.Atoi(vars["request_id"])
	theme_id, _ := strconv.Atoi(vars["theme_id"])
	magazine_id, _ := strconv.Atoi(vars["magazine_id"])

	var request Request
	var magazine Magazine
	var theme Themes
	var advertiser Advertiser
	var show Contract_page

	db, err := sql.Open("mysql", "root:bitard671K-On@tcp(127.0.0.1:3306)/mydb")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	req_res, err := db.Query(fmt.Sprintf("SELECT * FROM `requests` WHERE id = %d", request_id))
	if err != nil {
		panic(err)
	}
	for req_res.Next() {

		err = req_res.Scan(&request.Id, &request.Advertisers_id, &request.Theme_id, &request.Name, &request.Content, &request.Timestamp)
		if err != nil {
			panic(err)
		}
	}

	mag_res, err := db.Query(fmt.Sprintf("SELECT * FROM `magazines` WHERE id = %d", magazine_id))
	if err != nil {
		panic(err)
	}
	for mag_res.Next() {

		err = mag_res.Scan(&magazine.Id, &magazine.Name, &magazine.Tel, &magazine.Email, &magazine.Price)
		if err != nil {
			panic(err)
		}
	}

	th_res, err := db.Query(fmt.Sprintf("SELECT * FROM `themes` WHERE id = %d", theme_id))
	if err != nil {
		panic(err)
	}
	for th_res.Next() {

		err = th_res.Scan(&theme.Id, &theme.Name)
		if err != nil {
			panic(err)
		}
	}

	ad_res, err := db.Query(fmt.Sprintf("SELECT * FROM `advertisers` where id = %d", request.Advertisers_id))
	if err != nil {
		panic(err)
	}

	for ad_res.Next() {

		err = ad_res.Scan(&advertiser.Id, &advertiser.Name, &advertiser.Tel, &advertiser.Email)
		if err != nil {
			panic(err)
		}
	}

	if r.Method == "POST" {
		log.Println("пост метод после подтверждения")
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
		}
		var new_contract_table Contract_table
		new_contract_table.Timestamp = time.Now().Format("2006-01-02 15:04:05")
		new_contract_table.Duration, _ = strconv.Atoi(r.FormValue("duration"))
		new_contract_table.Amount, _ = strconv.Atoi(r.FormValue("amount"))
		new_contract_table.Request_id = request.Id
		new_contract_table.Magazines_id = magazine.Id
		// show.Duration, _ = strconv.Atoi(r.FormValue("duration"))
		// show.Amount, _ = strconv.Atoi(r.FormValue("amount"))
		// check, _ := strconv.Atoi(r.FormValue("price"))
		// check2 := new_contract_table.Duration
		if new_contract_table.Duration != 0 {
			new_contract_table.Amount = new_contract_table.Amount * new_contract_table.Duration
		}

		log.Println("price", r.FormValue("price"), "amount", new_contract_table.Amount, "duration", new_contract_table.Duration)

		log.Println("price", r.FormValue("price"), "amount", new_contract_table.Amount, "duration", new_contract_table.Duration)

		if r.FormValue("editorial") == "on" {
			new_contract_table.Editorial = 1
			new_contract_table.Amount += 2000
		} else {
			new_contract_table.Editorial = 0
		}

		db, err := sql.Open("mysql", "root:bitard671K-On@tcp(127.0.0.1:3306)/mydb")
		if err != nil {
			panic(err)
		}

		defer db.Close()

		insert, err := db.Query(fmt.Sprintf("INSERT INTO `contracts` (`timestamp`, `duration`, `requests_id`,`magazines_id`, `amount`, `editorial`) VALUES( '%s','%d','%d','%d','%d','%d')", new_contract_table.Timestamp, new_contract_table.Duration, new_contract_table.Request_id, new_contract_table.Magazines_id, new_contract_table.Amount, new_contract_table.Editorial))
		if err != nil {
			panic(err)
		}
		defer insert.Close()
		log.Println("Insert")
		defer http.Redirect(w, r, "/contracts/", 301)

	} else {

		show.Timestamp = time.Now().Format("2006-01-02 15:04:05")
		show.Request_Name = request.Name
		show.Request_Content = request.Content
		show.Request_timestamp = request.Timestamp
		show.Theme_name = theme.Name
		show.Advertisers_name = advertiser.Name
		show.Advertisers_Tel = advertiser.Tel
		show.Advertisers_email = advertiser.Email
		show.Magazines_name = magazine.Name
		show.Magazine_tel = magazine.Tel
		show.Magazine_price = magazine.Price

		tmpl, _ := template.ParseFiles("src/new_contract_page.html")
		tmpl.Execute(w, show)
	}

}

func handleRequest() {
	router := mux.NewRouter()
	router.HandleFunc("/", home_page)
	router.HandleFunc("/requests/", requests_page)
	router.HandleFunc("/advertisers/", advertisers_page)
	router.HandleFunc("/magazines/", magazines_page)
	router.HandleFunc("/contracts/", contracts_page)
	router.HandleFunc("/newrequest/", new_requests_page)
	router.HandleFunc("/newadvertiser/", new_advertiser_page)
	router.HandleFunc("/requests/{request_id:[0-9]+}", choosen_request_page)
	router.HandleFunc("/requests/{request_id:[0-9]+}/{theme_id:[0-9]+}", relevant_page)
	router.HandleFunc("/requests/{request_id:[0-9]+}/{theme_id:[0-9]+}/{magazine_id:[0-9]+}", new_contract_page)
	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}

func main() {
	handleRequest()
}
