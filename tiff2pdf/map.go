package tiff2pdf

import "sync"

type MapWrapper struct {
	wrappedMap sync.Map
}

// Load wraps sync.Maps Load func in a type safe format
func (m *MapWrapper) Load(key int) (value *fd, ok bool) {
	loaded, ok := m.wrappedMap.Load(key)
	if ok {
		return loaded.(*fd), ok
	}
	return nil, ok
}

// Store wraps sync.Maps Store func in a type safe format
func (m *MapWrapper) Store(key int, value *fd) {
	m.wrappedMap.Store(key, value)
}

// Delete wraps sync.Maps Delete func in a type safe format
func (m *MapWrapper) Delete(key int) {
	m.wrappedMap.Delete(key)
}

//***************************************************************
// Initial MapWrapper implementation used to migrate from
// regular map to sync.Map by proviidng same interface
//***************************************************************

// type MapWrapper struct {
// 	wrappedMap map[int]*fd
// }

// // func (m *MapWrapper) init() {
// // 	if m.wrappedMap == nil {
// // 		m.wrappedMap = make(map[int]*fd)
// // 	}
// // }

// func (m *MapWrapper) Load(key int) (value *fd, ok bool) {
// 	// m.init()
// 	val, ok := m.wrappedMap[key]
// 	return val, ok
// }

// func (m *MapWrapper) Store(key int, value *fd) {
// 	// m.init()
// 	m.wrappedMap[key] = value
// }

// func (m *MapWrapper) Delete(key int) {
// 	// m.init()
// 	// m.wrappedMap[key] = nil
// 	delete(m.wrappedMap, key)
// }

//***************************************************************
