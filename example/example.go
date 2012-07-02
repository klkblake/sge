package main

import (
	"time"
)

import (
	"github.com/klkblake/Go-SDL/sdl"
	"github.com/klkblake/s3dm"
	"github.com/klkblake/sge"
)

func main() {
	assets := sge.NewAssets()
	world := sge.NewWorld()
	view := sge.NewView("Example", 800, 600, 0.1, 10000)
	defer view.Close()
	// Test our shaders
	shader := assets.ShaderProgram("default", "defaultCube")
	texture := assets.TextureCubeMap(&[6]string{
		"skybox-pos-x.png",
		"skybox-neg-x.png",
		"skybox-pos-y.png",
		"skybox-neg-y.png",
		"skybox-pos-z.png",
		"skybox-neg-z.png",
	})
	world.Skybox = sge.NewSkybox(texture, shader, 10000)
	view.SetBackgroundColor(0.25, 0.5, 0.75)
	ticker := time.Tick(time.Second / 60)
	keystate := sdl.GetKeyState()
	cameraChanged := false
	cameraRotated := false
	last := time.Now()
mainloop:
	for {
		yaw := float64(0)
		pitch := float64(0)
		select {
		case t := <-ticker:
			delta := t.Sub(last)
			last = t
			if keystate[sdl.K_UP] == 1 {
				pitch = 0.05
				cameraRotated = true
			}
			if keystate[sdl.K_DOWN] == 1 {
				pitch = -0.05
				cameraRotated = true
			}
			if keystate[sdl.K_LEFT] == 1 {
				yaw = 0.05
				cameraRotated = true
			}
			if keystate[sdl.K_RIGHT] == 1 {
				yaw = -0.05
				cameraRotated = true
			}
			if keystate[sdl.K_w] == 1 {
				view.Camera.Position.Z -= 0.01
				cameraChanged = true
			}
			if keystate[sdl.K_s] == 1 {
				view.Camera.Position.Z += 0.01
				cameraChanged = true
			}
			if keystate[sdl.K_a] == 1 {
				view.Camera.Position.X -= 0.01
				cameraChanged = true
			}
			if keystate[sdl.K_d] == 1 {
				view.Camera.Position.X += 0.01
				cameraChanged = true
			}
			if cameraRotated {
				cameraRotated = false
				cameraChanged = true
				view.Camera.Rotation = s3dm.AxisAngle(s3dm.V3{0, 1, 0}, yaw).Mul(view.Camera.Rotation).Mul(s3dm.AxisAngle(s3dm.V3{1, 0, 0}, pitch))
			}
			if cameraChanged {
				cameraChanged = false
				view.Update()
			}
			world.Update(delta.Seconds())
			world.Render(view)
			sge.FlushGL()
		case event := <-sdl.Events:
			switch event.(type) {
			case sdl.QuitEvent:
				break mainloop
			}
		}
	}
}
