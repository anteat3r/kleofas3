package src

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/types"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var DoNotRedirectClient = http.Client{
  CheckRedirect: func(
    req *http.Request, via []*http.Request,
  ) error {
    return http.ErrUseLastResponse
  },
}

func LoginUser(
  dao *daos.Dao,
  user *models.Record,
  username string,
  password string,
) error {
  req, err := http.NewRequest(
    "POST",
    "https://bakalari.gchd.cz/bakaweb/api/login",
    strings.NewReader("client_id=ANDR&grant_type=password&username=" + username + "&password=" + password),
  )
  if err != nil { return err }
  req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
  resp, err := http.DefaultClient.Do(req)
  if err != nil { return err }

  body, err := io.ReadAll(resp.Body)
  if err != nil { return err }

  res := BakaLoginResponse{}
  err = json.Unmarshal(body, &res)
  if err != nil { return err }

  user.Set(USERS_ACCESS_TOKEN, res.AccessToken)
  user.Set(USERS_REFRESH_TOKEN, res.RefreshToken)
  
  date, _ := types.ParseDateTime(time.Now().Add(
    time.Second * time.Duration(res.ExpiresIn)))
  user.Set(USERS_ACCESS_TOKEN_EXPIRES, date)

  err = dao.Save(user)
  if err != nil { return err }

  req2, err := http.NewRequest(
    "GET",
    "https://bakalari.gchd.cz/bakaweb/login",
    strings.NewReader(""),
  )
  if err != nil { return err }

  resp2, err := http.DefaultClient.Do(req2)
  if err != nil { return err }

  rawcookie := resp2.Header.Get("Set-Cookie")
  rwckl := strings.Split(rawcookie, ";")
  if len(rwckl) < 2 { return errors.New("invalid cookie 1") }
  cookie := rwckl[0]

  payl := "username=" + username + "&password=" + password + "&persistent=true&returnUrl="

  req3, err := http.NewRequest(
    "POST",
    "https://bakalari.gchd.cz/bakaweb/Login",
    strings.NewReader(payl),
  )
  if err != nil { return err }

  req3.Header.Add("Cookie", cookie)
  req3.Header.Add(
    "Content-Type",
    "application/x-www-form-urlencoded",
  )
  req3.Header.Add(
    "Accept",
    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/png,image/svg+xml,*/*;q=0.8",
  )
  resp3, err := DoNotRedirectClient.Do(req3)
  if err != nil { return err }

  rawcok2 := resp3.Header.Get("Set-Cookie")
  rwckl2 := strings.Split(rawcok2, ";")
  if len(rwckl2) < 2 { return errors.New("invalid cookie") }
  cookie2 := rwckl2[0]

  user.Set(USERS_COOKIE, cookie2)
  nextweek, err := types.ParseDateTime(
    time.Now().Add(time.Hour * time.Duration(24 * 6.5)),
  )
  if err != nil { return err }
  user.Set(USERS_COOKIE_EXPIRES, nextweek)

  err = dao.SaveRecord(user)

  return err
}

func WebQuery(
  dao *daos.Dao,
  user *models.Record,
  endpoint string,
) (string, error) {
  if user.GetString("cookie") == "" {
    err := UpdateField(dao, user, "valid_login", false)
    if err != nil { return "", err }

    err = CredsExpiredNotif(user)
    if err != nil { return "", err }

    return "", CookieNotPresentError{user: user.Id}
  }

  if user.GetDateTime("cookie_expires").Time().Before(time.Now()) {
    err := UpdateField(dao, user, "valid_login", false)
    if err != nil { return "", err }

    err = CredsExpiredNotif(user)
    if err != nil { return "", err }

    return "", CookieExpiredError{user: user.Id}
  }

  req, err := http.NewRequest("GET", "https://bakalari.gchd.cz/bakaweb/" + endpoint, nil)
  if err != nil { return "", err }
  req.Header.Set("cookie", user.GetString("cookie"))

  resp, err := http.DefaultClient.Do(req)
  if err != nil { return "", err }

  if resp.StatusCode != 200 { return "", BakaRequestFailedError{code: resp.StatusCode, user: user.Id} }
  res, err := io.ReadAll(resp.Body)
  if err != nil { return "", err }

  return string(res), nil
}

func TimeTableQuery(
  dao *daos.Dao,
  user *models.Record,
  time, ttype, name string,
) (string, error) {
  caser := cases.Title(language.English)
  return WebQuery(
    dao, user,
    "Timetable/Public/" + caser.String(time) + "/" + caser.String(ttype) + "/" + name,
  )
}
