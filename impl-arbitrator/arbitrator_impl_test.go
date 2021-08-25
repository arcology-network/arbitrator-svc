package arbitrator_test

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

func TestStringAppend(t *testing.T) {
	begin := time.Now()
	var str string
	for i := 0; i < 50000; i++ {
		str += "blcc://eth1.0/accounts/0000111122223333444455556666777788889999/storage/containers/balance/$0000111122223333444455556666777788889999"
	}
	t.Log(len(str))
	t.Log(time.Duration(time.Since(begin)))
}

func TestStringAppend2(t *testing.T) {
	begin := time.Now()
	var buf bytes.Buffer
	for i := 0; i < 50000; i++ {
		fmt.Fprintf(&buf, "%s", "blcc://eth1.0/accounts/0000111122223333444455556666777788889999/storage/containers/balance/$0000111122223333444455556666777788889999")
	}
	str := buf.String()
	t.Log(len(str))
	t.Log(time.Duration(time.Since(begin)))
}
