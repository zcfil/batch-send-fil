<template>
  <div class="app-container">
    <el-form :model="configForm" v-loading="Loading" :rules="rules" ref="configForm" label-width="100px" class="demo-ruleForm">
  <el-form-item label="手续费设置：" prop="charge">
    <el-input  v-model="configForm.charge" 
 maxlength='8' oninput="if(isNaN(value)) { value = null } if(value.indexOf('.')>0){value=value.slice(0,value.indexOf('.')+8)}" 
 style="width: 100px;" placeholder="手续费百分比" :disabled="chargeBool" >

 </el-input> %
  <el-button icon="el-icon-edit" style="margin-left:50px;" size="mini" @click="editCharge('ruleForm')">{{cedit}}</el-button>
  </el-form-item>
  <el-form-item style="margin-top:30%;">
    <el-button type="primary" @click="onSubmit()">提交</el-button>
    <el-button @click="resetForm()">重置</el-button>
  </el-form-item>
</el-form>
  </div>
</template>

<script>
import { GetConfig,SetConfig} from '@/api/finance'
import { exportMnemonic,exportPrivateKey} from '@/api/wallet'
import AES from '@/common/AES'
export default {
  name: "Setting",
  created() {
    GetConfig().then(response => {
      this.Loading = false
      //this.configForm = response.data
      if (String(response.data.charge)==="undefined"){
        this.configForm.charge = 0
      }else{
        this.configForm.charge = Number(response.data.charge) * 100
      }
      //this.configForm1 = response.data
      this.configForm1.charge = this.configForm.charge 
    }),
    this.Mnemonic()

  },
  data() {
      var validatecharge = (rule, value, callback) => {
        // var n = value.indexOf("%")
        // var num = 1
        // if (n!=-1){
        //   num = value.substring(0, n)
        //   num = Number(num)/100
        // }
        // var num = 0
        // num = value
        // console.log(typeof num)
        if ( Number(value)>99){
            callback(new Error('手续费设置过大'));
        }
        // if (!Number.isInteger(Number(value))) {
        //     callback(new Error('请输入数字值'));
        //   } 
      };
    return {
      configForm: {
          charge: 0,
        },
      configForm1: {
          charge: 0,
        },
        configForm2: {
          charge: 0,
        },
      Loading: true,
      chargeBool: true,
      rules: {
          charge: [
            { validator: validatecharge, trigger: 'blur' }
          ],
        },
      cedit: '编辑',
        // charge: [
        //   {validator: validateCharge,trigger:'blur'}
        // ],
    }
  },
  methods: {
    onSubmit() {
      this.Loading = true
      this.configForm.charge = this.configForm.charge /100
      SetConfig(this.configForm).then(response => {
        GetConfig().then(response => {
        this.Loading = false
        this.chargeBool = true
        this.cedit= '编辑'
        //this.configForm = response.data
        if (String(response.data.charge)==="undefined"){
          this.configForm.charge = 0
        }else{
          this.configForm.charge = Number(response.data.charge) * 100
        }
        this.configForm1.charge = this.configForm.charge
    })
      })
    },
    onCancel() {
      this.$message({
        message: 'cancel!',
        type: 'warning'
      })
    },
    editCharge(){
      this.chargeBool = !this.chargeBool
      if (!this.chargeBool){
        this.cedit = '取消编辑'
      }else{ 
        this.cedit= '编辑'
      }
    },
    Mnemonic(){
      // var param = {
      //   wallet : 'f13djvobwob2les2mkqwttg6vkem2ozvofohef46q',

      // }
      // exportPrivateKey(param).then(response => {
      //    console.log(AES.decrypt(response.data.list))

      // })
    }
  }
}
</script>

<style scoped>
.line{
  text-align: center;
}
</style>

