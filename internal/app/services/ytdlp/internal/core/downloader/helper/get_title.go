package helper

func GetTitleFast(url string) (string, error) {
	info, err := GetInfoFast(url)
	if err != nil {
		return "", err
	}
	return info.Title, nil
}
