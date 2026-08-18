package packs

import (
	"errors"
	"math"
	"slices"

	"github.com/rs/zerolog"
)

type PackResult struct {
	PackSize int `json:"packSize"`
	Count    int `json:"count"`
}

func GetPacksForOrder(logger *zerolog.Logger, itemsOrdered int, packSizes []int) ([]PackResult, error) {
	orderedPackSizes := orderPackSizes(packSizes)
	smallestPack := orderedPackSizes[len(orderedPackSizes)-1]

	// worst possible solution is using smallest pack size only
	maxResult := PackResult{PackSize: smallestPack, Count: int(math.Ceil(float64(itemsOrdered) / float64(smallestPack)))}
	maxResultTotalItems := maxResult.Count * maxResult.PackSize

	packResultCombinations := map[int][]PackResult{}

	for numberOfItems := 1; numberOfItems <= maxResultTotalItems; numberOfItems++ {
		for _, packSize := range orderedPackSizes {
			// NOTE: break once best combination has been found, no need to calculate for further packsizes

			// if packsize matches items
			if numberOfItems == packSize {
				packResult := PackResult{PackSize: packSize, Count: 1}
				packResults := []PackResult{packResult}
				packResultCombinations[numberOfItems] = packResults
				break
			}

			// does (current packsize - item amount), exist? if so, we can append
			if packResultCombinations[numberOfItems-packSize] != nil {
				packResults := make([]PackResult, len(packResultCombinations[numberOfItems-packSize]))
				copy(packResults, packResultCombinations[numberOfItems-packSize])
				isExistingPack := false

				for i := range packResults {
					if packResults[i].PackSize == packSize {
						packResults[i].Count++
						isExistingPack = true
						break
					}
				}

				if !isExistingPack {
					packResults = append(packResults, PackResult{PackSize: packSize, Count: 1})
				}

				packResultCombinations[numberOfItems] = packResults
				break
			}
		}
	}

	for i := itemsOrdered; i <= maxResultTotalItems; i++ {
		if packResultCombinations[i] != nil {
			return packResultCombinations[i], nil
		}
	}

	return []PackResult{}, errors.New("no valid pack size combinations found for requested order")
}

func orderPackSizes(packSizes []int) []int {
	tempPackSizes := make([]int, len(packSizes))
	copy(tempPackSizes, packSizes)
	slices.Sort(tempPackSizes)
	slices.Reverse(tempPackSizes)
	return tempPackSizes
}
