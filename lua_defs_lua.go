package lua

/*
#include <lua.h>
 #include <lauxlib.h>
 #include <lualib.h>

*/
import "C"

type LuaValType int

const (
	LUA_TNIL           = LuaValType(C.LUA_TNIL)
	LUA_TNUMBER        = LuaValType(C.LUA_TNUMBER)
	LUA_TBOOLEAN       = LuaValType(C.LUA_TBOOLEAN)
	LUA_TSTRING        = LuaValType(C.LUA_TSTRING)
	LUA_TTABLE         = LuaValType(C.LUA_TTABLE)
	LUA_TFUNCTION      = LuaValType(C.LUA_TFUNCTION)
	LUA_TUSERDATA      = LuaValType(C.LUA_TUSERDATA)
	LUA_TTHREAD        = LuaValType(C.LUA_TTHREAD)
	LUA_TLIGHTUSERDATA = LuaValType(C.LUA_TLIGHTUSERDATA)
)

const (
	LUA_OK            = C.LUA_OK
	LUA_VERSION       = C.LUA_VERSION
	LUA_RELEASE       = C.LUA_RELEASE
	LUA_VERSION_NUM   = C.LUA_VERSION_NUM
	LUA_COPYRIGHT     = C.LUA_COPYRIGHT
	LUA_AUTHORS       = C.LUA_AUTHORS
	LUA_MULTRET       = C.LUA_MULTRET
	LUA_REGISTRYINDEX = C.LUA_REGISTRYINDEX
	LUA_YIELD         = C.LUA_YIELD
	LUA_ERRRUN        = C.LUA_ERRRUN
	LUA_ERRSYNTAX     = C.LUA_ERRSYNTAX
	LUA_ERRMEM        = C.LUA_ERRMEM
	LUA_ERRERR        = C.LUA_ERRERR
	LUA_TNONE         = C.LUA_TNONE
	LUA_MINSTACK      = C.LUA_MINSTACK
	LUA_GCSTOP        = C.LUA_GCSTOP
	LUA_GCRESTART     = C.LUA_GCRESTART
	LUA_GCCOLLECT     = C.LUA_GCCOLLECT
	LUA_GCCOUNT       = C.LUA_GCCOUNT
	LUA_GCCOUNTB      = C.LUA_GCCOUNTB
	LUA_GCSTEP        = C.LUA_GCSTEP
	LUA_GCSETPAUSE    = C.LUA_GCSETPAUSE
	LUA_GCSETSTEPMUL  = C.LUA_GCSETSTEPMUL
	LUA_HOOKCALL      = C.LUA_HOOKCALL
	LUA_HOOKRET       = C.LUA_HOOKRET
	LUA_HOOKLINE      = C.LUA_HOOKLINE
	LUA_HOOKCOUNT     = C.LUA_HOOKCOUNT
	LUA_MASKCALL      = C.LUA_MASKCALL
	LUA_MASKRET       = C.LUA_MASKRET
	LUA_MASKLINE      = C.LUA_MASKLINE
	LUA_MASKCOUNT     = C.LUA_MASKCOUNT
	LUA_ERRFILE       = C.LUA_ERRFILE
	LUA_NOREF         = C.LUA_NOREF
	LUA_REFNIL        = C.LUA_REFNIL
	LUA_FILEHANDLE    = C.LUA_FILEHANDLE
	LUA_COLIBNAME     = C.LUA_COLIBNAME
	LUA_TABLIBNAME    = C.LUA_TABLIBNAME
	LUA_IOLIBNAME     = C.LUA_IOLIBNAME
	LUA_OSLIBNAME     = C.LUA_OSLIBNAME
	LUA_STRLIBNAME    = C.LUA_STRLIBNAME
	LUA_MATHLIBNAME   = C.LUA_MATHLIBNAME
	LUA_DBLIBNAME     = C.LUA_DBLIBNAME
	LUA_LOADLIBNAME   = C.LUA_LOADLIBNAME
)

type LUA_TableType int8
