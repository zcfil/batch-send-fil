<template>
  <div class="login-container">
    <el-card class="content" shadow="always">
    
     <div style="margin-top: 30px;">
       <h1 style="text-align:center">验证助记词</h1>
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
      <el-tag v-for="item in items"
    :key="item"
    closable
    style="margin-top: 3%;margin-left:4%;width:12%;height:34px;text-align:center;"
    :disable-transitions="true"
    @close="handleClose(item)">
    {{item}}</el-tag>
    

    </div>
    <div style="margin-top: 6px;">
      <span style="font-size:14px;color:gray;margin-left:5%;">请点击下方的助记词，并确保顺序与您记忆的助记词一致。</span>
    </div>
    <div style="margin-top: 10px;margin-left:4%;height:180px">
       <el-tag v-for="item in choices"
    :key="item"
    effect="dark" 
    type="info"
    style="margin-top: 3%;margin-left:3%;width:12%;height:34px;text-align:center;"
    @click="handleChoice(item)"
    :disable-transitions="true">
    {{item}}</el-tag>
    </div>
    
    <!-- <el-button type="primary" style="width:90%;margin-top: 20%;margin-left: 5%;" plain>重新生成</el-button><br> -->
    <el-button type="primary" style="width:90%;margin-bottom: 3%;margin-left: 5%;"  @click="subClick" >下一步</el-button>


    </el-card>
    
  </div>
</template>

<script>
import AES from '@/common/AES'
import { newWallet } from '@/api/wallet'
export default {
  name: 'Login',
  data() {
    
    return {
       random : Math.random()*1000*1000*1000*1000*1000*1000 +"",
       items: [],
      choices: []
    }
  },

  watch: {
    // random: {
    //     immediate: false,
    //     handler(newValue,oldValue) { 

    //       console.log(newValue+"---------------"+oldValue)
    //         if (newValue !== oldValue) {
    //              console.log("has listened data!")
    //              this.$router.go(-1)
    //         }
    //     }
    // }
  },
  created(){
    // this.random=Math.random()*1000*1000*1000*1000*1000*1000 +""
    //   console.log(this.random)
    // this.choices = this.$route.params.mnemonic
    this.getMnemonic()
  },
  // destroyed(){
  //   this.$router.push({path: '/wallet/index/index'})
  // },
  mounted(){
      if (this.$route.params.num!==1){
        this.$router.replace({path: '/wallet/index/index'})
      }
  }
  ,
  methods: {
    handleClose(tag) {
        this.items.splice(this.items.indexOf(tag), 1);      //删除一个元素
        this.choices.push(tag)
      },
    handleChoice(tag) {
        this.choices.splice(this.choices.indexOf(tag), 1);
        this.items.push(tag)
      },
    RandomNum(Min, Max) {
      var Range = Max - Min;
      var Rand = Math.random();
      var num = Min + Math.floor(Rand * Range); //舍去
      return num;
    },
    //打乱顺序 
    getMnemonic(){
      console.log(this.$route.params.mnemonic)
      var mnemonic = [...this.$route.params.mnemonic]
      var max = 12
      var min  = 0
      for (var i=0;i<12;i++){
        var index = this.RandomNum(min,max)
        this.choices.push(mnemonic[index])
        mnemonic.splice(mnemonic.indexOf(mnemonic[index]), 1);
        max--
      }
    },
    subClick(){
      if (this.$route.params.mnemonic.toString()!==this.items.toString()){
        this.$message({
          message:'助记词验证失败，请检查您的助记词填写是否正确',
          type: 'error',
        })
        return
      }
     
      var param = {
        mnemonic : AES.encrypt(this.$route.params.mnemonic.join(' '))
      }
    
      newWallet(param).then(response => {
        this.$message({
          message:'创建成功：'+response.data,
          type: 'success',
          showClose: true,
          duration:1000*30,
        })
        this.$router.push({path: '/wallet/backups/index'})
      })

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
// div{text-align:center} 
</style>

