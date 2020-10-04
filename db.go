package f

import (
	"context"
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	ge "github.com/og/x/error"
	"reflect"
	"time"
)

type Database struct {
	Core *sqlx.DB
	onlyReadDataSourceName DataSourceName
}
func (database Database) Close() {
	ge.Check(database.Core.Close())
}
func (database Database) GetDataSourceName () DataSourceName {
	return database.onlyReadDataSourceName
}
func NewDatabase(dataSourceName DataSourceName) (database Database, err error) {
	db, err := sqlx.Connect(dataSourceName.DriverName, dataSourceName.GetString())
	if err != nil {
		return database, err
	}
	database = Database{Core: db,}
	database.onlyReadDataSourceName = dataSourceName
	return database, nil
}

func (database *Database) OneQB(modelPtr Model, has *bool, qb QB) After {
	return database.coreOneQB(sqlOpt{}, modelPtr, has, qb)
}
func (database *Database) OneQBContext(ctx context.Context, modelPtr Model, has *bool, qb QB) After {
	return database.coreOneQB(sqlOpt{Context: ctx}, modelPtr, has, qb)
}
func (database *Database) TxOneQB(tx *Tx, modelPtr Model, has *bool, qb QB) After {
	return database.coreOneQB(sqlOpt{Tx: tx.core}, modelPtr,has, qb)
}
func (database *Database) TxOneQBContext(ctx context.Context, tx *Tx, modelPtr Model, has *bool, qb QB) After {
	return database.coreOneQB(sqlOpt{Tx: tx.core,Context: ctx}, modelPtr,has, qb)
}
type sqlOpt struct {
	Tx *sqlx.Tx
	Context context.Context
}
func (opt sqlOpt) CtxOrBackground() context.Context {
	if opt.Context == nil {
		return context.Background()
	} else {
		return opt.Context
	}
}
func (database *Database) coreOneQB(opt sqlOpt, modelPtr Model, has *bool, qb QB) (after After) {
	_, isPtr := getPtrElem(modelPtr)
	if !isPtr {
		panic("db.OneID() or db.OneQB()  arg `modelPtr` must be a ptr")
	}
	scanModelMakeSQLSelect(reflect.ValueOf(modelPtr).Elem().Type(), &qb)
	query, values := qb.BindModel(modelPtr).GetSelect()
	after.ActualSQL = append(after.ActualSQL, query)
	var row *sqlx.Row
	if opt.Tx != nil  {
		row = opt.Tx.QueryRowxContext(opt.CtxOrBackground(), query, values...)
	} else {
		row = database.Core.QueryRowxContext(opt.CtxOrBackground(),query, values...)
	}
	err := row.StructScan(modelPtr)
	if err == sql.ErrNoRows { *has = false ; return}
	if err != nil {panic(err)}
	*has = true
	return
}



func (database *Database) OneID(modelPtr Model, has *bool, id interface{}) After {
	return database.OneQB(modelPtr, has, QB{
		Where:And("id", id),
	})
}

func (database *Database) TxOneID(tx *Tx, modelPtr Model, has *bool, id interface{}) After {
	return database.TxOneQB(tx, modelPtr, has, QB{
		Where:And("id", id),
	})
}
func (database *Database) CountQB(modelPtr Model, qb QB) (count int)  {
	return database.coreCountQB(sqlOpt{}, modelPtr, qb)
}
func (database *Database) CountQBContext(ctx context.Context, modelPtr Model, qb QB) (count int)  {
	return database.coreCountQB(sqlOpt{Context: ctx}, modelPtr, qb)
}
func (database *Database) TxCountQB(tx *Tx, modelPtr Model, qb QB) (count int) {
	return database.coreCountQB(sqlOpt{Tx: tx.core}, modelPtr, qb)
}
func (database *Database) TxCountQBContext(ctx context.Context, tx *Tx, modelPtr Model, qb QB) (count int) {
	return database.coreCountQB(sqlOpt{Tx: tx.core, Context: ctx}, modelPtr, qb)
}
func (database *Database) coreCountQB(opt sqlOpt, modelPtr Model, qb QB) (count int)  {
	qb.Count = true
	query, values := qb.BindModel(modelPtr).GetSelect()
	var row *sqlx.Row
	if opt.Tx != nil {
		row = opt.Tx.QueryRowxContext(opt.CtxOrBackground(),query, values...)
	} else {
		row = database.Core.QueryRowxContext(opt.CtxOrBackground(),query, values...)
	}
	err := row.Scan(&count)
	if err != nil {panic(err)}
	return
}

func (database *Database) ListQB(modelListPtr interface{}, qb QB) {
	database.coreListQB(sqlOpt{}, modelListPtr, qb)
}
func (database *Database) ListQBContext(ctx context.Context, modelListPtr interface{}, qb QB) {
	database.coreListQB(sqlOpt{Context: ctx}, modelListPtr, qb)
}
func (database *Database) TxListQB(tx *Tx, modelListPtr []Model, qb QB) {
	database.coreListQB(sqlOpt{Tx: tx.core}, modelListPtr, qb)
}
func (database *Database) TxListQBContext(ctx context.Context, tx *Tx, modelListPtr []Model, qb QB) {
	database.coreListQB(sqlOpt{Tx: tx.core, Context:ctx}, modelListPtr, qb)
}
func (database *Database) coreListQB(opt sqlOpt, modelListPtr interface{}, qb QB) {
	elementType := reflect.TypeOf(modelListPtr).Elem()
	reflectItemValue := reflect.MakeSlice(elementType, 1,1).Index(0)
	scanModelMakeSQLSelect(reflectItemValue.Type(), &qb)
	query, values := qb.BindModel(reflectItemValue.Addr().Interface().(Model)).GetSelect()
	if qb.Table == "" {
		tableName := reflectItemValue.MethodByName("TableName").Call([]reflect.Value{})[0].String()
		qb.Table = tableName
	}
	if opt.Tx != nil {
		err := opt.Tx.SelectContext(opt.CtxOrBackground(),modelListPtr, query, values...)
		ge.Check(err)
	} else {
		err := database.Core.SelectContext(opt.CtxOrBackground(),modelListPtr, query, values...)
		ge.Check(err)
	}
	return
}

func (database *Database) coreCreate(opt sqlOpt, modelPtr Model) {
	value, isPtr := getPtrElem(modelPtr)
	if !isPtr {
		panic("db.OneID() or db.OneQB()  arg `modelPtr` must be a ptr")
	}
	reflect.ValueOf(modelPtr).MethodByName("BeforeCreate").Call([]reflect.Value{})
	typeValue := reflect.TypeOf(modelPtr).Elem()
	fieldLen := value.NumField()
	insertData := map[Column]interface{}{}
	for i:=0;i<fieldLen;i++{
		item := value.Field(i)
		itemType := typeValue.Field(i)
		dbName := itemType.Tag.Get("db")
		if dbName == "" {
			continue
		}
		value := item.Interface()
		insertData[Column(dbName)] = value
	}
	createdAtValue := value.FieldByName("CreatedAt")
	if createdAtValue.IsValid() {
		createdAtType, _ := typeValue.FieldByName("CreatedAt")
		insertData[Column(createdAtType.Tag.Get("db"))] = time.Now()
		createdAtValue.Set(reflect.ValueOf(time.Now()))
	}
	updatedAtValue := value.FieldByName("UpdatedAt")
	if updatedAtValue.IsValid() {
		updatedAtType, _ := typeValue.FieldByName("UpdatedAt")
		insertData[Column(updatedAtType.Tag.Get("db"))] = time.Now()
		updatedAtValue.Set(reflect.ValueOf(time.Now()))
	}
	query, values := QB{
		Insert: insertData,
	}.BindModel(modelPtr).GetInsert()
	var result sql.Result
	if opt.Tx != nil {
		newResult, err := opt.Tx.ExecContext(opt.CtxOrBackground(),query, values...) ; ge.Check(err)
		result = newResult
	} else {
		newResult, err := database.Core.ExecContext(opt.CtxOrBackground(),query, values...) ; ge.Check(err)
		result = newResult
	}

	lastInsertID, err := result.LastInsertId() ; ge.Check(err)
	if  lastInsertID != 0 {
		value.FieldByName("ID").SetInt(lastInsertID)
	}

}

func (database *Database) Create(modelPtr Model) {
	database.coreCreate(sqlOpt{}, modelPtr)
}
func (database *Database) CreateContext(ctx context.Context, modelPtr Model) {
	database.coreCreate(sqlOpt{Context: ctx}, modelPtr)
}
func (database *Database) TxCreate(tx *Tx, modelPtr Model) {
	database.coreCreate(sqlOpt{Tx: tx.core}, modelPtr)
}
func (database *Database) TxCreateContext(ctx context.Context, tx *Tx, modelPtr Model) {
	database.coreCreate(sqlOpt{Tx: tx.core, Context: ctx}, modelPtr)
}

func (database *Database) coreDeleteQB(opt sqlOpt, modelPtr Model, qb QB) {
	_, isPtr := getPtrElem(modelPtr)
	if !isPtr {
		panic("db.DeleteQB() or db.TxDeleteQB()  arg `modelPtr` must be a ptr, eg: db.DeleteQB(&user, qb) db.TxDeleteQB(tx, &user, qb) ")
	}
	if len(qb.Update) == 0 {
		qb.Update = map[Column]interface{}{}
	}
	qb.Update["deleted_at"] = time.Now()
	query, values := qb.BindModel(modelPtr).GetUpdate()
	var result sql.Result
	if opt.Tx != nil {
		sqlResult, err := opt.Tx.ExecContext(opt.CtxOrBackground(),query, values...) ; ge.Check(err)
		result = sqlResult
	} else {
		sqlResult, err := database.Core.ExecContext(opt.CtxOrBackground(),query, values...) ; ge.Check(err)
		result = sqlResult
	}
	_, err := result.LastInsertId() ; ge.Check(err)
}
func (database *Database) DeleteQB(modelPtr Model, qb QB) {
	database.coreDeleteQB(sqlOpt{}, modelPtr, qb)
}
func (database *Database) DeleteQBContext(ctx context.Context, modelPtr Model, qb QB) {
	database.coreDeleteQB(sqlOpt{Context: ctx}, modelPtr, qb)
}
func (database *Database) TxDeleteQB(tx *Tx,modelPtr Model, qb QB) {
	database.coreDeleteQB(sqlOpt{Tx: tx.core}, modelPtr, qb)
}
func (database *Database) TxDeleteQBContext(ctx context.Context, tx *Tx,modelPtr Model, qb QB) {
	database.coreDeleteQB(sqlOpt{Tx: tx.core, Context:ctx}, modelPtr, qb)
}


func (database *Database) coreDelete(opt sqlOpt, modelPtr Model) {
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
	var result sql.Result
	if opt.Tx != nil {
		sqlResult, err := opt.Tx.ExecContext(opt.CtxOrBackground(),query, values...) ; ge.Check(err)
		result = sqlResult
	} else {
		sqlResult, err := database.Core.ExecContext(opt.CtxOrBackground(),query, values...) ; ge.Check(err)
		result = sqlResult
	}
	_, err := result.LastInsertId() ; ge.Check(err)
}
func (database *Database) Delete(modelPtr Model) {
	database.coreDelete(sqlOpt{}, modelPtr)
}
func (database *Database) DeleteContext(ctx context.Context, modelPtr Model) {
	database.coreDelete(sqlOpt{Context: ctx,}, modelPtr)
}
func (database *Database) TxDelete(tx *Tx,modelPtr Model, qb QB) {
	database.coreDelete(sqlOpt{Tx: tx.core}, modelPtr)
}
func (database *Database) TxDeleteContext(ctx context.Context, tx *Tx,modelPtr Model, qb QB) {
	database.coreDelete(sqlOpt{Tx: tx.core, Context:ctx}, modelPtr)
}

func (database *Database) Update(modelPtr Model) {
	database.coreUpdate(sqlOpt{}, modelPtr)
}
func (database *Database) UpdateContext(ctx context.Context, modelPtr Model) {
	database.coreUpdate(sqlOpt{Context: ctx}, modelPtr)
}
func (database *Database) TxUpdate(tx *Tx, modelPtr Model) {
	database.coreUpdate(sqlOpt{Tx: tx.core}, modelPtr)
}
func getPtrElem(ptr interface{}) (value reflect.Value, isPtr bool) {
	v := reflect.ValueOf(ptr)
	if v.Kind() != reflect.Ptr {
		isPtr = false
		return
	}
	value = v.Elem()
	isPtr = true
	return
}
func (database *Database) coreUpdate (opt sqlOpt, modelPtr Model) {
	rValue, isPtr := getPtrElem(modelPtr)
	if !isPtr {
		panic("db.Update() or db.TxUpdate()  arg `modelPtr` must be a ptr, eg: db.Update(&user) db.TxUpdate(tx, &user) ")
	}
	typeValue := reflect.TypeOf(modelPtr).Elem()
	fieldLen := rValue.NumField()
	updateData := map[Column]interface{}{}
	var id interface{}
	for i:=0;i<fieldLen;i++{
		item := rValue.Field(i)
		itemType := typeValue.Field(i)
		dbName := itemType.Tag.Get("db")
		value := item.Interface()
		if dbName == "id" {
			if item.IsZero() {
				panic(errors.New("db.Update(&model) or db.TxUpdate(&model) model.id is zero"))
			}
			id = value
			continue
		}
		updateData[Column(dbName)] = value
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
	var result sql.Result
	if opt.Tx != nil {
		newResult, err := opt.Tx.ExecContext(opt.CtxOrBackground(),query, values...) ; ge.Check(err)
		result = newResult
	} else {
		newResult, err := database.Core.ExecContext(opt.CtxOrBackground(),query, values...) ; ge.Check(err)
		result = newResult
	}
	lastInsertID, err := result.LastInsertId() ; ge.Check(err)
	if  lastInsertID != 0 {
		rValue.FieldByName("ID").SetInt(lastInsertID)
	}
}

func (database Database) Tx() *Tx {
	tx, err := database.Core.Beginx() ; ge.Check(err)
	return newTx(tx)
}