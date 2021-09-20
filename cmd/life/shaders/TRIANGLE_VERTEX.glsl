#version 410

in vec2 vp;
in float red;

out float fred;

void main() {
    fred = red;
    gl_Position = vec4(vp, 0.0, 1.0);
}
