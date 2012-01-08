package main

import "time"

import "sge"
import "s3dm"
import "atom/sdl"

func main() {
	assets := sge.NewAssets()
	world := sge.NewWorld()
	view := sge.NewView(800, 600, 0.1, 10000)
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
	ticker := time.Tick(int64(1e9) / 60)
	keystate := sdl.GetKeyState()
	yaw := float64(0)
	pitch := float64(0)
	cameraChanged := false
mainloop:
	for {
		select {
		case t := <-ticker:
			if keystate[sdl.K_UP] == 1 {
				pitch += 3
				if pitch > 90 {
					pitch = 90
				}
				cameraChanged = true
			}
			if keystate[sdl.K_DOWN] == 1 {
				pitch -= 3
				if pitch < -90 {
					pitch = -90
				}
				cameraChanged = true
			}
			if keystate[sdl.K_LEFT] == 1 {
				yaw += 3
				if yaw > 360 {
					yaw -= 360
				}
				cameraChanged = true
			}
			if keystate[sdl.K_RIGHT] == 1 {
				yaw -= 3
				if yaw < 0 {
					yaw += 360
				}
				cameraChanged = true
			}
			if keystate[sdl.K_w] == 1 {
				view.Camera.MoveLocal(s3dm.NewV3(0, 0, -0.01))
			}
			if keystate[sdl.K_s] == 1 {
				view.Camera.MoveLocal(s3dm.NewV3(0, 0, 0.01))
			}
			if keystate[sdl.K_a] == 1 {
				view.Camera.MoveLocal(s3dm.NewV3(-0.01, 0, 0))
			}
			if keystate[sdl.K_d] == 1 {
				view.Camera.MoveLocal(s3dm.NewV3(0.01, 0, 0))
			}
			if cameraChanged {
				cameraChanged = false
				view.Camera.SetIdentity()
				view.Camera.RotateGlobal(yaw, s3dm.NewV3(0, 1, 0))
				view.Camera.RotateGlobal(pitch, s3dm.NewV3(1, 0, 0))
			}
			world.Update(t)
			world.Render(view)
		case event := <-sdl.Events:
			switch e := event.(type) {
			case sdl.QuitEvent:
				break mainloop
			}
		}
	}
}
