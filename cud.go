package f

type CUDData struct {
	ExistIDList []string
	UpdateIDList []string
}
type CUDOutput struct {
	DeleteIDList []string
}
func CUD(data CUDData) (output CUDOutput) {
	updateIDMap := map[string/* id */]bool{}
	for _, updateID := range data.UpdateIDList {
		updateIDMap[updateID] = true
	}
	for _, existID := range data.ExistIDList {
		hasUpdate := updateIDMap[existID]
		if !hasUpdate {
		output.DeleteIDList = append(output.DeleteIDList, existID)
		}
	}
	return
}
