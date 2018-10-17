package utilx

import "time"

const (
	csvFormat = "02-01-2006 15:04:05"
)

func DateToCSVString(d int64) string {
	t := time.Unix(d, 0)
	t = t.UTC()
	return t.Format(csvFormat)

}

func CSVStringToDate(date string) (int64, error) {
	t, err := time.Parse(csvFormat, date)
	if err != nil {
		return -1, err
	}
	return t.Unix(), nil
}
