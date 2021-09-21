#version 410

uniform vec2 tilesize = vec2(2.0, 2.0);

in float idx;
in float cell;

out float fred;

void main() {
    float y = floor((idx+0.1)/(tilesize.x+1));
    float x = idx - y*(tilesize.x+1);

    x *= 2.0/tilesize.x;
    y *= 2.0/tilesize.y;

    gl_Position = vec4((x-1.0), (-y+1.0), 0.0, 1.0);

    fred = cell/(tilesize.x*tilesize.y-1);
}
