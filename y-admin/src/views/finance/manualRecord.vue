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
    <el-table
      :data="list"
      border
      fit
      highlight-current-row
      style="margin-top:20px;"
    >
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
import { ManualList } from '@/api/finance'

export default {
    name:"ManualRecord",
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
      list:[],
      total:0,
      value: '',
      pageParams:{
          pageNum: 1,
          pageSize: 10,
          status:'',
      },
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
    }
  },
  created() {
    this.auto()
  },
  methods: {
    auto() {
      ManualList(this.pageParams).then(response => {
        this.list = response.data.list
        for (var i=0;i<this.list.length;i++) {
            if (this.list[i].status==0||this.list[i].status==2) {
                this.list[i].charge1="0"
            }else {
                this.list[i].charge1=this.list[i].amount-this.list[i].filcount
            }
        }
        this.total = response.data.total
      })
    },
    recordList() {
      this.pageParams.status = this.value
      this.auto()
    },
    handleSizeChange(param){
      this.pageParams.pageSize = param
      this.auto()
    },
    handleCurrentChange(param){
      this.pageParams.pageNum = param
      this.auto()
    },
  }
}
</script>