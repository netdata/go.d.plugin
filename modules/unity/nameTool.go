package unity

import (
	"time"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"io/ioutil"
)

//Session struct
type Session struct {
	server   string
	token    string
	http     http.Client
}

type LunQueryResult struct {
	Base    string    `json:"@base"`
	Updated time.Time `json:"updated"`
	Links   []struct {
		Rel  string `json:"rel"`
		Href string `json:"href"`
	} `json:"links"`
	Entries []struct {
		Content struct {
			Id   string `json:"id"`
			Name string `json:"name"`
		} `json:"content"`
	} `json:"entries"`
}

func URL(server string, URI string) string {
	return "https://" + server + URI
}

func EncodeCredentials(username string, password string) string {
	return base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
}

func NewSession(server string, username string, password string) (*Session, error) {
	if server == "" || username == "" || password == "" {
		return nil, errors.New("Missing server (Unity IP), username or password")
	}

	var httpClient http.Client
	cookieJar, _ := cookiejar.New(nil)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient = http.Client{Transport: tr, Jar: cookieJar} //(insecure)

	var req *http.Request

	req, _ = http.NewRequest("GET", URL(server, "/api/types/system/instances"), nil)

	req.Header.Set("Authorization", "Basic "+EncodeCredentials(username, password))
	req.Header.Set("X-EMC-REST-CLIENT", "true")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	token := resp.Header.Get("Emc-Csrf-Token")

	return &Session{server, token, httpClient}, nil
}

//Request purpose is to send a Rest API request to the Unity array.
func (session *Session) Request(method string, URI string, fields string, engineering bool, resp interface{}) error {
	if method == "" || URI == "" {
		return errors.New("Missing method or URI")
	}

	endpoint := URL(session.server, URI)

	var req *http.Request

	req, _ = http.NewRequest(method, endpoint, nil)

	if method == "POST" {
		req.Header.Set("Emc-Csrf-Token", session.token)
	}
	
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-EMC-REST-CLIENT", "true")

	a := req.URL.Query()
	if fields != "" {
		a.Add("fields", fields)
	}
	if engineering == true {
		a.Add("visibility", "Engineering")
		a.Add("per_page","99999")
		a.Add("page","1")
	}
	a.Add("compact", "true")

	req.URL.RawQuery = a.Encode()

	httpResp, err := session.http.Do(req)
	if err != nil {
		return err
	}

	defer httpResp.Body.Close()

	switch {
	case (httpResp.StatusCode == 200 || httpResp.StatusCode == 201 || httpResp.StatusCode == 202) && method == "GET":
		body, err := ioutil.ReadAll(httpResp.Body)
		if err != nil {
			return err
		}
		return json.Unmarshal(body, &resp)

	case httpResp.StatusCode == 204:
		return nil

	case httpResp.StatusCode == 422:
		return fmt.Errorf("HTTP status codes: %d, detail: %v", httpResp.StatusCode, httpResp.Body)

	default:
		return fmt.Errorf("HTTP status codes: %d", httpResp.StatusCode)
	}
}

func (session *Session) CloseSession() (err error) {
	err = session.Request("POST", "/api/types/loginSessionInfo/action/logout", "", false,nil)
	return err
}

func (u *Unity) getLuns() (map[string]map[string]string, error){
	var resp map[string]map[string]string
	resp = make(map[string]map[string]string)

	for _,server := range u.config.Servers{
		resp[server.Adress] = make(map[string]string)
		session, err := NewSession(server.Adress, u.config.Username, u.config.Password)

		if err!=nil{
			continue
		}

		fields := "id,name"
		var raw *LunQueryResult
		err = session.Request("GET", "/api/types/lun/instances", fields, true, &raw)

		if err!=nil{
			continue
		}
		
		session.CloseSession()

		for _,entry := range raw.Entries{
			resp[server.Adress][entry.Content.Id]=entry.Content.Name
		}
	}	
	return resp,nil
}
