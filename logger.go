package fasthttp

import (
	"bytes"
	pub "github.com/edunx/rock-public-go"
	"github.com/valyala/fasthttp"
	"strconv"
	"strings"
)

func (self *vlogger) New( val string ) {
	self.format = strings.Split( strings.TrimSpace(val) , ",")
	self.encode = self.Json
}

func vlogJsonTranslation(buff *bytes.Buffer , v []byte) {
	len := len(v)
	var ch byte
	for i := 0;i<len;i++ {
		ch = v[i]

		switch ch {
		case '"':
			buff.WriteByte('\\')
			buff.WriteByte('"')
		case '\\':
			buff.WriteByte('\\')
			buff.WriteByte('\\')
		case '\r':
			buff.WriteByte('\\')
			buff.WriteByte('r')
		case '\n':
			buff.WriteByte('\\')
			buff.WriteByte('n')
		case '\t':
			buff.WriteByte('\\')
			buff.WriteByte('t')
		default:
			buff.WriteByte(ch)
		}

	}
}

func vlogJsonStr( buff *bytes.Buffer , v string ) {
	buff.WriteByte('"')
	vlogJsonTranslation(buff , pub.S2B( v ))
	buff.WriteByte('"')
}

func vlogJsonBytes( buff *bytes.Buffer , v []byte) {
	buff.WriteByte('"')
	vlogJsonTranslation(buff , v)
	buff.WriteByte('"')
}

func vlogJsonInt( buff *bytes.Buffer , v int ) {
	buff.WriteString(strconv.Itoa( v ))
}

func vlogJsonKey( buff *bytes.Buffer , key string ) {
	buff.WriteByte('"')
	buff.WriteString(key)
	buff.WriteString("\":")
}


func (self *vlogger) Json( ctx *fasthttp.RequestCtx ) []byte {
	buff := new( bytes.Buffer )

	buff.WriteByte('{')

	l := len(self.format)
	var key string
	for i := 0 ; i< l ; i++ {
		key = self.format[i]
		if i != 0 {
			buff.WriteByte(',')
		}

		switch self.format[i] {
		case "time":
			vlogJsonKey(buff , key)
			vlogJsonStr( buff , ctx.Time().Format("2006-01-02 15:04:05.00"))
		case "remote_addr":
			vlogJsonKey(buff , key)
			vlogJsonStr( buff , ctx.RemoteIP().String())
		case "path":
			vlogJsonKey(buff , key)
			vlogJsonBytes( buff , ctx.Request.URI().Path())
		case "query":
			vlogJsonKey(buff , key)
			vlogJsonBytes( buff , ctx.Request.URI().QueryString())
		case "host":
			vlogJsonKey(buff , key)
			vlogJsonBytes( buff , ctx.Host())
		case "server_addr":
			vlogJsonKey(buff , key)
			vlogJsonStr( buff , ctx.LocalIP().String())
		case "server_port":
			vlogJsonKey(buff , key)
			vlogJsonInt( buff , httpServerPort( ctx ))
		case "body":
			vlogJsonKey(buff , key)
			vlogJsonBytes( buff , ctx.Request.Body() )
		case "status":
			vlogJsonKey(buff , key)
			vlogJsonInt( buff , ctx.Response.StatusCode() )
		case "content-length":
			vlogJsonKey(buff , key)
			vlogJsonInt( buff , ctx.Request.Header.ContentLength() )
		case "referer":
			vlogJsonKey(buff , key)
			vlogJsonBytes(buff , ctx.Referer())
		case "ua":
			vlogJsonKey(buff , key)
			vlogJsonBytes( buff , ctx.Request.Header.UserAgent())
		case "region_cityid":
			vlogJsonKey(buff , key)
			vlogJsonInt(buff , httpRegionCityId( ctx ))
		case "region_info":
			vlogJsonKey(buff , key)
			vlogJsonBytes(buff , httpRegionInfoRaw( ctx ))
		case "sent":
			vlogJsonKey(buff , key)
			vlogJsonInt(buff , ctx.Response.Header.ContentLength())
		case "headers":
			vlogJsonKey(buff , key)
			buff.WriteByte('{')
			hi := 0

			ctx.Request.Header.VisitAll(func(key []byte , value []byte){
				if hi != 0 {
					buff.WriteByte(',')
				}
				vlogJsonKey( buff , pub.B2S( key ) )
				vlogJsonBytes(buff , value )
				hi++
			})
			buff.WriteString("}")

		default:
			size := len(key)
			if size > 4 && key[:4] == "arg_" {
				vlogJsonKey(buff , key[4:])
				vlogJsonBytes(buff , ctx.QueryArgs().Peek( key[4:] ))
				continue
			}

			if size > 5 && key[:5] == "post_" {
				vlogJsonKey(buff , key[5:])
				vlogJsonBytes(buff , ctx.Request.PostArgs().Peek( key[5:]))
				continue
			}

			if size > 5 && key[:5] == "http_" {
				vlogJsonKey(buff , key[5:])
				vlogJsonBytes(buff , ctx.Request.Header.Peek(key[5:]))
				continue
			}

			if size > 7 && key[:7] == "cookie_" {
				vlogJsonKey(buff , key[7:])
				vlogJsonBytes(buff , ctx.Request.Header.Cookie( key[7:]) )
				continue
			}

			if size > 6 && key[:6] == "param_" {
				obj := ctx.UserValue( key[6:])
				if val , ok := obj.(string); ok {
					vlogJsonKey(buff , key[6:])
					vlogJsonStr(buff , val)
					continue
				}
				vlogJsonStr(buff , "nil")
				continue
			}

			vlogJsonKey(buff , key)
			vlogJsonStr(buff , "nil")
		}
	}

	buff.WriteByte('}')
	return buff.Bytes()
}

func (self *vlogger) Raw(ctx *fasthttp.RequestCtx) []byte {
	return nil
}