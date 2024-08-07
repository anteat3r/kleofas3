package main

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/anteat3r/kleofas3/src"
	"github.com/labstack/echo/v5"

	"github.com/joho/godotenv"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/cron"
)

func main() {
  godotenv.Load()
  log.SetFlags(log.LstdFlags | log.Lshortfile)
  app := pocketbase.New()

  app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
    scheduler := cron.New()

    scheduler.MustAdd(
      "tick",
      // TODO: set to normal
      "* 7-20 * * *",
      src.Tick(app.Dao()),
    )

    scheduler.MustAdd(
      "srcsreload",
      "1 18 * * 1-5",
      src.ReloadSources(app.Dao()),
    )

    e.Router.GET(
      "/*",
      apis.StaticDirectoryHandler(os.DirFS("web"), false),
    )
    e.Router.GET(
      "/load",
      func(c echo.Context) error {
        src.ReloadSources(app.Dao())()
        return nil
      },
      // apis.RequireRecordAuth("users"),
    )
    e.Router.GET("/test/:user", func(c echo.Context) error {
      tm := time.Now()
      username := c.PathParam("user")
      user, _ := app.Dao().FindFirstRecordByData("users", "username", username)
      var wg sync.WaitGroup
      for name, endp := range map[string]string{
        "marks": "marks",
        "events": "events/my",
        "timetable": "timetable/actual",
        "absence": "absence/student",
        "nexttimetable": "timetable/actual?date=" +
          time.Now().Add(time.Hour * time.Duration(24 * 7)).Format("2006-01-02"),
      } {
        wg.Add(1)
        _, res, err := src.BakaQuery(app.Dao(), user, "GET", endp, "")
        if err != nil { continue }
        src.StoreData(app.Dao(), user.Username(), name, res)
        // src.LogInfo(sc, res, err)
        wg.Done()
      }
      wg.Wait()
      src.LogInfo(time.Now().Sub(tm))
      return c.String(200, "")
    })

    scheduler.Start()

    return nil
  })

  if err := app.Start(); err != nil {
    log.Fatal(err)
  }
}
