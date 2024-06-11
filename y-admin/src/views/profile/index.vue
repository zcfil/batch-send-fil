<template>
  <!-- <div slot="header" class="clearfix">
    <el-card class="box-card">
        <div class="head"></div>
    </el-card>
    <el-button style="float: right; padding: 3px 0" type="text">操作按钮</el-button>
  </div> -->
  <div>

<el-tabs type="border-card" class="right-card">
  
  <!-- <el-tab-pane label="Wallet List">
    <wallets title="传递的参数"  :dysncParams="dysncParamsText">
    </wallets>
  </el-tab-pane> -->
  <el-tab-pane label="系统设置">
    <setting title="传递的参数"  :dysncParams="dysncParamsText">
    </setting>
  </el-tab-pane>
  <el-tab-pane label="修改登录密码"> 
    <password title="传递的参数"  :dysncParams="dysncParamsText">
    </password>
  </el-tab-pane>
  <!-- <el-tab-pane label="ChartCard">
    <chart-card title="传递的参数"  :dysncParams="dysncParamsText">
    </chart-card>
  </el-tab-pane> -->
</el-tabs>

  </div>
  
</template>

<script>
import ChartCard from '@/components/ChartCard'
import Wallets from '@/components/Wallets/wallets'
import Setting from '@/components/User/setting'
import Password from '@/components/User/password'
import { memberinfo } from '@/api/user'
//import Setting from '@/components/User/setting.vue'
export default {
  components:{
    ChartCard,
    Wallets,
    Setting,
    Password,
    Setting,
  },
  data() {
    return {
      imageurl: './2.jpg',
      dysncParamsText:'动态参数',
      form: {
        name: '',
        region: '',
        date1: '',
        date2: '',
        delivery: false,
        type: [],
        resource: '',
        desc: ''
      },
      info:{}
    }
  },
   mounted() {
     this.fetchData()
   },
  methods: {
    async fetchData(){
     try{
      let res = await memberinfo();
      this.info = res.data||{}
     }catch(e){
     }finally{
     }
    },
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
 .text {
    font-size: 14px;
  }

  .item {
    padding: 18px 0;
  }

  .head-left{
    width: 80px;
    height: 80px;
    float: left;

  }
   .head-right{
    width: 240px;
    height: 80px;
    float: right;
  }
  .text {
    font-size: 14px;
  }

  .item {
    margin-bottom: 18px;
  }

  .clearfix:before,
  .clearfix:after {
    display: table;
    content: "";
  }
  .clearfix:after {
    clear: both
  }

  .box-card {
    margin-top: 40px;
    margin-left: 38px;
    width: 780px;
    float: left;
  }
  .right-card{
    width: 70%;
    float: left;
    margin-top: 24px;
    margin-left: 30px;
    height: 600px;
  }
</style>

