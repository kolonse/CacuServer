package lib

import (
	"testing"
	"time"
)

func TestTimeFormat(t *testing.T) {
	t.Log(TimeFormat(time.Now(), "YYYY-MM-DD HH:mm:ss"))
	t.Log(TimeFormat(time.Now(), "YYYY-MM-DD 00:00:00"))
	t.Log(TimeFormat(time.Now(), "YYYY-MM-00 00:00:00"))
}
