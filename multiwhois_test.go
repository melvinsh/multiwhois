package multiwhois

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestQueryDomainInfo(t *testing.T) {
	tests := []struct {
		domain      string
		expiration  time.Time
		isAvailable bool
	}{
		// .com
		{
			domain:      uuid.New().String() + ".com",
			isAvailable: true,
		},
		{
			domain:      "google.com",
			expiration:  time.Date(2028, 9, 14, 04, 00, 00, 0, time.UTC),
			isAvailable: false,
		},
		// .info
		{
			domain:      uuid.New().String() + ".info",
			isAvailable: true,
		},
		{
			domain:      "google.info",
			expiration:  time.Date(2023, 7, 31, 23, 57, 50, 0, time.UTC),
			isAvailable: false,
		},
		// .be
		{
			domain:      uuid.New().String() + ".be",
			isAvailable: true,
		},
		{
			domain:      "google.be",
			isAvailable: false,
		},
		// .nl
		{
			domain:      uuid.New().String() + ".nl",
			isAvailable: true,
		},
		{
			domain:      "google.nl",
			isAvailable: false,
		},
	}

	for _, test := range tests {
		info, err := QueryDomainInfo(test.domain)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if info.Domain != test.domain {
			t.Errorf("expected domain %s, got %s", test.domain, info.Domain)
		}

		if info.IsAvailable != test.isAvailable {
			t.Errorf("expected availability for %s to be %t, got %t\n%s", test.domain, test.isAvailable, info.IsAvailable, info.FullWhois)
		}

		if info.Expiration != test.expiration {
			t.Errorf("expected expiration for %s to be %s, got %s\n%s", test.domain, test.expiration, info.Expiration, info.FullWhois)
		}
	}
}
