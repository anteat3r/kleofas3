package src

import (
	"encoding/json"

	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/tools/types"
)

func Tick(dao *daos.Dao) func() {
  return func () {
    records, err := dao.FindRecordsByFilter(
      "sources",
      `id != ""`,
      "+last_fetched",
      1, 0,
    )
    if err != nil { LogError(err); return }
    if len(records) == 0 { LogInfo("no sources? :("); return }
    record := records[0]

    users, err := dao.FindRecordsByFilter(
      "users",
      "valid_login = true",
      "+last_used",
      1, 0,
    )
    if err != nil { LogError(err); return }
    if len(users) == 0 { LogInfo("no users? :("); return }
    user := users[0]

    datares := ""

    if record.GetString("type") == "allevents" {
      stat, res, err := BakaQuery(dao, user, "GET", "events/all", "")
      if err != nil { LogError(err); return }
      if stat != 200 {
        err := UpdateField(dao, user, "valid_login", false)
        if err != nil { LogError(err); return }
        return
      }
      datares = res

    } else {

      res, err := TimeTableQuery(
        dao, user, "actual",
        record.GetString("type"),
        record.GetString("name"),
      )
      if err != nil { LogError(err); return }

      tmtb, err := ParseTimeTableWeb(res)
      if err != nil { LogError(err); return }
      // LogInfo(tmtb)

      dataresb, err := json.Marshal(tmtb)
      if err != nil { LogError(err); return }
      datares = string(dataresb)
    }

    err = UpdateField(
      dao,
      record,
      "last_fetched",
      types.NowDateTime(),
    )
    if err != nil { LogError(err); return }

    err = UpdateField(
      dao,
      user,
      "last_used",
      types.NowDateTime(),
    )
    if err != nil { LogError(err); return }

    err = StoreData(
      dao,
      record.GetString("name"),
      record.GetString("type"),
      datares,
    )
    if err != nil { LogError(err) }
  }
}
