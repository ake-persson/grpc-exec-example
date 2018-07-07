package ts

import (
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
)

type Ts struct {
	time.Time
}

func Now() Ts {
	return Ts{time.Now()}
}

func Seconds(secs int64) Ts {
	return Ts{time.Unix(secs, 0)}
}

func Nanos(nanos int64) Ts {
	return Ts{time.Unix(0, nanos)}
}

func Timestamp(ts timestamp.Timestamp) Ts {
	secs := ts.GetSeconds()
	fnanos := int64(ts.GetNanos())
	nanos := secs*int64(time.Second) + fnanos
	return Ts{time.Unix(0, nanos)}
}

func (t Ts) Timestamp() timestamp.Timestamp {
	secs := t.Unix()
	nanos := t.UnixNano()
	fnanos := nanos - secs*int64(time.Second)
	return timestamp.Timestamp{Seconds: secs, Nanos: int32(fnanos)}
}

func (t Ts) String() string {
	return t.Format(time.RFC3339Nano)
}
