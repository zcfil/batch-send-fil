<template>
  <div class="login-container">
      <div class="title-container">
        <h2 class="title">绑定 MFA</h2>
      </div>
<span class="p1">安全令牌验证 请按照以下步骤完成绑定操作 </span>

<div class="link-top"></div>

<div style="margin-top: 20px;"><span class="p2">使用手机扫描以下二维码下载应用 </span>

<div style="margin-top: 20px;margin-left: 36%;"><canvas id="QRCode_header" ref="QRCode_header"></canvas>

<canvas id="QRCode_header1" ref="QRCode_header1" style=";margin-left: 12%;"></canvas></div>
<!-- <div style="margin-top: 20px;margin-right: 80%;"></div> -->



<div style="margin-top: 10px;"> <span class="p3">Android手机下载（阿里云） </span><span class="p4">iPhone手机下载（谷歌） </span></div>

<button style="margin-top: 80px;margin-left: 46%;width:9%;height:30px;background-color: rgb(60, 140, 200);font-weight: bold;" @click="handleVerify">下一步</button>
</div>

  </div>
</template>

<script>
// import { validUsername } from '@/utils/validate'
// import { register } from '@/api/user'
import QRCode from "qrcode"

export default {
  name: 'dow',
  data() {
    return {
      QRCodeiP:"https://apps.apple.com/cn/app/id388497605",
      QRCodeAn:"https://hd.m.aliyun.com/act/download.html"
    }
  },
  methods: {
      optsfun:function(){
          let opts = {
             errorCorrectionLevel: "H",//容错级别
             type: "image/png",//生成的二维码类型
             //quality: 0.3,//二维码质量
             //margin: 12,//二维码留白边距
             width: 200,//宽
             height: 180,//高
            //  text: "http://www.xxx.com",//二维码内容
             color: {
                 dark: "#333333",//前景色
                 light: "#fff"//背景色
             }
         };
        //  this.QRCodeMsg = "https://hd.m.aliyun.com/act/download.html"; //生成的二维码为URL地址js
         let msg =  this.$refs.QRCode_header    //
         let msg1 = this.$refs.QRCode_header1
         // 将获取到的数据（val）画到msg（canvas）上
         QRCode.toCanvas(
           msg, this.QRCodeAn, opts, 
         );
         QRCode.toCanvas(
           msg1, this.QRCodeiP, opts, 
         );
      },
      handleVerify() {
          this.$router.push({ path: '/verify' })
      }
  
  },
  created(){
         
  },
  mounted(){
    this.optsfun()
  }
   
}
</script>

<style lang="scss">
/* 修复input 背景不协调 和光标变色 */
/* Detail see https://github.com/PanJiaChen/vue-element-admin/pull/927 */

$bg:#a8b1bb;
$light_gray:#fff;
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
      margin-left: 41%;
      text-align: center;
      font-weight: bold;
    }
    .p2 {
      font-size: 18px;
      color: rgb(20, 20, 20);
      margin-left: 42%;
      text-align: center;
      font-weight: bold;
      
    }
  .p3 {
      font-size: 16px;
      color: rgb(78, 68, 68);
      margin-left: 36%;
      text-align: center;
     // font-weight: bold;
      margin-top: 20px;
      
    }
    .p4 {
      font-size: 16px;
      color: rgb(78, 68, 68);
      margin-left: 8%;
      text-align: center;
      //font-weight: bold;
      margin-top: 20px;
      
    }
}
</style>
