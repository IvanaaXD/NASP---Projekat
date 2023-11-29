package bloom_filter

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math"
)

// Struktura bloom filtera
type BloomFilter struct {
	M             uint
	niz           []byte
	hash_funkcije []HashWithSeed // hash funkcije
}

// Kreiranje novog bloom filtera:
//
//	expectedElements - ocekivani broj elemenata
//	falsePositiveRate - tolerancija na greske
func NewBloomFilter(expectedElements int, falsePositiveRate float64) *BloomFilter {

	// Racunanje broja bitova (m) i broja hes funkcija (k)
	m := CalculateM(expectedElements, falsePositiveRate)
	k := CalculateK(expectedElements, m)

	// Zauzimanje broja bajtova za niz izracunati niz bitova (m)
	bytess := int(math.Ceil(float64(m) / 8))
	niz := make([]byte, bytess)

	hashf := CreateHashFunctions(k)

	bf := BloomFilter{m, niz, hashf}
	return &bf
}

// Brisanje bloom filtera podrazumeva vracanje niza na pocetnu vrednost
// Zauzimamo novi niz, jer ce garbage collector da dealocira memoriju za prethodni niz
func (bf *BloomFilter) DeleteBloomFilter() {
	bf.niz = make([]byte, bf.M)
}

// Dodavanje elementa u bloom filter
func (bf *BloomFilter) Add(newElement string) {

	newEl := []byte(newElement)
	for _, heshF := range bf.hash_funkcije {
		newElHash := heshF.Hash(newEl)

		bitEl := newElHash % uint64(bf.M) // bit u nizu
		byteEl := int(bitEl / 8)          // bajt u nizu
		bitMask := 1 << (bitEl % 8)       // maska za postavljanje bita na 1

		bf.niz[byteEl] |= byte(bitMask)
	}
}

// Pretraga elementa u bloom filteru
func (bf *BloomFilter) Search(element string) bool {

	el := []byte(element)
	for _, heshF := range bf.hash_funkcije {
		elHash := heshF.Hash(el)

		bitEl := elHash % uint64(bf.M) // bit u nizu
		byteEl := bitEl / 8            // bajt u nizu
		bitMask := 1 << (bitEl % 8)    // maska sa vrednoscu 1 na potrebnom bitu

		// AND operacija elementa sa maskom za proveru da li je element u bloom filteru
		if bf.niz[byteEl]&byte(bitMask) == 0 {
			return false
		}
	}
	return true
}

// Serijalizacija bloom-filtera
func (b BloomFilter) Save() []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(b)

	if err != nil {
		panic(err)
	}

	return buffer.Bytes()
}

// Deserijalizacija bloom-filtera
func Load(data []byte) *BloomFilter {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)

	bf := &BloomFilter{}
	err := decoder.Decode(bf)
	if err != nil {
		panic(err)
	}

	return bf
}

func Test() {
	fmt.Println("Test Bloom filter:\n")

	bf := NewBloomFilter(30, 2)
	bf.Add("Bojan")
	bf.Add("Mićo")
	bf.Add("Katarina")
	bf.Add("Milica")
	bf.Add("Miloš")
	fmt.Println("\nNemanja ? ", bf.Search("Nemanja"))
	fmt.Println("Jovo ? ", bf.Search("Jovo"))
	bf.Add("Branko")
	bf.Add("Gaga")
	bf.Add("Djuro")
	bf.Add("Suncica")
	bf.Add("Ljupka")
	bf.Add("Krinka")
	bf.Add("Djole")
	bf.Add("Mirjana")
	bf.Add("Jovo")
	bf.Add("Dado")
	bf.Add("Mira")
	fmt.Println("\nNemanja ? ", bf.Search("Nemanja"))
	fmt.Println("Jovo ? ", bf.Search("Jovo"))
}
