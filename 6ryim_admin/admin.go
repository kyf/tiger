package main

import (
	"gopkg.in/mgo.v2/bson"
)

const (
	ADMIN_TABLE = "cc_admin"
)

type Admin struct {
	Id       bson.ObjectId `json:"id" bson:"_id"`
	UserName string        `json:"user" bson:"user"`
	Password string        `json:"pwd" bson:"pwd"`
	Opid     string        `json:"opid" bson:"opid"`
	mgo      *Mongo
}

func NewAdmin(user, pwd, opid string, mgo *Mongo) *Admin {
	return &Admin{bson.NewObjectId(), user, pwd, opid, mgo}
}

func (adm *Admin) add() error {
	data := bson.M{
		"_id":  adm.Id,
		"user": adm.UserName,
		"pwd":  adm.Password,
		"opid": adm.Opid,
	}
	err := adm.mgo.Add(ADMIN_TABLE, data)
	if err != nil {
		return err
	}
	return nil
}

func getAdminByName(name string, mgo *Mongo) (*Admin, error) {
	data := bson.M{
		"user": name,
	}

	adm := new(Admin)
	err := mgo.Find(ADMIN_TABLE, data).Limit(1).One(adm)
	if err != nil {
		return nil, err
	}
	return adm, nil
}

func (adm *Admin) checkUniq() (bool, error) {
	data := bson.M{
		"user": adm.UserName,
		"pwd":  adm.Password,
	}
	num, err := adm.mgo.Find(ADMIN_TABLE, data).Count()
	if err != nil {
		return false, err
	}
	if num > 0 {
		return false, nil
	} else {
		return true, nil
	}
}

func (adm *Admin) checkValid() (bool, error) {
	data := bson.M{
		"user": adm.UserName,
		"pwd":  adm.Password,
	}
	num, err := adm.mgo.Find(ADMIN_TABLE, data).Count()
	if err != nil {
		return false, err
	}
	if num > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func (adm *Admin) list() ([]Admin, error) {
	var result []Admin
	err := adm.mgo.Find(ADMIN_TABLE, nil).All(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (adm *Admin) remove(id string) error {
	data := bson.M{
		"_id": bson.ObjectIdHex(id),
	}
	err := adm.mgo.Remove(ADMIN_TABLE, data)
	if err != nil {
		return err
	}
	return nil
}

func (adm *Admin) edit(id string) error {
	data := bson.M{
		"_id": bson.ObjectIdHex(id),
	}

	newData := bson.M{
		"user": adm.UserName,
		"pwd":  adm.Password,
		"opid": adm.Opid,
	}

	err := adm.mgo.Update(ADMIN_TABLE, data, newData)
	if err != nil {
		return err
	}
	return nil
}
