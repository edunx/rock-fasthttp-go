package fasthttp

import (
	"fmt"
	"github.com/edunx/lua"
	pub "github.com/edunx/rock-public-go"
	"os"
)

func (self *vMCache) Path( v string)  {
	self.path = v
}

func (self *vMCache) filename(name string) string {
	return fmt.Sprintf("%s/%s.lua" , self.path , name)
}

func (self *vMCache) load( name string ) *vHandler {

	if self.L == nil {
		self.L = newState()
	}

	v , ok := self.pool.Load( name )
	if ok {
		return v.(vModule).handler
	}

	ret := self.require( name )
	self.store( name , ret)
	return ret
}

func (self *vMCache) store( name string , val *vHandler) {
	if val == nil {
		return
	}

	stat , err := os.Stat( self.filename( name ) )
	if os.IsNotExist( err ) {
		return
	}

	self.pool.Store(name , vModule{ val , stat.ModTime().Unix(), name })
}

func (self *vMCache) require( name string ) *vHandler {
	filename := self.filename( name )

	fn , err := self.L.LoadFile( filename )
	if err != nil {
		pub.Out.Err("load %s fail , err: %v" , filename , err)
		return nil
	}

	self.L.Push( fn )
	if e := self.L.PCall(0 , 1 , nil) ; e != nil {
		pub.Out.Err("pcall %s fail , err: %v" , filename , err)
		return nil
	}

	ret := self.L.Get(1)
	ud , ok := ret.(*lua.LightUserData)
	if !ok {
		pub.Out.Err("got %s ligthUserData fail, err: must be userdata " , filename)
		return nil
	}

	vh , ok := ud.Value.(*vHandler)
	if !ok {
		pub.Out.Err("got %s vhandler fail" , filename)
		return nil
	}

	return vh
}
