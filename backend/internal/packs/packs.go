package packs

func Packs(itemsOrdered int) map[int]int {
	packSizes := []int{5000, 2000, 1000, 500, 250}
	smallestPackSize := 250
	resultPacks := make(map[int]int)
	itemsRemaining := itemsOrdered

	for i, packSize := range packSizes {
		if itemsRemaining >= packSize {
			// between 251 & 499, package 1 pack of 500
			if packSize == smallestPackSize && itemsRemaining > packSize {
				resultPacks[packSizes[i-1]]++
				itemsRemaining = 0
				continue
			}

			numPacks := itemsRemaining / packSize
			resultPacks[packSize] = numPacks
			itemsRemaining -= numPacks * packSize
		}

		// between 1 & 249
		if packSize == smallestPackSize && itemsRemaining > 0 {
			resultPacks[smallestPackSize]++
			itemsRemaining = 0
		}
	}

	return resultPacks
}
