<template>
  <div class="app-container">
      <el-select v-model="value" placeholder="请选择" clearable  @change="recordList" >
    <el-option
      v-for="item in statusMap"
      :key="item.value"
      :label="item.label"
      :value="item.value"
      >
    </el-option>
  </el-select>
      <!-- <el-button type="primary" @click="auditUsers" style="margin-left: 20px;">搜索</el-button> -->
      <!-- <el-button type="danger" @click="refusedUsers" style="margin-left: 20px;">批量拒绝</el-button> -->

    

    <el-table
      v-loading="listLoading"
      :data="list"
      element-loading-text="Loading"
      border
      fit
      highlight-current-row
      style="margin-top:20px;"
      @selection-change="handleSelectionChange"

    >
    <el-table-column type="selection" width="55" align="center">
    </el-table-column>
      <el-table-column align="center" label="序号" width="60">
        <template slot-scope="scope">
          {{ scope.$index+1 }}
        </template>
      </el-table-column>
      <el-table-column label="用户ID" width="120" align="center">
        <template slot-scope="scope">
          {{ scope.row.user_id }}
        </template>
      </el-table-column>
      <el-table-column label="用户姓名" align="center">
        <template slot-scope="scope">
          <span>{{ scope.row.user_name }}</span>
        </template>
      </el-table-column>
      <el-table-column label="消息ID" :show-overflow-tooltip="true" width="200" align="center">
        <template slot-scope="scope">
          <span>{{ scope.row.cid }}</span>
        </template>
      </el-table-column>
      <el-table-column label="金钻" align="center">
        <template slot-scope="scope">
          {{ scope.row.amount }}
        </template>
      </el-table-column>
      <el-table-column label="实际到账" align="center">
        <template slot-scope="scope">
          {{ Number(scope.row.filcount).toFixed(6) }}
        </template>
      </el-table-column>
      <el-table-column label="设置手续费" width="100" align="center">
        <template slot-scope="scope">
          {{ Number(scope.row.charge1).toFixed(6) }}
        </template>
      </el-table-column>
      <el-table-column label="官方手续费" width="100" align="center">
        <template slot-scope="scope">
          <!-- {{ Number(scope.row.charge).toFixed(6) }} -->
          {{ (scope.row.charge)}}
        </template>
      </el-table-column>
      <el-table-column label="提钻账号" :show-overflow-tooltip="true" width="340" align="center">
        <template slot-scope="scope">
          {{ scope.row.address }}
        </template>
      </el-table-column>
      <el-table-column class-name="status-col" width="120" label="状态" align="center">
        <template slot-scope="scope">
          <el-tag :type="scope.row.status | statusType">{{ scope.row.status | statusFilter }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="类型" width="80" align="center">
        <template slot-scope="scope">
          {{ scope.row.coin_type }}
        </template>
      </el-table-column>
      <!-- <el-table-column align="center" prop="created_at" label="审核时间" width="200">
        <template slot-scope="scope">
          <i class="el-icon-time" />
          <span>{{ scope.row.create_time }}</span>
        </template>
      </el-table-column> -->
      <el-table-column align="center" prop="created_at" label="创建时间" width="240">
        <template slot-scope="scope">
          <i class="el-icon-time" />
          <span>{{ scope.row.create_time }}</span>
        </template>
      </el-table-column>
      <el-table-column align="center" prop="created_at" label="审核时间" width="240">
        <template slot-scope="scope">
          <i class="el-icon-time" />
          <span>{{ scope.row.update_time }}</span>
        </template>
      </el-table-column>

      <!-- <el-table-column label="操作" width="150" align="center">
      <template slot-scope="scope">
      <el-button
          size="mini"
          type="success"
          @click="auditUser(scope.row)">审核</el-button>
        <el-button
          size="mini"
          type="danger" 
          @click="refusedUser(scope.row)">拒绝</el-button>
      </template>
    </el-table-column> -->
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
import { ApplyList,Sends,Refuse } from '@/api/finance'
import { getToken} from '@/utils/auth'

export default {
  filters: {
    statusFilter(status) {
      const statusMap = {
        "0": '未审核',
        "1": '已通过',
        "2": '已拒绝',
        "3": '确认中'
      }
      return statusMap[status]
    },
    statusType(status) {
      const statusMap = {
        "1": 'success',
        "2": 'danger',
        "3": 'info'
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
          status:'',
          num:'',
      },
      value: '',
      statusMap: [
        {
        value: "1",
        label: "已通过"
      },
      {
        value: "2",
        label: "已拒绝"
      },
      {
        value: "0",
        label: "未审核"
      },
      {
        value: "3",
        label: "确认中"
      }],
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
      this.pageParams.num = this.$route.query.num
      ApplyList(this.pageParams).then(response => {
        this.list = response.data.list
        this.listLoading = false
        this.total = response.data.total
      })
    },
    recordList() {
      this.pageParams.status = this.value
      this.fetchData()
    },
    beforeRemove(file, fileList) {
        return this.$confirm(`确定移除 ${ file.name }？`);
      },
    handleSizeChange(param){
      this.pageParams.pageSize = param
      this.fetchData()
    },
    // handleExceed(files, fileList) {
    //     this.$message.warning(`当前限制选择 1 个文件，本次选择了 ${files.length} 个文件，共选择了 ${files.length + fileList.length} 个文件`);
    //   },
    handleCurrentChange(param){
      // alert(param)
      this.pageParams.pageNum = param
      this.fetchData()
    },
    auditUser(val){
      this.atitle = "sdfls"
      this.auditVisible = true
      this.$confirm(`确定通过 ${ val.user_name } 的提币审核，并转账？`, '通过审核', {
            cancelButtonText: '取消',
            confirmButtonText: '确定',
            type: 'success'
      }).then(() => {
        this.wparam.fil = val.amount
        this.wparam.to = val.address
        this.wparam.id = val.id
        Sends(this.wparam).then(() => {
            this.fetchData()
        })
      }).catch(() => {
        //console.log("取消")
      })
    },
    auditUsers(){
      // this.atitle = val.user_name 
      // this.auditVisible = true
      
      var names = " 【 ";
      var wparam = {
        fil:"",
        to:"",
        id:"",
      }
      for (var i=0;i<this.muSelect.length;i++){
        if (i!=0){
          names += "，"
          wparam.id += ","
          wparam.to += ","
          wparam.fil += ","
        }
        names += this.muSelect[i].user_name 
        wparam.id += this.muSelect[i].id 
        wparam.to += this.muSelect[i].address 
        wparam.fil += this.muSelect[i].amount
      }
      this.$confirm(`确定通过: ${ names } 】 的提币审核，并转账？`, '通过审核', {
            cancelButtonText: '取消',
            confirmButtonText: '确定',
            type: 'success'
      }).then(() => {
        Sends(wparam).then(() => {
            this.fetchData()
        })
      }).catch(() => {
        //console.log("取消")
      })
    },
    refusedUser(val){
      // this.atitle = val.user_name 
      // this.auditVisible = true
      this.$confirm(`确定拒绝 ${ val.user_name } 的提币申请?`, '提示', {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
      }).then(() => {
        const id = {
        "ids": val.id,
        "status": 2
        }
        Refuse(id).then(response => {
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
        names += this.muSelect[i].user_name 
        ids += this.muSelect[i].id 
      }
      this.$confirm(`确定拒绝: ${ names } 】 的提币申请？`, '提示', {
            cancelButtonText: '取消',
            confirmButtonText: '确定',
            type: 'warning'
      }).then(() => {
        const id = {
        "ids": ids,
        "status": 2
        }
        Refuse(id).then(response => {
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
