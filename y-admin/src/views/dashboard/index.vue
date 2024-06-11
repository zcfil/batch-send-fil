<template>
  <div class="dashboard-container" style="margin-top:0px;">
    <el-progress  type="circle" :percentage="circle" :show-text="false" :width="30" ></el-progress>
    <!-- <a class="login-a" href="/login">登录</a> -->
    <div id="div-text">
        Filecoin挖矿成本实时查询
    </div>
    <div class="dashboard-text">
      当前区块高度：{{ height1 }}
    </div>
    <div class="div-btn">
      <button :class="isXz==='1'?'btn11':'btn1'" @click="onclicks('1')">每32GB</button>
      <button :class="isXz==='2'?'btn22':'btn'" @click="onclicks('2')">每1TB</button>
      <button :class="isXz==='3'?'btn22':'btn'" @click="onclicks('3')">每1PB</button>
    </div>
  <el-card class="box-card1">
    <div class="box-text">
      封装扇区总成本
    </div>
    <div class="box-fil" style="color: rgb(194, 53, 49);">
      {{pleagelist1.totalfil}} FIL
    </div>
    <div class="box-rmb">
      {{(pleagelist1.totalfil*this.pleagelist1.cnytofil).toFixed(2)|NumberFormat}} ￥
    </div>
  </el-card>
  <el-card class="box-card">
    <div class="box-text">
      封装扇区质押费
    </div>
    <div class="box-fil" style="color: rgb(47, 69, 84);">
      {{pleagelist1.pleagefil}} FIL
    </div>
    <div class="box-rmb">
      {{(this.pleagelist1.pleagefil*this.pleagelist.cnytofil).toFixed(2)|NumberFormat}} ￥
    </div>
  </el-card>
  <el-card class="box-card">
    <div class="box-text">
      封装扇区Gas费
    </div>
    <div class="box-fil" style="color: rgb(97, 160, 168);">
      {{pleagelist1.gasfil}} FIL
    </div>
    <div class="box-rmb">
      {{(pleagelist1.gasfil*this.pleagelist.cnytofil).toFixed(2)|NumberFormat}} ￥
    </div>
  </el-card>
  
  <div id="gas">

  </div>

  </div>
</template>


<script>
import { mapGetters } from 'vuex'
import { getNowHeight,getFileType } from '@/utils/auth'
import { getpleagefil,getstatisticalfil } from '@/api/table'
import echarts from 'echarts'
import { number } from 'echarts/lib/export'

export default {
  name: 'Dashboard',
  filters: {
      formatN: function (s) {
        s = s + ''
        // eslint-disable-next-line no-useless-escape
        if (/[^0-9\.]/.test(s)) return ''
        s = s.replace(/^(\d*)$/, '$1.')
        s = (s + '00').replace(/(\d*\.\d\d)\d*/, '$1')
        s = s.replace('.', ',')
        const re = /(\d)(\d{3},)/
        while (re.test(s)) {
          s = s.replace(re, '$1,$2')
        }
        s = s.replace(/,(\d\d)$/, '.$1')
        return s.replace(/^\./, '0.')
      },
  },
  data(){
    return{
      height1:',',
      circle: 0,
      index: 0,
      flag: false,
      pleagelist:{
         gasfil:1,
         pleagefil:1,
         totalfil:1,
         cnytofil:1,
      },
      pleagelistt:{
         gasfil:1,
         pleagefil:1,
         totalfil:1,
         cnytofil:1,
      },
      pleagelistp:{
         gasfil:1,
         pleagefil:1,
         totalfil:1,
         cnytofil:1,
      },
      pleagelist1:{},
      listLoading:false,
      isXz: '2',
      historyFIL:null,
      pleagefil:[],
      totalfil:[],
      cnytofil:[],
      gasfil:[],
      create_time:[],
      option:{},
      
    }
  },
  components:{

  },

  computed: {
    ...mapGetters([
      'name'
    ])
  },
  created(){
    this.NowHeight()
    //this.$nextTick(()=>{
      // self.setInterval(()=>{
      //   let flag = true
      //    if (this.circle == 100){
      //     this.circle = 0
      //     if (this.index==5){
      //       this.index=0
      //     }
      //     this.index += 1
      //     this.getFIL()
      //     this.NowHeight()
      //   }
      //   this.circle += 1
      //   this.myEchars()
      // },300);
      self.setTimeout(()=>{
        let flag = true
         if (this.circle == 100){
          this.circle = 0
          if (this.index==5){
            this.index=0
          }
          this.index += 1
          this.getFIL()
          this.NowHeight()
        }
        this.circle += 1
        this.myEchars()
      },300);
    //})
    //this.listLoading = true
     this.getFIL()
    
  },
  // beforeCreate(){
  //   getstatisticalfil().then(response => {

  //       this.historyFIL = response.data
  //       for (var i=0;i<response.data.length;i++){

  //         this.pleagefil.push(Number(response.data[i].pleagefil)) 
  //         this.gasfil.push(Number(response.data[i].gasfil))
  //         this.totalfil.push(Number(response.data[i].totalfil))
  //         this.cnytofil.push(Number(response.data[i].cnytofil))
  //         this.create_time.push(response.data[i].create_time)
  //       }
  //       console.log("云构：",this.pleagefil)
        
  //     })
  // },
  mounted(){
    this.getHistory()  //获取数据
  },
  methods: {

 myEchars(){
      let that =this
      var p = [1,2,3,4,5,6,7,8]
      //let echarts1 = ;
      let ech = require('echarts').init(document.getElementById("gas")) 
       this.option = {
        title: {
            // text: '32GB'
        },
        tooltip: {
            trigger: 'axis'
        },
        legend: {
            data: [ '总成本（32GB）','质押费（32GB）', 'Gas费（32GB）'],
        },
        grid: {
            left: '3%',
            right: '4%',
            bottom: '3%',
            containLabel: true
        },
        toolbox: {
            feature: {
                //saveAsImage: {}
            }
        },
        xAxis: {
            type: 'category',
            boundaryGap: false,
            data: this.create_time,
            // axisLabel:{  
            // //  interval: 0,  //控制坐标轴刻度标签的显示间隔.设置成 0 强制显示所有标签。设置为 1，隔一个标签显示一个标签。设置为2，间隔2个标签。以此类推
            // rotate:-82,//倾斜度 -90 至 90 默认为0 
            // textStyle:{ 
            //     fontWeight:"bold",  //加粗
            //     color:"#000000"   //黑色
            // },                 
          // },
        },
        yAxis: {
            type: 'value'
        },
        series: [
            {
                name: '总成本（32GB）',
                type: 'line',
                stack: '总量',
                data: this.totalfil,
            },
            {
                name: '质押费（32GB）',
                type: 'line',
                stack: '总量',
                data: this.pleagefil,
            },
            {
                name: 'Gas费（32GB）',
                type: 'line',
                stack: '总量',
                data: this.gasfil
            },
            // {
            //     name: '汇率（￥）',
            //     type: 'line',
            //     stack: '总量',
            //     data: [],
            // }
        ]}
        ech.setOption(this.option)
    } ,
    NowHeight(){
      let genesis = new Date('2020/08/25 06:00:00')
      let height = (Date.parse(new Date())-Date.parse(genesis))/30/1000
      this.height1 = parseInt(height);
    },
    clock(){
      this.circle++ 
    },
    onclicks(val){
      // console.log(val)
      //alert("pb")
      switch(val){
        case '1':this.pleagelist1 = this.pleagelist;  this.isXz = '1'; break
        case '2':this.pleagelist1 = this.pleagelistt; this.isXz = '2' ;break
        case '3':this.pleagelist1 = this.pleagelistp; this.isXz = '3';break
        
      }
 
    },
    getFIL(){
      getpleagefil().then(response => {

        this.pleagelist.pleagefil = Number(response.data[this.index].pleagefil).toFixed(4)
        this.pleagelist.gasfil = Number(response.data[this.index].gasfil).toFixed(4)
        this.pleagelist.totalfil = Number(response.data[this.index].totalfil).toFixed(4)
        this.pleagelist.cnytofil = Number(response.data[this.index].cnytofil)

        this.pleagelistt.pleagefil = Number(response.data[this.index].pleagefil*32).toFixed(4)
        this.pleagelistt.gasfil = Number(response.data[this.index].gasfil*32).toFixed(4)
        this.pleagelistt.totalfil = Number(response.data[this.index].totalfil*32).toFixed(4)
        this.pleagelistt.cnytofil = Number(response.data[this.index].cnytofil)

     
        this.pleagelistp.pleagefil = Number(response.data[this.index].pleagefil*32*1024).toFixed(2)
        this.pleagelistp.gasfil = Number(response.data[this.index].gasfil*32*1024).toFixed(2)
        this.pleagelistp.totalfil = Number(response.data[this.index].totalfil*32*1024).toFixed(2)
        this.pleagelistp.cnytofil = Number(response.data[this.index].cnytofil)

        this.onclicks(this.isXz)
        //this.listLoading = false
      })
    },
    getHistory(){
      getstatisticalfil().then(response => {

        this.historyFIL = response.data
        for (var i=0;i<response.data.length;i++){

          this.pleagefil.push(Number(response.data[i].pleagefil)) 
          this.gasfil.push(Number(response.data[i].gasfil))
          this.totalfil.push(Number(response.data[i].totalfil))
          this.cnytofil.push(Number(response.data[i].cnytofil))
          this.create_time.push(response.data[i].create_time)
        }

        this.flag = true
      })
    }
  },
}
</script>

<style lang="scss" scoped>
// .progress{
//   float: left;
//   width: 100px;
// }
.login-a{
  font-size: 15px;
  font-weight:normal;
  float: right;
  color: blue;
  //margin-top: -40px;
}
a:hover{

color: red;
text-decoration:underline;

}
.dashboard {
  &-container {
    margin: 30px;
  }
  &-text {
    font-size: 30px;
    line-height: 46px;
  }
}
.dashboard-text{
  text-align:center;
  font-size: 28px;
  margin-top: 15px;
  font-weight:bold;
  color: rgb(156, 132, 9);
}
#div-text{
  text-align:center;
  font-size: 40px;
  margin-top: 0px;
  font-weight:bold;
  color: rgb(86, 89, 92);
}
.div-btn {
    height: 160px;
  }
.btn1{
  box-shadow: 5px 5px 5px #888888;
  width: 9%;
  height: 60px;
  margin-left: 33%;  
  border-radius:6px;
  border: none; 
  color: red;
  font-size: 20px;
  margin-top: 50px; 
}
.btn11{
  box-shadow: 5px 5px 5px #888888;
  width: 9%;
  height: 60px;
  margin-left: 33%;  
  border-radius:6px;
  color: black;
  border: 3px solid black;
  font-size: 19px;
  margin-top: 50px; 
}
.btn{
  box-shadow: 5px 5px 5px #888888;
  border: none; 
  border-radius:6px;
  width: 9%;
  height: 50px;
  margin-left: 32px;  
  margin-top: 50px; 
  color: red;
  font-size: 20px;
}
.btn22{
  box-shadow: 5px 5px 5px #888888;
  border-radius:6px;
  width: 9%;
  height: 60px;
  margin-left: 32px;  
  margin-top: 50px; 
  color: black;
  border: 3px solid black;
  font-size: 19px;
}
.box-card1 {
    float: left;
    width: 20%;
    height: 180px;
    margin-left: 16%;
  }
.box-card {
    float: left;
    width: 20%;
    height: 180px;
    margin-left: 40px;
  }
.box-text{
  text-align:center;
  font-size: 20px;
  margin-top: 15px;
  font-weight:bold;
  color: rgb(86, 89, 92);
}
.box-fil{
  text-align:center;
  font-size: 28px;
  margin-top: 20px;
  //font-weight:bold;
}
.box-rmb{
  text-align:right ;
  //font-size: 28px;
  margin-top: 20px;
  //font-weight:bold;
  color: rgb(212, 130, 101)
}
#gas{
  margin-top: 20px;
  float: left;
  height: 360px;
  width:100%;
  margin-left: 0%;
}
.div-gas{
  float: left;
  height: 200px;
  width:100%;
}
@media screen and (max-width:600px) {
    .dashboard {
  &-container {
    margin: 30px;
  }
  &-text {
    font-size: 30px;
    line-height: 46px;
  }
}
.dashboard-text{
  text-align:center;
  font-size: 28px;
  margin-top: 15px;
  font-weight:bold;
  color: rgb(156, 132, 9);
}
.div-text{
  text-align:center;
  font-size: 40px;
  margin-top: 15px;
  font-weight:bold;
  color: rgb(86, 89, 92);
}
.div-btn {
    height: 100px;
  }
.btn1{
  box-shadow: 1px 1px 1px #888888;
  width: 30%;
  height: 40px;
  margin-left: 2.5%;  
  border-radius:6px;
  border: none; 
  color: red;
  font-size: 14px;
  margin-top: 30px; 
  float: left;
  // text-align:center;
}
.btn11{
  box-shadow: 5px 5px 5px #888888;
  text-align:center;
  width: 30%;
  height: 40px;
  margin-left: 2.5%;  
  border-radius:6px;
  color: orange;
  border: 3px solid orange;
  font-size: 15px;
  margin-top: 30px; 
  float: left;
}
.btn{
  text-align:center;
  box-shadow: 1px 1px 1px #888888;
  border: none; 
  border-radius:6px;
  width: 30%;
  height: 40px;
  margin-left: 2.5%;  
  margin-top: 30px; 
  color: red;
  font-size: 14px;
  float: left;
}
.btn22{
  box-shadow: 5px 5px 5px #888888;
  border-radius:6px;
  width: 30%;
  height: 40px;
  margin-left: 2.5%;  
  margin-top: 30px; 
  color: orange;
  border: 3px solid orange;
  font-size: 15px;
  text-align:center;
  float: left;
}
.box-card1 {
  margin-top: 1px; 
    float: left;
    width: 96%;
    height: 130px;
    margin-left: 2%;
  }
.box-card {
  margin-top: 20px; 
    float: left;
    width: 96%;
    height: 130px;
    margin-left: 2%;
  }
.box-text{
  text-align:center;
  font-size: 20px;
  margin-top: 1px;
  font-weight:bold;
  color: rgb(86, 89, 92);
}
.box-fil{
  text-align:center;
  font-size: 28px;
  margin-top: 10px;
  //font-weight:bold;
}
.box-rmb{
  text-align:right ;
  //font-size: 28px;
  margin-top: 16px;
  //font-weight:bold;
  color: rgb(41, 49, 49);
}
}
</style>
