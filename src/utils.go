package src

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
)

func logMsg(msg string, v ...any) {
  fileres := ""
  pc := make([]uintptr, 10)
  n := runtime.Callers(3, pc)
  var frames *runtime.Frames
  if n == 0 {
    goto fileend
  }
  pc = pc[:n]
  frames = runtime.CallersFrames(pc)
  for {
    frame, more := frames.Next()
    if strings.Contains(frame.File, "asm_amd64") { break }
    fileres += fmt.Sprintf("%v:%v -> ", frame.File, frame.Line)
    if !more { break }
  }

  fileend:
  res := ""
  for _, e := range v {
    res += fmt.Sprintf("%v %T, ", e, e)
  }
  fmt.Printf(
    "%v %v%v: %v\n\n",
    strings.Split(time.Now().String(), ".")[0],
    fileres,
    msg,
    res,
  )
}

func LogError(v ...any) {
  logMsg("ERROR", v...)
}

func LogInfo(v ...any) {
  logMsg("INFO", v...)
}

func NewRecord(
  dao *daos.Dao,
  collname string,
) (*models.Record, error) {
  collection, err := dao.FindCollectionByNameOrId(collname)
  if err != nil { return nil, err }
  return models.NewRecord(collection), nil
}

func UpdateField(
  dao *daos.Dao,
  record *models.Record,
  field string,
  value any,
) error {
  record.Set(field, value)
  return dao.SaveRecord(record)
}

func StoreData(
  dao *daos.Dao,
  owner string,
  ttype string,
  data string,
) error {
  datarec, err := dao.FindFirstRecordByFilter(
    "data",
    "owner = {:owner} && type = {:type}",
    dbx.Params{"owner": owner, "type": ttype},
  )
  if err == sql.ErrNoRows {
    datarec, err = NewRecord(dao, "data")
    if err != nil { return err }

    datarec.Set("owner", owner)
    datarec.Set("type", ttype)
  } else if err != nil { return err }

  datarec.Set("data", data)

  return dao.SaveRecord(datarec)
}

func SendNotif(
  user *models.Record,
  body string,
) error {
  s := &webpush.Subscription{}
  json.Unmarshal([]byte(user.GetString("vapid")), s)

  _, err := webpush.SendNotification(
    []byte(body), s, 
    &webpush.Options{
      Subscriber: "example@example.com",
      VAPIDPublicKey: os.Getenv("VAPID_PUBLIC"),
      VAPIDPrivateKey: os.Getenv("VAPID_PRIVATE"),
      TTL: 30,
  })
  return err
}
