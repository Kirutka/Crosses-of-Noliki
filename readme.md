# 🎮 Крестики-нолики на Go [![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://golang.org/)

<p align="center">
  <img src="https://github.com/user-attachments/assets/6a5b9aa1-c1d8-4c26-8d88-3514e7b319dd" alt="Game Preview" width="600">
</p>


**Классическая игра с ИИ и продвинутыми алгоритмами на Golang**

---

## 🌟 Особенности
- 🕹️ **Два режима игры**: против друга или ИИ
- 🧠 **Умный ИИ** с алгоритмом минимакс (сложный уровень)
- 🎨 **Графический интерфейс** на библиотеке Ebiten
- ⚙️ **Настраиваемый размер поля** (3x3 по умолчанию)
- 🌈 **Градиентный фон** и плавная анимация
- 🔄 **Система рестарта** игры

---

## 🚀 Быстрый старт

### Требования
- Go 1.21+
- [Ebiten](https://ebiten.org/) (`go get -u github.com/hajimehoshi/ebiten/v2`)

```bash
# Клонировать репозиторий
git clone lagman https://github.com/Kirutka/Lagman.git
cd Крестики-нолики

# Запустить игру
go run main.go
```

---

## 🕹️ Управление
- **1** - режим "два игрока"
- **2** - игра против ИИ
- **ЛКМ** - сделать ход
- **R** - перезапустить игру

**Сложность ИИ** (после выбора режима 2):
- **1** - случайные ходы
- **2** - продвинутый алгоритм

---

## 🛠 Технологии
[![Ebiten](https://img.shields.io/badge/-Ebiten-2C2D2F?logo=go&logoColor=white)](https://ebiten.org/)
[![Go](https://img.shields.io/badge/-Go-00ADD8?logo=go&logoColor=white)](https://golang.org/)


---

## 📌 Особенности реализации
```go
// Пример алгоритма минимакс
func minimax(board [][]int, depth int, isMaximizing bool) int {
    if checkWinner(board) == 2 { return 1 }
    if checkWinner(board) == 1 { return -1 }
    if isBoardFull(board) { return 0 }
    
    // Рекурсивный поиск оптимального хода
    if isMaximizing {
        bestScore := -1
        for y := 0; y < boardSize; y++ {
            for x := 0; x < boardSize; x++ {
                if board[y][x] == 0 {
                    board[y][x] = 2
                    score := minimax(board, depth+1, false)
                    board[y][x] = 0
                    bestScore = max(bestScore, score)
                }
            }
        }
        return bestScore
    }
    // ... полный код в репозитории
}
```

---

## 📮 Контакты
**Кирилл** - начинающий Go-разработчик  
[![Telegram](https://img.shields.io/badge/-@Lesnoy_umorust-2CA5E0?style=flat&logo=telegram)](https://t.me/Lesnoy_umorust)
[![VK](https://img.shields.io/badge/-VK-0077FF?style=flat&logo=vk)](https://vk.com/id549536760)
[![Email](https://img.shields.io/badge/-kirillzaporozec48@gmail.com-D14836?style=flat&logo=gmail)](mailto:kirillzaporozec48@gmail.com)


