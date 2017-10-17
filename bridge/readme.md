# Cross-compilation bridge

This folder holds the tooling required to cross-compile the SDK into other languages. Implementation status for each language is tracked in the main project readme.

## Building the bridge

From this `bridge` folder:

```
./make.sh
```

This will build the SDK, confirm the generated bridge is valid, and compile the C header files along with any other language frontends.

## Developer notes

### Useful cgo functions

```go
// Free C heap memory. Include stdlib.h to use
func C.free(unsafe.Pointer)

// Malloc a C string on the C heap, must be freed.
func C.CString(string) *C.char

// Malloc a C array on the C heap, must be freed.
func C.CBytes([]byte) unsafe.Pointer

// C string to Go string
func C.GoString(*C.char) string

// C data with explicit length to Go string
func C.GoStringN(*C.char, C.int) string

// C data with explicit length to Go []byte
func C.GoBytes(unsafe.Pointer, C.int) []byte
```

### Type chart

This chart combines documentation from [Cgo](https://golang.org/cmd/cgo) and [Ctypes](https://docs.python.org/2.7/library/ctypes.html).
Some comparisons may be duplicated.

C type                     | Golang CGO      | Python Ctype     | Python type
---------------------------|-----------------|------------------|-------------------------------
_Bool                      |                 | c_bool           | bool (1)
char                       | C.char          | c_char OR c_byte | 1-character string OR int/long
wchar_t                    |                 | c_wchar          | 1-character unicode string
signed char                | C.schar         |                  |
unsigned char              | C.uchar         | c_ubyte          | int/long
short                      | C.short         | c_short          | int/long
unsigned short             | C.ushort        | c_ushort         | int/long
int                        | C.int           | c_int            | int/long
unsigned int               | C.uint          | c_uint           | int/long
long                       | C.long          | c_long           | int/long
unsigned long              | C.ulong         | c_ulong          | int/long
long long                  | C.longlong      | c_longlong       | int/long
unsigned long long         | C.ulonglong     | c_ulonglong      | int/long
float                      | C.float         | c_float          | float
double                     | C.double        | c_double         | float
complex float              | C.complexfloat  |                  |
complex double             | C.complexdouble |                  |
__int128_t                 | [16]byte        |                  |
__uint128_t                | [16]byte        |                  |
long double                |                 | c_longdouble     | float
char * (NUL terminated)    |                 | c_char_p         | string or None
wchar_t * (NUL terminated) |                 | c_wchar_p        | unicode or None
void *                     | unsafe.Pointer  | c_void_p         | int/long or None


To access a struct, union, or enum type directly, prefix it with struct_, union_, or enum_, as in C.struct_stat.
The size of any C type T is available as C.sizeof_T, as in C.sizeof_struct_stat.

### Useful links

https://golang.org/cmd/cgo/
https://docs.python.org/2.7/library/ctypes.html

https://blog.heroku.com/see_python_see_python_go_go_python_go

https://blog.golang.org/c-go-cgo
