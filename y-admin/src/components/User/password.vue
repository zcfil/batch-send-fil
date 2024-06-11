<template>
  <div class="app-container">
    <el-form ref="param" :rules="rules" :model="param" label-width="140px">
      <el-form-item label="原密码" prop="oldPwd">
        <el-input type="password" v-model="param.oldPwd"  placeholder="请输入原密码" />
      </el-form-item>
      <el-form-item label="新密码" prop="newPwd" >
        <el-input type="password" v-model="param.newPwd" placeholder="请输入新密码" />
      </el-form-item>
       <el-form-item label="确认新密码" prop="Confirm" >
        <el-input type="password" v-model="param.Confirm" placeholder="确认新密码" />
      </el-form-item>
    
      <el-form-item>
        <el-button type="primary" @click="onSubmit" :disabled="false">提交</el-button>
        <el-button @click="onCancel">重置</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script>
import { resetpwd } from '@/api/user'
export default {
  name: "Password",
  data() {
    var validatecharge = (rule, value, callback) => {
        console.log(value)
        if (value.length<6){
          this.btn = true
          callback(new Error('密码至少为6位'));
        }else if ( value== "" ){
          this.btn = true
          callback(new Error('密码不能为空'));
        }else if (this.param.newPwd!==this.param.Confirm){
          this.btn = true
          callback(new Error('新密码两次输入不一致！'));
        }else{
           this.btn = false
        }
      };
      var validatecharge1 = (rule, value, callback) => {
        console.log(value)
        if (value.length<6){
          this.btn = true
          callback(new Error('密码至少为6位'));
        }else if ( value== "" ){
          this.btn = true
          callback(new Error('密码不能为空'));
        }
      };
    return {
      btn: true,
      param:{
        oldPwd:'',
        newPwd:'',
        Confirm:'',
      },
      rules: {
          newPwd: [
            { validator: validatecharge1, trigger: 'blur' }
          ],
          Confirm: [
            { validator: validatecharge, trigger: 'blur' }
          ],
        },
    }
  },
  methods: {
    onSubmit() {
       if (this.param.newPwd.length<6){
          this.$message({
          message: '密码至少为6位!',
          type: 'error'
          
          })
          return
        }else if ( this.param.newPwd== "" ){
          this.$message({
          message: '密码不能为空!',
          type: 'error'
          
          })
          return
        }else if (this.param.newPwd!==this.param.Confirm){
          this.$message({
          message: '新密码两次输入不一致!',
          type: 'error'
          
       })
          return
        }
      resetpwd(this.param).then(res=>{
        
        if (res.msg==="success"){
          this.$message({
          message: '修改成功！!',
          type: 'success'
          
       })
       this.onCancel()
        }
        
      })
    },
    onCancel() {
      // this.$message({
      //   message: 'cancel!',
      //   type: 'warning'
      // })
      this.param.oldPwd = ""
      this.param.newPwd = ""
      this.param.Confirm = ""
    }
  }
}
</script>

<style scoped>
.line{
  text-align: center;
}
</style>

