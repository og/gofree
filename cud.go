package f

type CUDIDList struct {
	ExistIDList []string
	UpdateIDList []string
}
type CUDInfo struct {
	DeleteIDList []string
}
func CUD(data CUDIDList) (output CUDInfo) {
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
