package keykeeper

import (
	"encoding/base64"
	"encoding/binary"
	"strings"
)

type CustomID interface {
	ToUint32(string) (uint32, error)
	FromUint32(uint32) string
	GetNextUint32() uint32
}

type Base64ID struct {
	CustomID
	counter uint32
}

func NewBase64ID(start uint32) *Base64ID {
	bid := new(Base64ID)
	bid.counter = start
	return bid
}

func (bid *Base64ID) FromUint32(id uint32) string {
	var bytes = make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, id)
	b64 := base64.StdEncoding.EncodeToString(bytes)
	b64 = strings.ReplaceAll(b64, "/", "_")
	b64 = strings.ReplaceAll(b64, "+", "-")
	b64 = strings.TrimLeft(b64, "A")
	return strings.TrimRight(b64, "=")
}

func (bid *Base64ID) ToUint32(shortID string) (uint32, error) {
	shortID = strings.ReplaceAll(shortID, "_", "/")
	shortID = strings.ReplaceAll(shortID, "-", "+")

	for len(shortID) < 6 {
		shortID = "A" + shortID
	}

	for len(shortID)%4 != 0 {
		shortID += "="
	}

	bytes, err := base64.StdEncoding.DecodeString(shortID)

	if err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint32(bytes), nil
}

func (bid *Base64ID) GetNextUint32() uint32 {
	next := bid.counter
	bid.counter += 1
	return next
}
