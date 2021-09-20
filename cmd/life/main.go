package main

import (
	"runtime"

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
	vbo := mkvbo(triangle)

	// creating position vao
	vao := mkvao()
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	vertAttribPtr(prog, "vp\x00", 2, gl.FLOAT, 12, 0)
	vertAttribPtr(prog, "red\x00", 1, gl.FLOAT, 12, 8)

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(prog)

		gl.BindVertexArray(vao)
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(triangle)/3))
		glfw.PollEvents()
		window.SwapBuffers()
	}
}
