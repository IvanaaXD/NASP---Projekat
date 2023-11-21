package CountMinSketch

import (
	"bytes"
	"encoding/gob"
)

type CountMinSketch struct {
	m             uint           //kolone
	k             uint           //redovi
	table         [][]uint64     //tabela
	hashFunctions []HashWithSeed //hash f-je
}

func CreateCMS(epsilon, delta float64) *CountMinSketch {

	mParam := CalculateM(epsilon)
	kParam := CalculateK(delta)

	table := make([][]uint64, kParam)
	for i := 0; i < len(table); i++ {
		table[i] = make([]uint64, mParam)
	}

	hashFunctions := CreateHashFunctions(kParam)

	cms := CountMinSketch{m: mParam, k: kParam, table: table, hashFunctions: hashFunctions}
	return &cms

}

func (cms CountMinSketch) addItem(K []byte) {

	for row := 0; row < len(cms.hashFunctions); row++ {
		hashFunction := cms.hashFunctions[row]

		hash := hashFunction.Hash(K)
		col := hash % uint64(cms.m)

		cms.table[row][col] += 1
	}
}

func (cms CountMinSketch) getFrequency(K []byte) uint64 {

	minNum := ^uint64(0)

	for row := 0; row < len(cms.hashFunctions); row++ {
		hashFunction := cms.hashFunctions[row]

		hash := hashFunction.Hash(K)
		col := hash % uint64(cms.m)

		num := cms.table[row][col]

		if num < minNum {
			minNum = num
		}
	}

	return minNum
}

func (cms *CountMinSketch) GobEncode() ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	err := encoder.Encode(cms)
	if err != nil {
		panic("Error while encoding")
	}

	return buffer.Bytes(), nil
}

func GobDecode(data []byte) *CountMinSketch {
	var buffer bytes.Buffer
	decoder := gob.NewDecoder(&buffer)

	buffer.Write(data)

	cms := &CountMinSketch{}
	err := decoder.Decode(cms)
	if err != nil {
		panic("Error while decoding")
	}

	return cms
}
