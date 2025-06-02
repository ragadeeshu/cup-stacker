package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"
)

type cup struct {
	ordinal int
	flipped bool
}
type stack []cup

func canStack(cup, cupToStack cup) bool {
	if !cup.flipped {
		if !cupToStack.flipped {
			return cupToStack.ordinal == cup.ordinal-1 ||
				cupToStack.ordinal == cup.ordinal+1 ||
				cupToStack.ordinal == cup.ordinal-4
		} else {
			return cupToStack.ordinal == cup.ordinal-3 ||
				cupToStack.ordinal == cup.ordinal+3
		}
	} else {
		if !cupToStack.flipped {
			return cupToStack.ordinal == cup.ordinal-1 ||
				cupToStack.ordinal == cup.ordinal+1
		} else {
			return cupToStack.ordinal == cup.ordinal+1 ||
				cupToStack.ordinal == cup.ordinal-1 ||
				cupToStack.ordinal == cup.ordinal+4
		}
	}
}

func main() {
	port := flag.String("port", "8080", "Port to bind to")
	flag.Parse()

	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.Error(w, "404 not found.", http.StatusNotFound)
			return
		}
		if r.Method != "GET" {
			http.Error(w, "Method is not supported.", http.StatusMethodNotAllowed)
			return
		}
		if len(r.URL.Query()) == 0 {
			fmt.Fprint(w, "Please provide stack")
			return
		} else if stack := r.URL.Query().Get("stack"); len(r.URL.Query()) == 1 && stack != "" {
			cupStack, err := parseStack(stack)
			if err != nil {
				http.Error(w, "Unable to parse stack", http.StatusBadRequest)
				return
			}
			if isValidStack(cupStack) {
				renderStack(cupStack, w)
				return
			} else {
				http.Error(w, "Not a valid stack", http.StatusBadRequest)
				return
			}
		} else {
			http.Error(w, "Only stack is supported.", http.StatusBadRequest)
			return
		}
	}

	http.HandleFunc("/", handler)
	fmt.Println("ALIVE")
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

func renderStack(stack stack, w http.ResponseWriter) {
	var mat [][]rune
	highestOrdinal := -1
	for _, cup := range stack {
		if cup.ordinal > highestOrdinal {
			highestOrdinal = cup.ordinal
		}
	}
	cols := highestOrdinal*2 + 13
	rowOffset := 0
	for stackIndex, cup := range slices.Backward(stack) {
		rows := cup.ordinal + 3
		cupOffset := 0
		if stackIndex != len(stack)-1 {
			cupOffset = calculateOffset(stack[stackIndex], stack[stackIndex+1])
		}
		rows += cupOffset
		for range rows {
			mat = append(mat, make([]rune, cols))
		}
		colOffset := highestOrdinal - cup.ordinal
		renderCup(mat, cup, rowOffset+cupOffset, colOffset)
		rowOffset += rows
	}
	for i := range mat {
		fmt.Fprintln(w, strings.TrimRight(string(mat[i]), "\x00"))
	}
}

func calculateOffset(cup, previousCup cup) int {
	if !previousCup.flipped && cup.flipped {
		return 0
	}
	if previousCup.flipped && !cup.flipped {
		return -1
	}
	if diff := cup.ordinal - previousCup.ordinal; diff == 4 || diff == -4 {
		return 0
	}
	if !cup.flipped && cup.ordinal < previousCup.ordinal {
		return -(cup.ordinal + 2)
	}
	if cup.flipped && cup.ordinal > previousCup.ordinal {
		return -(cup.ordinal + 1)
	}
	return -1
}

func renderCup(mat [][]rune, cup cup, rowOffset, colOffset int) {
	for r := range 3 + cup.ordinal {
		for c := range colOffset {
			renderRune(mat, rowOffset+r, c, ' ')
		}
	}
	topMidWidth := (cup.ordinal-1)*2 + 1
	if !cup.flipped {
		renderRune(mat, rowOffset+0, colOffset+0, ' ')
		renderRune(mat, rowOffset+0, colOffset+1, ' ')
		renderRune(mat, rowOffset+0, colOffset+2, ' ')
		renderRune(mat, rowOffset+0, colOffset+3, ' ')
		renderRune(mat, rowOffset+0, colOffset+4, '╭')
		renderRune(mat, rowOffset+0, colOffset+5, '─')
		renderRune(mat, rowOffset+0, colOffset+6, '╮')
		for i := range topMidWidth {
			renderRune(mat, rowOffset+0, colOffset+7+i, ' ')
		}
		renderRune(mat, rowOffset+0, colOffset+7+topMidWidth, '╭')
		renderRune(mat, rowOffset+0, colOffset+8+topMidWidth, '─')
		renderRune(mat, rowOffset+0, colOffset+9+topMidWidth, '╮')
		renderRune(mat, rowOffset+1, colOffset+0, ' ')
		renderRune(mat, rowOffset+1, colOffset+1, ' ')
		renderRune(mat, rowOffset+1, colOffset+2, '┌')
		renderRune(mat, rowOffset+1, colOffset+3, '─')
		renderRune(mat, rowOffset+1, colOffset+4, '┘')
		renderRune(mat, rowOffset+1, colOffset+5, ' ')
		renderRune(mat, rowOffset+1, colOffset+6, '└')
		for i := range topMidWidth {
			renderRune(mat, rowOffset+1, colOffset+7+i, '─')
		}
		renderRune(mat, rowOffset+1, colOffset+7+topMidWidth, '┘')
		renderRune(mat, rowOffset+1, colOffset+8+topMidWidth, ' ')
		renderRune(mat, rowOffset+1, colOffset+9+topMidWidth, '└')
		renderRune(mat, rowOffset+1, colOffset+10+topMidWidth, '─')
		renderRune(mat, rowOffset+1, colOffset+11+topMidWidth, '┐')

		for i := range cup.ordinal {
			renderRune(mat, rowOffset+i+2, colOffset+0, ' ')
			renderRune(mat, rowOffset+i+2, colOffset+1, ' ')
			renderRune(mat, rowOffset+i+2, colOffset+2, '│')
			for j := range cup.ordinal*2 + 7 {
				renderRune(mat, rowOffset+i+2, colOffset+3+j, ' ')
			}
			renderRune(mat, rowOffset+i+2, colOffset+10+cup.ordinal*2, '│')
		}
		renderRune(mat, rowOffset+2+cup.ordinal, colOffset+0, '┌')
		renderRune(mat, rowOffset+2+cup.ordinal, colOffset+1, '─')
		renderRune(mat, rowOffset+2+cup.ordinal, colOffset+2, '┤')
		for j := range cup.ordinal*2 + 7 {
			renderRune(mat, rowOffset+2+cup.ordinal, colOffset+3+j, ' ')
		}
		renderRune(mat, rowOffset+2+cup.ordinal, colOffset+10+cup.ordinal*2, '├')
		renderRune(mat, rowOffset+2+cup.ordinal, colOffset+11+cup.ordinal*2, '─')
		renderRune(mat, rowOffset+2+cup.ordinal, colOffset+12+cup.ordinal*2, '┐')
	} else {
		renderRune(mat, rowOffset+0, colOffset+0, '└')
		renderRune(mat, rowOffset+0, colOffset+1, '─')
		renderRune(mat, rowOffset+0, colOffset+2, '┤')
		for j := range cup.ordinal*2 + 7 {
			renderRune(mat, rowOffset+0, colOffset+3+j, ' ')
		}
		renderRune(mat, rowOffset+0, colOffset+10+cup.ordinal*2, '├')
		renderRune(mat, rowOffset+0, colOffset+11+cup.ordinal*2, '─')
		renderRune(mat, rowOffset+0, colOffset+12+cup.ordinal*2, '┘')
		for i := range cup.ordinal {
			renderRune(mat, rowOffset+i+1, colOffset+0, ' ')
			renderRune(mat, rowOffset+i+1, colOffset+1, ' ')
			renderRune(mat, rowOffset+i+1, colOffset+2, '│')
			for j := range cup.ordinal*2 + 7 {
				renderRune(mat, rowOffset+i+1, colOffset+3+j, ' ')
			}
			renderRune(mat, rowOffset+i+1, colOffset+10+cup.ordinal*2, '│')
		}
		renderRune(mat, rowOffset+1+cup.ordinal, colOffset+0, ' ')
		renderRune(mat, rowOffset+1+cup.ordinal, colOffset+1, ' ')
		renderRune(mat, rowOffset+1+cup.ordinal, colOffset+2, '└')
		renderRune(mat, rowOffset+1+cup.ordinal, colOffset+3, '─')
		renderRune(mat, rowOffset+1+cup.ordinal, colOffset+4, '┐')
		renderRune(mat, rowOffset+1+cup.ordinal, colOffset+5, ' ')
		renderRune(mat, rowOffset+1+cup.ordinal, colOffset+6, '┌')
		for i := range topMidWidth {
			renderRune(mat, rowOffset+1+cup.ordinal, colOffset+7+i, '─')
		}
		renderRune(mat, rowOffset+1+cup.ordinal, colOffset+7+topMidWidth, '┐')
		renderRune(mat, rowOffset+1+cup.ordinal, colOffset+8+topMidWidth, ' ')
		renderRune(mat, rowOffset+1+cup.ordinal, colOffset+9+topMidWidth, '┌')
		renderRune(mat, rowOffset+1+cup.ordinal, colOffset+10+topMidWidth, '─')
		renderRune(mat, rowOffset+1+cup.ordinal, colOffset+11+topMidWidth, '┘')

		renderRune(mat, rowOffset+2+cup.ordinal, colOffset+0, ' ')
		renderRune(mat, rowOffset+2+cup.ordinal, colOffset+1, ' ')
		renderRune(mat, rowOffset+2+cup.ordinal, colOffset+2, ' ')
		renderRune(mat, rowOffset+2+cup.ordinal, colOffset+3, ' ')
		renderRune(mat, rowOffset+2+cup.ordinal, colOffset+4, '╰')
		renderRune(mat, rowOffset+2+cup.ordinal, colOffset+5, '─')
		renderRune(mat, rowOffset+2+cup.ordinal, colOffset+6, '╯')
		for i := range topMidWidth {
			renderRune(mat, rowOffset+2+cup.ordinal, colOffset+7+i, ' ')
		}
		renderRune(mat, rowOffset+2+cup.ordinal, colOffset+7+topMidWidth, '╰')
		renderRune(mat, rowOffset+2+cup.ordinal, colOffset+8+topMidWidth, '─')
		renderRune(mat, rowOffset+2+cup.ordinal, colOffset+9+topMidWidth, '╯')
	}

}

func renderRune(mat [][]rune, row, col int, r rune) {
	if mat[row][col] == 0 || mat[row][col] == ' ' {
		mat[row][col] = r
	}
}

func isValidStack(cupStack stack) bool {
	if len(cupStack) < 1 {
		return false
	}
	if len(cupStack) == 1 {
		return true
	}
	for i := 1; i < len(cupStack); i++ {
		if !canStack(cupStack[i-1], cupStack[i]) {
			return false
		}
	}
	return true
}

func parseStack(stack string) (parsedCups stack, err error) {
	for _, c := range strings.Split(stack, ",") {
		flipped := false
		if strings.HasSuffix(c, "f") {
			flipped = true
			c = strings.TrimSuffix(c, "f")
		}
		ordinal, err := strconv.Atoi(c)
		if err != nil {
			return nil, err
		}
		parsedCups = append(parsedCups, cup{ordinal: ordinal, flipped: flipped})
	}
	return
}
