package magic

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// Magic handles special "magic" commands that provide fun and entertaining features
type Magic struct {
	// Add any configuration or state here if needed
}

// NewMagic creates a new Magic instance
func NewMagic() *Magic {
	return &Magic{}
}

// Execute processes a magic command and returns the result
func (m *Magic) Execute(command string) (string, error) {
	// Convert to lowercase for case-insensitive matching
	command = strings.ToLower(strings.TrimSpace(command))

	// Handle different magic commands
	switch command {
	case "dance":
		return m.Dance(), nil
	default:
		return fmt.Sprintf("Unknown magic command: %s\n\nAvailable magic commands:\n- dance: Shows a fun dance animation", command), nil
	}
}

// Dance displays a fun dance animation in the terminal
func (m *Magic) Dance() string {
	// Create a new random source with current time (Go 1.20+ approach)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Choose a random dance animation
	danceIndex := r.Intn(4)

	switch danceIndex {
	case 0:
		return m.stickFigureDance()
	case 1:
		return m.partyParrot()
	case 2:
		return m.discoTime()
	case 3:
		return m.robotDance()
	default:
		return m.stickFigureDance()
	}
}

// stickFigureDance returns a stick figure dance animation
func (m *Magic) stickFigureDance() string {
	dance := `
    🎵 Let's dance! 🎵

      o   \o/    o_/    o
     /|\   |    /|      |\_
     / \  / \   / \    / \

      \o\  |o|  /o/     \o/
       |   |    |        |
      / \  / \  / \     / \

     \o/   o/   \o      o/
      |   /|     |\    /|
     / \  / \   / \   / \

    💃 Dance like nobody's watching! 🕺
`
	return dance
}

// partyParrot returns a party parrot ASCII art
func (m *Magic) partyParrot() string {
	parrot := `
    🎉 Party Time! 🎉

        ／⌒ヽ
       / ^  ^ \
      ｜  ●   ｜
      \   3  /
       \    /
        U U

        ／⌒ヽ
       / >  < \
      ｜  ●   ｜
      \   o  /
       \    /
        U U

        ／⌒ヽ
       / ^  ^ \
      ｜  ●   ｜
      \   O  /
       \    /
        U U

    🦜 Party Parrot is dancing with you! 🦜
`
	return parrot
}

// discoTime returns a disco-themed animation
func (m *Magic) discoTime() string {
	disco := `
    🪩 Disco Time! 🪩

     ┏━━━┓
     ┃   ┃
     ┃ ● ┃
     ┃   ┃
     ┗━━━┛
      ╱ ╲

     * ✧ * ✧ * ✧ * ✧ * ✧ * ✧ *

     \\(^o^)/  \\(^o^)/  \\(^o^)/

     * ✧ * ✧ * ✧ * ✧ * ✧ * ✧ *

    🎵 Stayin' Alive, Stayin' Alive! 🎵
`
	return disco
}

// robotDance returns a robot dance animation
func (m *Magic) robotDance() string {
	robot := `
    🤖 Robot Dance Mode Activated! 🤖

      ╔═══╗
      ║ ■ ║
    ╔═╝   ╚═╗
    ║       ║
    ╚═╗   ╔═╝
      ║   ║
     ╔╝   ╚╗
     ║     ║
     ╚═╦═╦═╝
       ║ ║
       ╚═╝

    ⚡ Beep Boop... Executing dance.exe ⚡
`
	return robot
}
