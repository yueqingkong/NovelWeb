package net

import "log"

type ProxyPool struct {
}

func NewProxyPool() ProxyPool {
	return ProxyPool{}
}

func (pool ProxyPool) Load() {
	log.Print("[代理池]初始化...")


}
