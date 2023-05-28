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
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Licence     string `json:"licence"`
	Tel         string `json:"tel"`
	Price       string `json:"price`
	Description string `json:"description"`
	Address     string `json:"address"`
	P           int
	P2          int
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
	Id                int    `json:"id"`
	Timestamp         string `json:"timestamp"`
	Deadline          string `json:"deadline"`
	Status            string `json:"status"`
	Employe_id        int    `json:"employe_id"`
	Request_id        int
	Request_themes_id int
	Advertisers_id    int
	Magazines_id      int
	Themes_id         int `json:"themes_id"`
	Amount            string
}

type Contract_page struct {
	Id                  int    `json:"id"`
	Timestamp           string `json:"timestamp"`
	Deadline            string `json:"deadline"`
	Status              string `json:"status"`
	Employe_Name        string
	Employe_tel         string
	Employe_Role        string
	Request_Name        string
	Request_Content     string
	Request_timestamp   string
	Theme_name          string
	Advertisers_id      int
	Advertisers_name    string
	Advertisers_Tel     string
	Advertisers_email   string
	Advertisers_address string
	Magazines_name      string
	Magazine_licence    string
	Magazine_tel        string
	Magazine_price      string
	Magazine_address    string
	Amount              string
}

type Advertiser struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Tel      string `json:"tel"`
	Email    string `json:"email"`
	Position string `json:"position"`
	Address  string `json:"address"`
}

type Employe struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Role       string `json:"role"`
	Birth_date string `json:"birth_date"`
	Tel        string `json:"tel"`
	Address    string `json:"address"`
}
type Themes struct {
	Id          int
	Name        string
	Description string
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
		err = res.Scan(&advertiser.Id, &advertiser.Name, &advertiser.Tel, &advertiser.Email, &advertiser.Position, &advertiser.Address)
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
		err = res.Scan(&magazine.Id, &magazine.Name, &magazine.Email, &magazine.Licence, &magazine.Tel, &magazine.Price, &magazine.Description, &magazine.Address)
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
		err = res.Scan(&request.Id, &request.Name, &request.Content, &request.Timestamp, &request.Theme_id, &request.Advertisers_id)
		if err != nil {
			panic(err)
		}
		feed = append(feed, request)
	}
	log.Println(feed)
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
		err = res.Scan(&contract.Id, &contract.Timestamp, &contract.Deadline, &contract.Status, &contract.Employe_id, &contract.Request_id, &contract.Request_themes_id, &contract.Advertisers_id, &contract.Magazines_id, &contract.Themes_id, &contract.Amount)
		if err != nil {
			panic(err)
		}
		storage = append(storage, contract)
	}

	for _, contract := range storage {

		list.Deadline = contract.Deadline
		list.Status = contract.Status
		list.Id = contract.Id
		list.Timestamp = contract.Timestamp
		list.Amount = contract.Amount

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
			err = req_res.Scan(&request.Id, &request.Name, &request.Content, &request.Timestamp, &request.Theme_id, &request.Advertisers_id)
			if err != nil {
				panic(err)
			}
			list.Request_Name = request.Name
			list.Request_Content = request.Content
			list.Request_timestamp = request.Timestamp
			list.Advertisers_id = request.Advertisers_id
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
			err = mag_res.Scan(&magazine.Id, &magazine.Name, &magazine.Email, &magazine.Licence, &magazine.Tel, &magazine.Price, &magazine.Description, &magazine.Address)
			if err != nil {
				panic(err)
			}
			list.Magazines_name = magazine.Name
			list.Magazine_licence = magazine.Licence
			list.Magazine_tel = magazine.Tel
			list.Magazine_price = magazine.Price
			list.Magazine_address = magazine.Address
		}

		em_res, err := db.Query("SELECT * FROM `employe` WHERE id = 1")
		if err != nil {
			panic(err)
		}
		for em_res.Next() {
			var employe Employe
			err = em_res.Scan(&employe.Id, &employe.Name, &employe.Role, &employe.Birth_date, &employe.Tel, &employe.Address)
			if err != nil {
				panic(err)
			}
			list.Employe_Name = employe.Name
			list.Employe_Role = employe.Role
			list.Employe_tel = employe.Tel

		}

		th_res, err := db.Query(fmt.Sprintf("SELECT * FROM `themes` WHERE id = %d", contract.Themes_id))
		if err != nil {
			panic(err)
		}
		for th_res.Next() {

			var theme Themes
			err = th_res.Scan(&theme.Id, &theme.Name, &theme.Description)
			if err != nil {
				panic(err)
			}

			list.Theme_name = theme.Name
		}

		ad_res, err := db.Query(fmt.Sprintf("SELECT * FROM `advertisers` WHERE id = %d", contract.Advertisers_id))
		if err != nil {
			panic(err)
		}

		for ad_res.Next() {
			var advertiser Advertiser

			err = ad_res.Scan(&advertiser.Id, &advertiser.Name, &advertiser.Tel, &advertiser.Email, &advertiser.Position, &advertiser.Address)
			if err != nil {
				panic(err)
			}
			list.Advertisers_Tel = advertiser.Tel
			list.Advertisers_name = advertiser.Name
			list.Advertisers_email = advertiser.Email
			list.Advertisers_address = advertiser.Address

		}

		show = append(show, list)
	}

	log.Println("вывод контрактов до")

	tmpl, _ := template.ParseFiles("src/contracts_page.html")
	tmpl.Execute(w, show)
	log.Println("вывод контрактов после")
}

func new_requests_page(w http.ResponseWriter, r *http.Request) {
	log.Println("вход в функцию")
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
		err = res1.Scan(&theme.Id, &theme.Name, &theme.Description)
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
		new_request.Theme_id = 1
		temp := time.Now()
		new_request.Timestamp = temp.Format("2006-01-02 15:04:05")
		// for i, abc := range themes_names {
		// 	log.Println("в цикле", i, abc, new_request.Theme_id)
		// 	if abc.Name == r.FormValue("theme") {
		// 		new_request.Theme_id = i

		// 	}
		// 	log.Println("после условия цикле", i, abc, new_request.Theme_id)
		// }
		log.Println("прокинутые данные", new_request)
		db, err := sql.Open("mysql", "root:bitard671K-On@tcp(127.0.0.1:3306)/mydb")
		if err != nil {
			panic(err)
		}
		log.Println("check1")
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
		tmpl.Execute(w, "")
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
		new_advertiser.Position = r.FormValue("position")
		new_advertiser.Address = r.FormValue("address")
		log.Println(new_advertiser, new_advertiser.Name, r.FormValue("name"))

		db, err := sql.Open("mysql", "root:bitard671K-On@tcp(127.0.0.1:3306)/mydb")
		if err != nil {
			panic(err)
		}
		log.Println("do bazi")
		defer db.Close()
		log.Println("posle bazi")
		insert, err := db.Query(fmt.Sprintf("INSERT INTO `advertisers` (`name`, `tel`, `email`, `position`, `address`) VALUES( '%s','%s','%s','%s','%s')", new_advertiser.Name, new_advertiser.Tel, new_advertiser.Email, new_advertiser.Position, new_advertiser.Address))
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

func employes_page(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:bitard671K-On@tcp(127.0.0.1:3306)/mydb")
	if err != nil {
		panic(err)
	}

	defer db.Close()
	res, err := db.Query("SELECT * FROM `employe`")
	if err != nil {
		panic(err)
	}
	staff := []Employe{}
	for res.Next() {
		var employe Employe
		err = res.Scan(&employe.Id, &employe.Name, &employe.Role, &employe.Birth_date, &employe.Tel, &employe.Address)
		if err != nil {
			panic(err)
		}
		staff = append(staff, employe)
	}

	tmpl, _ := template.ParseFiles("src/employes_page.html")
	tmpl.Execute(w, staff)
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

		err = res.Scan(&request.Id, &request.Name, &request.Content, &request.Timestamp, &request.Theme_id, &request.Advertisers_id)
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
			err = res.Scan(&magazine.Id, &magazine.Name, &magazine.Email, &magazine.Licence, &magazine.Tel, &magazine.Price, &magazine.Description, &magazine.Address)
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
	var employe Employe
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

		err = req_res.Scan(&request.Id, &request.Name, &request.Content, &request.Timestamp, &request.Theme_id, &request.Advertisers_id)
		if err != nil {
			panic(err)
		}
	}

	mag_res, err := db.Query(fmt.Sprintf("SELECT * FROM `magazines` WHERE id = %d", magazine_id))
	if err != nil {
		panic(err)
	}
	for mag_res.Next() {

		err = mag_res.Scan(&magazine.Id, &magazine.Name, &magazine.Email, &magazine.Licence, &magazine.Tel, &magazine.Price, &magazine.Description, &magazine.Address)
		if err != nil {
			panic(err)
		}
	}

	em_res, err := db.Query("SELECT * FROM `employe` WHERE id = 1")
	if err != nil {
		panic(err)
	}
	for em_res.Next() {

		err = em_res.Scan(&employe.Id, &employe.Name, &employe.Role, &employe.Birth_date, &employe.Tel, &employe.Address)
		if err != nil {
			panic(err)
		}
	}

	th_res, err := db.Query(fmt.Sprintf("SELECT * FROM `themes` WHERE id = %d", theme_id))
	if err != nil {
		panic(err)
	}
	for th_res.Next() {

		err = th_res.Scan(&theme.Id, &theme.Name, &theme.Description)
		if err != nil {
			panic(err)
		}
	}

	ad_res, err := db.Query(fmt.Sprintf("SELECT * FROM `advertisers` where id = %d", request.Advertisers_id))
	if err != nil {
		panic(err)
	}

	for ad_res.Next() {

		err = ad_res.Scan(&advertiser.Id, &advertiser.Name, &advertiser.Tel, &advertiser.Email, &advertiser.Position, &advertiser.Address)
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
		new_contract_table.Deadline = r.FormValue("deadline")
		new_contract_table.Status = "разрешено"
		new_contract_table.Amount = r.FormValue("amount")
		new_contract_table.Employe_id = employe.Id
		new_contract_table.Request_id = request.Id
		new_contract_table.Request_themes_id = request.Theme_id
		new_contract_table.Advertisers_id = advertiser.Id
		new_contract_table.Magazines_id = magazine.Id
		new_contract_table.Themes_id = theme.Id
		show.Deadline = r.FormValue("deadline")
		show.Amount = r.FormValue("amount")

		check, _ := strconv.Atoi(new_contract_table.Amount)
		check2, _ := strconv.Atoi(new_contract_table.Deadline)
		log.Println("price", r.FormValue("price"), "amount", new_contract_table.Amount, "deadline", new_contract_table.Deadline)
		new_contract_table.Amount = strconv.Itoa(check * check2)
		log.Println("price", r.FormValue("price"), "amount", new_contract_table.Amount, "deadline", new_contract_table.Deadline)
		db, err := sql.Open("mysql", "root:bitard671K-On@tcp(127.0.0.1:3306)/mydb")
		if err != nil {
			panic(err)
		}

		defer db.Close()

		insert, err := db.Query(fmt.Sprintf("INSERT INTO `contracts` (`timestamp`, `deadline`, `status`, `employe_id`, `requests_id`, `requests_themes_id`,`advertisers_id`,`magazines_id`,`themes_id`, `amount`) VALUES( '%s','%s','%s','%d','%d','%d','%d','%d','%d', '%s')", new_contract_table.Timestamp, new_contract_table.Deadline, new_contract_table.Status, new_contract_table.Employe_id, new_contract_table.Request_id, new_contract_table.Request_themes_id, new_contract_table.Advertisers_id, new_contract_table.Magazines_id, new_contract_table.Themes_id, new_contract_table.Amount))
		if err != nil {
			panic(err)
		}
		defer insert.Close()
		log.Println("Insert")
		defer http.Redirect(w, r, "/contracts/", 301)

	} else {

		show.Timestamp = time.Now().Format("2006-01-02 15:04:05")
		show.Employe_Name = employe.Name
		show.Employe_tel = employe.Tel
		show.Employe_Role = employe.Role
		show.Request_Name = request.Name
		show.Request_Content = request.Content
		show.Request_timestamp = request.Timestamp
		show.Theme_name = theme.Name
		show.Advertisers_name = advertiser.Name
		show.Advertisers_Tel = advertiser.Tel
		show.Advertisers_email = advertiser.Email
		show.Advertisers_address = advertiser.Address
		show.Magazines_name = magazine.Name
		show.Magazine_licence = magazine.Licence
		show.Magazine_tel = magazine.Tel
		show.Magazine_price = magazine.Price
		show.Magazine_address = magazine.Address

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
	router.HandleFunc("/employee/", employes_page)
	router.HandleFunc("/requests/{request_id:[0-9]+}", choosen_request_page)
	router.HandleFunc("/requests/{request_id:[0-9]+}/{theme_id:[0-9]+}", relevant_page)
	router.HandleFunc("/requests/{request_id:[0-9]+}/{theme_id:[0-9]+}/{magazine_id:[0-9]+}", new_contract_page)
	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}

func main() {
	handleRequest()
}
