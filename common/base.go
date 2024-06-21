package common

func ListContainersStr(s string, list []string) bool {
	for _, val := range list {
		if val == s {
			return true
		}
	}
	return false
}
