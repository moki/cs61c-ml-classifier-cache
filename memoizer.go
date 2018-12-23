package memoizer

import (
	"hash/crc64"

	proj5 "github.com/lambda7xx/sp18-proj5"
)

// Memoizer is simplest possible implementation that does anything interesting.
// This doesn't even do memoization, it just proxies requests between the client
// and the classifier. You will need to improve this to use the cache effectively.
func Memoizer(memHandle proj5.MnistHandle, classHandle proj5.MnistHandle, cacheHandle proj5.CacheHandle) {
	for req := range memHandle.ReqQ {
		crc64Table := crc64.MakeTable(crc64.ECMA)
		hash := crc64.Checksum(req.Val, crc64Table)

		cacheHandle.ReqQ <- proj5.CacheReq{
			Write: false,
			Id:    req.Id,
			Key:   hash,
		}

		if res := <-cacheHandle.RespQ; !res.Exists {
			classHandle.ReqQ <- req
			resp := <-classHandle.RespQ

			cacheHandle.ReqQ <- proj5.CacheReq{
				Write: true,
				Key:   hash,
				Val:   resp.Val,
				Id:    req.Id,
			}

			memHandle.RespQ <- resp
		} else {
			memHandle.RespQ <- proj5.MnistResp{
				Val: res.Val,
				Id:  res.Id,
			}
		}
	}
}
