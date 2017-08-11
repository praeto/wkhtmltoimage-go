package api

/*
#cgo CFLAGS: -I/usr/local/include
#cgo LDFLAGS: -L/usr/local/lib -lwkhtmltox -Wall -ansi -pedantic -ggdb
#include <stdbool.h>
#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <wkhtmltox/image.h>
extern void progress_changed_cb(void*, const int);
extern void error_cb(void*, char *msg);
extern void warning_cb(void*, char *msg);
extern void phase_changed_cb(void*);
extern void finished_cb(void*, const int);
static void setup_callbacks(wkhtmltoimage_converter * c) {
  wkhtmltoimage_set_progress_changed_callback(c, (wkhtmltoimage_int_callback)progress_changed_cb);
  wkhtmltoimage_set_error_callback(c, (wkhtmltoimage_str_callback)error_cb);
  wkhtmltoimage_set_warning_callback(c, (wkhtmltoimage_str_callback)warning_cb);
  wkhtmltoimage_set_phase_changed_callback(c, (wkhtmltoimage_void_callback)phase_changed_cb);
  wkhtmltoimage_set_finished_callback(c, (wkhtmltoimage_int_callback)finished_cb);
}
*/
import "C"

import (
	"unsafe"
)

type GlobalSettings struct {
	s *C.wkhtmltoimage_global_settings
}

type Converter struct {
	c               *C.wkhtmltoimage_converter
	ProgressChanged func(*Converter, int)
	Error           func(*Converter, string)
	Warning         func(*Converter, string)
	Phase           func(*Converter)
	Finished func(*Converter, int)
	quiet			bool
}

var converterMap map[unsafe.Pointer]*Converter

func init() {
	converterMap = map[unsafe.Pointer]*Converter{}
	C.wkhtmltoimage_init(C.false)
}

func NewGlobalSettings() *GlobalSettings {
	return &GlobalSettings{s: C.wkhtmltoimage_create_global_settings()}
}

func (gs *GlobalSettings) Set(name, value string) {
	c_name := C.CString(name)
	c_value := C.CString(value)
	defer C.free(unsafe.Pointer(c_name))
	defer C.free(unsafe.Pointer(c_value))
	C.wkhtmltoimage_set_global_setting(gs.s, c_name, c_value)
}

func (gs *GlobalSettings) NewConverter(html string, quiet bool) *Converter {
	cHtml := C.CString(html)
	defer C.free(unsafe.Pointer(cHtml))
	c := &Converter{c: C.wkhtmltoimage_create_converter(gs.s, cHtml), quiet: quiet}
	C.setup_callbacks(c.c)

	return c
}

//export progress_changed_cb
func progress_changed_cb(p unsafe.Pointer, i C.int) {
	conv := converterMap[p]
	if conv.ProgressChanged != nil && !conv.quiet {
		conv.ProgressChanged(conv, int(i))
	}
}

//export error_cb
func error_cb(p unsafe.Pointer, msg *C.char) {
	conv := converterMap[p]
	if conv.Error != nil && !conv.quiet {
		conv.Error(conv, C.GoString(msg))
	}
}

//export warning_cb
func warning_cb(p unsafe.Pointer, msg *C.char) {
	conv := converterMap[p]
	if conv.Warning != nil && !conv.quiet {
		conv.Warning(conv, C.GoString(msg))
	}
}

//export phase_changed_cb
func phase_changed_cb(p unsafe.Pointer) {
	conv := converterMap[p]
	if conv.Phase != nil && !conv.quiet {
		conv.Phase(conv)
	}
}

//export finished_cb
func finished_cb(c unsafe.Pointer, s C.int) {
	conv := converterMap[c]
	if conv.Finished != nil && !conv.quiet {
		conv.Finished(conv, int(s))
	}
}

func (converter *Converter) Convert() int {

	// To route callbacks right, we need to save a reference
	// to the converter object, base on the pointer.
	converterMap[unsafe.Pointer(converter.c)] = converter
	status := C.wkhtmltoimage_convert(converter.c)
	delete(converterMap, unsafe.Pointer(converter.c))
	if status != C.int(0) {
		return converter.ErrorCode()
	}
	return 0
}

func (converter *Converter) Output() (int64, string) {
	defer converter.Destroy()
	cc := C.CString("")
	ccc := (**C.uchar)(unsafe.Pointer(&cc))
	ll := C.wkhtmltoimage_get_output(converter.c, ccc)
	co := C.GoStringN((*C.char)(unsafe.Pointer(*ccc)), C.int(ll))
	return int64(ll), co
}

func (converter *Converter) ErrorCode() int {
	return int(C.wkhtmltoimage_http_error_code(converter.c))
}

func (converter *Converter) CurrentPhase() (int, string) {
	cpic := C.wkhtmltoimage_current_phase(converter.c)
	cpi := int(cpic)
	cps := C.GoString(C.wkhtmltoimage_phase_description(converter.c, cpic))
	return cpi, cps
}

func (converter *Converter) Destroy() {
	C.wkhtmltoimage_destroy_converter(converter.c)
}
