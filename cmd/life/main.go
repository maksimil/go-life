package main

import (
	"fmt"
	"log"
	"runtime"
	"strings"
	"unsafe"

	gl "github.com/go-gl/gl/v3.2-core/gl"
	glfw "github.com/go-gl/glfw/v3.3/glfw"
)

const (
	width  = 500
	height = 500
)

var (
	triangle = []float32{
		-0.5, -0.5, 1,
		0.5, -0.5, 0,
		0.5, 0.5, 0,
	}
)

func main() {
	runtime.LockOSThread()

	window := initGLFW()
	defer glfw.Terminate()

	initOpenGL()

	// creating shader program
	vertexshader, err := compileShader(TRIANGLE_VERTEX, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragmentshader, err := compileShader(TRIANGLE_FRAGMENT, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexshader)
	gl.AttachShader(prog, fragmentshader)
	gl.LinkProgram(prog)

	// creating vertex buffer
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(triangle), gl.Ptr(triangle), gl.STATIC_DRAW)

	// creating position vao
	var vao uint32
	vp := getAttribLoc(prog, "vp\x00")
	red := getAttribLoc(prog, "red\x00")
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.EnableVertexAttribArray(uint32(vp))
	gl.VertexAttribPointer(getAttribLoc(prog, "vp\x00"), 2, gl.FLOAT, false, 12, nil)
	gl.EnableVertexAttribArray(uint32(red))
	gl.VertexAttribPointer(uint32(red), 1, gl.FLOAT, false, 12, unsafe.Pointer(uintptr(8)))

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(prog)

		gl.BindVertexArray(vao)
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(triangle)/3))
		glfw.PollEvents()
		window.SwapBuffers()
	}
}

func getAttribLoc(prog uint32, name string) uint32 {
	return uint32(gl.GetAttribLocation(prog, gl.Str(name)))
}

func initGLFW() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 2)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "Conway's Game of Life", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	return window
}

// initOpenGL initializes OpenGL and returns an intiialized program.
func initOpenGL() {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)
}

func compileShader(source string, stype uint32) (uint32, error) {
	shader := gl.CreateShader(stype)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength)+1)
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}
