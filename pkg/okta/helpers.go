package okta

// UpdateStringsSet method removes the toRemove values from the source then adds the toAdd values to the source if they are not already presented.
func UpdateStringsSet(source, toRemove, toAdd []string) []string {
	// result slice with the sufficient capacity
	newSlice := make([]string, 0, len(source)+len(toAdd))

	// for fast access: map of the objects to remove
	toRemoveMap := map[string]struct{}{}
	for _, value := range toRemove {
		toRemoveMap[value] = struct{}{}
	}

	// for fast access: a map of all existing values
	allValuesMap := map[string]struct{}{}

	// copy the values from the original slice to the new slice if they are not in the toRemoveMap
	for _, value := range source {
		if _, ok := toRemoveMap[value]; !ok {
			newSlice = append(newSlice, value)
			allValuesMap[value] = struct{}{}
		}
	}

	// copy to the new slice values from toAdd if they are not already presented in the new slice
	for _, value := range toAdd {
		if _, ok := allValuesMap[value]; !ok {
			newSlice = append(newSlice, value)
			allValuesMap[value] = struct{}{}
		}
	}

	return newSlice
}
