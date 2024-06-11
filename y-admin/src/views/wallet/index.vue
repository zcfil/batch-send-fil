<template>
  <div class="login-container" style="">
    <el-card class="content" shadow="always">
    
     <div style="margin-top: 30px;text-align:center">
       <h1 >助记词生成</h1>
     </div>
     <!-- <el-input
  type="textarea"
  :rows="10"
  placeholder="请输入内容"
  v-model="textarea"
  resize="none"
  readonly="true"
  >
</el-input> -->
    <div style="margin-top: 30px;margin-left:5%;border:1px solid #000;height:140px;width:90%;">
      <el-tag v-for="(item,index) in items"
    :key="index"
    effect="dark" 
    type="info"
    :disable-transitions="true"
    style="margin-top: 3%;margin-left:4%;width:12%;height:34px;text-align:center">
    {{item}}</el-tag>
    
    </div>
    <div style="margin-top: 10px;text-align:center">
      <span style="font-size:13px;color:gray;">请依次抄写并妥善保管，助记词是您恢复钱包的唯一手段。一经丢失，无法找回。请勿使用截屏、拍照等方式保存。</span>
    </div>
    <div style="text-align:center"> 
    <el-button type="primary" style="width:90%;margin-top: 20%;" plain @click="onAnew">重新生成</el-button><br>
    <!-- <el-button type="primary" style="width:90%;margin-top: 3%;"  @click="$router.push({path: '/wallet/index/wallet-1'})">我已备份，下一步</el-button> -->
    <el-button type="primary" style="width:90%;margin-top: 3%;"  @click="$router.push({name: 'walletCreate', params: {num: 1,mnemonic: items}})">我已备份，下一步</el-button>
    </div>

    </el-card>
    
  </div>
</template>

<script>
import AES from '@/common/AES'
import { getMnemonic,exportPrivateKey } from '@/api/wallet'
import router from '@/router'
export default {
  name: 'Login',
  data() {
    
    return {
       items: [
          { type: '', label: '标签一' },
          { type: 'success', label: '标签二' },
          { type: 'info', label: '标签三' },
          { type: 'danger', label: '标签四' },
          { type: 'warning', label: '标签五' },
          { type: 'warning', label: '标签1' },
          { type: 'warning', label: '标签2' },
          { type: 'warning', label: '标签3' },
          { type: 'warning', label: '标签4' },
          { type: 'warning', label: '标签5' },
          { type: 'warning', label: '标签6' },
          { type: 'warning', label: '标签7' },
          
        ]
    }
  },
  created(){
    //getMnemonic()
    this.Mnemonic()
  },
  watch: {
    random: {
        immediate: false,
        handler(newValue,oldValue) {
            if (newValue !== oldValue) {
                // console.log("has listened data!")
                router.go(-1)
            }
        }
    }
  },
  methods: {
    
    Mnemonic(){
      // var param = {
      //   wallet : 'f13djvobwob2les2mkqwttg6vkem2ozvofohef46q',

      // }
      // exportPrivateKey(param).then(response => {
      //    console.log(AES.decrypt(response.data.list))

      // })
      getMnemonic().then(Response =>{
        var res = (AES.decrypt(Response.data)).split(' ')
        this.items = res
      } )
    },
    onAnew(){
      this.Mnemonic()
    }
  }
}
</script>

<style lang="scss">
.content{
  background-color:rgb(245, 245, 247);
  height:700px;
  width:50%;
  margin-left:24%;
  margin-top: 30px;
}
</style>

