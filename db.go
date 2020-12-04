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
	result, execErr := storage.ExecContext(ctx, query, values...) ; if execErr != nil {return execErr}
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

func (db *Database) UpdateData(ctx context.Context,  modelPtr Model, updateData Data) error {
	return coreUpdateData(ctx, db.Core, modelPtr, updateData)
}
func (tx *Tx) UpdateData(ctx context.Context,  modelPtr Model, updateData Data) error {
	return coreUpdateData(ctx, tx.Core, modelPtr, updateData)
}
func coreUpdateData (ctx context.Context, storage Storager, modelPtr Model,  updateData Data) error {
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
	_, err := storage.ExecContext(ctx, query, values...) ; if err != nil {return err}
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
func (db *Database) QueryRowScan(ctx context.Context, has *bool, qb QB,  dest ...interface{}) error {
	return coreQueryRowScan(ctx, db.Core, has, qb, dest...)
}
func (tx *Tx) QueryRowScan(ctx context.Context, qb QB, has *bool,  dest ...interface{}) error {
	return coreQueryRowScan(ctx, tx.Core, has, qb, dest...)
}
func coreQueryRowScan(ctx context.Context, storage Storager, has *bool, qb QB, dest ...interface{}) error {
	query, values := qb.GetSelect()
	row := storage.QueryRowxContext(ctx, query, values...)
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