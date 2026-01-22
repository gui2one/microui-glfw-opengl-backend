#version 460
uniform sampler2D uTexture;
in vec2 fUvs;
in vec3 fColor;
out vec4 frag_color;
void main() {

    vec4 tex = texture(uTexture, fUvs);
    
    vec3 debugColor = vec3(0.8,0.0,0.2); 
    float sharpAlpha = smoothstep(0.2, 0.8, tex.a);
    if (tex.a < 0.5){
    frag_color = vec4(debugColor, 1.0);
    } else{

    frag_color = vec4(fColor, tex.a);
    }
    // frag_color = vec4(fColor, 1.0);
}