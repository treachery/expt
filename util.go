package expt

import (
	"strconv"
	"strings"
)

func parseModSpec(spec string) (low, high int, err error) {
	inx := strings.Index(spec, "-")
	if inx <= 0 || len(spec) <= inx {
		err = ErrSpec
		return
	}
	if low, err = strconv.Atoi(spec[:inx]); err != nil {
		return
	}
	if high, err = strconv.Atoi(spec[inx+1:]); err != nil {
		return
	}
	if low < 0 || high > 99 || (low > high) {
		err = ErrSpec
		return
	}
	return
}
