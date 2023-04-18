package util

import (
	"strconv"
)

func DecimalToBinary16(value uint16) ([8]bool, [8]bool) {
	val := strconv.FormatInt(int64(value), 2)

	str := val
	for i := len([]rune(val)); i < 16; i++ {
		str = "0" + str
	}

	hbits := [8]bool{}
	lbits := [8]bool{}
	for i := 0; i < 16; i++ {
		if i < 8 {

            if str[i] == '0'{
                hbits[i] = false
            }else{
                hbits[i] = true
            }

		}else{

            if str[i] == '0'{
                lbits[i-8] = false
            }else{
                lbits[i-8] = true
            }
        }
	}

    return hbits, lbits

}
