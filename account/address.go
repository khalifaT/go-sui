package account

import (
	"encoding/hex"
	"fmt"
)

func IsValidAddress(address string) (bool, error) {
	if has0xPrefix(address) {
		address = address[2:]
	}
	addressBytes, err := hex.DecodeString(address)
	if err == nil && len(addressBytes)*2 == (ADDRESS_LENGTH) {
		return true, nil
	} else if err == nil {
		return false, fmt.Errorf("sui length address is not respected\n")
	} else {
		return false, err
	}

}
func has0xPrefix(str string) bool {
	return len(str) >= 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X')
}
