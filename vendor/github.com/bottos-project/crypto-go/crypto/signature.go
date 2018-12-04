// Copyright 2017~2022 The Bottos Authors
// This file is part of the Bottos Chain library.
// Created by Rocket Core Team of Bottos.

//This program is free software: you can distribute it and/or modify
//it under the terms of the GNU General Public License as published by
//the Free Software Foundation, either version 3 of the License, or
//(at your option) any later version.

//This program is distributed in the hope that it will be useful,
//but WITHOUT ANY WARRANTY; without even the implied warranty of
//MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//GNU General Public License for more details.

//You should have received a copy of the GNU General Public License
// along with bottos.  If not, see <http://www.gnu.org/licenses/>.

/*
 * file description:  signature
 * @Author:
 * @Date:   2017-12-06
 * @Last Modified by:
 * @Last Modified time:
 */

package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"math/big"

	secp256k1 "github.com/bottos-project/crypto-go/crypto/secp256k1"
)

func GenerateKey() (pubkey, seckey []byte) {
	key, err := ecdsa.GenerateKey(secp256k1.S256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	return elliptic.Marshal(secp256k1.S256(), key.X, key.Y), PaddedBigBytes(key.D, 32)
}

func Sign(msg, seckey []byte) ([]byte, error) {
	sign, err := secp256k1.Sign(msg, seckey)
	return sign[:len(sign)-1], err
}

func VerifySign(pubkey, msg, sign []byte) bool {
	return secp256k1.VerifySignature(pubkey, msg, sign)
}

func PaddedBigBytes(bigint *big.Int, n int) []byte {
	if bigint.BitLen()/8 >= n {
		return bigint.Bytes()
	}
	ret := make([]byte, n)
	secp256k1.ReadBits(bigint, ret)
	return ret
}
