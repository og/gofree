package scan

import (
	"github.com/pkg/errors"
	"reflect"
)

type RelationItem struct {
	FieldIndex int
	TableName string
	DBTag map[int]string `note:"index is structFieldIndex , value is  dbFieldName struct tag 'db'"`
}
type Relation struct {
	Single []RelationItem
	Many []RelationItem
}
func getStructRelation(value reflect.Value, index int) (mapRelationData RelationItem) {
	tableNameValue := value.MethodByName("TableName")
	if !tableNameValue.IsValid() {
		panic(errors.WithStack(errors.Errorf("%T should have method TableName() string {}", value.Interface(), )))
	}
	tableNameOutoput := tableNameValue.Call([]reflect.Value{})
	tableName := tableNameOutoput[0].String()
	if tableName == "" {
		panic(errors.WithStack(errors.Errorf("%T  TableName() can not return empty string \"\"", value.Interface(), )))
	}
	mapRelationData = RelationItem{
		FieldIndex: index,
		TableName:  tableName,
		DBTag:      map[int]string{},
	}
	reflectTableLen := value.NumField()
	reflectTableType := value.Type()
	for subIndex := 0; subIndex < reflectTableLen; subIndex++ {
		dbTag := reflect.StructTag(reflectTableType.Field(subIndex).Tag).Get("db")
		if dbTag != "" {
			mapRelationData.DBTag[subIndex] = dbTag
		}
	}
	return mapRelationData
}
func GetSliceRelation(slicePtr interface{}) (relation Relation) {
	newItemValue := reflect.MakeSlice(reflect.TypeOf(slicePtr), 1,1).Index(0)
	return getRelationByReflectValue(newItemValue)
}
func GetRelation(value interface{}) (relation Relation) {
	return getRelationByReflectValue(reflect.ValueOf(value))
}
func getRelationByReflectValue(relationValue reflect.Value) (relation Relation) {
	reflectItemLen := relationValue.NumField()
	for i:=0;i<reflectItemLen;i++ {
		reflectTable := relationValue.Field(i)
		switch reflectTable.Kind().String() {
			case "struct":
				relation.Single = append(relation.Single, getStructRelation(reflectTable, i))
			case "slice":
				reflectManyItem := reflect.MakeSlice(reflectTable.Type(), 1, 1).Index(0)
				relation.Many = append(relation.Many, getStructRelation(reflectManyItem, i))
		}
	}
	return
}