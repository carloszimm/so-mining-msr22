package types

import "time"

const layout string = "2006-01-02 15:04:05"

type DateTime struct {
	time.Time
}

// Convert the internal date as CSV string
func (date *DateTime) MarshalCSV() (string, error) {
	if !date.IsZero() {
		return date.Time.Format(layout), nil
	} else {
		return "", nil
	}
}

// Convert the CSV string as internal date
func (date *DateTime) UnmarshalCSV(csv string) (err error) {
	if csv != "" {
		date.Time, err = time.Parse(layout, csv)
	} else {
		date.Time, err = time.Time{}, nil
	}
	return err
}

type Post struct {
	Id                    int      `csv:"Id"`
	PostTypeId            int8     `csv:"PostTypeId"`
	AcceptedAnswerId      int      `csv:"AcceptedAnswerId"`
	ParentId              int      `csv:"ParentId"`
	CreationDate          DateTime `csv:"CreationDate"`
	DeletionDate          DateTime `csv:"DeletionDate"`
	Score                 int      `csv:"Score"`
	ViewCount             int      `csv:"ViewCount"`
	Body                  string   `csv:"Body"`
	OwnerUserId           int      `csv:"OwnerUserId"`
	OwnerDisplayName      string   `csv:"OwnerDisplayName"`
	LastEditorUserId      int      `csv:"LastEditorUserId"`
	LastEditorDisplayName string   `csv:"LastEditorDisplayName"`
	LastEditDate          DateTime `csv:"LastEditDate"`
	LastActivityDate      DateTime `csv:"LastActivityDate"`
	Title                 string   `csv:"Title"`
	Tags                  string   `csv:"Tags"`
	AnswerCount           int      `csv:"AnswerCount"`
	CommentCount          int      `csv:"CommentCount"`
	FavoriteCount         int      `csv:"FavoriteCount"`
	ClosedDate            DateTime `csv:"ClosedDate"`
	CommunityOwnedDate    DateTime `csv:"CommunityOwnedDate"`
	ContentLicense        string   `csv:"ContentLicense"`
}
