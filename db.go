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
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
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

func (db *Database) OneQB(ctx context.Context, modelPtr Model, has *bool, qb QB) error {
	return coreOneQB(ctx, db.Core, modelPtr, has, qb)
}
func (tx *Tx) OneQB(ctx context.Context, modelPtr Model, has *bool, qb QB) error {
	return coreOneQB(ctx, tx.Core, modelPtr, has, qb)
}
func coreOneQB(ctx context.Context, storage Storager, modelPtr Model, has *bool, qb QB) error {
	qb.Limit = 1
	scanModelMakeSQLSelect(reflect.ValueOf(modelPtr).Elem().Type(), &qb)
	query, values := qb.BindModel(modelPtr).GetSelect()
	row := storage.QueryRowxContext(ctx, query, values...)
	err := row.StructScan(modelPtr)
	if err == sql.ErrNoRows { *has = false ; return nil}
	if err != nil { return err }
	*has = true
	return nil
}
func (db *Database) OneID(ctx context.Context, modelPtr Model, has *bool, id interface{}) error {
	return coreOneQB(ctx, db.Core, modelPtr, has, QB{
		Where:And("id", id),
	})
}
func (tx *Tx) OneID(ctx context.Context, modelPtr Model, has *bool, id interface{}) error {
	return coreOneQB(ctx, tx.Core, modelPtr, has, QB{
		Where:And("id", id),
	})
}
func (tx *Tx) OneIDLock(ctx context.Context, modelPtr Model, has *bool, id interface{}, lock SelectLock) error {
	return coreOneQB(ctx, tx.Core, modelPtr, has, QB{
		Where:And("id", id),
		Lock: lock,
	})
}

func (db *Database) Count(ctx context.Context, modelPtr Model, qb QB) (count int, err error) {
	return coreCountQB(ctx, db.Core, modelPtr, qb)
}
func (tx *Tx) Count(ctx context.Context, modelPtr Model, qb QB) (count int, err error) {
	return coreCountQB(ctx, tx.Core, modelPtr, qb)
}
func coreCountQB(ctx context.Context, storage Storager, modelPtr Model, qb QB) (count int, err error)  {
	qb.Count = true
	query, values := qb.BindModel(modelPtr).GetSelect()
	row := storage.QueryRowxContext(ctx,query, values...)
	err = row.Scan(&count)
	if err != nil { return }
	return
}

func (db *Database) ListQB(ctx context.Context, modelListPtr interface{}, qb QB) error {
	return coreListQB(ctx, db.Core, modelListPtr, qb)
}
func (tx *Tx) ListQB(ctx context.Context, modelListPtr interface{}, qb QB) error {
	return coreListQB(ctx, tx.Core, modelListPtr, qb)
}
func coreListQB(ctx context.Context, storage Storager, modelListPtr interface{}, qb QB) error {
	elementType := reflect.TypeOf(modelListPtr).Elem()
	reflectItemValue := reflect.MakeSlice(elementType, 1,1).Index(0)
	scanModelMakeSQLSelect(reflectItemValue.Type(), &qb)
	query, values := qb.BindModel(reflectItemValue.Addr().Interface().(Model)).GetSelect()
	if qb.Table == "" {
		tableName := reflectItemValue.MethodByName("TableName").Call([]reflect.Value{})[0].String()
		qb.Table = tableName
	}
	return storage.SelectContext(ctx,modelListPtr, query, values...)
}
func (db *Database) Create(ctx context.Context, modelPtr Model) error {
	return coreCreate(ctx, db.Core , modelPtr)
}
func (tx *Tx) Create(ctx context.Context, modelPtr Model) error {
	return coreCreate(ctx, tx.Core , modelPtr)
}
func coreCreate(ctx context.Context, storage Storager,modelPtr Model) error {
	value, _ := getPtrElem(modelPtr)
	modelPtr.BeforeCreate()
	typeValue := reflect.TypeOf(modelPtr).Elem()
	insertData := map[Column]interface{}{}
	var idValue reflect.Value
	var hasPrimaryKeyAutoIncrement bool
	for i:=0;i<value.NumField();i++{
		item := value.Field(i)
		itemType := typeValue.Field(i)
		dbName, hasDBName := itemType.Tag.Lookup("db")
		if !hasDBName {
			continue
		}
		if dbName == "" {
			log.Print(`Maybe you forget set db:"name"` + itemType.Name)
			continue
		}
		if dbName == "id" {
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
		}
		insertData[Column(dbName)] = item.Interface()
	}
	nowTime := time.Now()
	{
		createdAtValue := value.FieldByName("CreatedAt")
		if createdAtValue.IsValid() {
			createdAtType, _ := typeValue.FieldByName("CreatedAt")
			timeValue := nowTime
			{
				createdAtTime := createdAtValue.Interface().(time.Time)
				if createdAtTime.IsZero() {
					createdAtValue.Set(reflect.ValueOf(timeValue))
				} else {
					timeValue = createdAtTime
				}
			}
			insertData[Column(createdAtType.Tag.Get("db"))] = timeValue

		}
	}
	{
		updatedAtValue := value.FieldByName("UpdatedAt")
		if updatedAtValue.IsValid() {
			updatedAtType, _ := typeValue.FieldByName("UpdatedAt")
			timeValue := nowTime
			{
				updatedAtTime := updatedAtValue.Interface().(time.Time)
				if updatedAtTime.IsZero() {
					updatedAtValue.Set(reflect.ValueOf(timeValue))
				} else {
					timeValue = updatedAtTime
				}
			}
			insertData[Column(updatedAtType.Tag.Get("db"))] = timeValue
		}
	}
	query, values := QB{
		Insert: insertData,
	}.BindModel(modelPtr).GetInsert()
	var result sql.Result
	result, err := storage.ExecContext(ctx, query, values...) ; if err != nil {return err}
	if hasPrimaryKeyAutoIncrement {
		lastInsertID, err := result.LastInsertId() ; ge.Check(err)
		if lastInsertID == 0 {
			panic(errors.New(modelPtr.TableName() + " does not support AutoIncrement or LastInsertId()"))
		}
		switch idValue.Type().Kind() {
		case reflect.Uint,reflect.Uint8,reflect.Uint16,reflect.Uint32,reflect.Uint64:
			idValue.SetUint(uint64(lastInsertID))
		case reflect.Int,reflect.Int8,reflect.Int16,reflect.Int32,reflect.Int64:
			idValue.SetInt(lastInsertID)
		case reflect.String:
			idValue.SetString(gconv.Int64String(lastInsertID))
		default:
			panic(errors.New(typeValue.Name() + ".ID type must be uint or int or string"))
		}
	}
	return nil
}

func (db *Database) DeleteQB(ctx context.Context, modelPtr Model, qb QB) error {
	return coreDeleteQB(ctx, db.Core, modelPtr, qb)
}
func (tx *Tx) DeleteQB(ctx context.Context, modelPtr Model, qb QB) error {
	return coreDeleteQB(ctx, tx.Core, modelPtr, qb)
}
func coreDeleteQB(ctx context.Context, storage Storager, modelPtr Model, qb QB) error {
	_, isPtr := getPtrElem(modelPtr)
	if !isPtr {
		panic("db.DeleteQB() or db.TxDeleteQB()  arg `modelPtr` must be a ptr, eg: db.DeleteQB(&user, qb) db.TxDeleteQB(tx, &user, qb) ")
	}
	if len(qb.Update) == 0 {
		qb.Update = map[Column]interface{}{}
	}
	qb.Update["deleted_at"] = time.Now()
	query, values := qb.BindModel(modelPtr).GetUpdate()
	_, err := storage.ExecContext(ctx, query, values...)
	if err != nil {return err}
	return nil
}

func (db *Database) Delete(ctx context.Context, modelPtr Model) error {
	return coreDelete(ctx, db.Core, modelPtr)
}
func (tx *Tx) Delete(ctx context.Context, modelPtr Model) error {
	return coreDelete(ctx, tx.Core, modelPtr)
}
func  coreDelete(ctx context.Context, storage Storager, modelPtr Model) (error) {
	rValue, isPtr := getPtrElem(modelPtr)
	if !isPtr {
		panic("db.Delete() or db.TxDelete()  arg `modelPtr` must be a ptr, eg: db.Delete(&user) db.TxDelete(tx, &user) ")
	}
	idValue := rValue.FieldByName("ID")
	if idValue.IsZero() {
		panic(errors.New("db.Update(&model) or db.TxUpdate(&model) model.id is zero"))
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
	_, err := storage.ExecContext(ctx,query, values...) ; if err != nil {return err}
	return nil
}

func (db *Database) Update(ctx context.Context,  modelPtr Model) error {
	return coreUpdate(ctx, db.Core, modelPtr, false, nil)
}
func (tx *Tx) Update(ctx context.Context,  modelPtr Model) error {
	return coreUpdate(ctx, tx.Core, modelPtr, false, nil)
}
func (db *Database) UpdateData(ctx context.Context,  modelPtr Model, data Data) error {
	return coreUpdate(ctx, db.Core, modelPtr, true, data)
}
func (tx *Tx) UpdateData(ctx context.Context,  modelPtr Model, data Data) error {
	return coreUpdate(ctx, tx.Core, modelPtr, true, data)
}
func coreUpdate (ctx context.Context, storage Storager, modelPtr Model, useUpdateData bool, updateData Data) error {
	rValue, isPtr := getPtrElem(modelPtr)
	if !isPtr { panic("db.Update() or db.TxUpdate()  arg `modelPtr` must be a ptr, eg: db.Update(&user) db.TxUpdate(tx, &user) ") }
	typeValue := reflect.TypeOf(modelPtr).Elem()
	fieldLen := rValue.NumField()
	var id interface{}
	var findID bool
	for i:=0;i<fieldLen;i++{
		item := rValue.Field(i)
		itemType := typeValue.Field(i)
		dbName, hasDBName := itemType.Tag.Lookup("db")
		if !hasDBName {
			continue
		}
		if dbName == "" {
			log.Print(`Maybe you forget set db:"name"` + itemType.Name)
			continue
		}
		value := item.Interface()
		if dbName == "id" {
			if item.IsZero() {
				panic(errors.New("update model.ID is zero"))
			}
			findID = true
			id = value
			continue
		}
		if useUpdateData {
			rValue.Set(reflect.ValueOf(updateData[Column(dbName)]))
		} else {
			updateData[Column(dbName)] = value
		}
	}
	if !findID {
		panic(errors.New("update not found id"))
	}
	updatedAtValue := rValue.FieldByName("UpdatedAt")
	if updatedAtValue.IsValid() {
		updatedAtType, _ := typeValue.FieldByName("UpdatedAt")
		updateData[Column(updatedAtType.Tag.Get("db"))] = time.Now()
		updatedAtValue.Set(reflect.ValueOf(time.Now()))
	}

	query, values := QB{
		Where: And("id", id),
		Update: updateData,
	}.BindModel(modelPtr).GetUpdate()
	result, err := storage.ExecContext(ctx, query, values...) ; if err != nil {return err}
	_, err = result.LastInsertId(); if err != nil {return err}
	return nil
}
func (db *Database) Transaction(ctx context.Context, transaction func(ftx *Tx) error) (error) {
	return db.TransactionOpts(ctx, nil, transaction)
}
func (db *Database) TransactionOpts(ctx context.Context, opts *sql.TxOptions, transaction func(*Tx) error) (txError error) {
	sqlxtx, err := db.Core.BeginTxx(ctx, opts) ; if err != nil {return err}
	tx := newTx(sqlxtx)
	defer func() {
		r := recover()
		if r != nil {
			rollbackErr := tx.Rollback() ; _= rollbackErr // 此时可以忽略 rollback 的错误
			panic(r)
		}
	}()
	err = transaction(tx)
	if err != nil {
		ge.Check(tx.Rollback())
		return err
	} else {
		return tx.commit()
	}
}
func (db *Database) Close() error {
	return db.Core.Close()
}