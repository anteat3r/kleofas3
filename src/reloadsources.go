package src

import (
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/tools/types"
)

func ReloadSources(dao *daos.Dao) func() {
  return func() {
    users, err := dao.FindRecordsByFilter(
      COLLECTION_USERS,
      USERS_VALID_LOGIN + " = true",
      "+last_used",
      1, 0,
    )
    if err != nil { LogError(err); return }
    if len(users) == 0 { return }
    user := users[0]

    html, err := WebQuery(dao, user, TIMETABLE_PUBLIC)
    if err != nil { LogError(err); return }

    srcs, err := ParseSourcesWeb(html)
    if err != nil { LogError(err); return }

    dao.RunInTransaction(func(txDao *daos.Dao) error {
      sources, err := txDao.FindRecordsByFilter(
        COLLECTION_SOURCES,
        `id != ""`,
        "-created",
        0, 0,
      )
      if err != nil { LogError(err); return nil }

      for _, src := range sources {
        err := txDao.DeleteRecord(src)
        if err != nil { LogError(err) }
      }

      for _, src := range srcs.teachers {
        rec, err := NewRecord(txDao, COLLECTION_SOURCES)
        if err != nil { LogError(err); continue }
        rec.Set(SOURCES_NAME, src.id)
        rec.Set(SOURCES_TYPE, "teacher")
        rec.Set(SOURCES_LAST_FETCHED, types.NowDateTime())
        err = txDao.SaveRecord(rec)
        if err != nil { LogError(err); continue }
      }

      for _, src := range srcs.rooms {
        rec, err := NewRecord(txDao, "sources")
        if err != nil { LogError(err); continue }
        rec.Set(SOURCES_NAME, src.id)
        rec.Set(SOURCES_TYPE, "room")
        rec.Set(SOURCES_LAST_FETCHED, types.NowDateTime())
        err = txDao.SaveRecord(rec)
        if err != nil { LogError(err); continue }
      }

      for _, src := range srcs.classes {
        rec, err := NewRecord(txDao, "sources")
        if err != nil { LogError(err); continue }
        rec.Set(SOURCES_NAME, src.id)
        rec.Set(SOURCES_TYPE, "class")
        rec.Set(SOURCES_LAST_FETCHED, types.NowDateTime())
        err = txDao.SaveRecord(rec)
        if err != nil { LogError(err); continue }
      }

      // TODO: add valid baka creds for this to work
      // rec, err := NewRecord(txDao, "sources")
      // if err != nil { LogError(err) }
      // rec.Set(SOURCES_NAME, "")
      // rec.Set(SOURCES_TYPE, "allevents")
      // rec.Set(SOURCES_LAST_FETCHED, types.NowDateTime())
      // err = txDao.SaveRecord(rec)
      // if err != nil { LogError(err) }

      return nil
    })
  }
}
