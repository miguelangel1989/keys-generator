package main

import (
	"math/big"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
)

var largestBitcoinSeed = new(big.Int).SetBytes([]byte{
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFE,
	0xBA, 0xAE, 0xDC, 0xE6, 0xAF, 0x48, 0xA0, 0x3B, 0xBF, 0xD2, 0x5E, 0x8C, 0xD0, 0x36, 0x41, 0x40,
})

var one = big.NewInt(1)

type key struct {
	private      string
	compressed   string
	uncompressed string
}

func generateBitcoinKeys(firstSeed *big.Int, amount int) (keys []key) {
	var padded [32]byte

	bitcoinKeys := make([]key, 0, amount)

	for i := 0; i < amount; i++ {
		// Check to make sure we're not out of range
		if firstSeed.Cmp(largestBitcoinSeed) > 0 {
			break
		}

		// Copy firstSeed value's bytes to padded slice
		copy(padded[32-len(firstSeed.Bytes()):], firstSeed.Bytes())

		// Get private and public keys
		privKey, public := btcec.PrivKeyFromBytes(btcec.S256(), padded[:])

		// Get compressed and uncompressed addresses for public key
		caddr, _ := btcutil.NewAddressPubKey(public.SerializeCompressed(), &chaincfg.MainNetParams)
		uaddr, _ := btcutil.NewAddressPubKey(public.SerializeUncompressed(), &chaincfg.MainNetParams)

		// Encode addresses
		wif, _ := btcutil.NewWIF(privKey, &chaincfg.MainNetParams, false)

		bitcoinKeys = append(bitcoinKeys, key{
			private:      wif.String(),
			compressed:   caddr.EncodeAddress(),
			uncompressed: uaddr.EncodeAddress(),
		})

		firstSeed.Add(firstSeed, one)
	}

	return bitcoinKeys
}