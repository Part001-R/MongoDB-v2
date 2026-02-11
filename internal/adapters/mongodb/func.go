package mongodb

func IsExistsCollection(listNames []string, name string) bool {

	for _, n := range listNames {

		if n == name {
			return true
		}
	}

	return false
}
