package dashboard

import "time"

func mustTestTime(value string) time.Time {
	t, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		panic(err)
	}
	return t
}

func sampleRecords() []LogRecord {
	return []LogRecord{
		{
			TS:         mustTestTime("2026-03-14T09:04:00Z"),
			Level:      "INFO",
			Path:       "/api/orders",
			Status:     200,
			DurationMS: 50,
		},
		{
			TS:         mustTestTime("2026-03-14T09:00:00Z"),
			Level:      "WARN",
			Path:       "/api/login",
			Status:     401,
			DurationMS: 150,
		},
		{
			TS:         mustTestTime("2026-03-14T09:07:00Z"),
			Level:      "ERROR",
			Path:       "/api/reports",
			Status:     503,
			DurationMS: 600,
		},
		{
			TS:         mustTestTime("2026-03-14T09:02:00Z"),
			Level:      "INFO",
			Path:       "/api/login",
			Status:     200,
			DurationMS: 300,
		},
		{
			TS:         mustTestTime("2026-03-14T09:06:00Z"),
			Level:      "INFO",
			Path:       "/api/orders",
			Status:     201,
			DurationMS: 20,
		},
	}
}
