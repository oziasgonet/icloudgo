package main

import "fmt"
import "net/http"
import "github.com/satori/go.uuid"
import "strings"
import "encoding/json"
import "bytes"
import "io/ioutil"

const (
	BASE_URL  = "https://www.icloud.com"
	SETUP_URL = "https://setup.icloud.com/setup/ws/1"
	LOGIN_URL = SETUP_URL + "/login"
)

var (
	clientBuildNumber = "v1"
	clientId          = strings.ToUpper(uuid.NewV1().String())
	contactsUrl       string
	dsid              string
	WEBAUTH_TOKEN     string
	WEBAUTH_USER      string
)

// struct for json parsing
type Info struct {
	DsInfo struct {
		Dsid string `json:"dsid"`
	} `json:"dsInfo"`
	Webservices struct {
		Contacts struct {
			Url string `json:"url"`
		} `'json:"contacts"`
	} `json:"webservices"`
}

func login(apple_id, password string) {
	json_str := `{"apple_id":"` + apple_id + `","password":"` + password +
		`","extended_login":"false"}`
	b := []byte(json_str)
	url := LOGIN_URL + "?clientBuildNumber=" + clientBuildNumber +
		"&clientId=" + clientId
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	client := &http.Client{}
	req.Header.Set("Host", "setup.icloud.com")
	req.Header.Set("Origin", BASE_URL)

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()
	reqBytes, _ := ioutil.ReadAll(resp.Body)
	WEBAUTH_TOKEN = resp.Header["Set-Cookie"][5]
	WEBAUTH_USER = resp.Header["Set-Cookie"][6]
	var f Info
	json.Unmarshal(reqBytes, &f)
	contactsUrl = f.Webservices.Contacts.Url
	// Remove the port from the url
	contactsUrl = contactsUrl[:len(contactsUrl)-4]
	dsid = f.DsInfo.Dsid
}

func getContacts() {
	client := &http.Client{}
	url := contactsUrl + "/co/startup" + "?clientBuildNumber=" +
		clientBuildNumber + "&clientId=" + clientId + "&clientVersion=2.1&" +
		"dsid=" + dsid + "&locale=en_US&order=last%2Cfirst"
	resp, err := http.NewRequest("GET", url, nil)
	host := strings.Split(contactsUrl, "//")[1]

	resp.Header.Set("Host", host)
	resp.Header.Set("Origin", BASE_URL)
	resp.Header.Set("Cookie", WEBAUTH_TOKEN+";"+WEBAUTH_USER)

	res, err := client.Do(resp)
	if err != nil {
		fmt.Println(err)
		return
	}
	
	defer res.Body.Close()
	respBytes, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(respBytes))
}
