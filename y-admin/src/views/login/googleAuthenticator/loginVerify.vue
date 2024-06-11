<template>
<div class="login-container">
    <div class="title-container">
        <h2 class="title">验证 MFA</h2>
    </div>
    <div class="link-top"></div>
    <div style="margin-top: 20px;margin-left: 40%;"><span>账号保护已开启，请根据提示完成以下操作</span></div>
    <div class="verify">
        <img src="./3.png" style="margin-top: 30px;margin-left: 28%;">
        <br>
        <br>
        <span style="margin-left: 66px;font-size: 15px;color:rgb(46, 44, 44);">请打开手机Google Authenticator应用，输入6位动态码</span>
        <div style="margin-top: 40px;margin-left: 28%;"> <input v-model="form.code" style="width:216px;height:30px;" placeholder="6位数字"></div>
        <button style="margin-top: 20px;margin-left: 28%;width:216px;height:30px;background-color: rgb(60, 140, 200);font-weight: bold;" @click="handleVerify">下一步</button>
    </div>

</div>

</template>

<script>
import { verifyCode } from '@/api/user'

export default {

  data() {
    return {
        form:{
          code:'',
      },
    }
  },
  methods: {
      handleVerify() {
        verifyCode(this.form).then(response => {
          if (response.code==0){
              this.$router.push({ path: this.redirect || '/' }) 
              this.$message('验证成功!')
          }
        })
      }
  
  },
  created(){
         
  },
  mounted(){
    
  }
   
}
</script>

<style lang="scss">
/* 修复input 背景不协调 和光标变色 */
/* Detail see https://github.com/PanJiaChen/vue-element-admin/pull/927 */

$bg:#a8b1bb;
$light_gray:rgb(46, 44, 44);
$cursor: rgb(37, 98, 177);

@supports (-webkit-mask: none) and (not (cater-color: $cursor)) {
  .login-container .el-input input {
    color: $cursor;
  }
}

/* reset element-ui css */
</style>

<style lang="scss" scoped>
$bg:#f2f2f3;
$light_gray:rgb(24, 1, 1);

.login-container {
  min-height: 100%;
  width: 100%;
  background-color: $bg;
  overflow: hidden;

.link-top {
            width: 40%;
            height: 1px;
            border-top: solid #ACC0D8 1px;
            margin-top: 20px;
            margin-left: 30%;
        }


  .title-container {
    position: relative;

    .title {
      font-size: 30px;
      color: $light_gray;
      margin: 30px auto 40px auto;
      text-align: center;
      font-weight: bold;
    }
  }
   .verify {
      width: 30%;
      margin-left: 35%;
      margin-top: 30px;
      background-color: white;
      height: 600px;
    }
}
</style>
