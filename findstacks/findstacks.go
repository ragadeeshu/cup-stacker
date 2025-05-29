package main

import (
	"fmt"
	"slices"
)

type cup struct {
	ordinal int
	flipped bool
}
type stack []cup

func findStacks(numCups int) []stack {
	var cupsToStack []int
	for i := 1; i <= numCups; i++ {
		cupsToStack = append(cupsToStack, i)
	}
	return stackRemainingCups([]cup{}, cupsToStack)
}

func stackRemainingCups(stackedCups []cup, cupsToStack []int) (foundStacks []stack) {
	for _, cupOrdinalToStack := range cupsToStack {
		for _, flipped := range []bool{false, true} {
			if len(stackedCups) == 0 || canStack(stackedCups[len(stackedCups)-1], cupOrdinalToStack, flipped) {
				freshlyStackedCups := append(stackedCups, cup{cupOrdinalToStack, flipped})
				if len(cupsToStack) == 1 {
					foundStacks = append(foundStacks, freshlyStackedCups)
				} else {
					remainingCups := slices.Clone(cupsToStack)
					remainingCups = slices.DeleteFunc(remainingCups, func(cup int) bool { return cup == cupOrdinalToStack })
					foundStacks = append(foundStacks, stackRemainingCups(freshlyStackedCups, remainingCups)...)
				}
			}
		}
	}
	return
}

func canStack(cup cup, cupOrdinalToStack int, flipped bool) bool {
	if !cup.flipped {
		if !flipped {
			return cupOrdinalToStack == cup.ordinal-1 ||
				cupOrdinalToStack == cup.ordinal+1 ||
				cupOrdinalToStack == cup.ordinal-4
		} else {
			return cupOrdinalToStack == cup.ordinal-3 ||
				cupOrdinalToStack == cup.ordinal+3
		}
	} else {
		if !flipped {
			return cupOrdinalToStack == cup.ordinal-1 ||
				cupOrdinalToStack == cup.ordinal+1
		} else {
			return cupOrdinalToStack == cup.ordinal+1 ||
				cupOrdinalToStack == cup.ordinal-1 ||
				cupOrdinalToStack == cup.ordinal+4
		}
	}
}

func main() {

	stacks := findStacks(8)
	fmt.Println(len(stacks))
	fmt.Println(stacks)

}
