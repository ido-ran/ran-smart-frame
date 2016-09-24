package smartframe

import (
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
    "google.golang.org/api/plus/v1"
    oldappengine "appengine"
    "google.golang.org/appengine"
    "google.golang.org/appengine/log"
    newuser "google.golang.org/appengine/user"
    "html/template"

    "net/http"
    "time"
    "appengine/datastore"
    "appengine/user"
    "fmt"
    "appengine/urlfetch"
    "io/ioutil"
    "strings"
    "encoding/json"
)

type UserInfo struct {
        UserID  string
        Email string
        DisplayName string
        GoogleAccessToken string
        LastGoogleAuthenticationTime time.Time
        FirstAuthenticationTime time.Time
}

type FollowUser struct {
        UserID string
        Email string
}

type MediaResponse struct {
  Media []MediaInfo
}

type MediaInfo struct {
  Type string
  URL string
  Timestamp string
}

var cached_templates = template.Must(template.ParseGlob("templates/*.html"))

var conf = &oauth2.Config{
    ClientID:     "53043632999-resi4cfbi53q4q6gplp46g757jnjb87d.apps.googleusercontent.com",       // Replace with correct ClientID
    ClientSecret: "IMkpURmmDD_7LYEtuuYzfWlH",   // Replace with correct ClientSecret
    //RedirectURL:  "https://ran-smart-frame.appspot.com/oauth2callback",
    RedirectURL:  "http://localhost:8080/oauth2callback",
    Scopes: []string{
        "https://www.googleapis.com/auth/userinfo.email",
        "https://picasaweb.google.com/data",
        "profile",
        "email",
    },
    Endpoint: google.Endpoint,
}

func init() {
    http.HandleFunc("/", handleRoot)
    http.HandleFunc("/authorize", handleAuthorize)
    http.HandleFunc("/oauth2callback", handleOAuth2Callback)
    http.HandleFunc("/photos", handleGetPhotos)

    http.HandleFunc("/login", handleConsoleLogin)
    http.HandleFunc("/follow", handleFollowUser)
    http.HandleFunc("/me", handleMe)
}

func handleMe(w http.ResponseWriter, r *http.Request) {
  ctx := appengine.NewContext(r)
  u, err := newuser.CurrentOAuth(ctx, "")
  if err != nil {
          http.Error(w, "OAuth Authorization header required", http.StatusUnauthorized)
          return
  }
  fmt.Fprintf(w, `%s`, u)
}

func handleGetPhotos(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  oldc := oldappengine.NewContext(r)

  u, err := user.CurrentOAuth(oldc, "https://picasaweb.google.com/data")
  if err != nil {
      http.Error(w, "OAuth Authorization header required " + err.Error(), http.StatusUnauthorized)
      return
  }
  if u == nil {
      fmt.Fprintf(w, `You are not logged in yet, please <a href='/authorize'>Login</a>`, u)
      return
  }

  authHeader := r.Header.Get("Authorization")
  if !strings.HasPrefix(authHeader, "Bearer") {
    http.Error(w, "Authentication header with bearer not found", http.StatusUnauthorized)
    return
  }

  AccessToken := strings.Split(authHeader, " ")[1] // Take the token

  var picasaUrl = "https://picasaweb.google.com/data/feed/api/user/113997888652562329648/albumid/6220969463722583649?alt=json&access_token=" + AccessToken
  urlfetchClient := urlfetch.Client(oldc)
  res, err := urlfetchClient.Get(picasaUrl)
  if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
  }
  jsonBuffer, err := ioutil.ReadAll(res.Body)
  res.Body.Close()
  if err != nil {
    http.Error(w, "Fail to get Picasa data", http.StatusInternalServerError)
  }

  // Read the Picasa response into struct for easy parsing
  var picasaRes PicasaResponse
  err = json.Unmarshal(jsonBuffer, &picasaRes)

  var mediaResp MediaResponse

  if err != nil {
    log.Errorf(c, "Prase Picasa response Error: %v", err)
    http.Error(w, "Fail to parse Picasa data", http.StatusUnauthorized)
  } else {
    for _,entry := range picasaRes.Feed.Entry {
      if entry.OriginalVideo.Type != "" {
        //fmt.Printf("video %s\n\n", entry.Title.T)
        // For now we ignore videos
      } else {
        for _,media := range entry.MediaGroup.MediaContent {
          mediaResp.Media = append(mediaResp.Media, MediaInfo{"photo", media.URL, entry.Timestamp.T})
          fmt.Printf("picture %s %s\n\n", media.Type, media.URL)

          // Break in case there is more than one media group for this photo
          break;
        }
      }
    }
  }

  b, err := json.Marshal(mediaResp)
  s := string(b)
  fmt.Fprintf(w, `%s`, s)
}

func handleFollowUser(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  oldc := oldappengine.NewContext(r)

  u, err := user.CurrentOAuth(oldc, "https://picasaweb.google.com/data")
  if err != nil {
      http.Error(w, "OAuth Authorization header required " + err.Error(), http.StatusUnauthorized)
      return
  }
  if u == nil {
      fmt.Fprintf(w, `You are not logged in yet, please <a href='/authorize'>Login</a>`, u)
      return
  }

  // Try to load the user from the database
  q := datastore.NewQuery("UserInfo").Filter("Email = ", u.Email)
  var userInfoResults []UserInfo
  userInfoKeys, err := q.GetAll(oldc, &userInfoResults)
  if err != nil {
      log.Errorf(c, "Fail to find user by email: %v", err)
  }

  if (len(userInfoResults) == 0) {
    // Get the user
    log.Infof(c, "The user was not found by email %v", u.Email)
  } else {
    userInfoKey := userInfoKeys[0]

    followUserID := r.FormValue("userid")
    followUserEmail := r.FormValue("email")

    newFollow := FollowUser {
      UserID: followUserID,
      Email: followUserEmail,
    }
    followKey := datastore.NewIncompleteKey(oldc, "FollowUser", userInfoKey)
    _, err = datastore.Put(oldc, followKey, &newFollow)
    if err != nil {
        log.Errorf(c, "Fail to create follow: %v", err)
    }

  }
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
  c := oldappengine.NewContext(r)
  u, err := user.CurrentOAuth(c, "https://picasaweb.google.com/data")
  if err != nil {
      http.Error(w, "OAuth Authorization header required " + err.Error(), http.StatusUnauthorized)
      return
  }
  if u == nil {
      fmt.Fprintf(w, `You are not logged in yet, please <a href='/authorize'>Login</a>`, u)
      return
  }

  fmt.Fprintf(w, `Welcome, %s!`, u)
}

func handleAuthorize(w http.ResponseWriter, r *http.Request) {
    cookie := &http.Cookie{Name:"auth_initiate", Value:"frame", Expires:time.Now().Add(356*24*time.Hour)}
    http.SetCookie(w, cookie)

    url := conf.AuthCodeURL("")
    http.Redirect(w, r, url, http.StatusFound)
}

func handleConsoleLogin(w http.ResponseWriter, r *http.Request) {
    cookie := &http.Cookie{Name:"auth_initiate", Value:"console", Expires:time.Now().Add(356*24*time.Hour)}
    http.SetCookie(w, cookie)

    url := conf.AuthCodeURL("")
    http.Redirect(w, r, url, http.StatusFound)
}

func handleOAuth2Callback(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    code := r.FormValue("code")
    tok, err := conf.Exchange(c, code)
    if err != nil {
        log.Errorf(c, "%v", err)
        http.Redirect(w, r, "app/authcomplete.html?error", http.StatusFound)
        return
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

    oldc := oldappengine.NewContext(r)

    userEmail := person.Emails[0].Value

    // Try to load the user from the database
    q := datastore.NewQuery("UserInfo").Filter("Email = ", userEmail)
    var userInfoResults []UserInfo
    userInfoKeys, err := q.GetAll(oldc, &userInfoResults)
    if err != nil {
        log.Errorf(c, "Fail to find user by email: %v", err)
    }

    var userInfo UserInfo
    var userInfoKey *datastore.Key
    if (len(userInfoResults) > 0) {
      // Get the user
      log.Infof(c, "Authorize exist user %v", userInfo)
      userInfo = userInfoResults[0]
      userInfoKey = userInfoKeys[0]
    } else {
      // Create new UserInfo
      log.Infof(c, "Authorize new user %v", userInfo)
      userInfoKey = datastore.NewIncompleteKey(oldc, "UserInfo", nil)
      userInfo = UserInfo{
        UserID: person.Id,
        Email: person.Emails[0].Value,
        DisplayName: person.DisplayName,
        FirstAuthenticationTime: time.Now(),
      }
    }

    userInfo.LastGoogleAuthenticationTime = time.Now()
    userInfo.GoogleAccessToken = tok.AccessToken

    _, err = datastore.Put(oldc, userInfoKey, &userInfo)
    if err != nil {
        log.Errorf(c, "Fail to update user: %v", err)
    }

    cookie := &http.Cookie{Name:"accessToken", Value:tok.AccessToken, Expires:time.Now().Add(356*24*time.Hour), HttpOnly:true}
    http.SetCookie(w, cookie)

    http.Redirect(w, r, "app/authcomplete.html?accesstoken=" + tok.AccessToken, http.StatusFound)
}
