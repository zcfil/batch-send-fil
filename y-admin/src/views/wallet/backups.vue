<template>
<div>
  <div style="margin-left:1%;margin-top:1%;">
  <el-button type="primary" @click="ImportVisible=true" style="margin-left:1%;">导入密钥</el-button>
  <el-button type="primary" @click="ImportVisibleMne=true" style="margin-left:1%;">导入助记词</el-button>
  </div>
  
<el-dialog
  width="33%"
  title="导入密钥"
  :visible.sync="ImportVisible"
  >
  <div style="height:120px;">
  <el-input
  type="textarea"
  :rows="4"
  placeholder="请输入内容"
  v-model="PrivateKey"
  resize="none">
  </el-input>


  </div>
  <span slot="footer" class="dialog-footer">
    <el-button @click="ImportVisible = false">取 消</el-button>
    <el-button type="primary" @click="onImportKey">确 定</el-button>
  </span>
  </el-dialog>

<el-dialog
  width="33%"
  title="导入助记词"
  :visible.sync="ImportVisibleMne"
  >
  <div style="height:120px;">
  <el-input
  type="textarea"
  :rows="4"
  placeholder="请输入内容"
  v-model="mnemonic"
  @input="onChange"
  resize="none">
  </el-input>


  </div>
  <span slot="footer" class="dialog-footer">
    <el-button @click="ImportVisibleMne = false">取 消</el-button>
    <el-button type="primary" @click="onImportKeyMne">确 定</el-button>
  </span>
  </el-dialog>
  

   <el-table
      :data="tableData"
      style="width: 90%;margin-left:2%;margin-top:1%;">
      <el-table-column
        fixed="left"
        prop="address"
        label="钱包地址"
        width="640"
        >
      </el-table-column>
      <el-table-column
        prop="balance"
        width="240"
        label="余额（FIL）"
      > 
      </el-table-column>
      <el-table-column
        prop="isdefault"
        label="是否默认钱包"
      > 
      <template slot-scope="scope">
          <div v-if="scope.row.isdefault==='1'">默认地址</div>
         <div v-else></div>
        </template>
      
      </el-table-column>
      <el-table-column
        fixed="right"
        label="操作"
        width="280">
        <template slot-scope="scope"> 
          <el-button @click="handleDefault(scope.row)" type="text" size="small">设为默认钱包</el-button>
          <el-button @click="MnemonicClick(scope.row)" type="text" size="small">查看助记词</el-button>
          <el-button @click="PrivateClick(scope.row)" type="text" size="small">查看私钥</el-button>
          <el-button @click="handleDelete(scope.row)" type="text" size="small">删除</el-button>
      </template>
      </el-table-column>
    </el-table>

    <el-dialog
  title="查看助记词"
  :visible.sync="MnemonicVisible"
  width="40%"
  >
  <div style="height:380px;">
  <div style="margin-left:5%;border:1px solid #000;height:234px;width:90%;">
      <el-tag v-for="(item,index) in items"
    :key="index"
    effect="dark" 
    type="info"
    :disable-transitions="true"
    style="margin-top: 3%;margin-left:3%;width:13%;height:34px;text-align:center">
    {{item}}</el-tag>

    </div>
    <div style="margin-top: 10px;text-align:center">
      <span style="font-size:13px;color:gray;">请依次抄写并妥善保管，助记词是您恢复钱包的唯一手段。一经丢失，无法找回。请勿使用截屏、拍照等方式保存。</span>
    </div>
    </div>
  </el-dialog>

  <el-dialog
  title="查看私钥"
  :visible.sync="PrivateVisible"
  width="40%">
  <div style="height:200px;">
  <el-input
  type="textarea"
  :rows="4"
  placeholder="请输入内容"
  v-model="PrivateText"
  readonly
  resize="none">
  </el-input>
  <div style="margin-top: 10px;">
      <span style="font-size:13px;color:gray;">请依次抄写并妥善保管，私钥是您恢复钱包的唯一手段。一经丢失，无法找回。请勿使用截屏、拍照等方式保存。</span>
  </div>
  </div>
  </el-dialog>
  
</div>
</template>

<script>
import { walletList,exportMnemonic,importWallet,exportPrivateKey,delWallet,setWallet,importMnemonic } from '@/api/wallet'
import AES from '@/common/AES'
export default {
  name: 'Login',
  data() {
    
    return {
      tableData: [
          
        ],
        dialogVisible:false,
        MnemonicVisible:false,
        PrivateVisible:false,
        ImportVisible:false,
        ImportVisibleMne:false,
        PrivateText:"",
        mnemonic:"",
        PrivateKey:"",
        items: [
          
        ],
        flage:true,
    }
  },
  created(){
    walletList().then(response => {
      this.tableData = response.data.list
    })
  },
  watch: {

  },
  methods: {
    handleDefault(val){
      this.$confirm(`确定将 ${ val.address } 钱包设置为默认钱包?`, '提示', {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
      }).then(() => {
        const id = {
        "wallet": val.address,
        }
        setWallet(id).then(response => {
            walletList().then(response => {
              this.tableData = response.data.list
            })
        })
      }).catch(() => {
        //console.log("取消")
      })
    },
    handleDelete(val){
      this.$confirm(` ${ val.address } ?`, '确定删除', {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'error'
      }).then(() => {
        const id = {
        "wallet": val.address,
        }
        delWallet(id).then(response => {
            walletList().then(response => {
              this.tableData = response.data.list
            })
        })
      }).catch(() => {
        //console.log("取消")
      })
    },
    MnemonicClick(val){
      var param = {
        wallet: val.address
      }
      this.items = ""
      exportMnemonic(param).then(response => {
        if (response.data===""){
          this.items = "未通过助记词导入钱包"
          return
        }
        // this.items = AES.decrypt(response.data).split(' ')
        this.items = AES.decrypt(response.data).split(' ')
      })
      this.MnemonicVisible = true
    },
    PrivateClick(val){
      this.PrivateVisible = true
      var param = {
        wallet: val.address
      }
      this.PrivateText = ""
      exportPrivateKey(param).then(response => {
          this.PrivateText = AES.decrypt(response.data) 
        })
    },
    onImportKeyMne(){
      this.mnemonic = this.mnemonic.replace("\n"," ")
      this.mnemonic = this.mnemonic.replace(/(^\s*)|(\s*$)/g,"");
      var param = {
        mnemonic:AES.encrypt(this.mnemonic)
      }
      console.log(this.mnemonic)
      this.ImportVisible = false
      this.$confirm(`确定导入密钥?`, '提示', {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
      }).then(() => {
        importMnemonic(param).then(response => {
            this.$message({
              message:"导入成功！",
              type:'success'
              })
            walletList().then(response => {
              this.tableData = response.data.list
            })
        })
      }),

      this.mnemonic = ""
    },
     onChange(){
      var str = this.mnemonic.split(" ")
      console.log(str)
      var len = str.length
      if (str[len-1]===""){
        len--
      }
      if (len>=12&&this.mnemonic.indexOf("\n")==-1){
        this.mnemonic = ''
        for (var i=0;i<str.length;i++){
           if (i == 11&&str[i]!==" "&&str[i]!==""&&str[i]!=="\n"&&this.flage){
             this.flage = false
             this.mnemonic += str[i]+ '\n'
           }else{
              this.mnemonic += str[i] +" "
           }
        }
      }
      if (len>=24){
        this.flage = true
      }
      // console.log(this.mnemonic)
    },

   onImportKey(){

      var param = {
        private_key:AES.encrypt(this.PrivateKey)
      }
      this.$confirm(`确定导入密钥?`, '提示', {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
      }).then(() => {
        importWallet(param).then(response => {
            this.ImportVisibleMne = false
            this.$message({
              message:"导入成功！",
              type:'success'
              })
            walletList().then(response => {
              this.tableData = response.data.list
            })
        })
      })

      this.PrivateKey = ""
    },
  }
}
</script>

<style lang="scss">
/* 修复input 背景不协调 和光标变色 */
/* Detail see https://github.com/PanJiaChen/vue-element-admin/pull/927 */

</style>

