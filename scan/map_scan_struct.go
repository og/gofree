package scan

import (
	"reflect"
)

type Scan struct {
	ValuePtr reflect.Value
}
func New(valuePtr interface{}) Scan {
	return Scan{
		ValuePtr: reflect.ValueOf(valuePtr),
	}
}
func (scan *Scan) MapScanSlice(data map[string]interface{}, relationData Relation) {

}
func (scan *Scan) MapScanStruct(data map[string]interface{}, relationData Relation) {
	for _,relation := range relationData.Single {
		for structFieldIndex, dbFieldName := range relation.DBTag {
			dataKey := relation.TableName + "." + dbFieldName
			dbValue, has := data[dataKey]
			if !has {
			} else {
				scan.ValuePtr.Elem().Field(relation.FieldIndex).Field(structFieldIndex).Set(reflect.ValueOf(dbValue))
			}

		}
	}
	for _,relation := range relationData.Many {
		newManyItem := reflect.MakeSlice(scan.ValuePtr.Elem().Field(relation.FieldIndex).Type(), 1,1).Index(0)
		for structFieldIndex, dbFieldName := range relation.DBTag {
			dataKey := relation.TableName + "." + dbFieldName
			dbValue, has := data[dataKey]
			if !has {
			} else {
				newManyItem.Field(structFieldIndex).Set(reflect.ValueOf(dbValue))
			}

		}
		scan.ValuePtr.Elem().Field(relation.FieldIndex).Set(reflect.Append(scan.ValuePtr.Elem().Field(relation.FieldIndex), newManyItem))
	}

}