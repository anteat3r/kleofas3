package src

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/types"
)

func BakaQuery(
  dao *daos.Dao,
  user *models.Record,
  method, endpoint, body string,
) (int, string, error) {
  if user.GetString(USERS_ACCESS_TOKEN) == "" {
    user.Set(USERS_VALID_LOGIN, false)
    err := dao.Save(user)
    if err != nil { return 0, "", err }
    err = CredsExpiredNotif(user)
    if err != nil { return 0, "", err }
    return 0, "", TokenNotPresentError{user: user.Id}
  }

  attempts := 0

  var resp *http.Response
  var req *http.Request
  var bodybuffer *strings.Reader
  var err error

  if user.GetDateTime(USERS_ACCESS_TOKEN_EXPIRES).Time().Before(time.Now()) {
    goto try_refresh
  }

  try_access:
    bodybuffer = strings.NewReader(body)
    if method == "GET" {
      bodybuffer = nil
    }
    req, err = http.NewRequest(
      method,
      "https://bakalari.gchd.cz/bakaweb/api/3/" + endpoint,
      bodybuffer,
    )
    req.Header.Set("Authorization", "Bearer " + user.GetString(USERS_ACCESS_TOKEN))
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    if err != nil { return 0, "", err }
    resp, err = http.DefaultClient.Do(req)
    if err != nil { return 0, "", err }

    if resp.StatusCode != 401 {
      res, err := io.ReadAll(resp.Body)
      if err != nil { return 0, "", err }
      resp.Body.Close()
      return resp.StatusCode, string(res), nil
    }
  
  try_refresh:
    req2, err := http.NewRequest(
      "POST",
      "https://bakalari.gchd.cz/bakaweb/api/login",
      strings.NewReader("client_id=ANDR&grant_type=refresh_token&refresh_token=" +
                        user.GetString(USERS_REFRESH_TOKEN)),
    )
    req2.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    if err != nil { return 0, "", err }
    resp2, err := http.DefaultClient.Do(req)
    if err != nil { return 0, "", err }

    if resp2.StatusCode == 200 {
      res, err := io.ReadAll(resp2.Body)
      if err != nil { return 0, "", err }
      resp2.Body.Close()
      jres := bakaLoginResponse{}
      json.Unmarshal(res, &jres)

      user.Set(USERS_ACCESS_TOKEN, jres.AccessToken)
      user.Set(USERS_REFRESH_TOKEN, jres.RefreshToken)
      date, _ := types.ParseDateTime(time.Now().Add(time.Second * time.Duration(jres.ExpiresIn)))
      user.Set(USERS_ACCESS_TOKEN_EXPIRES, date)

      err = dao.Save(user)
      if err != nil { return 0, "", err }
      if attempts > 0 { return 0, "", BakaLoginFailedError{user: user.Id} }
      attempts++
      goto try_access
    }

  return 0, "", nil
}

type bakaLoginResponse struct {
  UserId string `json:"bak:UserId"`
  AccessToken string `json:"access_token"`
  TokenType string `json:"token_type"`
  ExpiresIn int `json:"expires_in"`
  Scope string `json:"scope"`
  RefreshToken string `json:"refresh_token"`
  ApiVersion string `json:"bak:ApiVersion"`
  AppVersion string `json:"bak:AppVersion"`
}
