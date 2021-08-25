package storage

import (
	"errors"
	"fmt"
	"gimli/Internal/gRPC/domain"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
)

type DatabasePtr *gorm.DB

var db *gorm.DB

type OriginURL struct {
	model gorm.Model
	URL   string `gorm:"column:url"`
	ID    int    `gorm:"column:id"`
}

type ShortURL struct {
	model    gorm.Model
	Url      string `gorm:"column:url"`
	OriginID int    `gorm:"column:origin_id"`
}

var Running bool = false

const (
	host           = "127.0.0.1"
	port           = 5432
	user           = "postgres"
	password       = "Nata2010"
	dbname         = "gimli_db"
	shortURLTable  = "short_url"
	originURLTable = "origin_url"
)

func SetupConnection() *gorm.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	var err error
	db, err = gorm.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return "public." + defaultTableName
	}
	db.LogMode(true)
	Running = true
	return db
}

func CloseConnection() {
	db.Close()
}

func OriginToFull(o OriginURL) domain.FullURL {
	return domain.FullURL{
		Url: o.URL}
}

func ShortURLToShortPath(s ShortURL) domain.ShortPath {
	return domain.ShortPath{
		Path: s.Url}
}

func getOriginByURL(URL string) (*OriginURL, error) {
	origin := OriginURL{}
	db.Table(originURLTable).Where("url = ?", URL).First(&origin)
	if db.Error != nil {
		return nil, db.Error
	}
	return &origin, nil
}

func insertOriginURL(urlStr string) (OriginURL, error) {
	originURL := OriginURL{}
	db.Table(originURLTable).Last(&originURL)
	insert := OriginURL{
		ID:  originURL.ID + 1,
		URL: urlStr,
	}
	db.Table(originURLTable).Create(&insert)
	return insert, db.Error
}

func insertShortURL(shortURLStr string, originURLId int) (ShortURL, error) {
	insert := ShortURL{
		Url:      shortURLStr,
		OriginID: originURLId}
	db.Table(shortURLTable).Create(&insert)
	return insert, db.Error
}

//InsertPairURL added originURL and shortURL to database. If originURL already exists, returns (nil, error)
func InsertPairURL(originURL, shortURL string) (*ShortURL, error) {
	origin, err := getOriginByURL(originURL)
	if origin.URL != "" {
		return nil, errors.New("URL = " + originURL + " already exists in Service.")
	}
	*origin, err = insertOriginURL(originURL)
	if err != nil {
		log.Println("Cannot insertOriginURL: ", origin, "\nwith error ", err)
		panic(err)
	}
	short, err := insertShortURL(shortURL, origin.ID)
	return &short, err
}

//GetOriginByShort returns OriginURL matching shortUrlStr. If it is not exists, returns (nil, err)
func GetOriginByShort(shortUrlStr string) (*OriginURL, error) {
	shortURL := ShortURL{}
	originURL := OriginURL{}
	dbErr := db.Table(shortURLTable).Where("url = ?", shortUrlStr).First(&shortURL).Error
	if errors.Is(dbErr, gorm.ErrRecordNotFound){
		log.Println(shortUrlStr, " not found ")
		return nil, errors.New(shortUrlStr + " not found ")
	}
	dbErr = db.Table(originURLTable).Where("id = ?", shortURL.OriginID).First(&originURL).Error
	if errors.Is(dbErr, gorm.ErrRecordNotFound){
		log.Println(shortUrlStr, " not found ")
		return nil, errors.New(shortUrlStr + " not found ")
	}
	return &originURL, db.Error
}
