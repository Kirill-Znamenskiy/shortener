package crypto

import (
	"crypto/aes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
)

// GenerateSecretKey generate secret key with size len.
func GenerateSecretKey(size int) (ret []byte, err error) {
	ret = make([]byte, size)
	_, err = rand.Read(ret)
	if err != nil {
		ret = nil
		return
	}

	return
}

// MakeUserUUIDSign make userUUID sign with HMAC algorithm, use SHA256.
// If no error return slice of bytes with length 32, because of (256/8=32).
func MakeUserUUIDSign(userUUID *uuid.UUID, secretKey []byte) (ret []byte, err error) {

	h := hmac.New(sha256.New, secretKey)
	_, err = h.Write((*userUUID)[:])
	if err != nil {
		return
	}
	ret = h.Sum(nil)

	return
}

// EncryptMessage encrypt clearMsg message with AES.
// Return ecryptedMsg, same length as clearMsg.
func EncryptMessage(clearMsg []byte, secretKey []byte) (ecryptedMsg []byte, err error) {

	aesblock, err := aes.NewCipher(secretKey)
	if err != nil {
		return
	}

	bs := aesblock.BlockSize()

	if len(clearMsg)%bs != 0 {
		err = fmt.Errorf("invalid clearMsg size")
		return
	}

	bcnt := len(clearMsg) / bs

	ecryptedMsg = make([]byte, len(clearMsg))
	for i := 0; i < bcnt; i++ {
		aesblock.Encrypt(ecryptedMsg[bs*i:bs*(i+1)], clearMsg[bs*i:bs*(i+1)])
	}

	return
}

// DecryptMessage decrypt encryptedMsg message, encrypted with AES.
// Return clearMsg, same length as encryptedMsg.
func DecryptMessage(encryptedMsg []byte, secretKey []byte) (clearMsg []byte, err error) {
	aesblock, err := aes.NewCipher(secretKey)
	if err != nil {
		return
	}

	bs := aesblock.BlockSize()

	if len(encryptedMsg)%bs != 0 {
		err = fmt.Errorf("invalid msg size")
		return
	}

	bcnt := len(encryptedMsg) / bs

	clearMsg = make([]byte, len(encryptedMsg))
	for i := 0; i < bcnt; i++ {
		aesblock.Decrypt(clearMsg[bs*i:bs*(i+1)], encryptedMsg[bs*i:bs*(i+1)])
	}

	return
}

func EncryptAndSignUserUUID(userUUID *uuid.UUID, secretKey []byte) (ret []byte, err error) {
	//log.Printf("EncryptAndSignUserUUID(secretKey=%x)", secretKey)

	//log.Printf("userUUID: %x", userUUID)

	sign, err := MakeUserUUIDSign(userUUID, secretKey)
	if err != nil {
		return
	}

	aux := make([]byte, 0, len(*userUUID)+len(sign))
	aux = append(aux, (*userUUID)[:]...)
	aux = append(aux, sign...)

	//log.Printf("userUUIDSign: %x", aux)

	aux, err = EncryptMessage(aux, secretKey)
	if err != nil {
		return
	}

	//log.Printf("userUUIDSign ecrypted: %x", aux)

	ret = make([]byte, hex.EncodedLen(len(aux)))
	hex.Encode(ret, aux)

	return
}

func DecryptSignedUsedUUID(msg []byte, secretKey []byte) (ret *uuid.UUID, err error) {
	//log.Printf("DecryptSignedUsedUUID(secretKey=%x)", secretKey)

	// 16 for uuid + 32 for sha256 sign
	validSize := 16 + 32

	hexedMsg := msg

	if len(hexedMsg) != hex.EncodedLen(validSize) {
		err = fmt.Errorf("invalid sign format")
		return
	}

	msg = make([]byte, validSize)
	n, err := hex.Decode(msg, hexedMsg)
	if err != nil {
		return
	}
	if n != validSize {
		err = fmt.Errorf("invalid sign format")
		return
	}

	//log.Printf("msg ecrypted: %x", msg)

	msg, err = DecryptMessage(msg, secretKey)
	if err != nil {
		return
	}

	//log.Printf("msg: %x", msg)

	userUUID := uuid.UUID{}
	copy(userUUID[:], msg[:len(userUUID)])
	msgSign := msg[len(userUUID):]

	validSign, err := MakeUserUUIDSign(&userUUID, secretKey)
	if err != nil {
		return
	}

	//log.Printf("userUUID: %x", userUUID)
	//log.Printf("msgSign: %x", msgSign)
	//log.Printf("validSign: %x", validSign)

	if !hmac.Equal(msgSign, validSign) {
		err = fmt.Errorf("invalid sign")
		return
	}

	ret = &userUUID
	return
}
