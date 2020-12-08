package f

import (
	"context"
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	gconv "github.com/og/x/conv"
	ge "github.com/og/x/error"
	"log"
	"reflect"
	"time"
)

type Database struct {
	Core *sqlx.DB
	onlyReadDataSourceName DataSourceName
}
type Storager interface {
	QueryRowxContext(cTransaction context.Context, query string, args ...interface{}) *sqlx.Row
	SelectContext(cTransaction context.Context, dest interface{}, query string, args ...interface{}) error
	ExecContext(cTransaction context.Context, query string, args ...interface{}) (sql.Result, error)
}
func (database Database) DataSourceName () DataSourceName {
	return database.onlyReadDataSourceName
}
func NewDatabase(dataSourceName DataSourceName) (database Database, err error) {
	db, err := sqlx.Connect(dataSourceName.DriverName, dataSourceName.GetString())
	if err != nil {
		return database, err
	}
	database = Database{Core: db}
	database.onlyReadDataSourceName = dataSourceName
	return database, nil
}

func (db *Database) One(cTransaction context.Context, modelPtr Model, has *bool, qb QB) error {
	return coreOne(cTransaction, db.Core, modelPtr, has, qb)
}
func (Transaction *Transaction) One(cTransaction context.Context, modelPtr Model, has *bool, qb QB) error {
	return coreOne(cTransaction, Transaction.Core, modelPtr, has, qb)
}
func coreOne(cTransaction context.Context, storage Storager, modelPtr Model, has *bool, qb QB) error {
	qb.Limit = 1
	scanModelMakeSQLSelect(reflect.ValueOf(modelPtr).Elem().Type(), &qb)
	query, values := qb.BindModel(modelPtr).GetSelect()
	row := storage.QueryRowxContext(cTransaction, query, values...)
	err := row.StructScan(modelPtr)
	if err == sql.ErrNoRows { *has = false ; return nil}
	if err != nil { return err }
	*has = true
	return nil
}
func (db *Database) Relation(cTransaction context.Context, relationPtr RelationModel, has *bool, qb QB) error {
	return coreRelation(cTransaction, db.Core, relationPtr, has, qb)
}
func (Transaction *Transaction) Relation(cTransaction context.Context, relationPtr RelationModel, has *bool, qb QB) error {
	return coreRelation(cTransaction, Transaction.Core, relationPtr, has, qb)
}
func coreRelation(cTransaction context.Context, storage Storager, relationPtr RelationModel, has *bool, qb QB) error {
	qb.Limit = 1
	scanModelMakeSQLSelect(reflect.ValueOf(relationPtr).Elem().Type(), &qb)
	qb.Table = relationPtr.TableName()
	query, values := qb.GetSelect()
	row := storage.QueryRowxContext(cTransaction, query, values...)
	err := row.StructScan(relationPtr)
	if err == sql.ErrNoRows { *has = false ; return nil}
	if err != nil { return err }
	*has = true
	return nil
}
func (db *Database) OneID(cTransaction context.Context, modelPtr Model, has *bool, id interface{}) error {
	return coreOne(cTransaction, db.Core, modelPtr, has, QB{
		Where:And("id", id),
	})
}
func (Transaction *Transaction) OneID(cTransaction context.Context, modelPtr Model, has *bool, id interface{}) error {
	return coreOne(cTransaction, Transaction.Core, modelPtr, has, QB{
		Where:And("id", id),
	})
}
func (Transaction *Transaction) OneIDLock(cTransaction context.Context, modelPtr Model, has *bool, id interface{}, lock SelectLock) error {
	return coreOne(cTransaction, Transaction.Core, modelPtr, has, QB{
		Where:And("id", id),
		Lock: lock,
	})
}

func (db *Database) Count(cTransaction context.Context, modelPtr Model,count *int, qb QB) (err error) {
	return coreCount(cTransaction, db.Core, modelPtr, count,  qb)
}
func (Transaction *Transaction) Count(cTransaction context.Context, modelPtr Model,count *int, qb QB) ( err error) {
	return coreCount(cTransaction, Transaction.Core, modelPtr, count, qb)
}
func coreCount(cTransaction context.Context, storage Storager, modelPtr Model, count *int, qb QB)( err error)  {
	qb.Count = true
	query, values := qb.BindModel(modelPtr).GetSelect()
	row := storage.QueryRowxContext(cTransaction,query, values...)
	err = row.Scan(count)
	if err != nil { return }
	return
}

func (db *Database) List(cTransaction context.Context, modelListPtr interface{}, qb QB) error {
	return coreList(cTransaction, db.Core, modelListPtr, qb)
}
func (Transaction *Transaction) List(cTransaction context.Context, modelListPtr interface{}, qb QB) error {
	return coreList(cTransaction, Transaction.Core, modelListPtr, qb)
}
func coreList(cTransaction context.Context, storage Storager, modelListPtr interface{}, qb QB) error {
	elementType := reflect.TypeOf(modelListPtr).Elem()
	reflectItemValue := reflect.MakeSlice(elementType, 1,1).Index(0)
	scanModelMakeSQLSelect(reflectItemValue.Type(), &qb)
	query, values := qb.BindModel(reflectItemValue.Addr().Interface().(Model)).GetSelect()
	if qb.Table == "" {
		tableName := reflectItemValue.MethodByName("TableName").Call([]reflect.Value{})[0].String()
		qb.Table = tableName
	}
	return storage.SelectContext(cTransaction,modelListPtr, query, values...)
}
func (db *Database) Create(cTransaction context.Context, modelPtr Model) error {
	return coreCreate(cTransaction, db.Core , modelPtr)
}
func (Transaction *Transaction) Create(cTransaction context.Context, modelPtr Model) error {
	return coreCreate(cTransaction, Transaction.Core , modelPtr)
}
func coreCreate(cTransaction context.Context, storage Storager,modelPtr Model) error {
	value, _ := getPtrElem(modelPtr)
	modelPtr.BeforeCreate()
	typeValue := reflect.TypeOf(modelPtr).Elem()
	insertSort := []sortData{}
	var idValue reflect.Value
	var hasPrimaryKeyAutoIncrement bool
	for i:=0;i<value.NumField();i++{
		item := value.Field(i)
		itemType := typeValue.Field(i)
		dbName, hasDBName := itemType.Tag.Lookup("db")
		if !hasDBName { continue }
		switch dbName {
		case "":
			log.Print(`Maybe you forget set db:"name"` + itemType.Name)
			continue
		case "id":
			idValue = item
			autoIncrementValue, has := itemType.Tag.Lookup("dbAutoIncrement")
			if has {
				switch autoIncrementValue {
				case "true":
					hasPrimaryKeyAutoIncrement = true
					continue
				case "false":
				default:
					panic(errors.New(`dbAutoIncrement muse be dbAutoIncrement:"true" or dbAutoIncrement:"false" can not be dbAutoIncrement:"`+ autoIncrementValue + `"`))
				}
			}
		case "created_at","updated_at":
			insertSort = append(insertSort, sortData{Column(dbName), fillTime(time.Now(), item)})
			continue
		case "deleted_at":
			if item.IsZero() {
				continue
			}
		}
		insertSort = append(insertSort, sortData{Column(dbName), item.Interface()})
	}
	query, values := QB{
		Table: modelPtr.TableName(),
		insertSort: insertSort,
	}.GetInsert()
	result, execErr := storage.ExecContext(cTransaction, query, values...) ; if execErr != nil {return execErr}
	bindInsertID(bindInsertIDData{
		Result: result,
		IDValue: idValue,
		HasPrimaryKeyAutoIncrement: hasPrimaryKeyAutoIncrement,
		Model: modelPtr,
	})
	return nil
}
type bindInsertIDData struct {
	Result sql.Result
	IDValue reflect.Value
	HasPrimaryKeyAutoIncrement bool
	Model Model
}
func bindInsertID(data bindInsertIDData) {
	if data.HasPrimaryKeyAutoIncrement {
		lastInsertID, err := data.Result.LastInsertId() ; ge.Check(err)
		if lastInsertID == 0 {
			panic(errors.New(data.Model.TableName() + " does not support AutoIncrement or LastInsertId()"))
		}
		switch data.IDValue.Type().Kind() {
		case reflect.Uint,reflect.Uint8,reflect.Uint16,reflect.Uint32,reflect.Uint64:
			data.IDValue.SetUint(uint64(lastInsertID))
		case reflect.Int,reflect.Int8,reflect.Int16,reflect.Int32,reflect.Int64:
			data.IDValue.SetInt(lastInsertID)
		case reflect.String:
			data.IDValue.SetString(gconv.Int64String(lastInsertID))
		default:
			panic(errors.New(data.Model.TableName() + " ID type must be uint or int or string"))
		}
	}
}
func fillTime (timeValue time.Time, item reflect.Value) (sqlValue interface{}) {
	if item.IsZero() {
		sqlValue = timeValue
		item.Set(reflect.ValueOf(timeValue))
	} else {
		sqlValue = item.Interface()
	}
	return
}
func (db *Database) DeleteQB(cTransaction context.Context, modelPtr Model, qb QB) error {
	return coreDeleteQB(cTransaction, db.Core, modelPtr, qb)
}
func (Transaction *Transaction) DeleteQB(cTransaction context.Context, modelPtr Model, qb QB) error {
	return coreDeleteQB(cTransaction, Transaction.Core, modelPtr, qb)
}
func coreDeleteQB(cTransaction context.Context, storage Storager, modelPtr Model, qb QB) error {
	_, isPtr := getPtrElem(modelPtr)
	if !isPtr {
		panic("db.DeleteQB() or db.TransactionDeleteQB()  arg `modelPtr` must be a ptr, eg: db.DeleteQB(&user, qb) db.TransactionDeleteQB(Transaction, &user, qb) ")
	}
	if len(qb.Update) == 0 {
		qb.Update = map[Column]interface{}{}
	}
	qb.Update["deleted_at"] = time.Now()
	query, values := qb.BindModel(modelPtr).GetUpdate()
	_, err := storage.ExecContext(cTransaction, query, values...)
	if err != nil {return err}
	return nil
}

func (db *Database) Delete(cTransaction context.Context, modelPtr Model) error {
	return coreDelete(cTransaction, db.Core, modelPtr)
}
func (Transaction *Transaction) Delete(cTransaction context.Context, modelPtr Model) error {
	return coreDelete(cTransaction, Transaction.Core, modelPtr)
}
func coreDelete(cTransaction context.Context, storage Storager, modelPtr Model) (error) {
	rValue, isPtr := getPtrElem(modelPtr)
	if !isPtr {
		panic("db.Delete() or db.TransactionDelete()  arg `modelPtr` must be a ptr, eg: db.Delete(&user) db.TransactionDelete(Transaction, &user) ")
	}
	idValue := rValue.FieldByName("ID")
	if idValue.IsZero() {
		panic(errors.New("db.Update(&model) or db.TransactionUpdate(&model) model.id is zero"))
	}
	id := idValue.Interface()
	qb := QB{
		Where: And("id", id),
	}
	if len(qb.Update) == 0 {
		qb.Update = map[Column]interface{}{}
	}
	qb.Update["deleted_at"] = time.Now()
	query, values := qb.BindModel(modelPtr).GetUpdate()
	_, err := storage.ExecContext(cTransaction,query, values...) ; if err != nil {return err}
	return nil
}

func (db *Database) UpdateData(cTransaction context.Context,  modelPtr Model, updateData Data) error {
	return coreUpdateData(cTransaction, db.Core, modelPtr, updateData)
}
func (Transaction *Transaction) UpdateData(cTransaction context.Context,  modelPtr Model, updateData Data) error {
	return coreUpdateData(cTransaction, Transaction.Core, modelPtr, updateData)
}
func coreUpdateData (cTransaction context.Context, storage Storager, modelPtr Model,  updateData Data) error {
	rValue, _ := getPtrElem(modelPtr)
	typeValue := reflect.TypeOf(modelPtr).Elem()
	fieldLen := rValue.NumField()
	var id interface{}
	var findID bool
	for i:=0;i<fieldLen;i++{
		item := rValue.Field(i)
		itemType := typeValue.Field(i)
		dbName, hasDBName := itemType.Tag.Lookup("db")
		if !hasDBName { continue }
		value := item.Interface()
		switch dbName {
		case "":
			log.Print(`Maybe you forget set db:"name"` + itemType.Name)
			continue
		case "id":
			if dbName == "id" {
				if item.IsZero() {
					panic(errors.New("update model.ID is zero"))
				}
				findID = true
				id = value
				continue
			}
		}
		value, hasValue := updateData[Column(dbName)]
		if hasValue {
			item.Set(reflect.ValueOf(value))
		}
	}
	if !findID { panic(errors.New("update not found id")) }
	updatedAtValue := rValue.FieldByName("UpdatedAt")
	if updatedAtValue.IsValid() {
		nowTime := time.Now()
		updatedAtType, _ := typeValue.FieldByName("UpdatedAt")
		updateData[Column(updatedAtType.Tag.Get("db"))] = nowTime
		updatedAtValue.Set(reflect.ValueOf(nowTime))
	}
	query, values := QB{
		Table: modelPtr.TableName(),
		Where: And("id", id),
		Update: updateData,
	}.GetUpdate()
	_, err := storage.ExecContext(cTransaction, query, values...) ; if err != nil {return err}
	return nil
}
func (db *Database) Transaction(cTransaction context.Context, transaction func(fTransaction *Transaction) error) (error) {
	return db.TransactionOpts(cTransaction, nil, transaction)
}
func (db *Database) TransactionOpts(cTransaction context.Context, opts *sql.TxOptions, transaction func(*Transaction) error) (TransactionError error) {
	sqlxTx, err := db.Core.BeginTxx(cTransaction, opts) ; if err != nil {return err}
	Transaction := newTransaction(sqlxTx)
	defer func() {
		r := recover()
		if r != nil {
			rollbackErr := Transaction.Rollback() ; _= rollbackErr // 此时可以忽略 rollback 的错误
			panic(r)
		}
	}()
	err = transaction(Transaction)
	if err != nil {
		ge.Check(Transaction.Rollback())
		return err
	} else {
		return Transaction.commit()
	}
}
func (db *Database) Close() error {
	return db.Core.Close()
}
func (db *Database) QueryRowScan(cTransaction context.Context, has *bool, qb QB,  dest ...interface{}) error {
	return coreQueryRowScan(cTransaction, db.Core, has, qb, dest...)
}
func (Transaction *Transaction) QueryRowScan(cTransaction context.Context, qb QB, has *bool,  dest ...interface{}) error {
	return coreQueryRowScan(cTransaction, Transaction.Core, has, qb, dest...)
}
func coreQueryRowScan(cTransaction context.Context, storage Storager, has *bool, qb QB, dest ...interface{}) error {
	query, values := qb.GetSelect()
	row := storage.QueryRowxContext(cTransaction, query, values...)
	err := row.Scan(dest...)
	*has = true
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			*has = false
			return nil
		}
		return err
	}
	return nil
}
