// Code generated by "stringer -type=Rank"; DO NOT EDIT.

package deck

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Two-2]
	_ = x[Three-3]
	_ = x[Four-4]
	_ = x[Five-5]
	_ = x[Six-6]
	_ = x[Seven-7]
	_ = x[Eight-8]
	_ = x[Nine-9]
	_ = x[Ten-10]
	_ = x[Jack-11]
	_ = x[Queen-12]
	_ = x[King-13]
	_ = x[Ace-14]
	_ = x[JokerRank-15]
}

const _Rank_name = "TwoThreeFourFiveSixSevenEightNineTenJackQueenKingAceJokerRank"

var _Rank_index = [...]uint8{0, 3, 8, 12, 16, 19, 24, 29, 33, 36, 40, 45, 49, 52, 61}

func (i Rank) String() string {
	i -= 2
	if i < 0 || i >= Rank(len(_Rank_index)-1) {
		return "Rank(" + strconv.FormatInt(int64(i+2), 10) + ")"
	}
	return _Rank_name[_Rank_index[i]:_Rank_index[i+1]]
}
