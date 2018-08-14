package yourtime

// Timemark is the data structure for a timemark
type Timemark struct {
	TimemarkID int64  `json:"timemarkID"`
	Author     string `json:"author"`
	AuthorURL  string `json:"authorURL"`
	Timemark   int64  `json:"timemark"`
	Content    string `json:"content"`
	Votes      int64  `json:"votes"`
	Date       int64  `json:"date"`
}
