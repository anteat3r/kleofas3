package src

import "fmt"

const (
  COLLECTION_USERS = "users"
  COLLECTION_DATA = "data"
  COLLECTION_SOURCES = "sources"

  USERS_VALID_LOGIN = "valid_login"
  USERS_LAST_USED = "last_used"
  USERS_ACCESS_TOKEN = "access_token"
  USERS_REFRESH_TOKEN = "refresh_token"
  USERS_ACCESS_TOKEN_EXPIRES = "access_token_expires"
  USERS_COOKIE = "cookie"
  USERS_COOKIE_EXPIRES = "cookie_expires"

  DATA_OWNER = "owner"
  DATA_TYPE = "type"
  DATA_DATA = "data"

  SOURCES_NAME = "name"
  SOURCES_TYPE = "type"
  SOURCES_LAST_FETCHED = "last_fetched"
  SOURCES_DETAIL = "detail"
  
  BASE_URL = "https://bakalari.gchd.cz/bakaweb/"
  TIMETABLE_PUBLIC = "TimeTable/Public"
)

type CookieNotPresentError struct { user string }
func (e CookieNotPresentError) Error() string {
  return fmt.Sprintf("cookie not present in db for user %v\n", e.user)
}

type CookieExpiredError struct { user string }
func (e CookieExpiredError) Error() string {
  return fmt.Sprintf("cookie has expired for user %v\n", e.user)
}

type TokenNotPresentError struct { user string }
func (e TokenNotPresentError) Error() string {
  return fmt.Sprintf("token not present in db for user %v\n", e.user)
}

type BakaLoginFailedError struct { user string }
func (e BakaLoginFailedError) Error() string {
  return fmt.Sprintf("baka login failed for user %v\n", e.user)
}

type BakaRequestFailedError struct { user string; code int }
func (e BakaRequestFailedError) Error() string {
  return fmt.Sprintf("invalid status code %v for user %v\n", e.code, e.user)
}
