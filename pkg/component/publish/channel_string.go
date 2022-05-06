package publish

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[None-0]
	_ = x[Org-1]
}

const _Channel_name = "NoneOrg"

var _Channel_index = [...]uint8{0, 4, 7}

func (i Channel) String() string {
	if i < 0 || i >= Channel(len(_Channel_index)-1) {
		return "Channel(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Channel_name[_Channel_index[i]:_Channel_index[i+1]]
}
