#version 410

uniform uvec2 tilesize = uvec2(2, 2);

uniform sampler2D state;

in float idx;
in float cell;

out float fred;
out float alpha;

void main() {
    float size = tilesize.x*tilesize.y;

    float y = floor(idx/(tilesize.x+1));
    float x = idx - y*(tilesize.x+1);

    x *= 2.0/tilesize.x;
    y *= 2.0/tilesize.y;

    gl_Position = vec4((x-1.0), (-y+1.0), 0.0, 1.0);

    fred = cell/(size-1);

    float ty = floor(cell/tilesize.x);
    float tx = (cell-ty*tilesize.x)/tilesize.x;
    ty /= tilesize.y;
    tx += 1/(2*tilesize.x);
    ty += 1/(2*tilesize.y);

    alpha = texture(state, vec2(tx, ty)).x;
}
