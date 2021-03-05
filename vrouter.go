package fasthttp

import (
	"fmt"
	"github.com/edunx/lua"
	pub "github.com/edunx/rock-public-go"
	"github.com/fasthttp/router"
	"os"
	"time"
)

func(self *vRCache) Path( v string ) {
	stat , err := os.Stat( v )
	if os.IsNotExist( err ) {
		pub.Out.Err( "not found %s" , v )
		return
	}

	if !stat.IsDir() {
		pub.Out.Err("%s must be dir" , v)
		return
	}

	self.path = v
}

func (self *vRCache) filename( name string ) string {
	return fmt.Sprintf("%s/%s.lua" , self.path , name )
}

func (self *vRCache) NewState() *lua.LState {
	return newState()
}

func (self *vRCache) Compile( name string  ) *vRouter {

	filename := self.filename( name )
	stat , err := os.Stat( filename )
	if os.IsNotExist(err) {
		//pub.Out.Debug("not found %s vhost" , filename )
		return nil
	}

	//编译文件
	L := self.NewState()
	if e := L.DoFile( filename ); e != nil {
		pub.Out.Debug("load %s vhost fail , err: %v" , name , e)
		return nil
	}

	v := &vRouter{
		L:L,
		name: name,
		modTime: stat.ModTime().Unix(),
	}

	self.pool.Store(name , v)

	return v
}

func (self *vRCache) load( name string ) *router.Router {
	r , ok := self.pool.Load( name )
	if ok {
		return r.(*vRouter).L.GetExdata().(*router.Router)
	}

	v := self.Compile( name )
	if v == nil {
		return self.unknown.L.GetExdata().(*router.Router)
	}

	return v.L.GetExdata().(*router.Router)
}

func (self *vRCache) Unknown( name string ) {
	//判断文件是否存在
	v := self.Compile( name )
	if v == nil {
		return
	}
	self.unknown = v
}

func (self *vRCache) update( name string , obj *vRouter ) {
		filename := self.filename( name )
		stat , err := os.Stat( filename )
		if os.IsNotExist( err ) {
			self.pool.Delete( name )
			return
		}

		if obj.modTime == stat.ModTime().Unix() {
			return
		}

		self.Compile( name )
}

func (self *vRCache) sync() {
	tk := time.NewTicker( 1000 * time.Millisecond )
	for range tk.C {
		self.pool.Range(func(k interface{} , v interface{} ) bool {
			name := k.(string)
			obj := v.(*vRouter)
			self.update( name , obj)
			return false
		})

		if self.unknown == nil {
			continue
		}

		stat , err := os.Stat( self.filename( self.unknown.name ))
		if os.IsNotExist( err ) {
			pub.Out.Err("default not found")
			continue
		}

		if stat.ModTime().Unix() == self.unknown.modTime {
			continue
		}

		self.Unknown( self.unknown.name )
	}
}
