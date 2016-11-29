package icloudgo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"net/http"
	"strings"
)

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

// Login to iCloud
func Login(apple_id, password string) error {
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
		return err
	}
	defer resp.Body.Close()
	return parseLoginResponse(resp)
}

// Parses the HTTP response from an iCloud login attempt
func parseLoginResponse(r *http.Response) error {
	reqBytes, _ := ioutil.ReadAll(r.Body)

	if len(r.Header["Set-Cookie"]) < 5 {
		return errors.New("Login failed: Unable to retrieve webauth cookies")
	}
	WEBAUTH_TOKEN = r.Header["Set-Cookie"][5]
	WEBAUTH_USER = r.Header["Set-Cookie"][6]

	var f Info
	json.Unmarshal(reqBytes, &f)

	contactsUrl = f.Webservices.Contacts.Url
	if contactsUrl == "" {
		return errors.New("Login failed: Unable to retrieve contacts url")
	}
	// Remove the port from the url
	contactsUrl = contactsUrl[:len(contactsUrl)-4]
	dsid = f.DsInfo.Dsid

	return nil
}

// Retrieve the contacts from iCloud
func GetContacts() interface{} {
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
		return nil
	}

	defer res.Body.Close()
	respBytes, _ := ioutil.ReadAll(res.Body)
	var f interface{}
	json.Unmarshal(respBytes, &f)
	return f
}
