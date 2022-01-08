# zigbee-pi
一个物联网环境监测平台的 prototype，只是为了水一水期末论文。

利用 ZigBee 节点及传感器采集环境的温湿度，经过 ZigBee 协调器汇总信息后通过串口传送给树莓派，再通过树莓派使用 WebSocket 协议将数据发送到云服务器上，实现在公网下的环境温湿度可视化。云服务器端运行基于 Go 开发的服务端，利用 WebSocket 技术将数据发送给基于Vue3开发的前端。

## 系统总体架构
![image](https://user-images.githubusercontent.com/13096985/148629177-e792cbfe-ab5c-4cc2-9e7b-b370cbede351.png)

在本原型中，传感器部分采用 DHT11 温湿度传感器，ZigBee 节点使用基于 TI 公司研发的 CC2530 芯片并整合了 USB 转串口芯片芯片（CH340）的开发板，树莓派使用 Raspberry Pi 3 Model B 型号。云服务器预装 Ubuntu 20.04 操作系统，使用 Go 1.17、Gin 框架开发服务端，前端部分使用 Vue3 进行开发。

树莓派与云服务器的通信完全采用 WebSocket 协议，浏览器与云服务器的通信，除下载静态资源时使用 HTTP/2，获取监控数据时云采用 WebSocket 进行通信。
采用 WebSocket 作为主要的通信协议是因为它是一种长连接，可以在单个 TCP 上建立全双工的通信，具有较低开销，便于实现客户端与服务器进行实时的数据传输。在连接创建完成后，服务器与客户端之间交换数据时，用于协议控制的数据包头部相对较小。由于协议是全双工的，所以服务器可以随时主动给客户端下发数据。相比于 HTTP 请求需要等待客户端请求服务端才能响应，延迟明显减少。

为什么数据上报要用 WebSocket 呢？主要是因为我懒，其实这里应该用 UDP 的。

## 软件部分设计
### ZigBee终端及协调器的软件设计

ZigBee 终端需要实现对 DHT11 传感器数据的读取，DHT11 通过 CC2530 的 P0_7 接口进行数据传送。读取 DHT11 传送过来的温度高8位和湿度高8位，显示在 LCD 显示器上，然后通过 ZigBee 发送给协调器，协调器接收到数据后通过串口打印出接收到的数据，用于树莓派的读取。
  
### 树莓派上的软件设计
树莓派上使用 Python 语言进行开发，通过 pyserial 库获取由 ZigBee 协调器节点发送过来的串口消息。为了避免长时间轮询对后端服务器的压力，树莓派使用 Python 的 websocket-client库，通过 WebSocket 与后端服务器连接。
配合循环和 try…catch… 语句实现 WebSocket的自动重连，一旦程序启动，则会立即与后端接口进行握手，若因为网络问题或服务器调试导致链接断开，则每10秒尝试一次重连，直到再次握手成功。
当握手成功后，程序将每隔2秒读取一次串口缓冲区内的数据，转换成 json 后发送给后端接口。
  
### 服务端软件设计

服务端采用Go语言进行开发，得益于Go语言的特性，非常适合用于开发 WebSocket 服务端。在这个部分中，笔者使用了 Gin 作为 HTTP 框架，Gorilla/WebSocket作为WebSocket 框架，分别设计了两个接口，Client接口用于广播监控数据，Pi接口用于接收从树莓派上传的监控数据，数据的流动逻辑如图所示。
  
  ![image](https://user-images.githubusercontent.com/13096985/148629352-234d8e2c-fd9e-40b8-ad76-c6ab41ad694d.png)

创建一个 interface{} 类型的管道用于 Goroutine 与 Pi 接口的通信。
前端访问Client 接口的地址，等待客户端与服务器成功握手后，客户端即加入广播字典中。
创建一个 Goroutine，通过 select 语句阻塞，等待管道中出现消息后，向所有存在于广播字典的客户端套接字(socket)推送监控数据，即可实现广播效果，若推送给某个终端的过程中发生错误，则认为该终端已离线，需要将其移出广播字典。
 
树莓派启动监控后，通过 Pi 接口与服务器进行 WebSocket 握手，握手后需要发送一组 UUID 作为 Token 确保监控端口不会被恶意访问，该Token为服务端在配置文件中设置的 TrustedToken 的值。树莓派通过认证后，每隔2秒将监控数据发送给服务端。服务端收到监控数据后更新缓存内的监控数据并向管道发送一个int，通知Goroutine进行广播操作。
需要注意的是，缓冲区为一个结构体，定义如下：
```
var buffer struct {
	Temperature string `json:"temperature"`
	Humidity    string `json:"humidity"`
	mux         sync.Mutex
}
```
Pi 接口允许多个设备同时连接，为了避免出现多个 Goroutine 同时写入缓冲区，需要使用 `sync.Mutex` 作为互斥锁，在读写 buffer 前需要用 `buffer.mux.Lock()` 申请锁，访问结束后需要用 `buffer.mux.UnLock()` 释放锁。
### 前端软件设计
前端基于 Vue3 进行开发，使用 Ant Design Vue 作为UI组件库，为了 WebSocket 能自动重连，笔者使用`ReconnectingWebSocket`库进行数据传输。
 
访问前端首页后，将立即与后端 `/monitor` 接口握手，设置接收到消息的回调函数 `wsOnMessage(m)`，每当树莓派推送一组数据到服务器，服务器就会广播数据给客户端，客户端接收到的数据为 json，包含温度和湿度的数据。之后，前端将数据更新到视图上，由此可以实现监测数据的动态更新。
当前端与后端 `/monitor` 接口断开连接后，将自动尝试重连直至握手成功。

## 系统测试
启动服务器上的服务端软件，将ZigBee 协调器与树莓派连接如图所示。

![image](https://user-images.githubusercontent.com/13096985/148629444-75cd7881-6cc0-43cc-a9bc-e2c047aa929a.png)

![image](https://user-images.githubusercontent.com/13096985/148629445-6996ac01-eed8-43c6-a729-8eb790a035d6.png)


启动 ZigBee 终端节点与协调器节点，等待终端与协调器配对成功。ZigBee终端上的显示器会打印出数据。之后，观察服务端程序的日志输出.

![image](https://user-images.githubusercontent.com/13096985/148629504-21fcc7f0-e8cb-4255-9df7-902b8b59e199.png)

在树莓派上运行 Python 脚本，分别代表环境的温度和湿度。

![image](https://user-images.githubusercontent.com/13096985/148629511-ed32fc43-5068-43cc-8535-c641cf1f74b6.png)

启动浏览器访问前端页面，即刻获取监控数据，浏览器上的访问效果如图所示。

![image](https://user-images.githubusercontent.com/13096985/148629392-2c3779eb-2be5-4723-a871-015f42c6484e.png)

![image](https://user-images.githubusercontent.com/13096985/148629393-1ba75594-ff13-4bcd-8b28-d3f6bf3ecd00.png)

### 参考文献
[1] 王凯巍,孙康,张子伊,陈美娟,朱晓荣.基于ZigBee与树莓派的环境信息采集系统[J].实验科学与技术,2019,17(06):132-136.

[2] Nikhade S G. Wireless sensor network system using Raspberry Pi and zigbee for environmental monitoring applications[C]//2015 International Conference on Smart Technologies and Management for Computing, Communication, Controls, Energy and Materials (ICSTM). IEEE, 2015: 376-381.

[3] CC2530 datasheet (Rev. B)[DB/OL].
https://www.ti.com/lit/ds/symlink/cc2530.pdf. 2021-12-19 

[4] Temperature-Humidity-Sensor-Schematic[DB/OL].
https://www.waveshare.net/w/upload/1/14/Temperature-Humidity-Sensor-Schematic.pdf.2021-12-19

[5] ZIGBEE开发套件WiFi模块全部资料[DB/OL].
https://pan.baidu.com/s/1Yk7IaLiKJsjoMc9HT9Xlmw?pwd=j87f. 2021-12-19

[6] 数字温湿度传感器[DB/OL].
https://cdn-shop.adafruit.com/datasheets/DHT11-chinese.pdf. 2021-12-21


