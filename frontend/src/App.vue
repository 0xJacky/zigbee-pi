<template>
    <div class="container">
        <h1>宿舍</h1>
        <a-row type="flex" justify="center" :gutter="20">
            <a-col>
                <a-progress
                    :format="p=>p+'°C'"
                    :percent="temperature" type="dashboard"/>
                <p>温度</p>
            </a-col>
            <a-col>
                <a-progress :percent="humidity" type="dashboard"/>
                <p>湿度</p>
            </a-col>
        </a-row>
        <a-row type="flex" justify="center">
            <a-col>
                <p>© {{ year }} 0xJacky</p>
            </a-col>
        </a-row>
    </div>
</template>

<script>
import ReconnectingWebSocket from 'reconnecting-websocket'

export default {
    name: 'App',
    components: {},
    data() {
        return {
            humidity: 0,
            temperature: 0,
            year: new Date().getFullYear()
        }
    },
    created() {
        this.websocket = new ReconnectingWebSocket('wss://homework.jackyu.cn/zigbee-pi/api/monitor')
        this.websocket.onmessage = this.wsOnMessage
    },
    methods: {
        wsOnMessage(m) {
            const r = JSON.parse(m.data)
            this.humidity = r.humidity
            this.temperature = r.temperature
        }
    }
}
</script>

<style lang="less">
@media (prefers-color-scheme: dark) {
    @import '~ant-design-vue/dist/antd.dark.less';
}
#app {
    font-family: Avenir, Helvetica, Arial, sans-serif;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
    text-align: center;
    height: 100%;
    .container {
        margin-top: 120px;
    }
}
</style>
