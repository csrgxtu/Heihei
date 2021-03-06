package services

import (
  "github.com/huichen/wukong/engine"
  "github.com/huichen/wukong/types"
  "github.com/astaxie/beego"
  "gopkg.in/mgo.v2/bson"
  "Heihei/models"
  "errors"
)

var BookfulCollection = beego.AppConfig.String("BookfulCollection")

var (
	// searcher是线程安全的
	searcher = engine.Engine{}
)

// func init() {
//   beego.Info("初始化悟空")
//   // 初始化
//   searcher.Init(types.EngineInitOptions{
//     SegmenterDictionaries: "./data/dictionary.txt", StopTokenFile: "./data/stop_tokens.txt", UsePersistentStorage: true, PersistentStorageFolder: "./data", PersistentStorageShards: 20})
//
//   // defer searcher.Close()
//   beego.Info("悟空初始化完毕")
// }

// 这个方法将会很耗时间
func Indexer() (err error, rtv int64) {
  if CheckAndReconnect() != nil {
    return
  }

  var Book models.Bookful
  Iterator := Session.DB(DB).C(BookfulCollection).Find(nil).Iter()
  for Iterator.Next(&Book) {
    beego.Info(Book.WukongDocId)
    searcher.IndexDocument(Book.WukongDocId, types.DocumentIndexData{Content: Book.Title}, false)
  }

  // 等待索引刷新完毕
  searcher.FlushIndex()

  return
}

func Searcher(query string) (err error, rtv models.Bookful) {
  if CheckAndReconnect() != nil {
    return
  }

  beego.Info(searcher.Search(types.SearchRequest{Text: query}))
  var SearchRes = searcher.Search(types.SearchRequest{Text: query}).Docs
  if len(SearchRes) == 0 {
    err = errors.New("Server Internal Error")
    return
  }
  var WukongDocId = SearchRes[0].DocId
  var criteria = bson.M{"wukongdocid": WukongDocId, "clc_sort_num": bson.M{"$ne": ""}}
  err = Session.DB(DB).C(BookfulCollection).Find(criteria).One(&rtv)
  if err != nil {
    beego.Info(err)
    err = errors.New("Server Internal Error")
    return
  }

  return
}

func SearchBookful(query string) (err error, rtv []models.Bookful) {
  if CheckAndReconnect() != nil {
    return
  }

  var criteria = bson.M{"title": query}
  err = Session.DB(DB).C(BookfulCollection).Find(criteria).All(&rtv)
  if err != nil {
    criteria = bson.M{"title": bson.M{"$regex": bson.RegEx{".*" + query + "*.", ""}}}
    err = Session.DB(DB).C(BookfulCollection).Find(criteria).One(&rtv)
    if err != nil {
      beego.Info(err)
      err = errors.New("Server Internal Error")
      return
    }
  }

  return
}
