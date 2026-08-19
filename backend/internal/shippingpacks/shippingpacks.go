package shippingpacks

import (
	"errors"
	"math"
	"slices"

	"github.com/rs/zerolog"
)

type ShippingPack struct {
	PackSize int `json:"packSize"`
	Count    int `json:"count"`
}

func GetShippingPacksForOrder(logger *zerolog.Logger, itemsOrdered int, packSizes []int) ([]ShippingPack, error) {
	orderedPackSizes := orderPackSizes(packSizes)
	smallestPack := orderedPackSizes[len(orderedPackSizes)-1]

	maxShippingPacks := ShippingPack{PackSize: smallestPack, Count: int(math.Ceil(float64(itemsOrdered) / float64(smallestPack)))}
	maxShippingPacksTotalItems := maxShippingPacks.Count * maxShippingPacks.PackSize

	shippingPacksForItemAmount := map[int][]ShippingPack{}

	for numberOfItems := 1; numberOfItems <= maxShippingPacksTotalItems; numberOfItems++ {
		shippingPackCombinations := [][]ShippingPack{}

		for _, packSize := range orderedPackSizes {
			if numberOfItems == packSize {
				shippingPackCombinations = append(shippingPackCombinations, []ShippingPack{{PackSize: packSize, Count: 1}})
				continue
			}

			prevShippingPacks := shippingPacksForItemAmount[numberOfItems-packSize]
			if prevShippingPacks != nil {
				shippingPacks := make([]ShippingPack, len(prevShippingPacks))
				copy(shippingPacks, prevShippingPacks)
				isPresent := false

				for i := range shippingPacks {
					if shippingPacks[i].PackSize == packSize {
						shippingPacks[i].Count++
						isPresent = true
						break
					}
				}

				if !isPresent {
					shippingPacks = append(shippingPacks, ShippingPack{PackSize: packSize, Count: 1})
				}

				shippingPackCombinations = append(shippingPackCombinations, shippingPacks)
			}
		}

		if len(shippingPackCombinations) > 0 {
			bestCombination := shippingPackCombinations[0]
			for _, spc := range shippingPackCombinations {
				if getPackCount(spc) < getPackCount(bestCombination) {
					bestCombination = spc
				}
			}

			shippingPacksForItemAmount[numberOfItems] = bestCombination
		}
	}

	for i := itemsOrdered; i <= maxShippingPacksTotalItems; i++ {
		if shippingPacksForItemAmount[i] != nil {
			return shippingPacksForItemAmount[i], nil
		}
	}

	return []ShippingPack{}, errors.New("no valid pack size combinations found for requested order")
}

func orderPackSizes(packSizes []int) []int {
	tempPackSizes := make([]int, len(packSizes))
	copy(tempPackSizes, packSizes)
	slices.Sort(tempPackSizes)
	slices.Reverse(tempPackSizes)
	return tempPackSizes
}

func getPackCount(packs []ShippingPack) int {
	count := 0
	for _, p := range packs {
		count += p.Count
	}
	return count
}
