<template>
  <div class="app-container">
    <el-form ref="form" :model="form" label-width="120px">
      <el-form-item label="Activity name">
        <el-input v-model="form.name"  placeholder="lzj" />
      </el-form-item>
      <el-form-item label="Activity email" >
        <el-input v-model="form.name" placeholder="469230190@qq.com" />
      </el-form-item>
       <el-form-item label="Activity phone" >
        <el-input v-model="form.name" placeholder="19142755309" />
      </el-form-item>
      <el-form-item label="Activity zone">
        <el-select v-model="form.region" placeholder="深圳">
          <el-option label="北京" value="shanghai" />
          <el-option label="上海" value="beijing" />
        </el-select>
      </el-form-item>
    
      <el-form-item label="Activity form">
        <el-input v-model="form.desc" type="textarea" placeholder="备注" />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="onSubmit">Submit</el-button>
        <el-button @click="onCancel">Cancel</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script>
export default {
  name: "User",
  data() {
    return {
      form: {
        name: '',
        region: '',
        date1: '',
        date2: '',
        delivery: false,
        type: [],
        resource: '',
        desc: ''
      }
    }
  },
  methods: {
    onSubmit() {
      this.$message('submit!')
    },
    onCancel() {
      this.$message({
        message: 'cancel!',
        type: 'warning'
      })
    }
  }
}
</script>

<style scoped>
.line{
  text-align: center;
}
</style>

