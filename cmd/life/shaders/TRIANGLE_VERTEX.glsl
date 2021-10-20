#version 410

uniform uvec2 tilesize = uvec2(2, 2);

uniform sampler2D state;

in float idx;
in float cell;

out float alpha;

void main() {
    float tw = float(tilesize.x);
    float th = float(tilesize.y);

    float size = tw*th;

    float y = floor((idx+0.1)/(tw+1));
    float x = idx - y*(tw+1);

    x *= 2.0/tw;
    y *= 2.0/th;

    gl_Position = vec4((x-1.0), (-y+1.0), 0.0, 1.0);

    float ty = floor((cell+0.1)/tw);
    float tx = cell-ty*tw;

    tx += 0.5;
    tx /= tw;

    ty += 0.5;
    ty /= th;

    alpha = texture(state, vec2(tx, ty)).x*255;
}
