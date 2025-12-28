package main

import (
	"image/color"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var screenWidth, screenHeight int32 = 1200, 800

func getScreenSize() (int, int) {
	return int(rl.GetScreenWidth()), int(rl.GetScreenHeight())
}

func getMaxDistance() float32 {
	w, h := getScreenSize()
	return float32(math.Sqrt(float64(w*w + h*h)))
}

// Shape interface - all shapes must implement this
type Shape interface {
	Draw()
	RayIntersection(rayStart, rayDir rl.Vector2) (bool, rl.Vector2)
}

// Circle shape
type Circle struct {
	X, Y   float32
	Radius float32
	Color  color.RGBA
}

func (c Circle) Draw() {
	rl.DrawCircleV(rl.Vector2{X: c.X, Y: c.Y}, c.Radius, c.Color)
}

func (c Circle) RayIntersection(rayStart, rayDir rl.Vector2) (bool, rl.Vector2) {
	oc := rl.Vector2{
		X: rayStart.X - c.X,
		Y: rayStart.Y - c.Y,
	}

	a := rayDir.X*rayDir.X + rayDir.Y*rayDir.Y
	b := 2.0 * (oc.X*rayDir.X + oc.Y*rayDir.Y)
	discriminant := b*b - 4*a*(oc.X*oc.X+oc.Y*oc.Y-c.Radius*c.Radius)

	if discriminant < 0 {
		return false, rl.Vector2{}
	}

	t := (-b - float32(math.Sqrt(float64(discriminant)))) / (2.0 * a)
	if t < 0 {
		return false, rl.Vector2{}
	}

	return true, rl.Vector2{
		X: rayStart.X + rayDir.X*t,
		Y: rayStart.Y + rayDir.Y*t,
	}
}

// Rectangle shape
type Rectangle struct {
	X, Y, Width, Height float32
	Color               color.RGBA
}

func (r Rectangle) Draw() {
	rl.DrawRectangle(int32(r.X), int32(r.Y), int32(r.Width), int32(r.Height), r.Color)
}

func (r Rectangle) RayIntersection(rayStart, rayDir rl.Vector2) (bool, rl.Vector2) {
	var closestT float32 = math.MaxFloat32
	hasHit := false

	edges := []struct{ p1, p2 rl.Vector2 }{
		{rl.Vector2{X: r.X, Y: r.Y}, rl.Vector2{X: r.X + r.Width, Y: r.Y}},
		{rl.Vector2{X: r.X + r.Width, Y: r.Y}, rl.Vector2{X: r.X + r.Width, Y: r.Y + r.Height}},
		{rl.Vector2{X: r.X + r.Width, Y: r.Y + r.Height}, rl.Vector2{X: r.X, Y: r.Y + r.Height}},
		{rl.Vector2{X: r.X, Y: r.Y + r.Height}, rl.Vector2{X: r.X, Y: r.Y}},
	}

	for _, edge := range edges {
		hit, t := rayLineIntersection(rayStart, rayDir, edge.p1, edge.p2)
		if hit && t < closestT {
			closestT = t
			hasHit = true
		}
	}

	if hasHit {
		return true, rl.Vector2{
			X: rayStart.X + rayDir.X*closestT,
			Y: rayStart.Y + rayDir.Y*closestT,
		}
	}

	return false, rl.Vector2{}
}

// Line shape (wall)
type Line struct {
	Start, End rl.Vector2
	Color      color.RGBA
	Thickness  float32
}

func (l Line) Draw() {
	rl.DrawLineEx(l.Start, l.End, l.Thickness, l.Color)
}

func (l Line) RayIntersection(rayStart, rayDir rl.Vector2) (bool, rl.Vector2) {
	hit, t := rayLineIntersection(rayStart, rayDir, l.Start, l.End)
	if hit {
		return true, rl.Vector2{
			X: rayStart.X + rayDir.X*t,
			Y: rayStart.Y + rayDir.Y*t,
		}
	}
	return false, rl.Vector2{}
}

// Helper function: Ray-Line segment intersection
func rayLineIntersection(rayStart, rayDir, lineStart, lineEnd rl.Vector2) (bool, float32) {
	v1 := rl.Vector2{X: rayStart.X - lineStart.X, Y: rayStart.Y - lineStart.Y}
	v2 := rl.Vector2{X: lineEnd.X - lineStart.X, Y: lineEnd.Y - lineStart.Y}
	v3 := rl.Vector2{X: -rayDir.Y, Y: rayDir.X}

	dot := v2.X*v3.X + v2.Y*v3.Y
	if math.Abs(float64(dot)) < 0.000001 {
		return false, 0
	}

	t := (v2.X*v1.Y - v2.Y*v1.X) / dot
	u := (v1.X*v3.X + v1.Y*v3.Y) / dot

	if t >= 0 && u >= 0 && u <= 1 {
		return true, t
	}

	return false, 0
}

// Light source
var lightCircle = Circle{
	X:      150,
	Y:      150,
	Radius: 15,
	Color: color.RGBA{
		R: 255,
		G: 255,
		B: 0,
		A: 255,
	},
}

// Create a map with rooms and walls
var obstacles []Shape

func createMap() {
	obstacles = []Shape{}
	wallColor := color.RGBA{R: 100, G: 100, B: 120, A: 255}
	thickness := float32(8)

	// Outer boundary walls
	obstacles = append(obstacles, Line{
		Start:     rl.Vector2{X: 50, Y: 50},
		End:       rl.Vector2{X: 1150, Y: 50},
		Color:     wallColor,
		Thickness: thickness,
	})
	obstacles = append(obstacles, Line{
		Start:     rl.Vector2{X: 1150, Y: 50},
		End:       rl.Vector2{X: 1150, Y: 750},
		Color:     wallColor,
		Thickness: thickness,
	})
	obstacles = append(obstacles, Line{
		Start:     rl.Vector2{X: 1150, Y: 750},
		End:       rl.Vector2{X: 50, Y: 750},
		Color:     wallColor,
		Thickness: thickness,
	})
	obstacles = append(obstacles, Line{
		Start:     rl.Vector2{X: 50, Y: 750},
		End:       rl.Vector2{X: 50, Y: 50},
		Color:     wallColor,
		Thickness: thickness,
	})

	// Room 1 - Top Left
	obstacles = append(obstacles, Line{
		Start:     rl.Vector2{X: 50, Y: 300},
		End:       rl.Vector2{X: 350, Y: 300},
		Color:     wallColor,
		Thickness: thickness,
	})
	obstacles = append(obstacles, Line{
		Start:     rl.Vector2{X: 350, Y: 300},
		End:       rl.Vector2{X: 350, Y: 50},
		Color:     wallColor,
		Thickness: thickness,
	})

	// Room 2 - Top Middle (with doorway)
	obstacles = append(obstacles, Line{
		Start:     rl.Vector2{X: 450, Y: 50},
		End:       rl.Vector2{X: 450, Y: 200},
		Color:     wallColor,
		Thickness: thickness,
	})
	obstacles = append(obstacles, Line{
		Start:     rl.Vector2{X: 450, Y: 280},
		End:       rl.Vector2{X: 450, Y: 350},
		Color:     wallColor,
		Thickness: thickness,
	})
	obstacles = append(obstacles, Line{
		Start:     rl.Vector2{X: 450, Y: 350},
		End:       rl.Vector2{X: 750, Y: 350},
		Color:     wallColor,
		Thickness: thickness,
	})
	obstacles = append(obstacles, Line{
		Start:     rl.Vector2{X: 750, Y: 350},
		End:       rl.Vector2{X: 750, Y: 50},
		Color:     wallColor,
		Thickness: thickness,
	})

	// Room 3 - Top Right
	obstacles = append(obstacles, Line{
		Start:     rl.Vector2{X: 850, Y: 50},
		End:       rl.Vector2{X: 850, Y: 300},
		Color:     wallColor,
		Thickness: thickness,
	})
	obstacles = append(obstacles, Line{
		Start:     rl.Vector2{X: 850, Y: 300},
		End:       rl.Vector2{X: 1150, Y: 300},
		Color:     wallColor,
		Thickness: thickness,
	})

	// Room 4 - Bottom Left (with doorway)
	obstacles = append(obstacles, Line{
		Start:     rl.Vector2{X: 50, Y: 450},
		End:       rl.Vector2{X: 250, Y: 450},
		Color:     wallColor,
		Thickness: thickness,
	})
	obstacles = append(obstacles, Line{
		Start:     rl.Vector2{X: 330, Y: 450},
		End:       rl.Vector2{X: 400, Y: 450},
		Color:     wallColor,
		Thickness: thickness,
	})
	obstacles = append(obstacles, Line{
		Start:     rl.Vector2{X: 400, Y: 450},
		End:       rl.Vector2{X: 400, Y: 750},
		Color:     wallColor,
		Thickness: thickness,
	})

	// Room 5 - Bottom Middle
	obstacles = append(obstacles, Line{
		Start:     rl.Vector2{X: 500, Y: 500},
		End:       rl.Vector2{X: 700, Y: 500},
		Color:     wallColor,
		Thickness: thickness,
	})
	obstacles = append(obstacles, Line{
		Start:     rl.Vector2{X: 700, Y: 500},
		End:       rl.Vector2{X: 700, Y: 750},
		Color:     wallColor,
		Thickness: thickness,
	})

	// Room 6 - Bottom Right (with doorway)
	obstacles = append(obstacles, Line{
		Start:     rl.Vector2{X: 800, Y: 750},
		End:       rl.Vector2{X: 800, Y: 550},
		Color:     wallColor,
		Thickness: thickness,
	})
	obstacles = append(obstacles, Line{
		Start:     rl.Vector2{X: 800, Y: 470},
		End:       rl.Vector2{X: 800, Y: 400},
		Color:     wallColor,
		Thickness: thickness,
	})
	obstacles = append(obstacles, Line{
		Start:     rl.Vector2{X: 800, Y: 400},
		End:       rl.Vector2{X: 1150, Y: 400},
		Color:     wallColor,
		Thickness: thickness,
	})

	// Add some circular obstacles/pillars
	obstacles = append(obstacles, Circle{
		X:      250,
		Y:      150,
		Radius: 30,
		Color:  color.RGBA{R: 80, G: 80, B: 100, A: 255},
	})
	obstacles = append(obstacles, Circle{
		X:      600,
		Y:      200,
		Radius: 35,
		Color:  color.RGBA{R: 80, G: 80, B: 100, A: 255},
	})
	obstacles = append(obstacles, Circle{
		X:      950,
		Y:      180,
		Radius: 40,
		Color:  color.RGBA{R: 80, G: 80, B: 100, A: 255},
	})
	obstacles = append(obstacles, Circle{
		X:      200,
		Y:      600,
		Radius: 45,
		Color:  color.RGBA{R: 80, G: 80, B: 100, A: 255},
	})
	obstacles = append(obstacles, Circle{
		X:      550,
		Y:      650,
		Radius: 30,
		Color:  color.RGBA{R: 80, G: 80, B: 100, A: 255},
	})

	// Add some rectangular obstacles/furniture
	obstacles = append(obstacles, Rectangle{
		X:      900,
		Y:      500,
		Width:  80,
		Height: 50,
		Color:  color.RGBA{R: 120, G: 80, B: 60, A: 255},
	})
	obstacles = append(obstacles, Rectangle{
		X:      150,
		Y:      380,
		Width:  60,
		Height: 40,
		Color:  color.RGBA{R: 120, G: 80, B: 60, A: 255},
	})
}

var maxDistance = getMaxDistance()

func updateMaxDistance() {
	maxDistance = getMaxDistance()
}

func draw() {
	rl.BeginDrawing()

	// Draw all obstacles (walls, pillars, furniture)
	for _, obstacle := range obstacles {
		obstacle.Draw()
	}

	// Draw the light source
	lightCircle.Draw()

	numRays := 360
	// Draw rays
	for i := 0; i < numRays; i++ {
		angle := float64(360.0/float32(numRays)*float32(i)) * math.Pi / 180.0

		start := rl.Vector2{
			X: lightCircle.X + float32(math.Cos(angle))*lightCircle.Radius,
			Y: lightCircle.Y + float32(math.Sin(angle))*lightCircle.Radius,
		}

		rayDir := rl.Vector2{
			X: float32(math.Cos(angle)) * maxDistance,
			Y: float32(math.Sin(angle)) * maxDistance,
		}

		end := rl.Vector2{
			X: lightCircle.X + rayDir.X,
			Y: lightCircle.Y + rayDir.Y,
		}

		// Check collision with ALL obstacles
		closestDistance := maxDistance
		hasCollision := false
		var closestPoint rl.Vector2

		for _, obstacle := range obstacles {
			hit, intersection := obstacle.RayIntersection(start, rayDir)
			if hit {
				dx := intersection.X - start.X
				dy := intersection.Y - start.Y
				distance := float32(math.Sqrt(float64(dx*dx + dy*dy)))

				if distance < closestDistance {
					closestDistance = distance
					closestPoint = intersection
					hasCollision = true
				}
			}
		}

		if hasCollision {
			end = closestPoint
		}

		// Make rays semi-transparent yellow
		rl.DrawLineV(start, end, color.RGBA{R: 255, G: 255, B: 0, A: 30})
	}

	// Instructions
	rl.DrawText("Click and drag to move light", 10, 10, 20, rl.RayWhite)
	rl.DrawText("Press F11 for fullscreen", 10, 35, 20, rl.RayWhite)

	if rl.IsMouseButtonDown(rl.MouseButtonLeft) {
		lightCircle.X = float32(rl.GetMouseX())
		lightCircle.Y = float32(rl.GetMouseY())
	}
}

func main() {
	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.InitWindow(screenWidth, screenHeight, "2D Raycasting - Room Map")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	// Create the map
	createMap()

	for !rl.WindowShouldClose() {
		// Toggle fullscreen with F11
		if rl.IsKeyPressed(rl.KeyF11) {
			rl.ToggleFullscreen()
			updateMaxDistance()
		}

		// Update maxDistance if window was resized
		if rl.IsWindowResized() {
			updateMaxDistance()
		}

		rl.ClearBackground(color.RGBA{R: 20, G: 20, B: 30, A: 255})
		draw()
		rl.EndDrawing()
	}
}
