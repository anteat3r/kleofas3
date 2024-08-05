package src

import (
	"github.com/pocketbase/pocketbase/models"
)

func CredsExpiredNotif(user *models.Record) error {
  return SendNotif(user, `{
    "type": "notif",
    "name": "Vypršel refresh token, pls přihlas se znovu"
  }`)
}
