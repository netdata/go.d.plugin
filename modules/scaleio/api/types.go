package api

// For all 4xx and 5xx return codes, the body may contain an apiError instance
// with more specifics about the failure.
type apiError struct {
	Message        string
	HTTPStatusCode int
	ErrorCode      int
}

func (e apiError) Error() string {
	return e.Message
}

type Version struct {
	Major int64
	Minor int64
}

//type System struct {
//	ID                int64
//	DaysInstalled     int64
//	SystemVersionName string
//}

type Bwc struct {
	NumOccured      int64
	NumSeconds      int64
	TotalWeightInKb int64
}
