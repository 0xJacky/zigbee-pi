import {createApp} from 'vue'
import App from './App.vue'
import {Progress, Row, Col} from 'ant-design-vue'

const app = createApp(App)

app.use(Progress)
app.use(Row)
app.use(Col)

app.mount('#app')

export default app
