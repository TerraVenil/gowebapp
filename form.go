package main

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"os"

	"fmt"
)

const CookieName = "MTR"

const LoginRoute = "/login"
const HomeRoute = "/home"
const LogoutRoute = "/logout"

const LoginTemplate string = "www/pages/login"
const HomeTemplate string = "www/pages/home"
const NotFoundTemplate string = "www/pages/404"

const UserNameTest = "johndoe@gmail.com"
const PasswordTest = "123456"

type Login struct {
	UserId     int `json:"Userid"`
	Username   string
	Fname      string
	Lname      string
	Human      int
	Manager    int
	Officeid   int
	Officelug  string
	Officename string
}

type Customer struct {
	FirstName string `json:"Firstname"`
	LastName  string `json:"Lastname"`
	Address   string `json:"Address"`
	Phone     string `json:"Phone"`
	Cell      string `json:"Cell"`
}

type Log struct {
	Date      string
	Direction string
	CommType  string
}

type Communication struct {
	Login   Login
	CommLog []Log
}

type Client struct {
	UserName string
	Password string
}

var storageMutex sync.RWMutex
var sessionsStore map[string]Client
var csrfTokens []string

func authenticationHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(fmt.Sprintf("Authentication resource of %s", r.URL.Path))

		cookie, err := r.Cookie(CookieName)
		if err != nil || cookie == nil || cookie.Value == "" {
			log.Println("Cookie value is't present.")
			http.Redirect(w, r, LoginRoute, http.StatusFound)
		} else {
			log.Println("Cookie value is present.")
			storageMutex.RLock()
			sessionID, _ := url.QueryUnescape(cookie.Value)
			_, exists := sessionsStore[sessionID]
			storageMutex.RUnlock()
			if exists {
				next.ServeHTTP(w, r)
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
		}
	})
}

func tokenHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" && r.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
			storageMutex.RLock()
			csrfToken := r.FormValue("csrfToken")
			storageMutex.RUnlock()
			if contains(csrfTokens, csrfToken) {
				h.ServeHTTP(w, r)
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
		} else {
			h.ServeHTTP(w, r)
		}
	})
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	sessionsStore = make(map[string]Client)
	log.Println(fmt.Sprintf("Cout of sessions %d", len(sessionsStore)))

	csrfTokens = []string{}
	log.Println(fmt.Sprintf("Cout of csrfTokens %d", len(csrfTokens)))

	// Authentication access
	http.HandleFunc("/communication", GetCommunicationLogs) ///?cid=1234
	http.HandleFunc("/logout", LogoutHandler)
	http.Handle("/dashboard", authenticationHandler(http.HandlerFunc(DashboardHandler)))
	http.Handle("/home", authenticationHandler(http.HandlerFunc(HomeHandler)))

	// Anonymous access
	http.Handle("/login", tokenHandler(http.HandlerFunc(LoginHandler)))

	// Public access, should be hidden!!!
	http.Handle("/", http.FileServer(http.Dir("www")))

	http.ListenAndServe(":"+port, nil)
}

func GetCommunicationLogs(rw http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	pid := query.Get("pid")
	var communication *Communication
	if pid != "" {
		communication = &Communication{
			Login{37, "peter", "Peter", "Parker", 0, 0, 13, "myOffice", "MyOffice"},
			[]Log{Log{"2017-03-15T10:04:34.620415276-04:00", "Out", "SMS"}, Log{"2017-03-15T10:04:34.620415276-04:00", "Out", "SMS"}}}
	} else {
		communication = &Communication{}
	}

	js, err := json.Marshal(*communication)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(js)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		{
			log.Println(fmt.Sprintf("[GET] Get login template by %s.", LoginTemplate))
			renderTemplate(w, LoginTemplate)
		}
	case "POST":
		{
			log.Println("[POST] Parse login form.")

			err := r.ParseForm()
			if err != nil {
				log.Println(fmt.Sprintf("Parse form with erro %s.", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			userName := r.FormValue("username")
			password := r.FormValue("password")
			csrfToken := r.FormValue("csrfToken")

			log.Println(fmt.Sprintf("[POST] Username %s and password %s and CSRF token %s.", userName, password, csrfToken))

			if userName == UserNameTest && password == PasswordTest {
				storageMutex.Lock()
				csrfTokens = append(csrfTokens, csrfToken)
				storageMutex.Unlock()
				SessionStart(w, r, Client{UserName: userName, Password: password})
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("[GET] Call logout.")
	SessionClose(w, r)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(fmt.Sprintf("[GET] Get home template by path %s.", HomeTemplate))
	renderTemplate(w, HomeTemplate)
}

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("[GET] Get customer info.")

	customer := Customer{"Jone", "Land", "1234 Main St, Atlanta, GA 30236", "(770)222-6666", "223-559-3324"}
	js, err := json.Marshal(customer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func renderTemplate(w http.ResponseWriter, tmpl string) {
	hash := md5.New()
	io.WriteString(hash, strconv.FormatInt(time.Now().UnixNano(), 10))
	token := fmt.Sprintf("%x", hash.Sum(nil))
	storageMutex.Lock()
	csrfTokens = append(csrfTokens, token)
	storageMutex.Unlock()
	t, err := template.ParseFiles(tmpl+".html", "www/pages/minimalbase.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := struct{ CsrfToken string }{token}
	t.ExecuteTemplate(w, "base", data)
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func SessionStart(w http.ResponseWriter, r *http.Request, client Client) {
	sessionID := NewSessionId()
	storageMutex.Lock()
	sessionsStore[sessionID] = client
	storageMutex.Unlock()
	cookie := http.Cookie{Name: CookieName, Value: url.QueryEscape(sessionID), Path: "/", HttpOnly: true, MaxAge: int(365 * 24 * time.Hour)}
	http.SetCookie(w, &cookie)
}

func SessionClose(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie(CookieName)
	sessionID, _ := url.QueryUnescape(cookie.Value)
	storageMutex.Lock()
	delete(sessionsStore, sessionID)
	storageMutex.Unlock()
	cookieExpired := http.Cookie{Name: CookieName, Value: url.QueryEscape(sessionID), Path: "/", HttpOnly: true, Expires: time.Date(1970, time.November, 1, 1, 0, 0, 0, time.UTC), MaxAge: -1}
	http.SetCookie(w, &cookieExpired)
	http.Redirect(w, r, LoginRoute, http.StatusFound)
}

func NewSessionId() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func contains(source []string, value string) bool {
	for _, v := range source {
		if v == value {
			return true
		}
	}
	return false
}
