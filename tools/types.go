package tools

type Booking struct {
	PlaceID   int    `json:"place_id"`
	UserName  string `json:"user_name"`
	UserPhone string `json:"user_phone"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}
