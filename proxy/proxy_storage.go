package proxy

import "strconv"

var ProxyStorage = make(map[string]*ProxyInstance)

func RegisterProxy(proxy *ProxyInstance) bool {
	bindPortString := strconv.Itoa(proxy.BindPort)
	_, found := ProxyStorage[bindPortString]
	if found {
		return false
	} else {
		ProxyStorage[bindPortString] = proxy
		return true
	}
}

func PersistProxyConfig(proxy *ProxyInstance) error {
	//TODO
	return nil
}
