package native

/* #include <stdint.h>

size_t convert(unsigned long long integer) {
	return (size_t) integer;
}
*/
import "C"

func toInt64(in int64) C.size_t {
	return C.convert(C.ulonglong(in))
}

func intToInt64(in int) C.size_t {
	return toInt64(int64(in))
}
