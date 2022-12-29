package esms

func isLegalPosition(p string) bool {
	if p == "GK" {
		return true
	}

	if len(p) != 3 {
		return false
	}

	pos := fullPosToPosition(p)
	side := fullPosToSide(p)

	return isLegalSide(side) && (pos == "DF" || pos == "DM" || pos == "MF" || pos == "AM" || pos == "FW")
}

func isLegalSide(s string) bool {
	return s == "L" || s == "R" || s == "C"
}

// Given a full position (like DML), get only
// the position (DM)
//
func fullPosToPosition(p string) string {
	return p[:2]
}

// Given full position (like DML), get only
// the side (L)
//
func fullPosToSide(p string) string {
	return p[2:]
}
