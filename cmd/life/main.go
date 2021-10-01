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

const (
	tw = 2
	th = 2
)

func main() {
	// initialization
	runtime.LockOSThread()

	window := initGLFW()
	defer glfw.Terminate()

	initOpenGL()

	// creating shader program
	prog := func() uint32 {

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
		return prog
	}()

	// creating idicies for rendering
	idxs := func() []float32 {
		idxs := []float32{}

		for j := 0; j < th; j++ {
			for i := 0; i < tw; i++ {
				idxs = append(idxs,
					float32(i+j*tw), float32(i+j*(tw+1)),
					float32(i+j*tw), float32(i+j*(tw+1)+1),
					float32(i+j*tw), float32(i+(j+1)*(tw+1)+1),

					float32(i+j*tw), float32(i+j*(tw+1)),
					float32(i+j*tw), float32(i+(j+1)*(tw+1)),
					float32(i+j*tw), float32(i+(j+1)*(tw+1)+1),
				)
			}
		}
		return idxs
	}()

	vbo := mkvbo(idxs)

	vao := mkvao()
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	vertAttribPtr(prog, "idx\x00", 1, gl.FLOAT, 8, 4)
	vertAttribPtr(prog, "cell\x00", 1, gl.FLOAT, 8, 0)

	// settting tilesize uniform
	gl.UseProgram(prog)
	tsloc := gl.GetUniformLocation(prog, gl.Str("tilesize\x00"))
	gl.Uniform2ui(tsloc, uint32(tw), uint32(th))

	// redering loop
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(prog)

		gl.BindVertexArray(vao)
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(idxs)))
		glfw.PollEvents()
		window.SwapBuffers()
	}
}
