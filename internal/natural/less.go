package natural

func isdigit(v byte) bool {
	return '0' <= v && v <= '9'
}

func Less(a, b string) bool {
	ia, ib := 0, 0
	for ia < len(a) && ib < len(b) {
		ca, cb := a[ia], b[ib]

		da, db := isdigit(ca), isdigit(cb)
		if da != db {
			// one is digit, the other isn't
			return da
		}

		if !da { // && db
			// both letters
			if ca != cb {
				return ca < cb
			}
			ia++
			ib++
			continue
		}

		// digits
		for ; ia < len(a) && a[ia] == '0'; ia++ {
		}
		for ; ib < len(b) && b[ib] == '0'; ib++ {
		}

		// remember zero position
		nzia, nzib := ia, ib

		// advance over digits
		for ; ia < len(a) && isdigit(a[ia]); ia++ {
		}
		for ; ib < len(b) && isdigit(b[ib]); ib++ {
		}

		// different lengths, longer one is larger
		if lena, lenb := ia-nzia, ib-nzib; lena != lenb {
			return lena < lenb
		}

		// same length
		if sa, sb := a[nzia:ia], b[nzib:ib]; sa != sb {
			return sa < sb
		}

		// different no of zeros
		if nzia != nzib {
			return nzia < nzib
		}
	}

	// didn't find anything interesting
	return len(a) < len(b)
}
