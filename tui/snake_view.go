package tui

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	// Game constants
	snakeSymbol   = '█'
	foodSymbol    = '●'
	emptySymbol   = ' '
	borderSymbol  = '█'
	scoreTemplate = "Score: %d | High Score: %d | Press q to quit, p to pause"

	// Direction constants
	dirUp    = 0
	dirRight = 1
	dirDown  = 2
	dirLeft  = 3
)

// Coordinate represents a point on the game grid
type Coordinate struct {
	x, y int
}

// SnakeGame holds the state of the Snake game
type SnakeGame struct {
	snake          []Coordinate
	food           Coordinate
	direction      int
	nextDirection  int
	score          int
	highScore      int
	gameOver       bool
	paused         bool
	width          int
	height         int
	lastUpdateTime time.Time
	speed          time.Duration
}

// SnakeView component for the Snake game
type SnakeView struct {
	*tview.Box
	theme     *Theme
	game      *SnakeGame
	gameTimer *time.Timer
	stopChan  chan struct{}
}

// NewSnakeView creates a new Snake game view
func NewSnakeView(theme *Theme) *SnakeView {
	box := tview.NewBox().
		SetBorder(true).
		SetTitle(" SNAKE ").
		SetTitleAlign(tview.AlignCenter).
		SetBorderColor(theme.Colors.Primary).
		SetTitleColor(theme.Colors.Primary)

	game := &SnakeGame{
		snake:          []Coordinate{{10, 10}, {9, 10}, {8, 10}},
		direction:      dirRight,
		nextDirection:  dirRight,
		score:          0,
		highScore:      0,
		gameOver:       false,
		paused:         false,
		speed:          10 * time.Millisecond,
		lastUpdateTime: time.Now(),
		width:          20,
		height:         10,
	}

	sv := &SnakeView{
		Box:      box,
		theme:    theme,
		game:     game,
		stopChan: make(chan struct{}),
	}

	// Start the game loop
	sv.resetGame()
	sv.startGameLoop()

	return sv
}

// Draw renders the Snake game
func (sv *SnakeView) Draw(screen tcell.Screen) {
	sv.Box.Draw(screen)

	// Get the canvas dimensions
	x, y, width, height := sv.GetInnerRect()

	// Update game dimensions if needed
	if sv.game.width != width-2 || sv.game.height != height-2 {
		sv.game.width = width - 2 // -2 for inner padding
		sv.game.height = height - 2

		// Reset food position if outside bounds
		if sv.game.food.x >= sv.game.width || sv.game.food.y >= sv.game.height {
			sv.placeFood()
		}
	}

	// Draw borders
	borderStyle := tcell.StyleDefault.Foreground(sv.theme.Colors.Border)
	for i := 0; i < width; i++ {
		screen.SetContent(x+i, y, borderSymbol, nil, borderStyle)
		screen.SetContent(x+i, y+height-1, borderSymbol, nil, borderStyle)
	}
	for i := 0; i < height; i++ {
		screen.SetContent(x, y+i, borderSymbol, nil, borderStyle)
		screen.SetContent(x+width-1, y+i, borderSymbol, nil, borderStyle)
	}

	// Draw score
	scoreText := fmt.Sprintf(scoreTemplate, sv.game.score, sv.game.highScore)
	for i, c := range scoreText {
		if x+i+1 < x+width-1 {
			screen.SetContent(x+i+1, y+height-1, c, nil, tcell.StyleDefault.Foreground(sv.theme.Colors.Text))
		}
	}

	// Game over message
	if sv.game.gameOver {
		gameOverText := "GAME OVER - Press SPACE to restart"
		startX := x + (width-len(gameOverText))/2
		startY := y + height/2

		for i, c := range gameOverText {
			screen.SetContent(startX+i, startY, c, nil, tcell.StyleDefault.Foreground(sv.theme.Colors.Error))
		}
	}

	// Pause message
	if sv.game.paused && !sv.game.gameOver {
		pausedText := "PAUSED - Press P to continue"
		startX := x + (width-len(pausedText))/2
		startY := y + height/2

		for i, c := range pausedText {
			screen.SetContent(startX+i, startY, c, nil, tcell.StyleDefault.Foreground(sv.theme.Colors.Warning))
		}
	}

	// Draw food
	if !sv.game.gameOver {
		foodX, foodY := x+sv.game.food.x+1, y+sv.game.food.y+1
		if foodX > x && foodY > y && foodX < x+width-1 && foodY < y+height-1 {
			screen.SetContent(foodX, foodY, foodSymbol, nil, tcell.StyleDefault.Foreground(sv.theme.Colors.Error))
		}
	}

	// Draw snake
	snakeStyle := tcell.StyleDefault.Foreground(sv.theme.Colors.Primary)
	for _, segment := range sv.game.snake {
		segX, segY := x+segment.x+1, y+segment.y+1
		if segX > x && segY > y && segX < x+width-1 && segY < y+height-1 {
			screen.SetContent(segX, segY, snakeSymbol, nil, snakeStyle)
		}
	}
}

// InputHandler handles keyboard input for the Snake game
func (sv *SnakeView) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return sv.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		// Handle key input
		switch event.Key() {
		case tcell.KeyUp:
			if sv.game.direction != dirDown {
				sv.game.nextDirection = dirUp
			}
		case tcell.KeyRight:
			if sv.game.direction != dirLeft {
				sv.game.nextDirection = dirRight
			}
		case tcell.KeyDown:
			if sv.game.direction != dirUp {
				sv.game.nextDirection = dirDown
			}
		case tcell.KeyLeft:
			if sv.game.direction != dirRight {
				sv.game.nextDirection = dirLeft
			}
		case tcell.KeyRune:
			switch event.Rune() {
			case 'q', 'Q': // Quit the game view
				setFocus(nil) // This will make the dashboard switch away from the game
			case 'p', 'P': // Pause/resume the game
				sv.game.paused = !sv.game.paused
			case ' ': // Restart on space when game over
				if sv.game.gameOver {
					sv.resetGame()
				}
			}
		}
	})
}

// startGameLoop begins the game loop
func (sv *SnakeView) startGameLoop() {
	go func() {
		ticker := time.NewTicker(sv.game.speed)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if !sv.game.gameOver && !sv.game.paused {
					sv.updateGame()
				}
			case <-sv.stopChan:
				return
			}
		}
	}()
}

// updateGame advances the game state
func (sv *SnakeView) updateGame() {
	// Update direction
	sv.game.direction = sv.game.nextDirection

	// Get head position
	head := sv.game.snake[0]

	// Calculate new head position based on direction
	newHead := head
	switch sv.game.direction {
	case dirUp:
		newHead.y--
	case dirRight:
		newHead.x++
	case dirDown:
		newHead.y++
	case dirLeft:
		newHead.x--
	}

	// Check for collisions with walls
	if newHead.x < 0 || newHead.y < 0 || newHead.x >= sv.game.width || newHead.y >= sv.game.height {
		sv.game.gameOver = true
		return
	}

	// Check for collisions with self
	for _, segment := range sv.game.snake {
		if newHead.x == segment.x && newHead.y == segment.y {
			sv.game.gameOver = true
			return
		}
	}

	// Add new head
	sv.game.snake = append([]Coordinate{newHead}, sv.game.snake...)

	// Check if food was eaten
	if newHead.x == sv.game.food.x && newHead.y == sv.game.food.y {
		sv.game.score++
		if sv.game.score > sv.game.highScore {
			sv.game.highScore = sv.game.score
		}

		// Speed up the game slightly
		if sv.game.speed > 50*time.Millisecond {
			sv.game.speed -= 2 * time.Millisecond
		}

		// Place new food
		sv.placeFood()
	} else {
		// Remove tail if no food was eaten
		sv.game.snake = sv.game.snake[:len(sv.game.snake)-1]
	}
}

// placeFood positions food randomly on the grid
func (sv *SnakeView) placeFood() {
	for {
		// Generate random position
		x := rand.Intn(sv.game.width)
		y := rand.Intn(sv.game.height)

		// Check if position is occupied by snake
		occupied := false
		for _, segment := range sv.game.snake {
			if segment.x == x && segment.y == y {
				occupied = true
				break
			}
		}

		// If position is not occupied, place food there
		if !occupied {
			sv.game.food = Coordinate{x, y}
			break
		}
	}
}

// resetGame initializes a new game
func (sv *SnakeView) resetGame() {
	// Create initial snake (3 segments)
	sv.game.snake = []Coordinate{{10, 10}, {9, 10}, {8, 10}}
	sv.game.direction = dirRight
	sv.game.nextDirection = dirRight
	sv.game.score = 0
	sv.game.gameOver = false
	sv.game.paused = false
	sv.game.speed = 100 * time.Millisecond

	// Place initial food
	sv.placeFood()
}

// Stop stops the game loop
func (sv *SnakeView) Stop() {
	close(sv.stopChan)
}
