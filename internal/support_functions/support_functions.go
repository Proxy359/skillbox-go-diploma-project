package supportFunctions

var countruList = [15]string{"RU", "US", "GB", "FR", "BL", "AT", "BG", "DK", "CA", "ES", "CH", "TR", "PE", "NZ", "MC"}
var providersList = [3]string{"Topolo", "Rond", "Kildy"}

func CheckCountry(str string) bool {
	for i := 0; i < len(countruList); i++ {
		if str == countruList[i] {
			return true
		}
	}
	return false
}

func CheckProvider(str string) bool {
	for i := 0; i < len(providersList); i++ {
		if str == providersList[i] {
			return true
		}
	}
	return false
}
