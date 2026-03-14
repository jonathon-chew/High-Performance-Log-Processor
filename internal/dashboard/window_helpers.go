package dashboard

import (
	"slices"
	"time"
)

// groupRecordsByWindow should contain the shared "split records into contiguous
// time buckets" behavior used by windowed metric functions.
// It should return one WindowBucket per bucket, with each bucket containing
// the raw records that fall within that bucket's time range.
// Specific metric functions can then transform those buckets into request counts,
// level totals, status summaries, or path-based aggregates.
func groupRecordsByWindow(records []LogRecord, bucketSize BucketSize) []WindowBucket {
	// panic("not implemented")
	if len(records) == 0 {
		return []WindowBucket{}
	}

	slices.SortFunc(records, func(a, b LogRecord) int {
		if a.TS.Unix() > b.TS.Unix() {
			return 1
		} else if a.TS.Unix() < b.TS.Unix() {
			return -1
		} else {
			return 0
		}
	})

	endTime := records[0].TS.Add(time.Duration(bucketSize))
	var returnWindowBucket []WindowBucket
	var tempLogRecords []LogRecord

	// Loop through log records
	for _, record := range records {

		if len(tempLogRecords) == 0 {
			tempLogRecords = append(tempLogRecords, record)
		} else if record.TS.Unix() <= endTime.Unix() {
			tempLogRecords = append(tempLogRecords, record)
		} else {

			if len(tempLogRecords) > 0 {
				returnWindowBucket = append(returnWindowBucket, WindowBucket{
					Window: TimeWindow{
						Start: tempLogRecords[0].TS,
						End:   tempLogRecords[len(tempLogRecords)-1].TS,
					},
					Records: tempLogRecords,
				})
			}

			endTime = record.TS.Add(time.Duration(bucketSize))
			tempLogRecords = []LogRecord{record}
		}
	}

	if len(tempLogRecords) > 0 {
		returnWindowBucket = append(returnWindowBucket, WindowBucket{
			Window: TimeWindow{
				Start: tempLogRecords[0].TS,
				End:   tempLogRecords[len(tempLogRecords)-1].TS,
			},
			Records: tempLogRecords,
		})
	}

	return returnWindowBucket
}
