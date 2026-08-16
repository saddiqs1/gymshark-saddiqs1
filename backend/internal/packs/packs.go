package packs

import (
	"slices"
)

// TODO - revisit this eventually, logic is currently incorrect
func GetPacksForOrder(itemsOrdered int, packSizes []int) map[int]int {
	orderedPackSizes := orderPackSizes(packSizes)
	smallestPackSize := 250
	resultPacks := make(map[int]int)
	itemsRemaining := itemsOrdered

	for i, packSize := range orderedPackSizes {
		// TODO - make the logic for small pack sizes more generic, currently hardcoded for 250 and 500
		if packSize == smallestPackSize && itemsRemaining > packSize {
			// between 251 & 499, package 1 pack of 500
			resultPacks[orderedPackSizes[i-1]]++
		} else if itemsRemaining >= packSize {
			numPacks := itemsRemaining / packSize
			resultPacks[packSize] = numPacks
			itemsRemaining -= numPacks * packSize
		} else if packSize == smallestPackSize && itemsRemaining > 0 {
			// between 1 & 249
			resultPacks[smallestPackSize]++
		}
	}

	return resultPacks
}

func orderPackSizes(packSizes []int) []int {
	tempPackSizes := make([]int, len(packSizes))
	copy(tempPackSizes, packSizes)
	slices.Sort(tempPackSizes)
	slices.Reverse(tempPackSizes)
	return tempPackSizes
}
