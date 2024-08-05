package src

import (
	"encoding/json"
	"io"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
)

type LoadMsg struct {
  DataType string `json:"datatype"`
  Name string `json:"name"`
  Body string `json:"body"`
}

func LoadEnpoint(dao *daos.Dao) echo.HandlerFunc {
  return func (c echo.Context) error {
    user := c.Get(apis.ContextAuthRecordKey).(*models.Record)
    body, err := io.ReadAll(c.Request().Body)
    c.Request().Body.Close()
    if err != nil { return err }

    data := &LoadMsg{}
    err = json.Unmarshal(body, data)
    if err != nil { return err }

    datarecs, err := dao.FindRecordsByFilter(
      COLLECTION_DATA,
      DATA_OWNER + " = {:owner} && " + DATA_TYPE + " = {:type}",
      "updated",
      1, 0,
      dbx.Params{DATA_OWNER: data.Name, DATA_TYPE: data.DataType},
    )
    if err != nil { return err }

    var datarec *models.Record
    if len(datarecs) == 0 {
      datarec, err = NewRecord(dao, COLLECTION_DATA)
      datarec.Set(DATA_OWNER, data.Name)
      datarec.Set(DATA_TYPE, data.DataType)
    } else {
      datarec = datarecs[0]
    }

    err = UpdateField(dao, datarec, DATA_DATA, data.Body)
    if err != nil { return err }

    sources, err := dao.FindRecordsByFilter(
      COLLECTION_SOURCES,
      SOURCES_NAME + " = {:name} && " + SOURCES_TYPE + " = {:type}",
      "updated",
      1, 0,
      dbx.Params{SOURCES_NAME: data.Name, SOURCES_TYPE: data.DataType},
    )
    if err != nil { return err }
    if len(sources) == 0 { return nil }
    source := sources[0]

    err = UpdateField(dao, source, SOURCES_LAST_FETCHED, time.Now())
    if err != nil { return err }

    err = UpdateField(dao, user, USERS_LAST_USED, time.Now())
    if err != nil { return err }

    return nil
  }
}
