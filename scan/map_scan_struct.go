package scan

import (
	"github.com/pkg/errors"
	"reflect"
)

type Scan struct {
	StructValuePtr reflect.Value
}
func New(valuePtr interface{}) Scan {
	structValuePtr := reflect.ValueOf(valuePtr)
	return Scan{
		StructValuePtr: structValuePtr,
	}
}

func (scan *Scan) MapScanStruct(data map[string]interface{}, relationData Relation, structPtr interface{}) {
	for _,relation := range relationData.Single {
		for structFieldIndex, dbFieldName := range relation.DBTag {
			dataKey := relation.TableName + "." + dbFieldName
			dbValue, has := data[dataKey]
			if !has {
				panic(errors.WithStack(errors.New(dataKey + "is not found")))
			}
			scan.StructValuePtr.Elem().Field(relation.FieldIndex).Field(structFieldIndex).Set(reflect.ValueOf(dbValue))
		}
	}
	for _,relation := range relationData.Many {
		newManyItem := reflect.MakeSlice(scan.StructValuePtr.Elem().Field(relation.FieldIndex).Type(), 1,1).Index(0)
		for structFieldIndex, dbFieldName := range relation.DBTag {
			dataKey := relation.TableName + "." + dbFieldName
			dbValue, has := data[dataKey]
			if !has {
				panic(errors.WithStack(errors.New(dataKey + "is not found")))
			}
			newManyItem.Field(structFieldIndex).Set(reflect.ValueOf(dbValue))
		}
		scan.StructValuePtr.Elem().Field(relation.FieldIndex).Set(reflect.Append(scan.StructValuePtr.Elem().Field(relation.FieldIndex), newManyItem))
	}

}