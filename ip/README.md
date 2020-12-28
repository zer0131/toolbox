
用户配置的white_list中的元素有两种形态
1. ip段: 10.189.240.0-10.189.240.255 
2. 具体ip: 10.189.241.35 

代码示例：

// 这个操作需要放在hook.go中，传入的参数在app.ini中配置，例如：

// white_list[string_array] = 10.189.240.0-10.189.240.255,10.189.241.0-10.189.241.255

ip.Init(context.Background(), []string{"10.188.0.24", "10.188.1.25-10.188.255.255"})

ip.CheckIp(context.Background(), "10.188.0.24")