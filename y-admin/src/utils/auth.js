import Cookies from 'js-cookie'

const TokenKey = 'vue_admin_template_token'

export function getToken() {
  return Cookies.get(TokenKey)
}

export function setToken(token) {
  var millisecond = new Date().getTime();
  var expiresTime = new Date(millisecond + 60 * 1000 * 60*12);
  return Cookies.set(TokenKey, token,{
    expires: expiresTime,
    })
}

export function removeToken() {
  return Cookies.remove(TokenKey)
}

export function getTime(oldTime){
  return oldTime.toLocaleDateString().replace(/\//g, "-") + " " + oldTime.toTimeString().substr(0, 8)
}
export function getFileSizeByBit(size){
  let count = size
  let unit = ['B','KB','MB','GB']
  let k = 0
  for(let i = 0;i < 4; i++){
    if (count<1024){
      break
    }
    count = count/1024
    k++
  }
  if (k ==0){
    return Number(count)+" " +unit[k] 
  }
  return Number(count).toFixed(2)+" " +unit[k] 
}


export function getExpiryTime(expiry){

  let genesis = new Date('2020-08-25 06:00:00')
  let Expiry = Date.parse(genesis)+expiry*30*1000
  var oldTime = new Date(Expiry)
  return getTime(oldTime)
}

export function getNowHeight(){

  let genesis = new Date('2020-08-25 06:00:00')
  let height = (Date.parse(new Date())-Date.parse(genesis))/30/1000
  return parseInt(height);
}
export function getFileType(str){
  let MapFileType = new Map([
  ["bmp",1],["gif",1],["jpg",1],["pic",1],["png",1],["tif",1],
  ["txt",2],["doc",2],["hlp",2],["wps",2],["rtf",2],["html",2],["pdf",2],
  ["avi",3],["mpg",3],["mov",3],["swf",3],["rm",3],["mkv",3],["rmvb",3],["ogg",3],["mod",3],["wmv",3],["qt",3],["asf",3],["navi",3],["divx",3],["mpeg",3],["dat",3],["mp4",3],
  ["wav",4],["aif",4],["au",4],["mp3",4],["ram",4],["wma",4],["mmf",4],["amr",4],["aac",4],["flac",4],
  ])
  
  // return MapFileType.get('jpg')
  return MapFileType.get(str)
}