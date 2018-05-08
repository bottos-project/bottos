package common

func Filter(sources []string, filters []string) []string {
	var tmpList []string
	var has bool
	for _, src := range sources {
		has = false
		for _, filter := range filters {
			if src == filter {
				has = true
				break
			}
		}
		if has == false {
			tmpList = append(tmpList, src)
		}
	}
	return tmpList
}
