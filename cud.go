package f
type CreateUpdateDeleteInput struct {
	ExistIDList []string
	UpdateIDList []string
}
type CreateUpdateDeleteOutput struct {
	DeleteIDList []string
}
func CreateUpdateDelete(data CreateUpdateDeleteInput) (output CreateUpdateDeleteOutput) {
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
