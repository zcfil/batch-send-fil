<template>
  <div class="login-container">
      <div class="title-container">
        <h2 class="title">绑定 MFA</h2>
      </div>
<span class="p1">安全令牌验证 请按照以下步骤完成绑定操作 </span>

<div class="link-top"></div>

<div style="margin-top: 20px;"><span class="p2">使用手机 Google Authenticator 应用扫描以下二维码，获取6位验证码 </span>

<div style="margin-top: 20px;margin-left: 42%;"><canvas id="QRCode_header" ref="QRCode_header"></canvas></div>
<div style="margin-top: 20px;"> <span class="p3">或手动添加密钥：{{verify.googleSecret}} </span><br>
</div>
<div style="margin-top: 20px;margin-left: 44%;"> <span style="color:red;">验证码：<input v-model="form.code" style="width:120px;height:30px;"> </span></div>
<button style="margin-top: 80px;margin-left: 45%;width:9%;height:30px;background-color: rgb(60, 140, 200);font-weight: bold;" @click="handleVerify">下一步</button>



</div>

  </div>
</template>

<script>
// import { validUsername } from '@/utils/validate'
import { getVerificationCode,bindVerificationCode,verifyCode } from '@/api/user'
import QRCode from "qrcode"

export default {
  name: 'dow',
  data() {
    return {
      verify :'',
      form:{
          code:'',
          googleSecret:'',
      },
      QRCodeMsg:'',
    }
  },
  methods: {
      optsfun:function(){
          let opts = {
             errorCorrectionLevel: "H",//容错级别
             type: "image/png",//生成的二维码类型
             //quality: 0.3,//二维码质量
             //margin: 12,//二维码留白边距
             width: 300,//宽
             height: 280,//高
             //text: "otpauth://totp/BET:lzj1?secret=7NXKXWLCVYNFRZ6Y",//二维码内容
             color: {
                 dark: "#333333",//前景色
                 light: "#fff"//背景色
             }
         };
        //this.QRCodeMsg = "otpauth://totp/BET:lzj1?secret=7NXKXWLCVYNFRZ6Y"; //生成的二维码为URL地址js
        let msg = document.getElementById("QRCode_header"); //this.$refs.QRCode_header    //
         // 将获取到的数据（val）画到msg（canvas）上
        QRCode.toCanvas(
           msg, this.QRCodeMsg, opts, 
         );
      },
      handleVerify() {
          bindVerificationCode(this.form).then(response => {
          if (response.code==0){
              this.$router.push({ path: this.redirect || '/' }) 
              this.$message('绑定成功!')
          }
        })
      },
      handleBindVerificationCode() {
          this.$router.push({ path: '/download' })
      },
      handleVerifyCode() {
          this.$router.push({ path: '/download' })
      }
  
  },
  created(){
      getVerificationCode().then(response => {
        this.verify = response.data
        this.QRCodeMsg = response.data.url
        this.form.googleSecret = response.data.googleSecret
        console.log(this.QRCodeMsg)
      })
  },
  mounted(){
    self.setTimeout(()=>{
    this.optsfun()
    },100)
  }
   
}
</script>

<style lang="scss">
/* 修复input 背景不协调 和光标变色 */
/* Detail see https://github.com/PanJiaChen/vue-element-admin/pull/927 */

$bg:#a8b1bb;
$light_gray:#fff;
$cursor: #fff;


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

  .tips {
    font-size: 14px;
    color: #fff;
    margin-bottom: 10px;

    span {
      &:first-of-type {
        margin-right: 16px;
      }
    }
  }
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
      margin: 80px auto 40px auto;
      text-align: center;
      font-weight: bold;
    }
  
  }
  .p1 {
      font-size: 16px;
      color: rgb(112, 110, 110);
      margin-left: 42%;
      text-align: center;
      font-weight: bold;
    }
    .p2 {
      font-size: 18px;
      color: rgb(20, 20, 20);
      margin-left: 34%;
      text-align: center;
      font-weight: bold;
      
    }
  .p3 {
      font-size: 16px;
      color: rgb(78, 68, 68);
      margin-left: 42%;
      text-align: center;
      font-weight: bold;
      margin-top: 20px;
      
    }
}
</style>
