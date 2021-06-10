package wallet

/**
 * Copyright (C), 2019-2021
 * FileName: creator
 * Author:   LinGuanHong
 * Date:     2021/3/29 3:00
 * Description:
 */
// https://github.com/vsergeev/btckeygenie

import (
	"bytes"
	"fmt"
	"io"
	"math/big"
)

/******************************************************************************/
/* ECDSA Keypair Generation */
/******************************************************************************/

var eccSecp256k1 EllipticCurve

func init() {
	/* See Certicom's SEC2 2.7.1, pg.15 */
	/* secp256k1 elliptic curve parameters */
	eccSecp256k1.P, _ = new(big.Int).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F", 16)
	eccSecp256k1.A, _ = new(big.Int).SetString("0000000000000000000000000000000000000000000000000000000000000000", 16)
	eccSecp256k1.B, _ = new(big.Int).SetString("0000000000000000000000000000000000000000000000000000000000000007", 16)
	eccSecp256k1.G.X, _ = new(big.Int).SetString("79BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F81798", 16)
	eccSecp256k1.G.Y, _ = new(big.Int).SetString("483ADA7726A3C4655DA4FBFC0E1108A8FD17B448A68554199C47D08FFB10D4B8", 16)
	eccSecp256k1.N, _ = new(big.Int).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141", 16)
	eccSecp256k1.H, _ = new(big.Int).SetString("01", 16)
}

// PublicKey represents a Bitcoin public key.
type PublicKey struct {
	Point
}

// PrivateKey represents a Bitcoin private key.
type PrivateKey struct {
	PublicKey
	D *big.Int
}

func NewPrivateKey(d *big.Int) *PrivateKey {
	key := &PrivateKey{D: d}
	key.derive()
	return key
}

// derive derives a Bitcoin public key from a Bitcoin private key.
func (priv *PrivateKey) derive() (pub *PublicKey) {
	/* See Certicom's SEC1 3.2.1, pg.23 */

	/* Derive public key from Q = d*G */
	Q := eccSecp256k1.ScalarBaseMult(priv.D)

	/* Check that Q is on the curve */
	if !eccSecp256k1.IsOnCurve(Q) {
		panic("Catastrophic math logic failure in public key derivation.")
	}

	priv.X = Q.X
	priv.Y = Q.Y

	return &priv.PublicKey
}

// GenerateKey generates a public and private key pair using random source rand.
func GenerateKey(rand io.Reader) (priv PrivateKey, err error) {
	/* See Certicom's SEC1 3.2.1, pg.23 */
	/* See NSA's Suite B Implementer’s Guide to FIPS 186-3 (ECDSA) A.1.1, pg.18 */

	/* Select private key d randomly from [1, n) */

	/* Read N bit length random bytes + 64 extra bits  */
	b := make([]byte, eccSecp256k1.N.BitLen()/8+8)
	_, err = io.ReadFull(rand, b)
	if err != nil {
		return priv, fmt.Errorf("Reading random reader: %v", err)
	}

	d := new(big.Int).SetBytes(b)

	/* Mod n-1 to shift d into [0, n-1) range */
	d.Mod(d, new(big.Int).Sub(eccSecp256k1.N, big.NewInt(1)))
	/* Add one to shift d to [1, n) range */
	d.Add(d, big.NewInt(1))

	priv.D = d

	/* Derive public key from private key */
	priv.derive()

	return priv, nil
}

// ToBytes converts a Bitcoin private key to a 32-byte byte slice.
func (priv *PrivateKey) ToBytes() (b []byte) {
	d := priv.D.Bytes()

	/* Pad D to 32 bytes */
	padded_d := append(bytes.Repeat([]byte{0x00}, 32-len(d)), d...)

	return padded_d
}

// ToBytes converts a Bitcoin public key to a 33-byte byte slice with point compression.
func (pub *PublicKey) ToBytes() (b []byte) {
	/* See Certicom SEC1 2.3.3, pg. 10 */

	x := pub.X.Bytes()

	/* Pad X to 32-bytes */
	padded_x := append(bytes.Repeat([]byte{0x00}, 32-len(x)), x...)

	/* Add prefix 0x02 or 0x03 depending on ylsb */
	if pub.Y.Bit(0) == 0 {
		return append([]byte{0x02}, padded_x...)
	}

	return append([]byte{0x03}, padded_x...)
}

// ToBytesUncompressed converts a Bitcoin public key to a 65-byte byte slice without point compression.
func (pub *PublicKey) ToBytesUncompressed() (b []byte) {
	/* See Certicom SEC1 2.3.3, pg. 10 */

	x := pub.X.Bytes()
	y := pub.Y.Bytes()

	/* Pad X and Y coordinate bytes to 32-bytes */
	padded_x := append(bytes.Repeat([]byte{0x00}, 32-len(x)), x...)
	padded_y := append(bytes.Repeat([]byte{0x00}, 32-len(y)), y...)

	/* Add prefix 0x04 for uncompressed coordinates */
	return append([]byte{0x04}, append(padded_x, padded_y...)...)
}