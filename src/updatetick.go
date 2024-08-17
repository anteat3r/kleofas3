package src

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/types"
  log "github.com/anteat3r/golog"
)

func Tick(dao *daos.Dao) func() {
  return func () {
    err := TimeTableUpdate(dao)
    if err != nil { log.LogError(err) }
    users, err := dao.FindRecordsByFilter(
      COLLECTION_USERS,
      "valid_login = true && next_reload > @now",
      "-next_reload",
      0, 0,
    )
    if err != nil { log.LogError(err); return }
    for _, usr := range users {
      UserReload(dao, usr)
      dur := usr.GetInt("reload_interval")
      next_reload := time.Now().Add(
        time.Minute * time.Duration(dur),
      )
      err := UpdateField(dao, usr, "next_reload", next_reload)
      if err != nil { log.LogError(err) }
    }
  }
}

func UserReload(dao *daos.Dao, user *models.Record) error {
  nextweek := time.Now().Add(
    time.Hour * time.Duration(24 * 7)).Format("2006-01-02")
  endps := map[string]string{
    "marks": "marks",
    "events": "events/my",
    "timetable": "timetable/actual",
    "absence": "absence/student",
    "nexttimetable": "timetable/actual?date=" + nextweek,
  }
  uendps := user.GetStringSlice("endpoints")
  for _, uep := range uendps {
    ep := endps[uep]
    _, res, err := BakaQuery(dao, user, "GET", ep, "")
    if err != nil { log.LogError(err); continue }
    err = StoreData(dao, user.Id, uep, res)
    if err != nil { log.LogError(err) }
  }
  return nil
}

func TimeTableUpdate(dao *daos.Dao) error {
  records, err := dao.FindRecordsByFilter(
    "sources",
    `id != ""`,
    "+last_fetched",
    1, 0,
  )
  if err != nil { return err }
  if len(records) == 0 { return NoSourcesError{} }
  record := records[0]

  users, err := dao.FindRecordsByFilter(
    "users",
    "valid_login = true",
    "+last_used",
    1, 0,
  )
  if err != nil { return err }
  if len(users) == 0 { return NoUsersError{} }
  user := users[0]

  datares := ""

  if record.GetString("type") == "allevents" {
    stat, res, err := BakaQuery(dao, user, "GET", "events/all", "")
    if err != nil { return err }
    if stat != 200 {
      err := UpdateField(dao, user, "valid_login", false)
      if err != nil { return err }
      return ReqFailedError{}
    }
    datares = res

  } else {

    res, err := TimeTableQuery(
      dao, user, "actual",
      record.GetString("type"),
      record.GetString("name"),
    )
    if err != nil { return err }

    tmtb, err := ParseTimeTableWeb(res)
    if err != nil { return err }
    // LogInfo(tmtb)

    dataresb, err := json.Marshal(tmtb)
    if err != nil { return err }
    datares = string(dataresb)
  }

  err = UpdateField(
    dao,
    record,
    "last_fetched",
    types.NowDateTime(),
  )
  if err != nil { return err }

  err = UpdateField(
    dao,
    user,
    "last_used",
    types.NowDateTime(),
  )
  if err != nil { return err }

  err = StoreData(
    dao,
    record.GetString("name"),
    record.GetString("type"),
    datares,
  )
  if err != nil { return err }
  return nil
}

type NoSourcesError struct {}
func (e NoSourcesError) Error() string {
  return "no valid sources available"
}

type NoUsersError struct {}
func (e NoUsersError) Error() string {
  return "no valid users available"
}

type ReqFailedError struct {code int}
func (e ReqFailedError) Error() string {
  return fmt.Sprintf("req failed with code %v", e.code)
}
