#version 410

in float fred;

void main() {
    gl_FragColor = vec4(fred, 0, 1-fred, 1);
}
