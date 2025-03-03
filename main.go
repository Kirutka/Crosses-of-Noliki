package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 600 // Ширина
	screenHeight = 600 // Высота
)

var (
	boardSize     = 3     // Размер поля
	board         [][]int // Поле
	currentPlayer = 1     // Текущий игрок
	gameMode      = 0     // Режим игры (0 - выбор режима, 1 - два игрока, 2 - против ИИ)
	aiDifficulty  = 0     // Сложность ИИ (0 - выбор сложности, 1 - простой, 2 - сложный)
	gameOver      = false // Флаг завершения игры
	winner        = 0     // Победитель (0 - ничья, 1 или 2 - победитель)
)

func main() {
	rand.Seed(time.Now().UnixNano()) // Инициализация генератора случайных чисел

	board = make([][]int, boardSize) // Инициализация игрового поля
	for i := range board {
		board[i] = make([]int, boardSize) // Создание строк в поле
	}

	ebiten.SetWindowSize(screenWidth, screenHeight) // Установка размера окна
	ebiten.SetWindowTitle("Crosses-noliki")         // Установка заголовка окна

	if err := ebiten.RunGame(&Game{}); err != nil { // Запуск игры
		log.Fatal(err)
	}
}

type Game struct{}

func (g *Game) Update() error {
	if gameMode == 0 { // Выбор режима игры
		if inpututil.IsKeyJustPressed(ebiten.Key1) {
			gameMode = 1 // Режим "два игрока"
		}
		if inpututil.IsKeyJustPressed(ebiten.Key2) {
			gameMode = 2 // Режим "против ИИ"
		}
		return nil
	}

	if gameMode == 2 && aiDifficulty == 0 { // Выбор сложности ИИ
		if inpututil.IsKeyJustPressed(ebiten.Key1) {
			aiDifficulty = 1 // Простой ИИ
		}
		if inpututil.IsKeyJustPressed(ebiten.Key2) {
			aiDifficulty = 2 // Сложный ИИ
		}
		return nil
	}

	if gameOver { // Перезапуск игры
		if inpututil.IsKeyJustPressed(ebiten.KeyR) {
			resetGame() // Сброс игры
		}
		return nil
	}

	if gameMode == 1 || (gameMode == 2 && currentPlayer == 1) { // Ход игрока
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()         // Получение позиции курсора
			cellX := x / (screenWidth / boardSize)  // Вычисление ячейки по X
			cellY := y / (screenHeight / boardSize) // Вычисление ячейки по Y
			if cellX >= 0 && cellX < boardSize && cellY >= 0 && cellY < boardSize && board[cellY][cellX] == 0 {
				board[cellY][cellX] = currentPlayer // Запись хода в поле
				checkGameOver()                     // Проверка завершения игры
				currentPlayer = 3 - currentPlayer   // Смена игрока
			}
		}
	} else if gameMode == 2 && currentPlayer == 2 { // Ход ИИ
		aiMove()                          // Ход ИИ
		checkGameOver()                   // Проверка завершения игры
		currentPlayer = 3 - currentPlayer // Смена игрока
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	drawGradientBackground(screen) // Градиентный фон

	// Отрисовка сетки
	cellWidth := screenWidth / boardSize   // Ширина ячейки
	cellHeight := screenHeight / boardSize // Высота ячейки
	for i := 1; i < boardSize; i++ {
		// Отрисовка вертикальных линий
		ebitenutil.DrawLine(screen, float64(i*cellWidth), 0, float64(i*cellWidth), float64(screenHeight), color.Black)
		// Отрисовка горизонтальных линий
		ebitenutil.DrawLine(screen, 0, float64(i*cellHeight), float64(screenWidth), float64(i*cellHeight), color.Black)
	}

	// Отрисовка крестиков и ноликов
	for y := 0; y < boardSize; y++ {
		for x := 0; x < boardSize; x++ {
			if board[y][x] == 1 { // Отрисовка крестика
				drawCross(screen, x*cellWidth, y*cellHeight, cellWidth, cellHeight)
			} else if board[y][x] == 2 { // Отрисовка нолика
				drawCircle(screen, x*cellWidth, y*cellHeight, cellWidth, cellHeight)
			}
		}
	}

	if gameMode == 0 {
		ebitenutil.DebugPrintAt(screen, "Select the game mode: \n1 - two players \n2 - against AI", 10, 10)
	} else if gameMode == 2 && aiDifficulty == 0 {
		ebitenutil.DebugPrintAt(screen, "Select the complexity of AI: \n1 - simple \n2 - complex", 10, 10)
	} else if gameOver {
		if winner == 0 {
			ebitenutil.DebugPrintAt(screen, "Draw! Click R for restart.", 10, 10)
		} else {
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Player %d won! Click R for restart.", winner), 10, 10)
		}
	} else {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Player's course %d ", currentPlayer), 10, 10)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) { // Layout задаёт размеры окна
	return screenWidth, screenHeight
}

func drawGradientBackground(screen *ebiten.Image) { // Градиентный фон
	for y := 0; y < screenHeight; y++ {
		alpha := uint8(255 * y / screenHeight)                                                    // Вычисление прозрачности
		col := color.RGBA{100, 100, 200, alpha}                                                   // Цвет градиента
		vector.StrokeLine(screen, 0, float32(y), float32(screenWidth), float32(y), 1, col, false) // Отрисовка линии
	}
}

func drawCross(screen *ebiten.Image, x, y, width, height int) { // Отрисовка крестика
	padding := 20 // Отступ от краёв ячейки
	x1 := float32(x + padding)
	y1 := float32(y + padding)
	x2 := float32(x + width - padding)
	y2 := float32(y + height - padding)

	// Отрисовка двух линий крестика
	vector.StrokeLine(screen, x1, y1, x2, y2, 10, color.RGBA{255, 0, 0, 255}, true)
	vector.StrokeLine(screen, x2, y1, x1, y2, 10, color.RGBA{255, 0, 0, 255}, true)
}

func drawCircle(screen *ebiten.Image, x, y, width, height int) { // Отрисовка нолика
	cx := float32(x + width/2)      // Центр по X
	cy := float32(y + height/2)     // Центр по Y
	radius := float32(width/2 - 20) // Радиус нолика

	// Отрисовка круга
	vector.StrokeCircle(screen, cx, cy, radius, 10, color.RGBA{0, 0, 255, 255}, true)
}

func checkGameOver() { // Проверка завершения игры
	// Проверка строк и столбцов
	for i := 0; i < boardSize; i++ {
		if board[i][0] != 0 && allEqual(board[i]) {
			winner = board[i][0]
			gameOver = true
			return
		}
		if board[0][i] != 0 && allEqual(getColumn(i)) {
			winner = board[0][i]
			gameOver = true
			return
		}
	}

	// Проверка диагоналей
	if board[0][0] != 0 && allEqual(getDiagonal(true)) {
		winner = board[0][0]
		gameOver = true
		return
	}
	if board[0][boardSize-1] != 0 && allEqual(getDiagonal(false)) {
		winner = board[0][boardSize-1]
		gameOver = true
		return
	}

	// Проверка на ничью
	if isBoardFull(board) {
		gameOver = true
		winner = 0
	}
}

func aiMove() { // Ход ИИ
	if aiDifficulty == 1 {
		// Простой ИИ: случайный ход
		for {
			x := rand.Intn(boardSize)
			y := rand.Intn(boardSize)
			if board[y][x] == 0 {
				board[y][x] = 2
				break
			}
		}
	} else if aiDifficulty == 2 {
		// Сложный ИИ: минимизация проигрыша
		bestScore := -1
		bestMove := [2]int{-1, -1}
		for y := 0; y < boardSize; y++ {
			for x := 0; x < boardSize; x++ {
				if board[y][x] == 0 {
					board[y][x] = 2
					score := minimax(board, 0, false)
					board[y][x] = 0
					if score > bestScore {
						bestScore = score
						bestMove = [2]int{y, x}
					}
				}
			}
		}
		if bestMove[0] != -1 {
			board[bestMove[0]][bestMove[1]] = 2
		}
	}
}

// Минимакс для сложного ИИ
func minimax(board [][]int, depth int, isMaximizing bool) int {
	if checkWinner(board) == 2 {
		return 1
	}
	if checkWinner(board) == 1 {
		return -1
	}
	if isBoardFull(board) {
		return 0
	}

	if isMaximizing {
		bestScore := -1
		for y := 0; y < boardSize; y++ {
			for x := 0; x < boardSize; x++ {
				if board[y][x] == 0 {
					board[y][x] = 2
					score := minimax(board, depth+1, false)
					board[y][x] = 0
					if score > bestScore {
						bestScore = score
					}
				}
			}
		}
		return bestScore
	} else {
		bestScore := 1
		for y := 0; y < boardSize; y++ {
			for x := 0; x < boardSize; x++ {
				if board[y][x] == 0 {
					board[y][x] = 1
					score := minimax(board, depth+1, true)
					board[y][x] = 0
					if score < bestScore {
						bestScore = score
					}
				}
			}
		}
		return bestScore
	}
}

func checkWinner(board [][]int) int { // Проверка победителя
	// Проверка строк и столбцов
	for i := 0; i < boardSize; i++ {
		if board[i][0] != 0 && allEqual(board[i]) {
			return board[i][0]
		}
		if board[0][i] != 0 && allEqual(getColumnFromBoard(board, i)) {
			return board[0][i]
		}
	}

	// Проверка диагоналей
	if board[0][0] != 0 && allEqual(getDiagonalFromBoard(board, true)) {
		return board[0][0]
	}
	if board[0][boardSize-1] != 0 && allEqual(getDiagonalFromBoard(board, false)) {
		return board[0][boardSize-1]
	}

	return 0
}

func allEqual(slice []int) bool { // Проверяет, все ли элементы в слайсе (срезе) равны первому элементу
	for i := 1; i < len(slice); i++ {
		if slice[i] != slice[0] {
			return false
		}
	}
	return true
}

func getColumn(col int) []int { // Возвращает столбец игрового поля по указанному индексу col
	column := make([]int, boardSize)
	for i := 0; i < boardSize; i++ {
		column[i] = board[i][col]
	}
	return column
}

func getColumnFromBoard(board [][]int, col int) []int {
	column := make([]int, boardSize)
	for i := 0; i < boardSize; i++ {
		column[i] = board[i][col]
	}
	return column
}

func getDiagonal(main bool) []int { // Возвращает диагональ игрового поля
	diagonal := make([]int, boardSize)
	for i := 0; i < boardSize; i++ {
		if main {
			diagonal[i] = board[i][i]
		} else {
			diagonal[i] = board[i][boardSize-1-i]
		}
	}
	return diagonal
}

func getDiagonalFromBoard(board [][]int, main bool) []int {
	diagonal := make([]int, boardSize)
	for i := 0; i < boardSize; i++ {
		if main {
			diagonal[i] = board[i][i]
		} else {
			diagonal[i] = board[i][boardSize-1-i]
		}
	}
	return diagonal
}

func isBoardFull(board [][]int) bool { // Проверяет, заполнено ли игровое поле
	for y := 0; y < boardSize; y++ {
		for x := 0; x < boardSize; x++ {
			if board[y][x] == 0 {
				return false
			}
		}
	}
	return true
}

func resetGame() { // Сбрасывает игру
	board = make([][]int, boardSize)
	for i := range board {
		board[i] = make([]int, boardSize)
	}
	currentPlayer = 1
	gameOver = false
	winner = 0
}
