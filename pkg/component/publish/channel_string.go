package publish

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[None-0]
	_ = x[Org-1]
}

// channelName topic
const channelName = "NoneOrg"

// channelIndex index
var channelIndex = [...]uint8{0, 4, 7}

// String chanel string
func (i Channel) String() string {
	if i < 0 || i >= Channel(len(channelIndex)-1) {
		return "Channel(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return channelName[channelIndex[i]:channelIndex[i+1]]
}
