<template>
  <div class="app-container">
    
    <!-- <el-select style="margin-left: 20px;margin-right: 10px; width:220px;"  placeholder="文件类型">
        <el-option
      v-for="item in options"
      :key="item.value"
      :label="item.label"
      :value="item.value"
      >
        </el-option> 
      </el-select> -->
      <!-- <el-button type="primary" @click="ApplyList" >搜索</el-button> -->
 
      <el-button type="primary" @click="dialogVisible = true"  >导入审核列表</el-button>
      <el-button type="success" @click="auditUsers" style="margin-left: 20px;">批量审核</el-button>
      <el-button type="danger" @click="refusedUsers" style="margin-left: 20px;">批量拒绝</el-button>
      <!-- <el-input style="width: 220px; float: right" v-model="input" placeholder="请输入内容"></el-input> -->
      <span style="margin-left: 40px;color:red">钱包地址：{{wallet.address}}</span>
      <span style="margin-left: 20px;color:green">余额：{{wallet.balance}}　FIL</span>
  
      <el-dialog
        title="请选择文件"
        :visible.sync="dialogVisible"
        width="30%">
      <el-upload
        class="upload-demo"
        :action="base_api"
        :before-remove="beforeRemove"
        multiple
        :headers="uploadh"
        :limit="3"
        :on-exceed="handleExceed"
        :on-success="fetchData"
        :file-list="fileList">
      <el-button size="small" type="primary">批量导入</el-button>
  <!-- <div slot="tip" class="el-upload__tip">只能上传jpg/png文件，且不超过500kb</div> -->
      </el-upload>
 <!-- <span slot="footer" class="dialog-footer">
    <el-button @click="dialogVisible = false">取 消</el-button>
    <el-button type="primary" @click="dialogVisible = false">确 定</el-button>
  </span> -->

    </el-dialog>

    <el-dialog
        :title="atitle"
        :visible.sync="auditVisible"
        width="30%">

    </el-dialog>

    <el-table
      v-loading="listLoading"
      :data="list"
      element-loading-text="Loading"
      border
      fit
      highlight-current-row
      style="margin-top:20px;"
      @selection-change="handleSelectionChange"
      empty-text="暂无数据"
    >
    <el-table-column type="selection" width="55" align="center">
    </el-table-column>
       <el-table-column label="批次" align="center">
        <template slot-scope="scope">
           <!-- <a style="color:blue;font-size:16px;" @click="$router.push({name: 'MentionAudit', params: {num: scope.row.num}})">第 {{ scope.row.num }} 批次</a> -->
           <a style="color:blue;font-size:16px;" @click="$router.push({path: '/finance/index/audita?num='+scope.row.num})">第 {{ scope.row.num }} 批次</a>
        </template>
      </el-table-column>

      <el-table-column label="创建人" align="center">
        <template slot-scope="scope">
          <span>{{ scope.row.operator }}</span>
        </template>
      </el-table-column>

      <el-table-column class-name="status-col" width="120" label="状态" align="center">
        <template slot-scope="scope">
          <el-tag :type="scope.row.surplus | statusType">
            <div v-if="scope.row.surplus==='0'">审核完成</div>
            <div v-else>未审核</div>
          </el-tag>
        </template>
      </el-table-column>
    
      <el-table-column align="center" prop="created_at" label="创建时间" width="240">
        <template slot-scope="scope">
          <i class="el-icon-time" />
          <span>{{ scope.row.create_time }}</span>
        </template>
      </el-table-column>

      <el-table-column label="操作" width="200" align="center" fixed="right">
      <template slot-scope="scope">
      <el-button
          size="mini"
          type="success"
          @click="auditUser(scope.row)">全部通过</el-button>
        <el-button
          size="mini"
          type="danger" 
          @click="refusedUser(scope.row)">全部拒绝</el-button>
      </template>
    </el-table-column>
    </el-table>

    <div class="block" style="margin-top:10px;float:right;">
    <!-- layout="total, sizes, prev, pager, next, jumper" -->
    <el-pagination
      background
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
      :page-sizes="[10, 20, 50, 100]"
      :page-size="10"
      layout="sizes, prev, pager, next, jumper"
      :total="total">
    </el-pagination>
  </div>
  </div>

</template>

<script>
import { batchlist,batchSends,batchRefuse,walletbalance } from '@/api/finance'
import { getToken} from '@/utils/auth'

export default {
  filters: {
    statusFilter(status) {
      const statusMap = {
        "0": '审核完毕',
      }
      return statusMap[status]
    },
    statusType(status) {
      const statusMap = {
        "0": 'success',
      }
      return statusMap[status]
    }
  },
  data() {
    return {
      uploadh:{
        token:''
      },
      list: null,
      listLoading: true,
      fileList: undefined,
      pageParams:{
          pageNum: 1,
          pageSize: 10,
          status:0,
      },
      wallet:{
          address: "",
          balance: "0",
      },
      dialogVisible: false,
      auditVisible: false,
      total:0,
      atitle:"",
      base_api: process.env.VUE_APP_BASE_API+"yungo/uploadApply",
      wparam:{
        fil:"",
        to:"",
        id:"",
      },
      muSelect: []
    }
  },
  created() {
    this.fetchData()
    this.uploadh.token = getToken()
  },
  methods: {
    fetchData() {
      this.listLoading = true
      batchlist(this.pageParams).then(response => {
        this.list = response.data.list
        this.listLoading = false
        this.total = response.data.total
      })
      walletbalance().then(response => {
        this.wallet.address = response.data.address
        this.wallet.balance = response.data.balance
      })
    },
     BatchDetails(){
      //this.$router.push({ path:'/two.html'  })
      console.log("dskfhskd")
    },
    beforeRemove(file, fileList) {
        return this.$confirm(`确定移除 ${ file.name }？`);
      },
    handleSizeChange(param){
      this.pageParams.pageSize = param
      this.fetchData()
    },
    handleExceed(files, fileList) {
        this.$message.warning(`当前限制选择 1 个文件，本次选择了 ${files.length} 个文件，共选择了 ${files.length + fileList.length} 个文件`);
      },
    handleCurrentChange(param){
      // alert(param)
      this.pageParams.pageNum = param
      this.fetchData()
    },
    auditUser(val){
      // this.atitle = val.user_name 
      // this.auditVisible = true
      this.$confirm(`确定通过：第 ${ val.num } 批次 的提币审核，并转账？`, '通过审核', {
            cancelButtonText: '取消',
            confirmButtonText: '确定',
            type: 'success'
      }).then(() => {
        this.wparam.id = val.num
        batchSends(this.wparam).then(() => {
            this.fetchData()
        })
      }).catch(() => {
        //console.log("取消")
      })
    },
    auditUsers(){
      // this.atitle = val.user_name 
      // this.auditVisible = true
      if (this.muSelect.length==0){
        this.$message("未选中用户！")
        return
      }
      var names = " 【 ";
      var wparam = {
        id:"",
      }
      for (var i=0;i<this.muSelect.length;i++){
        if (i!=0){
          names += "，"
          wparam.id += ","
        }
        names += this.muSelect[i].num 
        wparam.id += this.muSelect[i].num 

      }
      this.$confirm(`确定通过：第${ names } 】批次 的提币审核，并转账？`, '通过审核', {
            cancelButtonText: '取消',
            confirmButtonText: '确定',
            type: 'success'
      }).then(() => {
        batchSends(wparam).then(() => {
            this.fetchData()
        })
      }).catch(() => {
        //console.log("取消")
      })
    },
    refusedUser(val){
      // this.atitle = val.user_name 
      // this.auditVisible = true
      this.$confirm(`确定拒绝：第 ${ val.num } 批次 的提币申请?`, '提示', {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
      }).then(() => {
        const id = {
        "ids": val.num,
        "status": 2
        }
        batchRefuse(id).then(response => {
            this.fetchData()
        })
      }).catch(() => {
        //console.log("取消")
      })
    },
    refusedUsers(val){
      // this.atitle = val.user_name 
      // this.auditVisible = true
      var names = " 【 ";
      var ids = ""
      
      for (var i=0;i<this.muSelect.length;i++){
        if (i!=0){
          names += "，"
          ids += ","
        }
        names += this.muSelect[i].num 
        ids += this.muSelect[i].num 
      }
      this.$confirm(`确定拒绝：第 ${ names } 】批次 的提币申请？`, '提示', {
            cancelButtonText: '取消',
            confirmButtonText: '确定',
            type: 'warning'
      }).then(() => {
        const id = {
        "ids": ids,
        "status": 2
        }
        batchRefuse(id).then(response => {
            this.fetchData()
        })
      }).catch(() => {
        //console.log("取消")
      })
    },
    handleSelectionChange(val){
      this.muSelect = val
    }
  }
}
</script>
