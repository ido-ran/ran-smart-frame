package hello

import (
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
    "google.golang.org/api/plus/v1"
    oldappengine "appengine"
    "google.golang.org/appengine"
    "google.golang.org/appengine/log"
    "html/template"

    "net/http"
    "appengine/urlfetch"
    "io/ioutil"
)

var cached_templates = template.Must(template.ParseGlob("templates/*.html"))

var conf = &oauth2.Config{
    ClientID:     "53043632999-resi4cfbi53q4q6gplp46g757jnjb87d.apps.googleusercontent.com",       // Replace with correct ClientID
    ClientSecret: "IMkpURmmDD_7LYEtuuYzfWlH",   // Replace with correct ClientSecret
    RedirectURL:  "https://ran-smart-frame.appspot.com/oauth2callback",
    Scopes: []string{
        "https://picasaweb.google.com/data",
        "profile",
    },
    Endpoint: google.Endpoint,
}

func init() {
    http.HandleFunc("/", handleRoot)
    http.HandleFunc("/authorize", handleAuthorize)
    http.HandleFunc("/oauth2callback", handleOAuth2Callback)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
    err := cached_templates.ExecuteTemplate(w, "notAuthenticated.html", nil)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
    }
}

func handleAuthorize(w http.ResponseWriter, r *http.Request) {
    //c := appengine.NewContext(r)
    url := conf.AuthCodeURL("")
    http.Redirect(w, r, url, http.StatusFound)
}

func handleOAuth2Callback(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    code := r.FormValue("code")
    tok, err := conf.Exchange(c, code)
    if err != nil {
        log.Errorf(c, "%v", err)
    }
    client := conf.Client(c, tok)

    // PLUS SERVICE CLIENT
    pc, err := plus.New(client)
    if err != nil {
        log.Errorf(c, "An error occurred creating Plus client: %v", err)
    }
    person, err := pc.People.Get("me").Do()
    if err != nil {
        log.Errorf(c, "Person Error: %v", err)
    }
    log.Infof(c, "Name: %v", person.DisplayName)

    log.Infof(c, "Access Token: %v", tok.AccessToken);

    var picasaUrl = "https://picasaweb.google.com/data/feed/api/user/113997888652562329648?kind=photo&tag=parentview&access_token=" + tok.AccessToken + "&alt=json"
    oldc := oldappengine.NewContext(r)
    urlfetchClient := urlfetch.Client(oldc)
    res, err := urlfetchClient.Get(picasaUrl)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json, err := ioutil.ReadAll(res.Body)
  	res.Body.Close()
  	if err != nil {
  		log.Errorf(c, "fail to read all body");
  	}
    log.Infof(c, "HTTP GET returned %v", json)
}
