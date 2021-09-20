package clock

// #cgo CFLAGS: -Wall -pthread -fPIC
// #include <stdlib.h>
// #include "dgt3000.h"
import "C"
import "unsafe"

// Get direct access to BCM2708/9
func Init() int {
	return int(C.dgt3000Init())
}

// Configure DGT3000 to on, central control and mode 25
func Configure() int {
	return int(C.dgt3000Configure())
}

// Send a wake command to the DGT3000
func Wake() int {
	return int(C.dgt3000Wake())
}

// Send set central control command to DGT3000
func SetCC() int {
	return int(C.dgt3000SetCC())
}

// Send set mode 25 to DGT3000
func Mode25() int {
	return int(C.dgt3000Mode25())
}

// Send end display to DGT3000 to clear display
func EndDisplay() int {
	return int(C.dgt3000EndDisplay())
}

// Return current button state
func GetButtonState() int {
	return int(C.dgt3000GetButtonState())
}

// Send set display command
func SetDisplay(msg string) int {
	cs := C.CString(msg)
	ret := int(C.dgt3000SetDisplay(cs))
	C.free(unsafe.Pointer(cs))
	return ret
}

// Send set and run command
func SetNRun(lr byte, lh byte, lm byte, ls byte, rr byte, rh byte, rm byte, rs byte) int {
	c_lr := (*C.char)(unsafe.Pointer(&lr))
	c_lh := (*C.char)(unsafe.Pointer(&lh))
	c_lm := (*C.char)(unsafe.Pointer(&lm))
	c_ls := (*C.char)(unsafe.Pointer(&ls))
	c_rr := (*C.char)(unsafe.Pointer(&rr))
	c_rh := (*C.char)(unsafe.Pointer(&rh))
	c_rm := (*C.char)(unsafe.Pointer(&rm))
	c_rs := (*C.char)(unsafe.Pointer(&rs))

	ret := int(C.dgt3000SetNRun(*c_lr, *c_lh, *c_lm, *c_ls, *c_rr, *c_rh, *c_rm, *c_rs))
	C.free(unsafe.Pointer(c_lr))
	C.free(unsafe.Pointer(c_lh))
	C.free(unsafe.Pointer(c_lm))
	C.free(unsafe.Pointer(c_ls))
	C.free(unsafe.Pointer(c_rr))
	C.free(unsafe.Pointer(c_rh))
	C.free(unsafe.Pointer(c_rm))
	C.free(unsafe.Pointer(c_rs))
	return ret
}

// Send set and run command to dgt3000 with current clock values
func Run(lr byte, rr byte) int {
	c_lr := (*C.char)(unsafe.Pointer(&lr))
	c_rr := (*C.char)(unsafe.Pointer(&lr))
	ret := int(C.dgt3000Run(*c_lr, *c_rr))
	C.free(unsafe.Pointer(c_lr))
	C.free(unsafe.Pointer(c_rr))
	return ret
}

// End and set display
func Display(text string, beep byte, ld byte, rd byte) int {
  c_text := C.CString(text)
  c_beep := (*C.char)(unsafe.Pointer(&beep))
	c_ld := (*C.char)(unsafe.Pointer(&ld))
	c_rd := (*C.char)(unsafe.Pointer(&rd))
	ret := int(C.dgt3000Display(c_text, *c_beep, *c_ld, *c_rd))
  C.free(unsafe.Pointer(c_beep))
	C.free(unsafe.Pointer(c_ld))
	C.free(unsafe.Pointer(c_rd))
	return ret
}

/*

// check for messages from dgt3000
void *dgt3000Receive(void *);

// wait for an Ack message
int dgt3000GetAck(char adr, char cmd, long long int timeOut);

// return last received time
void dgt3000GetTime(char time[]);

// return buttons pressed
int dgt3000GetButton(char *buttons, char *time);

*/

// turn off dgt3000
func Off(returnMode byte) int {
	cs := (*C.char)(unsafe.Pointer(&returnMode))
	ret := int(C.dgt3000Off(*cs))
	C.free(unsafe.Pointer(cs))
	return ret
}

// stop receiving
func Stop() int {
	return int(C.dgt3000Stop())
}
