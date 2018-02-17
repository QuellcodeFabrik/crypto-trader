package helper

func IsElementInArray(element string, list []string) bool {
    for _, listItem := range list {
        if listItem == element {
            return true
        }
    }
    return false
}