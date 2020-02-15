package core

import (
	"bufio"
	"bytes"
	"errors"

	"github.com/OguzhanE/saltyrtc-server-go/pkg/base"
	"github.com/ugorji/go/codec"
	"golang.org/x/crypto/nacl/box"
)

var (
	// ErrNotAllowedMessage occurs when you are trying to relay a message to invalid dest
	ErrNotAllowedMessage = errors.New("not allowed message")
	// ErrNotMatchedIdentities occurs when you identities dont match for two different source
	ErrNotMatchedIdentities = errors.New("identities dont match")
	// ErrNotAuthenticatedClient occurs when you are trying to encrypt message
	ErrNotAuthenticatedClient = errors.New("client is not authenticated")
	// ErrMessageTooShort occurs when the length of message less than expected
	ErrMessageTooShort = errors.New("message is too short")
	// ErrCantDecodePayload occurs when try to decode payload
	ErrCantDecodePayload = errors.New("cant decode payload")
	// ErrFieldNotExist occurs when a field should exist but it doesnt
	ErrFieldNotExist = errors.New("field doesnt exist")
	// ErrInvalidFieldValue occurs when fiel value is not valid
	ErrInvalidFieldValue = errors.New("invalid field value")
	// ErrCantDecryptPayload occurs when try to decrypt payload
	ErrCantDecryptPayload = errors.New("cant decrypt payload")
)

func encodePayload(payload interface{}) ([]byte, error) {
	b := new(bytes.Buffer)
	bw := bufio.NewWriter(b)
	h := new(codec.MsgpackHandle)
	h.WriteExt = true
	enc := codec.NewEncoder(bw, h)
	err := enc.Encode(payload)
	if err != nil {
		bw.Flush()
		return nil, err
	}
	err = bw.Flush()
	return b.Bytes(), err
}

func decodePayload(encodedPayload []byte) (PayloadUnion, error) {
	h := new(codec.MsgpackHandle)
	h.WriteExt = true
	h.ErrorIfNoField = true
	dec := codec.NewDecoderBytes(encodedPayload, h)
	v := PayloadUnion{}
	err := dec.Decode(&v)
	return v, err
}

func encryptPayload(client *Client, nonce []byte, encodedPayload []byte) ([]byte, error) {
	var nonceArr [base.NonceLength]byte
	copy(nonceArr[:], nonce[:base.NonceLength])
	return box.Seal(nil, encodedPayload, &nonceArr, &client.ClientKey, &client.ServerSessionBox.Sk), nil
}

func decryptPayload(client *Client, nonce []byte, data []byte) ([]byte, error) {
	var nonceArr [base.NonceLength]byte
	copy(nonceArr[:], nonce[:base.NonceLength])
	decryptedData, ok := box.Open(nil, data, &nonceArr, &client.ClientKey, &client.ServerSessionBox.Sk)
	if !ok {
		return nil, ErrCantDecryptPayload
	}
	return decryptedData, nil
}

// PayloadPacker ..
type PayloadPacker interface {
	Pack(client *Client, nonceReader NonceReader) ([]byte, error)
}

// PayloadUnion ..
type PayloadUnion struct {
	Type               base.MessageType `codec:"type"`
	Key                []byte           `codec:"key,omitempty"`
	YourCookie         []byte           `codec:"your_cookie,omitempty"`
	Subprotocols       []string         `codec:"subprotocols,omitempty"`
	PingInterval       uint32           `codec:"ping_interval,omitempty"`
	YourKey            []byte           `codec:"your_key,omitempty"`
	InitiatorConnected bool             `codec:"initiator_connected,omitempty"`
	Responders         []uint16         `codec:"responders,omitempty"`
	SignedKeys         []byte           `codec:"signed_keys,omitempty"`
	Id                 interface{}      `codec:"id,omitempty"`
	Reason             int              `codec:"reason,omitempty"`
}

// PayloadFieldError ..
type PayloadFieldError struct {
	Type  string
	Field string
	Err   error
}

// NewPayloadFieldError creates PayloadFieldError instance
func NewPayloadFieldError(payloadType string, field string, err error) *PayloadFieldError {
	return &PayloadFieldError{
		Type:  payloadType,
		Field: field,
		Err:   err,
	}
}

// Error ..
func (e *PayloadFieldError) Error() string {
	return e.Type + "." + e.Field + ": " + e.Err.Error()
}