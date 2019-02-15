package orm

import (
	"encoding/json"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

type Mongo struct {

}

// 解析成功的站点页面
type Website struct {
	WebsiteURL string `bson:"websiteurl" json:"websiteurl"`
	LastTop    int    `bson:"lasttop"  json:"ctr"` //上一次解析的页数
}

// 书本解析失败
type FailBook struct {
}

// 章节解析失败
type FailChapter struct {
}

type MBook struct {
	Identifier  string  `bson:"identifier" json:"identifier"`
	Name        string  `bson:"name" json:"name"`
	Domain      string  `bson:"domain" json:"domain"`
	Cover       string  `json:"cover"`
	Source      string  `json:"source"`
	Describe    string  `json:"describe"`
	Author      string  `json:"author"`
	Type        string  `json:"type"`
	Last_update string  `bson:"last_update" json:"last_update"`
	Language    string  `json:"language"`
	Source_ctr  int64   `json:"source_ctr"`
	Ctr         int64   `json:"ctr"`
	Score       float32 `json:"score"`
	Finish      string  `bson:"finish" json:"finish"` //章节解析完成时间
	UpTime      string  `bson:"uptime" json:"uptime"` //上传时间,保存的时候该时间为空。上传成功后，设置为更新章节时间
}

type MChapter struct {
	Identifier string `bson:"identifier" json:"identifier"`
	Idx        int    `json:"idx"`                //索引序列号
	Idx_name   string `json:"idx_name"`           //索引名，第一章，第二章等
	Title      string `bson:"title" json:"title"` //标题
	Content    string `json:"content"`            //内容
	Source     string `json:"source"`             //来源 crawler
	Domain     string `json:"domain"`
	UpTime     string `bson:"uptime" json:"uptime"` //上传时间,保存的时候该时间为空。上传成功后，设置为更新章节时间
}

type ResultMBook struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ResultMChapter struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var collection *mgo.Collection

func init() {
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}
	// defer session.Close()//?? 关闭session

	session.SetMode(mgo.Monotonic, true)
	collection = session.DB("test").C("person")
}

func (mongo Mongo)Insert(i interface{}) {
	err := collection.Insert(i)
	if err != nil {
		log.Println(err)
	}
}

func (mongo Mongo)WebSite(url string) Website {
	var website Website
	err := collection.Find(bson.M{"websiteurl": url, "lasttop": bson.M{"$exists": true}}).One(&website)
	if err != nil {
		log.Println(err)
	}
	return website
}

func (mongo Mongo)AllMBook() []MBook {
	var mBooks []MBook
	err := collection.Find(bson.M{"uptime": "", "last_update": bson.M{"$exists": true}}).All(&mBooks)
	if err != nil {
		log.Println(err)
	}

	return mBooks
}

func (mongo Mongo)AllMChapter() []MChapter {
	var mMChapters []MChapter
	err := collection.Find(bson.M{"uptime": "", "title": bson.M{"$exists": true}}).All(&mMChapters)
	if err != nil {
		log.Println(err)
	}

	return mMChapters
}

func (mongo Mongo)BookIdentify(idnetify string) MBook {
	var mBook MBook
	err := collection.Find(bson.M{"identifier": idnetify, "name": bson.M{"$exists": true}}).One(&mBook)
	if err != nil {
		log.Println(err)
	}

	return mBook
}

func (mongo Mongo)ChapterIdentifyName(identify string, title string) MChapter {
	var mMChapter MChapter
	err := collection.Find(bson.M{"identifier": identify, "title": title}).One(&mMChapter)
	if err != nil {
		log.Println(err)
	}

	return mMChapter
}

func (mongo Mongo)LocalBookFinsih(identify string) {
	_, err := collection.UpdateAll(bson.M{"identifier": identify}, bson.M{"$set": bson.M{"finish": time.Now().String()}})
	if err != nil {
		log.Println("UpdateBook: ", err)
	}
}

func (mongo Mongo)UpBookSuccess(identify string) {
	_, err := collection.UpdateAll(bson.M{"identifier": identify}, bson.M{"$set": bson.M{"uptime": time.Now().String()}})
	if err != nil {
		log.Println("UpdateBook: ", err)
	}
}

func(mongo Mongo) UpChapterSuccess(identify string, title string) {
	_, err := collection.UpdateAll(bson.M{"identifier": identify, "title": title}, bson.M{"$set": bson.M{"uptime": time.Now().String()}})
	if err != nil {
		log.Println("UpdateChapter: ", err)
	}
}

func (mongo Mongo)Upload() {
	var mBooks = AllMBook()

	log.Print("Upload mBooks: ", mBooks)
	for _, mbook := range mBooks {
		result := UpBook(mbook)

		var resultBook ResultMBook
		err := json.Unmarshal([]byte(result), &resultBook)
		if err != nil {
			log.Fatal(err, resultBook)
		}
		UpBookSuccess(mbook.Identifier)
	}

	UploadChapter()
}

func (mongo Mongo)UploadChapter() {
	var mChapters = AllMChapter()
	log.Print("Upload mChapters: ", mChapters)

	for _, mChapter := range mChapters {
		result := UpChapter(mChapter)

		var resultChapter ResultMChapter
		err := json.Unmarshal([]byte(result), &resultChapter)
		if err != nil {
			log.Fatal(err, mChapter)
		}
		UpChapterSuccess(mChapter.Identifier, mChapter.Title)
	}
}
