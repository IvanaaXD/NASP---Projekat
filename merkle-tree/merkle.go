package merkle_tree

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
)

// HELPER FUNCTIONS

func Hash(data []byte) []byte {
	hash := sha1.Sum(data)
	return hash[:]
}

type MerkleRoot struct {
	root *Node
}

type Node struct {
	data  []byte
	left  *Node
	right *Node
}

func (mr *MerkleRoot) String() string {
	return mr.root.String()
}
func (n *Node) String() string {
	return hex.EncodeToString(n.data[:])
}

// Serijalizacija merkle stabla po nivoima (od korena ka listovima)
func SerializeMerkleTree(root *Node, file *os.File) {
	if root == nil {
		return
	}

	fmt.Fprintln(file, root.String())
	SerializeMerkleTree(root.left, file)
	SerializeMerkleTree(root.right, file)
}

// Kreiranje merkle stabla
func BuildMerkleTree(data [][]byte) *MerkleRoot {
	if len(data) == 0 {
		return nil
	}

	// Niz listova Merkle stabla
	var nodes []*Node
	for _, d := range data {
		nodes = append(nodes, &Node{data: Hash(d)})
	}

	// Petlja za kreiranje merkle stabla
	for len(nodes) > 1 {
		var newNodes []*Node
		for i := 0; i < len(nodes); i += 2 {
			var left, right *Node

			// Levo podstablo
			left = nodes[i]

			// Provera da li postoji dovoljno čvorova za desno podstablo
			if i+1 < len(nodes) {
				right = nodes[i+1]
			} else {
				// Ako nema dovoljno čvorova, kreiraj prazan čvor (hash praznih podataka)
				right = &Node{data: Hash([]byte{})}
			}

			// Spajanje podataka čvorova i kreiranje novog čvora
			// '...' - dodavanje podataka niza desnog podstabla nakon podataka levog podstabla
			data := append(left.data[:], right.data[:]...)
			newNodes = append(newNodes, &Node{data: Hash(data), left: left, right: right})
		}
		// Prelazak na nivo roditelja
		nodes = newNodes
	}

	return &MerkleRoot{root: nodes[0]}
}

func Test() {
	data := [][]byte{[]byte("data1"), []byte("data2"), []byte("data3")}
	merkleRoot := BuildMerkleTree(data)

	file, _ := os.Create("merkle-tree/merkle_tree.txt")
	defer file.Close()

	SerializeMerkleTree(merkleRoot.root, file)

}
