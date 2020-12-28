package layer

// service
var serviceList = make(map[interface{}]interface{})

func RegisterService(k interface{}, v interface{}) {
	serviceList[k] = v
}

func ServiceList() map[interface{}]interface{} {
	return serviceList
}

// model
var modelList = make(map[interface{}]interface{})

func RegisterModel(k, v interface{}) {
	modelList[k] = v
}

func ModelList() map[interface{}]interface{} {
	return modelList
}

// model wrapper
var modelWrapperList = make(map[interface{}]interface{})

func RegisterModelWrapper(k, v interface{}) {
	modelWrapperList[k] = v
}

func ModelWrapperList() map[interface{}]interface{} {
	return modelWrapperList
}

// service wrapper
var serviceWrapperList = make(map[interface{}]interface{})

func RegisterServiceWrapper(k, v interface{}) {
	serviceWrapperList[k] = v
}

func ServiceWrapperList() map[interface{}]interface{} {
	return serviceWrapperList
}

// plugins会被注入到所有model层以上的所有对象中
// 存在不能划入model层的功能，希望被很多service或者hook或者handler共享
var pluginsList = make(map[interface{}]interface{})

func RegisterPlugins(k, v interface{}) {
	pluginsList[k] = v
}

func PluginsList() map[interface{}]interface{} {
	return pluginsList
}
