package src

import (
	"encoding/json"

	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
)

type BakaHour struct {
  Id int `json:"Id"`
  Caption string `json:"Caption"`
  BeginTime string `json:"BeginTime"`
  EndTime string `json:"EndTime"`
}

type BakaDay struct {
  Atoms []BakaAtom `json:"Atoms"`
  DayOfWeek int `json:"DayOfWeek"`
  Date string `json:"Date"`
  DayDescription string `json:"DayDescription"`
  DayType string `json:"DayType"`
}

type BakaAtom struct {
  HourId int `json:"HourId"`
  GroupIds []string `json:"GroupIds"`
  SubjectId string `json:"SubjectId"`
  TeacherId string `json:"TeacherId"`
  RoomId string `json:"RoomId"`
  CycleIds []string `json:"CycleIds"`
  Change BakaChange `json:"Change"`
  HomeworkIds []string `json:"HomeworkIds"`
  HomeWorks []any `json:"HomeWorks"`
  Theme string `json:"Theme"`
  Assistants []any `json:"Assistants"`
  Notice string `json:"Notice"`
}

type BakaChange struct {
  ChangeSubject string `json:"ChangeSubject"`
  Day string `json:"Day"`
  Hours string `json:"Hours"`
  ChangeType string `json:"ChangeType"`
  Description string `json:"Description"`
  Time string `json:"Time"`
  TypeAbbrev string `json:"TypeAbbrev"`
  TypeName string `json:"TypeName"`
}

type BakaTimeTable struct {
  Hours []BakaHour `json:"Hours"`
  Days []BakaDay `json:"Days"`
  Classes []BakaClassIdPair `json:"Classes"`
  Groups []BakaIdPair `json:"Groups"`
  Subjects []BakaIdPair `json:"Subjects"`
  Teachers []BakaIdPair `json:"Teachers"`
  Rooms []BakaIdPair `json:"Rooms"`
  Cycles []BakaIdPair `json:"Cycles"`
  Students []BakaIdPair `json:"Students"`
}

type BakaClassIdPair struct {
  Id string `json:"Id"`
  ClassId string `json:"ClassId"`
  Abbrev string `json:"Abbrev"`
  Name string `json:"Name"`
}

type BakaIdPair struct {
  Id string `json:"Id"`
  Abbrev string `json:"Abbrev"`
  Name string `json:"Name"`
}

func ProcessTimeTable(
  dao *daos.Dao,
  user *models.Record,
  tables string,
  oldtables string,
) error {
  table := BakaTimeTable{}
  oldtable := BakaTimeTable{}
  err := json.Unmarshal([]byte(tables), &table)
  if err != nil { return err }
  err = json.Unmarshal([]byte(oldtables), &oldtable)
  if err != nil { return err }

  // if 
  return nil
  
}
