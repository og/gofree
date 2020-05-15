package f_test

import (
	f "github.com/og/gofree"
	gtest "github.com/og/x/test"
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
type optionCUDOutout struct {
	UpdateOptionList []Option
	CreateOptionList []Option
	DeleteIDList []IDOption
}
func optionCUD (databaseOptionList []Option, requestOptionList []OptionDTO) (output optionCUDOutout) {
	output.DeleteIDList = []IDOption{}
	output.CreateOptionList = []Option{}
	output.UpdateOptionList = []Option{}
	cudOutput := f.CUD(f.CUDIDList{
		ExistIDList:  func() (idList []string){
			for _,option := range databaseOptionList {
				idList = append(idList, string(option.ID))
			}
			return
		}(),
		UpdateIDList: func() (idList []string) {
			for _,option := range requestOptionList {
				if option.IsCreate {
					output.CreateOptionList = append(output.CreateOptionList, Option{
						Name: option.Name,
					})
				} else {
					output.UpdateOptionList = append(output.UpdateOptionList, Option{
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
		output.DeleteIDList = append(output.DeleteIDList, IDOption(deleteID))
	}
	return
}
func TestCUD(t *testing.T) {
	as := gtest.NewAS(t)
	// cud
	{
		out := optionCUD([]Option{
			{ID: IDOption("1"), Name:"a",},
			{ID: IDOption("2"), Name:"b",},
			{ID: IDOption("3"), Name:"c",},
		}, []OptionDTO{
			// 删 1
			// 改 2 3
			{ID:IDOption("2"),Name:"bb"},
			{ID:IDOption("3"),Name:"cc"},
			// 增 4
			{Name:"dd", IsCreate: true},
		})
		as.Equal(out.UpdateOptionList, []Option{
			{IDOption("2"), "bb",},
			{IDOption("3"), "cc",},
		})
		as.Equal(out.CreateOptionList, []Option{
			{IDOption(""), "dd",},
		})
		as.Equal(out.DeleteIDList, []IDOption{
			IDOption("1"),
		})
	}
	// only delete
	{
		out := optionCUD([]Option{
			{ID: IDOption("1"), Name:"a",},
			{ID: IDOption("2"), Name:"b",},
			{ID: IDOption("3"), Name:"c",},
		}, []OptionDTO{

		})
		as.Equal(out.UpdateOptionList, []Option{

		})
		as.Equal(out.CreateOptionList, []Option{

		})
		as.Equal(out.DeleteIDList, []IDOption{
			IDOption("1"),IDOption("2"),IDOption("3"),
		})
	}
	// only update
	{
		out := optionCUD([]Option{
			{ID: IDOption("1"), Name:"a",},
			{ID: IDOption("2"), Name:"b",},
			{ID: IDOption("3"), Name:"c",},
		}, []OptionDTO{
			{ID: IDOption("1"), Name:"aa",},
			{ID: IDOption("2"), Name:"bb",},
			{ID: IDOption("3"), Name:"cc",},
		})
		as.Equal(out.UpdateOptionList, []Option{
			{IDOption("1"), "aa",},
			{IDOption("2"), "bb",},
			{IDOption("3"), "cc",},
		})
		as.Equal(out.CreateOptionList, []Option{

		})
		as.Equal(out.DeleteIDList, []IDOption{

		})
	}
	// only create
	{
		out := optionCUD([]Option{

		}, []OptionDTO{
			{Name:"a",IsCreate: true,},
			{Name:"b",IsCreate: true,},
			{Name:"c",IsCreate: true,},
		})
		as.Equal(out.UpdateOptionList, []Option{

		})
		as.Equal(out.CreateOptionList, []Option{
			{Name:"a",},
			{Name:"b",},
			{Name:"c",},
		})
		as.Equal(out.DeleteIDList, []IDOption{

		})
	}
	// delete create
	{
		out := optionCUD([]Option{
			{ID: IDOption("1"), Name:"a",},
			{ID: IDOption("2"), Name:"b",},
			{ID: IDOption("3"), Name:"c",},
		}, []OptionDTO{
			{Name:"e",IsCreate: true,},
			{Name:"d",IsCreate: true,},
		})
		as.Equal(out.UpdateOptionList, []Option{

		})
		as.Equal(out.CreateOptionList, []Option{
			{Name:"e",},
			{Name:"d",},
		})
		as.Equal(out.DeleteIDList, []IDOption{
			IDOption("1"),IDOption("2"),IDOption("3"),
		})
	}
	// update create
	{
		out := optionCUD([]Option{
			{ID: IDOption("1"), Name:"a",},
			{ID: IDOption("2"), Name:"b",},
			{ID: IDOption("3"), Name:"c",},
		}, []OptionDTO{
			{ID: IDOption("1"), Name:"aa",},
			{ID: IDOption("2"), Name:"bb",},
			{ID: IDOption("3"), Name:"cc",},
			{Name:"e",IsCreate: true,},
			{Name:"d",IsCreate: true,},
		})
		as.Equal(out.UpdateOptionList, []Option{
			{IDOption("1"), "aa",},
			{IDOption("2"), "bb",},
			{IDOption("3"), "cc",},
		})
		as.Equal(out.CreateOptionList, []Option{
			{Name:"e",},
			{Name:"d",},
		})
		as.Equal(out.DeleteIDList, []IDOption{

		})
	}
}
