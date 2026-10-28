package dtypes

type SSEEvent struct {
	Name string
	Data []byte

	Num   int
	Index string // dedup key: "userKey:dataHash"
}
