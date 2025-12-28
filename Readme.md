# 2D Raycasting Demo

A real-time 2D raycasting simulation built with Go and Raylib, demonstrating dynamic lighting and shadow effects in a room-based environment.

![Screenshot](screenshots/demo.png)

## Features

- **Real-time raycasting** with 360 rays emanating from a movable light source
- **Multiple room layout** with walls, doorways, and corridors
- **Various obstacle types**: walls (lines), pillars (circles), and furniture (rectangles)
- **Dynamic shadows** that respond to light position
- **Interactive light control** - drag the light source anywhere on the map
- **Resizable window** with fullscreen support
- **Collision detection** with multiple shape types

## Prerequisites

- Go 1.18 or higher
- Raylib dependencies for your platform

### Installing Raylib Dependencies

**Linux (Debian/Ubuntu):**
```bash
sudo apt-get install libgl1-mesa-dev libxi-dev libxcursor-dev libxrandr-dev libxinerama-dev
```

**macOS:**
```bash
brew install glfw
```

**Windows:**
No additional dependencies required - raylib-go includes pre-built binaries.

## Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/2d-raycasting.git
cd 2d-raycasting
```

2. Download dependencies:
```bash
go mod download
```

3. Run the project:
```bash
go run main.go
```

## Usage

- **Click and drag** with left mouse button to move the light source
- **Press F11** to toggle fullscreen mode
- **Resize the window** by dragging window borders

## How It Works

The raycasting algorithm shoots 360 rays from the light source in all directions. Each ray:
1. Checks for intersection with all obstacles in the scene
2. Finds the closest collision point
3. Stops at that point, creating realistic shadows

The project implements ray-shape intersection algorithms for:
- **Lines** (walls) - Ray-line segment intersection
- **Circles** (pillars) - Quadratic equation solution
- **Rectangles** (furniture) - Four edge intersection checks

## Project Structure

```
2d-raycasting/
├── main.go           # Main application code
├── go.mod            # Go module dependencies
├── go.sum            # Dependency checksums
├── README.md         # This file
└── screenshots/      # Screenshot folder
    └── demo.png      # Demo screenshot
```

## Dependencies

- [raylib-go](https://github.com/gen2brain/raylib-go) - Go bindings for Raylib
- [Raylib](https://www.raylib.com/) - Simple and easy-to-use library for game development

## Technical Details

- **Language:** Go 1.25.5
- **Graphics Library:** Raylib via raylib-go
- **Ray Count:** 360 rays per frame
- **Collision Detection:** Custom ray-shape intersection algorithms
- **Target FPS:** 60

## Performance

The simulation runs at 60 FPS with 360 rays being cast per frame. Performance may vary based on:
- Number of obstacles in the scene
- Window resolution
- Hardware capabilities

## Future Improvements

- [ ] Add textured walls
- [ ] Implement field of view restrictions
- [ ] Add multiple light sources
- [ ] Support for polygon obstacles
- [ ] Light intensity falloff with distance
- [ ] Color/brightness customization
- [ ] Save/load custom maps

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is open source and available under the [MIT License](LICENSE).

## Acknowledgments

- Built with [Raylib](https://www.raylib.com/)
- Inspired by classic raycasting techniques used in early 3D games

## Contact

Email: prashant2002singh915@gmail.com

Project Link: [https://github.com/prashant2002pd/raycast-go](https://github.com/prashant2002pd/raycast-go)
