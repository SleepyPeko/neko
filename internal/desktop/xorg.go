package desktop

import (
	"image"
	"regexp"
	"os/exec"

	"demodesk/neko/internal/types"
	"demodesk/neko/internal/desktop/xorg"
)

func (manager *DesktopManagerCtx) Move(x, y int) {
	xorg.Move(x, y)
}

func (manager *DesktopManagerCtx) Scroll(x, y int) {
	xorg.Scroll(x, y)
}

func (manager *DesktopManagerCtx) ButtonDown(code int) error {
	return xorg.ButtonDown(code)
}

func (manager *DesktopManagerCtx) KeyDown(code uint64) error {
	return xorg.KeyDown(code)
}

func (manager *DesktopManagerCtx) ButtonUp(code int) error {
	return xorg.ButtonUp(code)
}

func (manager *DesktopManagerCtx) KeyUp(code uint64) error {
	return xorg.KeyUp(code)
}

func (manager *DesktopManagerCtx) ResetKeys() {
	xorg.ResetKeys()
}

func (manager *DesktopManagerCtx) ScreenConfigurations() map[int]types.ScreenConfiguration {
	return xorg.ScreenConfigurations
}

func (manager *DesktopManagerCtx) SetScreenSize(size types.ScreenSize) error {
	mu.Lock()
	manager.emmiter.Emit("before_screen_size_change")

	defer func() {
		manager.emmiter.Emit("after_screen_size_change")
		mu.Unlock()
	}()

	return xorg.ChangeScreenSize(size.Width, size.Height, size.Rate)
}

func (manager *DesktopManagerCtx) GetScreenSize() *types.ScreenSize {
	return xorg.GetScreenSize()
}

func (manager *DesktopManagerCtx) SetKeyboardMap(kbd types.KeyboardMap) error {
	// TOOD: Use native API.
	cmd := exec.Command("setxkbmap", "-layout", kbd.Layout, "-variant", kbd.Variant)
	_, err := cmd.Output()
	return err
}

func (manager *DesktopManagerCtx) GetKeyboardMap() (*types.KeyboardMap, error) {
	// TOOD: Use native API.
	cmd := exec.Command("setxkbmap", "-query")
	res, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	kbd := types.KeyboardMap{}

	re := regexp.MustCompile(`layout:\s+(.*)\n`)
	arr := re.FindStringSubmatch(string(res))
	if len(arr) > 1 {
		kbd.Layout = arr[1]
	}

	re = regexp.MustCompile(`variant:\s+(.*)\n`)
	arr = re.FindStringSubmatch(string(res))
	if len(arr) > 1 {
		kbd.Variant = arr[1]
	}

	return &kbd, nil
}

func (manager *DesktopManagerCtx) SetKeyboardModifiers(mod types.KeyboardModifiers) {
	if mod.NumLock != nil {
		xorg.SetKeyboardModifier(xorg.KBD_NUM_LOCK, *mod.NumLock)
	}

	if mod.CapsLock != nil {
		xorg.SetKeyboardModifier(xorg.KBD_CAPS_LOCK, *mod.CapsLock)
	}
}

func (manager *DesktopManagerCtx) GetKeyboardModifiers() types.KeyboardModifiers {
	modifiers := xorg.GetKeyboardModifiers()

	NumLock := (modifiers & xorg.KBD_NUM_LOCK) != 0
	CapsLock := (modifiers & xorg.KBD_CAPS_LOCK) != 0

	return types.KeyboardModifiers{
		NumLock: &NumLock,
		CapsLock: &CapsLock,
	}
}

func (manager *DesktopManagerCtx) GetCursorImage() *types.CursorImage {
	return xorg.GetCursorImage()
}

func (manager *DesktopManagerCtx) GetScreenshotImage() *image.RGBA {
	return xorg.GetScreenshotImage()
}