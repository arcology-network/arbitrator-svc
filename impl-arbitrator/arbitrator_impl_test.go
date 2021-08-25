package arbitrator

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	codec "github.com/HPISTechnologies/common-lib/codec"
)

func TestStringAppend(t *testing.T) {
	t0 := time.Now()
	var str string
	for i := 0; i < 50000; i++ {
		str += "blcc://eth1.0/accounts/0000111122223333444455556666777788889999/storage/containers/balance/$0000111122223333444455556666777788889999"
	}
	t.Log(len(str))
	t.Log(time.Duration(time.Since(t0)))
}

func TestStringAppendWithBuffer(t *testing.T) {
	t0 := time.Now()
	var buf bytes.Buffer
	for i := 0; i < 50000; i++ {
		fmt.Fprintf(&buf, "%s", "blcc://eth1.0/accounts/0000111122223333444455556666777788889999/storage/containers/balance/$0000111122223333444455556666777788889999")
	}
	str := buf.String()
	t.Log(len(str))
	t.Log(time.Duration(time.Since(t0)))
}

func TestStringAppendCodecStrings(t *testing.T) {
	buf := make([]string, 50000)
	for i := 0; i < 50000; i++ {
		buf[i] = fmt.Sprintf("blcc://eth1.0/accounts/0000111122223333444455556666777788889999/storage/containers/balance/$0000111122223333444455556666777788889999")
	}

	t0 := time.Now()
	t.Log(len(codec.Strings(buf).Flatten()))
	t.Log(time.Duration(time.Since(t0)))
}
