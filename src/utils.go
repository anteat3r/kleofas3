package src

import (
	"database/sql"
	"encoding/json"
	"os"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
)

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
