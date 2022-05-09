package mysql

/*
Copyright 2022 QuanxiangCloud Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
     http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
import (
	"gorm.io/gorm"
	"strconv"

	"github.com/quanxiang-cloud/organizations/internal/logic/org/consts"
	"github.com/quanxiang-cloud/organizations/internal/models/octopus"
)

type extendRepo struct {
}

//NewExtendRepo new
func NewExtendRepo() octopus.ExtendRepo {
	return new(extendRepo)
}

func (e *extendRepo) Insert(db, tx *gorm.DB, tableName string, r map[string]interface{}) error {
	if tableName == "" {
		tableName = "default"
	}
	if !checkTableExist(db, tableName) {
		if err := createTable(db, tableName, &octopus.Extend{}); err != nil {
			return err
		}
	}
	return tx.Table(tableName).Create(r).Error
}

func (e *extendRepo) InsertList(db, tx *gorm.DB, tableName string, r []map[string]interface{}) (err error) {
	//TODO implement me
	panic("implement me")
}

func (e *extendRepo) UpdateByID(db, tx *gorm.DB, tableName string, extend *octopus.Extend, r map[string]interface{}) (err error) {
	if tableName == "" {
		tableName = "default"
	}
	if !checkTableExist(db, tableName) {
		if err := createTable(db, tableName, &octopus.Extend{}); err != nil {
			return err
		}
		return tx.Table(tableName).Create(r).Error
	}
	return tx.Table(tableName).Model(extend).Updates(r).Error
}

func (e *extendRepo) SelectList(db *gorm.DB, tableName string, status, page, limit int) (list []map[string]interface{}, total int64) {
	//TODO implement me
	panic("implement me")
}

func (e *extendRepo) SelectByID(db *gorm.DB, tableName string, id string) (res map[string]interface{}) {
	if tableName == "" {
		tableName = "default"
	}
	db = db.Table(tableName).Where("id=?", id)
	data := make(map[string]interface{})
	affected := db.Find(&data).RowsAffected
	if affected == 1 {
		return data
	}
	return nil
}

func (e *extendRepo) SelectByIDs(db *gorm.DB, tableName string, ids []string) (list []map[string]interface{}) {
	if tableName == "" {
		tableName = "default"
	}
	db = db.Table(tableName).Where("id in (?)", ids)
	res := make([]map[string]interface{}, 0)
	affected := db.Find(&res).RowsAffected
	if affected > 0 {
		return res
	}
	return nil
}

func checkTableExist(db *gorm.DB, tableName string) bool {
	if tableName == "" {
		tableName = "default"
	}
	return db.Migrator().HasTable(tableName)
}
func createTable(db *gorm.DB, tableName string, tb interface{}) error {
	if tableName == "" {
		tableName = "default"
	}
	return db.Table(tableName).Migrator().CreateTable(tb)
}

//-----------------------扩展字段相关处理---------------------------

//NewManageColumnRepo 初始化
func NewManageColumnRepo() octopus.ManageColumn {
	return new(extendRepo)
}

func (e *extendRepo) AddColumns(db *gorm.DB, tableName, columnName string, types string, len int, pointLen int) error {
	if tableName == "" {
		tableName = "default"
	}
	var sql = "alter table `" + tableName + "`"
	switch types {
	case consts.STRING:
		if len == 0 {
			sql = sql + " add " + columnName + " " + consts.DBColumns[types] + ";"
		}
		if len > 0 && len <= 255 {
			sql = sql + " add " + columnName + " " + consts.DBColumns[types] + "(" + strconv.Itoa(len) + ");"
		}
		if len > 255 {
			sql = sql + " add " + columnName + " text;"
		}
	case consts.INT, consts.INT64, consts.TIME:
		sql = sql + " add " + columnName + " " + consts.DBColumns[types] + "(" + strconv.Itoa(len) + ");"
	case consts.FLOAT:
		sql = sql + " add " + columnName + " " + consts.DBColumns[types] + "(" + strconv.Itoa(len) + "," + strconv.Itoa(pointLen) + ");"
	case consts.TEXT:
		sql = sql + " add " + columnName + " " + consts.DBColumns[types] + ";"

	}
	return db.Exec(sql).Error
}

func (e *extendRepo) DropColumns(db *gorm.DB, tableName, columnName string) error {
	if tableName == "" {
		tableName = "default"
	}
	return db.Exec("alter table `" + tableName + "` drop column " + columnName).Error
}

func (e *extendRepo) CreateTable(db *gorm.DB, tableName string) error {
	if tableName == "" {
		tableName = "default"
	}
	return createTable(db, tableName, &octopus.Extend{})
}

func (e *extendRepo) CheckTableExist(db *gorm.DB, tableName string) bool {
	if tableName == "" {
		tableName = "default"
	}
	return checkTableExist(db, tableName)
}
