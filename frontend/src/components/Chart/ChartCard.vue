<template>
    <a-card>
        <a-statistic
            :title="title"
            :value="value"
        />
    </a-card>
</template>

<script>
import app from '@/main'
import {Statistic} from 'ant-design-vue'

app.use(Statistic)
export default {
    name: 'ChartCard',
    props: ['title', 'value']
}
</script>

<style scoped>

</style>
