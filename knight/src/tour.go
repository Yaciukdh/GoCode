package main

import (
	"fmt"
	"os"
	"strconv"
)

/* global */
var boardWidth int = 4
var boardHeight int = 4


type board struct {

}

func printBoard(board[][]uint8){
	for i := 0; i < boardHeight; i++ {
		fmt.Println(board[i])
		}
}

func makeBoard() [][]uint8 {
	board := make([][]uint8, boardHeight)
	for i := range board {
		board[i] = make([]uint8, boardWidth)
	}

	return board
}

func main() {
	startPos := os.Args[1:]
	startX,_ := strconv.Atoi(startPos[0])
	startY,_ := strconv.Atoi(startPos[1])
	tempBoard := makeBoard()
	printBoard(tempBoard)
	fmt.Println(startPos)
	fmt.Println(startX)
	fmt.Println(startY)



}
