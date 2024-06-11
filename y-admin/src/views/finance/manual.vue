<template>
    <div class="manual">
      <el-button type="success" @click="newMoney">新增提币</el-button>
      <el-button type="success" @click="auditUsers" style="margin-left: 20px;">批量审核</el-button>
      <el-button type="danger" @click="refusedUsers" style="margin-left: 20px;">批量拒绝</el-button>
      <span style="margin-left: 40px;color:red">钱包地址：{{wallet.address}}</span>
      <span style="margin-left: 20px;color:green">余额：{{wallet.balance}}　FIL</span>
      <el-table
       :data="list"
       border
       fit
       highlight-current-row
       style="margin-top:20px;"
       @selection-change="handleSelectionChange"
       empty-text="暂无数据"
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
        <el-table-column label="消息ID" :show-overflow-tooltip="true" width="220" align="center">
          <template slot-scope="scope">
            <span>{{ scope.row.cid }}</span>
          </template>
        </el-table-column>
        <el-table-column label="金钻" align="center">
          <template slot-scope="scope">
            {{ scope.row.amount }}
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
        <el-table-column label="类型" width="110" align="center">
          <template slot-scope="scope">
            {{ scope.row.coin_type }}
          </template>
        </el-table-column>
        <el-table-column align="center" prop="created_at" label="创建时间" width="200">
          <template slot-scope="scope">
            <i class="el-icon-time" />
            <span>{{ scope.row.create_time }}</span>
          </template>
        </el-table-column>
  
        <el-table-column label="操作" width="150" align="center" fixed="right">
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
      </el-table-column>
      </el-table>
  
      <div class="block" style="margin-top:10px;float:right;">
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
      <el-dialog title="新增提钻" :visible.sync="open" width="500px">
          <el-form ref="form" :model="form" :rules="rules" label-width="80px">
            <el-form-item label="姓名" prop="user_name">
              <el-input v-model="form.user_name" placeholder="请输入姓名" />
            </el-form-item>
            <el-form-item label="提钻数" prop="amount">
              <el-input v-model="form.amount" placeholder="请输入提钻数" />
            </el-form-item>
            <el-form-item label="提钻地址" prop="address">
              <el-input v-model="form.address" placeholder="请输入提钻地址" />
            </el-form-item>
          </el-form>
          <div slot="footer" class="dialog-footer">
            <el-button type="primary" @click="submitForm">确 定</el-button>
            <el-button @click="open = false">取 消</el-button>
          </div>
        </el-dialog>
    </div>
</template>

<script>
import { walletbalance,ManualList,updateManual,sendManual,manualAdd } from '@/api/finance'
export default {
    name:'Manual',
    props: {},
    filters: {
    statusFilter(status) {
        const statusMap = {
          "0": '未审核',
          "1": '审核通过',
          "2": '审核拒绝'
        }
        return statusMap[status]
      },
      statusType(status) {
        const statusMap = {
          "1": 'success',
          "2": 'danger'
        }
        return statusMap[status]
      }
    },
    data() {
        return {
            wallet:{
                address: "",
                balance: "0",
            },
            muSelect: [],
            list:[],
            pageParams:{
                pageNum: 1,
                pageSize: 10,
                status:0,
            },
            total: 0,
            open:false,
            form:{},
            rules: {
              user_name: [
                { required: true, message: '姓名不能为空', trigger: 'blur' }
              ],
              amount: [
                { required: true, message: '提钻数不能为空', trigger: 'blur' }
              ],
              address: [
                { required: true, message: '提钻地址不能为空', trigger: 'blur' },
                {
                  pattern: /^(t1|t2|t3|f1|f2|f3)[a-zA-Z0-9]{38,}$/,
                  message: '请输入正确的地址',
                  trigger: 'blur'
                }
              ],
            },
        };
    },
    computed: {},
    created() {
        this.auto()
    },
    mounted() {},
    watch: {},
    methods: {
        auto() {
            ManualList(this.pageParams).then(response => {
              this.list = response.data.list
              this.total = response.data.total
            })
            walletbalance().then(response => {
              this.wallet.address = response.data.address
              this.wallet.balance = response.data.balance
            })
        },
        newMoney() {
          this.form={
            user_name:undefined,
            amount:undefined,
            address:undefined,
          }
          this.open=true
        },
        submitForm() {
          this.$refs['form'].validate(valid => {
            if (valid) {
              manualAdd(this.form).then(res=> {
                this.$message('新增成功')
                this.open = false
                this.auto()
              })
            }
          })
        },
        auditUser(val){
          var wparam = {
            fil:val.amount,
            to:val.address,
            id:val.id,
          }
          this.$confirm(`确定通过 【 ${ val.user_name } 】 的提币审核，并转账？`, '通过审核', {
                cancelButtonText: '取消',
                confirmButtonText: '确定',
                type: 'success'
          }).then(() => {
            sendManual(wparam).then(() => {
                this.auto()
            })
          }).catch(() => {
            //console.log("取消")
          })
        },
        auditUsers(){
          if (this.muSelect.length==0){
            this.$message("未选中用户！")
            return
          }
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
            sendManual(wparam).then(() => {
              this.auto()
            })
          }).catch(() => {
            //console.log("取消")
          })
        },
        refusedUser(val){
          this.$confirm(`确定拒绝 ${ val.user_name } 的提币申请?`, '提示', {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }).then(() => {
            const id = {
            "ids": val.id,
            "status": 2
            }
            updateManual(id).then(response => {
                this.$message("审核成功！")
                this.auto()
            })
          }).catch(() => {
            //console.log("取消")
          })
        },
        refusedUsers(){
          if (this.muSelect.length==0){
            this.$message("未选中用户！")
            return
          }
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
             updateManual(id).then(response => {
                 this.$message("审核成功！")
                 this.auto()
             })
           }).catch(() => {
             //console.log("取消")
           })
        },
        handleSelectionChange(val){
          this.muSelect = val
        },
        handleSizeChange(param){
          this.pageParams.pageSize = param
          this.auto()
        },
        handleCurrentChange(param){
          this.pageParams.pageNum = param
          this.auto()
        },
    },
    components: {},
};
</script>

<style scoped>
.manual {
    padding: 20px;
}
</style>
