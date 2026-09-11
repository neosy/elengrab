package qkeys

type QueryKey string

func (k QueryKey) String() string {
	return string(k)
}

func (k QueryKey) Short() string {
	return Keys.ShortKeyByKey(k.String())
}
