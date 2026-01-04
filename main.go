package main

import (
	"fmt"
	"log"
	"unicode"

	"github.com/AllMeatball/u64-remote/keyboard"
	"github.com/AllMeatball/u64-remote/server"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct{
	server server.U64Server
	keybuffer []byte
}

const (
	KEY_RUN_STOP = 3
)

var priv_ebiten_to_chrcode = map[ebiten.Key]byte {
	ebiten.KeyUp:    keyboard.KEY_CRSR_UP,
	ebiten.KeyDown:  keyboard.KEY_CRSR_DOWN,
	ebiten.KeyLeft:  keyboard.KEY_CRSR_LEFT,
	ebiten.KeyRight: keyboard.KEY_CRSR_RIGHT,

	ebiten.KeyDelete:    keyboard.KEY_DEL_INS,
	ebiten.KeyBackspace: keyboard.KEY_DEL_INS,

	ebiten.KeyEnter:     keyboard.KEY_RETURN,
	ebiten.KeyPageUp:    KEY_RUN_STOP,
}

func (g *Game) HandleKeys() error {
	if len(g.keybuffer) <= 0 { return nil }

	is_buffer_full, err := keyboard.IsBufferFull(g.server)

	if err != nil { return nil }
	if is_buffer_full { return nil }

	// fmt.Println(g.keybuffer)
	if len(g.keybuffer) >= keyboard.KEYBOARD_MAX_COUNT {
		key_slice := g.keybuffer[:keyboard.KEYBOARD_MAX_COUNT]
		g.keybuffer = g.keybuffer[keyboard.KEYBOARD_MAX_COUNT:]

		err := keyboard.TypeBytes(g.server, key_slice)
		if err != nil { return err }
	} else {
		err := keyboard.TypeBytes(g.server, g.keybuffer)
		g.keybuffer = nil // clear keys
		if err != nil { return err }
	}


	// fmt.Println(key_slice)

	//
	return nil
}

func (g *Game) Update() error {
	var runes []rune
	runes = ebiten.AppendInputChars(runes)

	for _, r := range runes {
		g.keybuffer = append(g.keybuffer, keyboard.UnicodeToPet(unicode.ToUpper(r)))
	}

	for key, c64_key := range priv_ebiten_to_chrcode {
		if inpututil.IsKeyJustPressed(key) {
			g.keybuffer = append(g.keybuffer, c64_key)
		}
	}

	return g.HandleKeys()
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, fmt.Sprintf("%2.f", ebiten.ActualFPS()))
}

func (g *Game) Layout(width, height int) (screen_width, screen_height int) {
	return 640, 480
}

func main() {
	creds := server.U64Creds{
		EnableMessageBox: true,
	}

	err := server.LoadCreds(&creds)
	if err != nil { log.Fatal(err) }

	server, err := server.NewU64Server(creds)
	if err != nil { log.Fatal(err) }

	game := &Game{
		server: server,
	}

	keyboard.StopTyping(server)

	// ebiten.SetVsyncEnabled(false)
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("U64 NetKeyboard")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
