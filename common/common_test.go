package common

import (
	"fmt"
	"testing"
	"time"
)

func TestDuration2Seconds(t *testing.T) {
	var tests = []struct {
		dur    string
		expect int
	}{
		{"1d", 60 * 60 * 24},
		{"1h", 60 * 60},
		{"10m", 60 * 10},
		{"2h4m10s", (60 * 120) + (4 * 60) + 10},
		{"88u", 0},
		{"19t77X what?4s", 4},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("duration-%s", tt.dur)
		t.Run(testname, func(t *testing.T) {
			seconds := duration2int(tt.dur)
			if seconds != tt.expect {
				t.Errorf("got %d, want %d", seconds, tt.expect)
			}
		})
	}
}

func TestIsExpired(t *testing.T) {
	var tests = []struct {
		expire string
		start  time.Time
		expect bool
	}{
		{"3s", time.Now(), true},
		{"1d", time.Now(), false},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("isexpired-%s-%s", tt.start, tt.expire)
		t.Run(testname, func(t *testing.T) {
			time.Sleep(5 * time.Second)
			got := IsExpired(tt.start, tt.expire)
			if got != tt.expect {
				t.Errorf("got %t, want %t", got, tt.expect)
			}
		})
	}
}

func TestUntaint(t *testing.T) {
	var tests = []struct {
		want    string
		input   string
		expect  string
		wanterr bool
	}{
		{`[^a-zA-Z0-9\-]`, "ab23-bb43-beef", "ab23-bb43-beef", false},
		{`[^a-zA-Z0-9\-]`, "`cat passwd`+ab23-bb43-beef", "catpasswdab23-bb43-beef", true},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("untaint-%s-%s", tt.want, tt.expect)
		t.Run(testname, func(t *testing.T) {
			untainted, err := Untaint(tt.input, tt.want)
			if untainted != tt.expect {
				t.Errorf("got %s, want %s", untainted, tt.expect)
			}
			if err != nil && !tt.wanterr {
				t.Errorf("got error: %s", err)
			}
		})
	}
}
