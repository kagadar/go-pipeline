package must

import (
	"errors"
	"testing"
)

func TestOk(t *testing.T) {
	defer func() {
		if got := recover(); got != nil {
			t.Errorf("Ok(nil) panic got %v, want <nil>", got)
		}
	}()
	Ok(nil)
}

func TestOk_Panic(t *testing.T) {
	err := errors.New("panic")
	defer func() {
		if got := recover(); got != err {
			t.Errorf("Ok(panic) panic got %v, want %v", got, err)
		}
	}()
	Ok(err)
}

func TestDo(t *testing.T) {
	defer func() {
		if got := recover(); got != nil {
			t.Errorf("Do() panic got %v, want <nil>", got)
		}
	}()
	if got := Do(1, nil); got != 1 {
		t.Errorf("Do(1,nil) got %d, want 1", got)
	}
}

func TestDo_Panic(t *testing.T) {
	err := errors.New("panic")
	defer func() {
		if got := recover(); got != err {
			t.Errorf("Do(panic) panic got %v, want %v", got, err)
		}
	}()
	Do(1, err)
}

func TestZeroErr(t *testing.T) {
	want := errors.New("fail")
	if got, err := ZeroErr(1, nil); got != 1 || err != nil {
		t.Errorf("ZeroErr(1,nil) got %d %v, want 1 <nil>", got, err)
	}
	if got, err := ZeroErr(1, want); got != 0 || err != want {
		t.Errorf("ZeroErr(1,fail) got %d %v, want 0 %v", got, err, want)
	}
}
