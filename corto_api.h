#ifndef CORTO_CODEC_H_
#define CORTO_CODEC_H_

#include <stdint.h>
#include <string.h>
#include <stdbool.h>

#ifdef __cplusplus
extern "C" {
#endif

#if defined(WIN32) || defined(WINDOWS) || defined(_WIN32) || defined(_WINDOWS)
#define FLYWAVE_CORTO_API __declspec(dllexport)
#else
#define FLYWAVE_CORTO_API
#endif

typedef struct _corto_encoder_t corto_encoder_t;

FLYWAVE_CORTO_API corto_encoder_t *
corto_new_encoder(uint32_t _nvert, uint32_t _nface, uint32_t entropy);

FLYWAVE_CORTO_API void corto_encoder_free(corto_encoder_t *enc);

FLYWAVE_CORTO_API _Bool corto_encoder_add_positions(corto_encoder_t *enc,
                                                    const float *buffer,
                                                    float q, float *o);

FLYWAVE_CORTO_API _Bool
corto_encoder_add_positions_index32(corto_encoder_t *enc, const float *buffer,
                                    const uint32_t *index, float q, float *o);

FLYWAVE_CORTO_API _Bool
corto_encoder_add_positions_index16(corto_encoder_t *enc, const float *buffer,
                                    const uint16_t *index, float q, float *o);

FLYWAVE_CORTO_API _Bool corto_encoder_add_positions_bits(corto_encoder_t *enc,
                                                         const float *buffer,
                                                         int bits);
                                                         
FLYWAVE_CORTO_API _Bool corto_encoder_add_positions_bits_index32(
    corto_encoder_t *enc, const float *buffer, uint32_t *index, int bits);

FLYWAVE_CORTO_API _Bool corto_encoder_add_positions_bits_index16(
    corto_encoder_t *enc, const float *buffer, uint16_t *index, int bits);

FLYWAVE_CORTO_API _Bool corto_encoder_add_normals_float(corto_encoder_t *enc,
                                                        const float *buffer,
                                                        int bits,
                                                        uint32_t prediction);
FLYWAVE_CORTO_API _Bool corto_encoder_add_normals_short(corto_encoder_t *enc,
                                                        const int16_t *buffer,
                                                        int bits,
                                                        uint32_t prediction);

FLYWAVE_CORTO_API _Bool corto_encoder_add_colors(corto_encoder_t *enc,
                                                 const unsigned char *buffer,
                                                 int rbits, int gbits,
                                                 int bbits, int abits);
FLYWAVE_CORTO_API _Bool corto_encoder_add_colors3(corto_encoder_t *enc,
                                                  const unsigned char *buffer,
                                                  int rbits, int gbits,
                                                  int bbits);

FLYWAVE_CORTO_API _Bool corto_encoder_add_uvs(corto_encoder_t *enc,
                                              const float *buffer, float q);

FLYWAVE_CORTO_API _Bool corto_encoder_add_attribute(
    corto_encoder_t *enc, const char *name, const char *buffer, uint32_t format,
    int components, float q, uint32_t strategy);

FLYWAVE_CORTO_API void corto_encoder_add_group(corto_encoder_t *enc,
                                               int end_triangle);

FLYWAVE_CORTO_API void
corto_encoder_add_group_props(corto_encoder_t *enc, int end_triangle,
                              char **props_keys, char **props_values, int len);

FLYWAVE_CORTO_API void corto_encoder_encode(corto_encoder_t *enc);

typedef struct _corto_decoder_t corto_decoder_t;

FLYWAVE_CORTO_API corto_decoder_t *corto_new_decoder(int len,
                                                     const uint8_t *input);

FLYWAVE_CORTO_API void corto_decoder_free(corto_decoder_t *dec);

FLYWAVE_CORTO_API _Bool corto_encoder_has_attr(corto_decoder_t *dec,
                                               const char *name);

FLYWAVE_CORTO_API _Bool corto_encoder_set_positions(corto_decoder_t *dec,
                                                    float *buffer);

FLYWAVE_CORTO_API _Bool corto_encoder_set_normals_float(corto_decoder_t *dec,
                                                        float *buffer);

FLYWAVE_CORTO_API _Bool corto_encoder_set_normals_short(corto_decoder_t *dec,
                                                        int16_t *buffer);

FLYWAVE_CORTO_API _Bool corto_encoder_set_uvs(corto_decoder_t *dec,
                                              float *buffer);

FLYWAVE_CORTO_API _Bool corto_encoder_set_colors(corto_decoder_t *dec,
                                                 uint8_t *buffer,
                                                 int components);

FLYWAVE_CORTO_API _Bool corto_encoder_set_attribute(corto_decoder_t *dec,
                                                    const char *name,
                                                    char *buffer,
                                                    uint32_t format);

FLYWAVE_CORTO_API void corto_encoder_set_index32(corto_decoder_t *dec,
                                                 uint32_t *buffer);

FLYWAVE_CORTO_API void corto_encoder_set_index16(corto_decoder_t *dec,
                                                 uint16_t *buffer);

FLYWAVE_CORTO_API void corto_encoder_decode(corto_decoder_t *dec);

#ifdef __cplusplus
}
#endif

#endif // CORTO_CODEC_H_