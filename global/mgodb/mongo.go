package mgodb

import (
	"fmt"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

var MongoSession *mgo.Session

type Options struct {
	Addr     string
	Database string
	Username string
	Password string
	Idle     int
}

func InitMongo(option Options) error {
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{option.Addr},
		Timeout:  60 * time.Second,
		Database: option.Database,
		Username: option.Username,
		Password: option.Password,
	}
	// Create a session which maintains a pool of socket connections
	// to our MongoDB.
	session, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		fmt.Println("mongodb init fail ")
		return err
	}
	session.SetPoolLimit(option.Idle)
	session.SetMode(mgo.Monotonic, true)
	MongoSession = session
	return nil
}

type MgoSelector struct {
	condition map[string]interface{}
	sort      string
	table     string
	limit     int
	skip      int
}

func NewMgoSelector() *MgoSelector {
	mgoselector := new(MgoSelector)
	return mgoselector
}

func (this *MgoSelector) Table(table string) *MgoSelector {
	this.table = table
	return this
}

func (this *MgoSelector) Where(condition map[string]interface{}) *MgoSelector {
	this.condition = condition
	return this
}

func (this *MgoSelector) Sort(sort string) *MgoSelector {
	this.sort = sort
	return this
}

func (this *MgoSelector) Limit(limit int) *MgoSelector {
	this.limit = limit
	return this
}

func (this *MgoSelector) Skip(skip int) *MgoSelector {
	this.skip = skip
	return this
}

func (this *MgoSelector) Find() (map[string]interface{}, error) {
	newsession := MongoSession.Copy()
	defer newsession.Close()
	res := map[string]interface{}{}
	err := newsession.DB("").C(this.table).Find(this.condition).One(&res)
	return res, err
}

func (this *MgoSelector) Pipe(query []bson.M) (map[string]int, error) {
	newsession := MongoSession.Copy()
	defer newsession.Close()
	res := map[string]int{}
	err := newsession.DB("").C(this.table).Pipe(query).One(&res)
	return res, err
}

func (this *MgoSelector) Select(res interface{}) error {
	newsession := MongoSession.Copy()
	defer newsession.Close()
	query := newsession.DB("").C(this.table).Find(this.condition)
	if this.limit != 0 {
		query = query.Limit(this.limit)
	}
	if this.skip != 0 {
		query = query.Skip(this.skip)
	}
	if this.sort != "" {
		query = query.Sort(this.sort)
	}

	err := query.All(res)
	return err
}

func (this *MgoSelector) InsertOne(insertData interface{}) error {
	newsession := MongoSession.Copy()
	defer newsession.Close()
	err := newsession.DB("").C(this.table).Insert(&insertData)
	if err != nil {
		return err
	}
	return nil
}

func (this *MgoSelector) UpInsert(selector interface{}, update interface{}) (*mgo.ChangeInfo, error) {
	newsession := MongoSession.Copy()
	defer newsession.Close()
	return newsession.DB("").C(this.table).Upsert(selector, update)
}

func (this *MgoSelector) RemoveOne(selector interface{}) error {
	newsession := MongoSession.Copy()
	defer newsession.Close()
	return newsession.DB("").C(this.table).Remove(selector)
}

func (this *MgoSelector) RemoveAll(selector interface{}) (*mgo.ChangeInfo, error) {
	newsession := MongoSession.Copy()
	defer newsession.Close()
	return newsession.DB("").C(this.table).RemoveAll(selector)
}

func (this *MgoSelector) Count(selector interface{}) (int, error) {
	newsession := MongoSession.Copy()
	defer newsession.Close()
	return newsession.DB("").C(this.table).Find(selector).Count()
}

func (this *MgoSelector) InsertMore(insertDatas []map[string]interface{}) error {
	newsession := MongoSession.Copy()
	defer newsession.Close()
	retMsg := ""
	for _, v := range insertDatas {
		err := newsession.DB("").C(this.table).Insert(&v)
		if err != nil {
			retMsg = retMsg + err.Error() + "|"
		}
	}
	if retMsg == "" {
		return nil
	} else {
		return errors.New(retMsg)
	}
}

func (this *MgoSelector) Update(selector interface{}, updateData interface{}) error {
	newsession := MongoSession.Copy()
	defer newsession.Close()
	err := newsession.DB("").C(this.table).Update(selector, updateData)
	return err
}

func (this *MgoSelector) Upsert(selector interface{}, updateData interface{}) error {
	newsession := MongoSession.Copy()
	defer newsession.Close()
	_, err := newsession.DB("").C(this.table).Upsert(selector, updateData)
	return err
}

func (this *MgoSelector) Remove(selector interface{}) error {
	newsession := MongoSession.Copy()
	defer newsession.Close()
	_, err := newsession.DB("").C(this.table).RemoveAll(selector)
	return err
}
