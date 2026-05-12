import CryptoJS from 'crypto-js'
const nowtime = () => {
    let showDate = new Date()
    let seperator = '-'
    let year = showDate.getFullYear()
    let month = showDate.getMonth() + 1
    let day = showDate.getDate()
    let hour = showDate.getHours()
    let min = showDate.getMinutes()
    let seconds = showDate.getSeconds()
    var strDate = ''
    if (month >= 1 && month <= 9) {
        month = "0" + month
    }
    if (day >= 0 && day <= 9) {
        day = "0" + day
    }
    if (hour >= 0 && hour <= 9) {
        hour = "0" + hour
    }
    if (min >= 0 && min <= 9) {
        min = "0" + min
    }
    if (seconds >= 0 && seconds <= 9) {
        seconds = "0" + seconds
    }


    let currentdate = year + month + day + hour + min + seconds
    return currentdate
}
const nowtime1 = () => {
    let showDate = new Date()
    let seperator = '-'
    let year = showDate.getFullYear()
    let month = showDate.getMonth() + 1
    let day = showDate.getDate()
    let hour = showDate.getHours()
    let min = showDate.getMinutes()
    let seconds = showDate.getSeconds()
    var strDate = ''
    if (month >= 1 && month <= 9) {
        month = "0" + month
    }
    if (day >= 0 && day <= 9) {
        day = "0" + day
    }
    if (hour >= 0 && hour <= 9) {
        hour = "0" + hour
    }
    if (min >= 0 && min <= 9) {
        min = "0" + min
    }
    if (seconds >= 0 && seconds <= 9) {
        seconds = "0" + seconds
    }


    let currentdate = year + seperator + month + seperator + day + ' ' + hour + ':' + min + ':' + seconds
    return currentdate
}
const commonjs = {
    login: function () {
        console.log("1111")
    },
    pageSize: 10,
    myKey: '4dogs.cn4dogs.cn',
    encryptCBC (msg, key) {
        key = this.PaddingLeft(key, 16),//保证key的长度为16byte,进行'0'补位
            key = CryptoJS.enc.Utf8.parse(key)
        // 加密结果返回的是CipherParams object类型
        // key 和 iv 使用同一个值
        var encrypted = CryptoJS.AES.encrypt(msg.toString(), key, {
            iv: key,
            mode: CryptoJS.mode.CBC,// CBC算法
            padding: CryptoJS.pad.Pkcs7 //使用pkcs7 进行padding 后端需要注意
        })
        // ciphertext是密文,toString()内传编码格式,比如Base64,这里用了16进制
        // 如果密文要放在 url的参数中 建议进行 base64-url-encoding 和 hex encoding, 不建议使用base64 encoding
        return encrypted.ciphertext.toString(CryptoJS.enc.Hex)  //后端必须进行相反操作
    },
    // 确保key的长度,使用 0 字符来补位
    // length 建议 16 24 32
    PaddingLeft (key, length) {
        let pkey = key.toString()
        let l = pkey.length
        if (l < length) {
            pkey = new Array(length - l + 1).join('0') + pkey
        } else if (l > length) {
            pkey = pkey.slice(length)
        }
        return pkey
    }, 
    decryptCBC (msg, key) {
        key = this.PaddingLeft(key, 16)//保证key的长度为16byte,进行'0'补位
        key = CryptoJS.enc.Utf8.parse(key)
        var encryptedHexStr = CryptoJS.enc.Hex.parse(msg.toString())
        var encryptedBase64Str = CryptoJS.enc.Base64.stringify(encryptedHexStr)
        var decryptedData = CryptoJS.AES.decrypt(encryptedBase64Str, key, {
            iv: key,
            mode: CryptoJS.mode.CBC,
            padding: CryptoJS.pad.Pkcs7
        })
        return decryptedData.toString(CryptoJS.enc.Utf8) 
    },
    risklevellist: [[0, '全部'], [1, '高危'], [2, '中危'], [3, '低危'], [4, '未发现'],],
    taskstatus: [[0, '全部'], [1, '待触发'], [2, '待执行'], [3, '运行中'], [4, '已完成'], [5, '暂停中']], //待触发  1   待执行  2     运行中   3     已完成   4    暂停中   5 
    nowtime,
    nowtime1,
    timeer: null,
    timermillisec: 50000, //wdh add 时间间隔，默认3000毫秒
}

export default commonjs