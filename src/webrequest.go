package src

import (
	"io"
	"net/http"
	"time"

	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

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
