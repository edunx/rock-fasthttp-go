package fasthttp

import (
	"fmt"
	"github.com/edunx/lua"
	pub "github.com/edunx/rock-public-go"
	"os"
	"time"
)

func (self *vMCache) Path( v string)  {
	self.path = v
}

func (self *vMCache) filename(name string) string {
	return fmt.Sprintf("%s/%s.lua" , self.path , name)
}

func (self *vMCache) load( name string ) *vHandler {

	if self.L == nil {
		self.L = newHttpThreadState()
	}

	v , ok := self.pool.Load( name )
	if ok {
		return v.(vModule).handler
	}

	ret := self.require( name )
	self.store( name , ret , nil)
	return ret
}

func (self *vMCache) store( name string , val *vHandler , stat os.FileInfo ) {
	if val == nil {
		return
	}

	var err error
	if stat != nil {
		goto DONE
	}

	stat , err = os.Stat( self.filename( name ) )
	if os.IsNotExist( err ) {
		return
	}

DONE:
	pub.Out.Debug("load %s vhandler success" ,  name )
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
		pub.Out.Err("got %s vhandler fail, err: must be lightuserdata , type: %T " , filename , ud)
		self.L.Pop(1)
		return nil
	}

	vh , ok := ud.Value.(*vHandler)
	if !ok {
		pub.Out.Err("got %s vhandler fail" , filename)
		self.L.Pop(1)
		return nil
	}

	self.L.Pop(1)
	pub.Out.Debug("require %s vhandler success" , name)
	return vh
}

func (self *vMCache) sync() {
	tk := time.NewTicker( 500 * time.Millisecond )
	for range tk.C {
		self.pool.Range(func(k interface{} , v interface{} ) bool {
			name := k.(string)
			obj := v.(vModule)
			filename := self.filename( name )
			stat , err := os.Stat( filename )
			if os.IsNotExist( err ) {
				self.pool.Delete( name )
				return false
			}

			if stat.ModTime().Unix() == obj.modTime {
				return  false
			}

			ret := self.require( name )
			if ret == nil {
				return false
			}

			self.store( name , ret , stat)
			obj.handler = nil

			return false
		})
	}
}
