package main

import (
	"fmt"
	"log"
	"os"

	"github.com/anteat3r/kleofas3/src"
	"github.com/labstack/echo/v5"

	"github.com/joho/godotenv"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/cron"
	"github.com/pocketbase/pocketbase/tools/types"
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
      username := c.PathParam("user")
      user, _ := app.Dao().FindFirstRecordByData("users", "username", username)
      data := user.Get("last_used").(types.DateTime)
      return c.String(200, fmt.Sprintf("%v %T\n", data.IsZero(), data))
    })

    scheduler.Start()

    return nil
  })

  if err := app.Start(); err != nil {
    log.Fatal(err)
  }
}
