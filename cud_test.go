package f_test

import (
	f "github.com/og/gofree"
	"log"
	"testing"
)
type IDOption string
type Option struct {
	ID IDOption
	Name string
}
type OptionDTO struct {
	ID IDOption
	Name string
	IsCreate bool
}
func TestCUD(t *testing.T) {
	updateOptionList := []Option{}
	createOptionList := []Option{}
	deleteIDList := []IDOption{}
	databaseOptionList := []Option{
		{ID: IDOption("1"), Name:"a",},
		{ID: IDOption("2"), Name:"b",},
		{ID: IDOption("3"), Name:"c",},
	}
	requestOptionList := []OptionDTO{
		// 删 1
		// 改 2 3
		{ID:IDOption("2"),Name:"bb"},
		{ID:IDOption("3"),Name:"bb"},
		// 增 4
		{Name:"dd", IsCreate: true},
	}
	cudOutput := f.CUD(f.CUDData{
		ExistIDList:  func() (idList []string){
			for _,option := range databaseOptionList {
				idList = append(idList, string(option.ID))
			}
			return
		}(),
		UpdateIDList: func() (idList []string) {
			for _,option := range requestOptionList {
				if option.IsCreate {
					createOptionList = append(createOptionList, Option{
						Name: option.Name,
					})
				} else {
					updateOptionList = append(updateOptionList, Option{
						ID:   option.ID,
						Name: option.Name,
					})
					idList = append(idList, string(option.ID))
				}
			}
			return
		}(),
	})
	for _,deleteID :=  range cudOutput.DeleteIDList {
		deleteIDList = append(deleteIDList, IDOption(deleteID))
	}
	log.Print("updateOptionList", updateOptionList)
	log.Print("createOptionList", createOptionList)
	log.Print("deleteIDList", deleteIDList)
}
