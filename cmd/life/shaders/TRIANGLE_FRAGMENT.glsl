#version 410

in float fred;
in float alpha;

void main() {
    gl_FragColor = vec4(fred, 0, 1-fred, alpha);
}
