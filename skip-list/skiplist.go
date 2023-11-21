package skip_list

import (
	"fmt"
	"math/rand"
)

// Struktura cvora skip liste
type SkipNode struct {
	Value int
	Level int
	Next  []*SkipNode
}

// Struktura skip liste
type SkipList struct {
	Header    *SkipNode
	maxHeight int
}

// Kreiranje cvora skip liste
// Next: make([]*SkipNode, level+1) - ide level+1 jer mora postojati veza ka sledecem visem nivou
func NewSkipNode(value, level int) *SkipNode {
	return &SkipNode{Value: value, Level: level, Next: make([]*SkipNode, level+1)}
}

// Kreiranje skip liste
func NewSkipList(maxLvl int) *SkipList {
	return &SkipList{Header: NewSkipNode(0, maxLvl), maxHeight: maxLvl}
}

// Dodavanje novog elementa u skip listu
func (sl *SkipList) Insert(value int) bool {

	// Ako element postoji ne dodaje se u skip listu
	if sl.Search(value) {
		fmt.Println("Element", value, "već postoji u skip listi.")
		return false
	}

	// Pravimo niz u kome se za svaki nivo skip liste nalazi cvor nakon kog treba ubaciti novi cvor
	update := make([]*SkipNode, sl.maxHeight)
	current := sl.Header

	// Petlja za popunjavanje update niza
	for i := sl.maxHeight - 1; i >= 0; i-- {
		// Ako je vrednost sledeceg cvora manja od vrednosti novog cvora
		for current.Next[i] != nil && current.Next[i].Value < value {
			current = current.Next[i]
		}
		update[i] = current
	}

	// Racunanje nivoa za novi cvor
	Level := sl.roll()

	// Podizanje nivoa header pokazivaca na najvisi nivo
	if Level > sl.Header.Level {
		sl.Header.Level = Level
	}

	// Kreiranje novog cvora
	newNode := NewSkipNode(value, Level)

	// Ubacivanje novog cvora na svaki njegov nivo
	for i := 0; i <= Level && i < len(update); i++ {
		newNode.Next[i] = update[i].Next[i]
		update[i].Next[i] = newNode
	}

	return true
}

// Pretraga elementa u skip listi
func (sl *SkipList) Search(value int) bool {
	current := sl.Header
	// Nalazenje pozicije gde element treba biti
	for i := sl.Header.Level - 1; i >= 0; i-- {
		for current.Next[i] != nil && current.Next[i].Value < value {
			current = current.Next[i]
		}
	}

	// Ako je element na predvidjenoj poziciji
	if current.Next[0] != nil && current.Next[0].Value == value {
		return true
	}

	return false
}

// Brisanje elementa iz skip liste
func (sl *SkipList) Delete(value int) bool {

	// Pravimo niz u kome se za svaki nivo skip liste nalazi cvor nakon kog treba izbaciti trazeni cvor
	update := make([]*SkipNode, sl.maxHeight)
	current := sl.Header

	// Nalazenje pozicija na svakom nivou za brisanje cvora
	for i := sl.maxHeight - 1; i >= 0; i-- {
		for current.Next[i] != nil && current.Next[i].Value < value {
			current = current.Next[i]
		}
		update[i] = current
	}

	// Ako element postoji
	if current.Next[0] != nil && current.Next[0].Value == value {
		// Element za brisanje na najnizem nivou
		deletedNode := current.Next[0]

		// Brisanje elementa na svakom nivou
		//   i < sl.maxHeight - da petlja ne bi pokušala pristupiti nivoima iznad trenutnog najvišeg nivoa u listi
		//   i < len(update)  - da petlja ne bi pokušala pristupiti indeksima van granica niza update
		//   i < len(deletedNode.Next) - da petlja ne bi pokušala pristupiti indeksima van granica niza deletedNode.Next
		for i := 0; i < sl.maxHeight && i < len(update) && i < len(deletedNode.Next); i++ {
			if update[i].Next[i] != deletedNode {
				break
			}
			update[i].Next[i] = deletedNode.Next[i]
		}

		// Oslobadjanje memorije za izbrisan čvor
		for i := 0; i < len(deletedNode.Next); i++ {
			deletedNode.Next[i] = nil
		}

		// Smanjivanje nivoa glave ako su svi viši nivoi prazni
		for sl.Header.Level > 1 && sl.Header.Next[sl.Header.Level-1] == nil {
			sl.Header.Level--
		}
	} else {
		fmt.Println("Element", value, "ne postoji u skip listi.")
		return false
	}

	return true
}

// Funkcija za nasumican broj nivoa cvora
func (sl *SkipList) roll() int {
	Level := 0
	// Dok se dobija glava (1) baca se novcic -> dodaje se nivo
	for ; rand.Intn(2) == 1; Level++ {
		if Level >= sl.maxHeight {
			return Level
		}
	}
	return Level
}

// Ispis skip liste
func (sl *SkipList) Print() {
	for i := sl.maxHeight - 1; i >= 0; i-- {
		current := sl.Header
		for current.Next[i] != nil {
			fmt.Printf("%d ", current.Next[i].Value)
			current = current.Next[i]
		}
		fmt.Println()
	}
	fmt.Println()
}

func Test() {
	fmt.Println("\nTest Skip list:\n")
	skipList := NewSkipList(4)

	// Dodavanje
	skipList.Insert(3)
	skipList.Insert(6)
	skipList.Insert(7)
	skipList.Insert(9)
	skipList.Insert(12)
	skipList.Insert(19)
	skipList.Insert(17)
	skipList.Insert(26)
	skipList.Insert(21)
	skipList.Insert(25)

	fmt.Println("Skip lista:")
	skipList.Print()

	// Pretraga
	searchValue := 6
	if skipList.Search(searchValue) {
		fmt.Printf("Element %d je pronađen u Skip listi.\n", searchValue)
	} else {
		fmt.Printf("Element %d nije pronađen u Skip listi.\n", searchValue)
	}

	// Brisanje
	deleteValue := 17
	obrisan := skipList.Delete(deleteValue)
	if obrisan {
		fmt.Printf("Element %d je obrisan iz Skip liste.\n", deleteValue)

		fmt.Println("Skip lista nakon brisanja:")
		skipList.Print()
	}

}
