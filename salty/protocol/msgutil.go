package protocol

import (
	"errors"

	"github.com/OguzhanE/saltyrtc-server-go/pkg/crypto/nacl"
)

// IsValidYourCookieBytes ..
func IsValidYourCookieBytes(pk interface{}) bool {
	if pk == nil {
		return false
	}
	b, ok := pk.([]byte)
	if ok && len(b) == CookieLength {
		return true
	}
	return false
}

// ParseYourCookie ..
func ParseYourCookie(pk interface{}) ([]byte, error) {
	if !IsValidYourCookieBytes(pk) {
		return nil, errors.New("invalid your_cookie")
	}
	b, _ := pk.([]byte)
	return b, nil
}

// IsValidSubprotocols ..
func IsValidSubprotocols(subprotocols interface{}) bool {
	if subprotocols == nil {
		return false
	}
	_, ok := subprotocols.([]string)
	return ok
}

// ParseSubprotocols ..
func ParseSubprotocols(subprotocols interface{}) ([]string, error) {
	if !IsValidSubprotocols(subprotocols) {
		return nil, errors.New("invalid subprotocols")
	}
	val, _ := subprotocols.([]string)
	return val, nil
}

// IsValidPingInterval ..
func IsValidPingInterval(pingInterval interface{}) bool {
	if pingInterval == nil {
		return false
	}
	v, ok := pingInterval.(int)
	if ok && v >= 0 {
		return true
	}
	return false
}

// ParsePingInterval ..
func ParsePingInterval(pingInterval interface{}) (int, error) {
	if !IsValidPingInterval(pingInterval) {
		return 0, errors.New("invalid ping_interval")
	}
	val, _ := pingInterval.(int)
	return val, nil
}

// IsValidYourKey ..
func IsValidYourKey(yourKey interface{}) bool {
	return nacl.IsValidBoxPkBytes(yourKey)
}

// ParseYourKey ..
func ParseYourKey(yourKey interface{}) ([KeyBytesSize]byte, error) {
	yourKeyBytes, err := nacl.ConvertBoxPkToBytes(yourKey)
	if err != nil {
		var tmpArr [KeyBytesSize]byte
		return tmpArr, err
	}
	return nacl.CreateBoxPkFromBytes(yourKeyBytes)
}

// IsValidAddressID checks whether id is a valid address
func IsValidAddressID(id interface{}) bool {
	if id == nil {
		return false
	}
	_, ok := id.(AddressType)
	return ok
}

// ParseAddressID parses id to address of type
func ParseAddressID(id interface{}) (AddressType, error) {
	if !IsValidAddressID(id) {
		return 0, errors.New("Invalid address id")
	}
	v, _ := id.(AddressType)
	return v, nil
}

// IsValidResponderAddressID returns true if id is a valid responder address
func IsValidResponderAddressID(id interface{}) bool {
	v, err := ParseAddressID(id)
	return err == nil && IsValidResponderAddressType(v)
}

// ParseResponderAddressID parses id as address of type
func ParseResponderAddressID(id interface{}) (AddressType, error) {
	if !IsValidResponderAddressID(id) {
		return 0, errors.New("Invalid responder address id")
	}
	v, _ := id.(AddressType)
	return v, nil
}

// IsValidReasonCode ..
func IsValidReasonCode(reason interface{}) bool {
	if reason == nil {
		return false
	}
	v, ok := reason.(int)

	if ok &&
		v == CloseCodeGoingAway ||
		v == CloseCodeSubprotocolError ||
		(v >= CloseCodePathFullError && v <= CloseCodeInvalidKey) {
		return true
	}
	return false
}

// ParseReasonCode ..
func ParseReasonCode(reason interface{}) (int, error) {
	if !IsValidReasonCode(reason) {
		return 0, errors.New("Invalid reason code")
	}
	v, _ := reason.(int)
	return v, nil
}
