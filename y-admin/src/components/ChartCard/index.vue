<template>
  <div>
    这里是内容{{title}}-----{{dysncParams}}
  </div>
</template>

<script>
  export default {
    name: "ChartCard",
    props: {
      title: {//参数
        type: String,
        default: '1'
      },
      dysncParams: {
         type: String,
        default: '0000'
      }
  
    }
  }
</script>

<style  scoped>
/*css */
</style>