// Copyright (c) 2015 Taco de Wolff.
// Use of this source code is governed by a MIT style
// license that can be found in https://github.com/tdewolff/parse/blob/master/LICENSE.md

package lexer_old

// AsIdentifierName returns true if a valid identifier name is given.
func AsIdentifierName(b []byte) bool {
	if len(b) == 0 || !identifierStartTable[b[0]] {
		return false
	}

	i := 1
	for i < len(b) {
		if identifierTable[b[i]] {
			i++
		} else {
			return false
		}
	}
	return true
}

// AsDecimalLiteral returns true if a valid decimal literal is given.
func AsDecimalLiteral(b []byte) bool {
	if len(b) == 0 || (b[0] < '0' || '9' < b[0]) && (b[0] != '.' || len(b) == 1) {
		return false
	} else if b[0] == '0' {
		return len(b) == 1
	}
	i := 1
	for i < len(b) && '0' <= b[i] && b[i] <= '9' {
		i++
	}
	if i < len(b) && b[i] == '.' && b[0] != '.' {
		i++
		for i < len(b) && '0' <= b[i] && b[i] <= '9' {
			i++
		}
	}
	return i == len(b)
}

//func CreateToken(name []byte) string {
//	return strconv.FormatUint(name)
//}
