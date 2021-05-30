#ifndef CORTO_CODEC_H
#define CORTO_CODEC_H

#include "decoder.h"
#include <string.h>

// If compiling with Visual Studio
#if defined(_MSC_VER)
#define EXPORT_API __declspec(dllexport)
#else
// Other platforms don't need this
#define EXPORT_API
#endif

#endif // !CORTO_UNITY_PLUGIN_H