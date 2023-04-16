package main

import (
	"github.com/dgraph-io/ristretto"
	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/golang/glog"
	"github.com/klauspost/compress/zstd"
	"net/http"
)

var debug = false
var input = loadContent()
var encoder, _ = zstd.NewWriter(nil)
var decoder, _ = zstd.NewReader(nil, zstd.WithDecoderConcurrency(0))

func main() {
	http.ListenAndServe(":8888", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		Cache()
		writer.Write([]byte("ok"))
	}))
}

var cache, cacheErr = ristretto.NewCache(&ristretto.Config{
	NumCounters: 1e7,     // number of keys to track frequency of (10M).
	MaxCost:     4 << 30, // maximum cost of cache (1GB).
	BufferItems: 64,      // number of keys per Get buffer.
})

func Cache() {
	if cacheErr != nil {
		panic(cacheErr)
	}

	key := gtime.TimestampNanoStr()

	output := Compress(input)
	if debug {
		glog.Info(key)
		glog.Info("input string len ", len(input))
		glog.Info("compress string len ", len(output))
		glog.Info("compress string hash is ", gmd5.MustEncrypt(output))
	}

	cache.Set(key, output, 1)
	cache.Wait()

	value, found := cache.Get(key)
	if !found {
		panic("missing value")
	}

	out := Decompress(value.([]byte))

	if debug {
		glog.Info("cached string len ", len(value.([]byte)))
		glog.Info("decompress string len ", len(out))
		glog.Info("decompress string hash is ", gmd5.MustEncrypt(out))
	}

	if string(input) != string(out) {
		glog.Info("比对失败")
	}
}

func loadContent() []byte {
	in := gfile.GetContents("js.js")
	return []byte(in)
}

func Compress(src []byte) []byte {
	return src
	//return encoder.EncodeAll(src, make([]byte, 0, len(src)))
}

func Decompress(src []byte) []byte {
	return src
	//ret, _ := decoder.DecodeAll(src, nil)
	//return ret
}
