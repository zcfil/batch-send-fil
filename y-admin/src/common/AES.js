
import CryptoJS from 'crypto-js';
 
export default {
    //随机生成指定数量的16进制key
    generatekey(num) {
        let library = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
        let key = "";
        for (var i = 0; i < num; i++) {
            let randomPoz = Math.floor(Math.random() * library.length);
            key += library.substring(randomPoz, randomPoz + 1);
        }
        return key;
    },
    getkey(){
        return "1234567890123456";
    },
    
    //加密
    encrypt(word) {
        let keyStr = '1234567890123456'; //判断是否存在ksy，不存在就用定义好的key
        var iv = CryptoJS.enc.Utf8.parse("1234567890123456");
        var key = CryptoJS.enc.Utf8.parse(keyStr);
        var srcs = CryptoJS.enc.Utf8.parse(word);
        var encrypted = CryptoJS.AES.encrypt(srcs, key, { iv:iv, mode: CryptoJS.mode.CTR, padding: CryptoJS.pad.NoPadding });
        return encrypted.ciphertext.toString()
    },
    //解密
    decrypt(word) {
        let keyStr = '1234567890123456';
        var iv = CryptoJS.enc.Utf8.parse("1234567890123456");
        var key = CryptoJS.enc.Utf8.parse(keyStr);
        var decrypt = CryptoJS.AES.decrypt(word, key, { iv:iv, mode: CryptoJS.mode.CTR, padding: CryptoJS.pad.NoPadding });
        return CryptoJS.enc.Utf8.stringify(decrypt).toString();
    }
    // encrypt(word) {
    //     var key = CryptoJS.enc.Utf8.parse("1234567890000000"); //16位
    //     var iv = CryptoJS.enc.Utf8.parse("1234567890000000");
    //     var encrypted = '';
    //     if (typeof(word) == 'string') {
    //         var srcs = CryptoJS.enc.Utf8.parse(word);
    //         encrypted = CryptoJS.AES.encrypt(srcs, key, {
    //             iv: iv,
    //             mode: CryptoJS.mode.CBC,
    //             padding: CryptoJS.pad.Pkcs7
    //         });
    //     } else if (typeof(word) == 'object') {//对象格式的转成json字符串
    //         data = JSON.stringify(word);
    //         var srcs = CryptoJS.enc.Utf8.parse(data);
    //         encrypted = CryptoJS.AES.encrypt(srcs, key, {
    //             iv: iv,
    //             mode: CryptoJS.mode.CBC,
    //             padding: CryptoJS.pad.Pkcs7
    //         })
    //     }
    //     return encrypted.ciphertext.toString();
    // },
    // decrypt(word) {
    //     var key = CryptoJS.enc.Utf8.parse("1234567890000000"); 
    //     var iv = CryptoJS.enc.Utf8.parse("1234567890000000");
    //     var encryptedHexStr = CryptoJS.enc.Hex.parse(word);
    //     var srcs = CryptoJS.enc.Base64.stringify(encryptedHexStr);
    //     var decrypt = CryptoJS.AES.decrypt(srcs, key, {
    //         iv: iv,
    //         mode: CryptoJS.mode.CBC,
    //         padding: CryptoJS.pad.Pkcs7
    //     });
    //     var decryptedStr = decrypt.toString(CryptoJS.enc.Utf8);
    //     return decryptedStr.toString();
    // }
}