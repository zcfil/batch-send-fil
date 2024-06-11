 <template>
    <el-table
      :data="tableData"
      style="width: 100%">
      <el-table-column
        prop="Address"
        label="Address"
        width="600">
      </el-table-column>
      <el-table-column
        prop="Balance"
        label="Balance">
      </el-table-column>
      <el-table-column label="操作">
        <el-button
          @click.native.prevent="deleteRow(scope.$index, tableData)"
          type="text"
          size="small">
          移除
        </el-button>
    </el-table-column>
    </el-table>
  </template>

  <script>
    export default {
      
      data() {
        return {
          tableData: [{
            Address: 'f3shcfsi7b557m2fuenen6fqd7pogfz2obt2fzkxzocqhvkotlpdyelevufoy6jfiscb4tjw4avjx5ibf7wtra',
            Balance: '81.830674494633663574 FIL'
          }, {
            Address: 'f3shcfsi7b557m2fuenen6fqd7pogfz2obt2fzkxzocqhvkotlpdyelevufoy6jfiscb4tjw4avjx5ibf7wtra',
            Balance: '81.830674494633663574 FIL'
          }, 
          {
            Address: 'f3shcfsi7b557m2fuenen6fqd7pogfz2obt2fzkxzocqhvkotlpdyelevufoy6jfiscb4tjw4avjx5ibf7wtra',
            Balance: '81.830674494633663574 FIL'
          },]
        }
      },
     methods: {
      deleteRow(index, rows) {
        rows.splice(index, 1);
      }
      },
      
    }
  </script>