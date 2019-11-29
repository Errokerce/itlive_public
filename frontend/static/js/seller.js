const app = new Vue({
  delimiters: ['&{', '}'],
  el: '#app',
  data() {
    return {
      bidTableVisible: false,
      dialogFormVisible: false,
      itemTableVisible: false,
      bidSusDialogVisible: false,
      bidFailDialogVisible: false,
      timeendDialogVisible: false,
      socket: null,
      name: name,
      count: 20,
      msg: '',
      chatData: [],
      itemID: 0,
      formitems: {
        id: '',
        itemName: 'pen',
        itemQuantity: '1',
        aSetQuantity: '1',
        itemPrice: '100',
        addPrice1: '200',
        addPrice2: '',
        addPrice3: '',
        itemText: ''
      },
      items: [],                  //登入的商品資訊
      tempitem: {
        id: '',
        itemName: 'pen',
        itemQuantity: '1',
        aSetQuantity: '1',
        itemPrice: '100',
        addPrice1: '200',
        addPrice2: '',
        addPrice3: '',
        itemText: ''
      },
      biditemSelect:
      {
        itemName: '',
        // itemQuantity: '',
        // aSetQuantity: '',
        // itemPrice: '',
        // addPrice1: '',
        // addPrice2: '',
        // addPrice3: '',
        // itemText: '',
        totalTime: ''
      },
      biditemForm: [],
      formLabelWidth: '100px',
      itembid: [],                //目前競標的商品
      totalTime: '',              //目前的競標時間
      record: '',                 //紀錄選擇的競標時間
      showselect: true,
      showBidBut: true,
      showConBut: false,
      showPauBut: false,
      showResBut: false,
      timer: null,
      times: [
        {
          id: '1',
          totalTime: (30),
          title: '00:30'

        },
        {
          id: '2',
          totalTime: (1 * 60),
          title: '01:00'
        },
        {
          id: '3',
          totalTime: (1 * 90),
          title: '01:30'

        },
        {
          id: '4',
          totalTime: (2 * 60),
          title: '02:00'

        },
        {
          id: '5',
          totalTime: (1 * 150),
          title: '延長時間'

        }],
      bidtable: [],             //競標者名字及最高價
      tempheight: [],
      expTime: ''
    }
  },
  mounted() {
    const socket = initSocket(name)
    this.setUpSocket(socket)
    this.socket = socket

    this.$nextTick(replaceStreamSrc())

    axios.get('http://localhost:8888/items')
      .then((res) => {
        this.items = res.data
      }).catch((err) => {
        console.log(err)
      })
  },
  watch: {
    socket(val) {
      if (!val) {
        this.socket = initSocket(Cookies.get('username'))
        this.setUpSocket(this.socket)
      }
    },
    chatData() {
      // 滚动到最底部
      this.$nextTick(() => {
        this.$refs.chattable.bodyWrapper.scrollTop = this.$refs.chattable.bodyWrapper.scrollHeight;
      })
    }
  },
  filters: {
    formatDate(val) {
      const date = new Date(val)
      const y = date.getFullYear()
      const m = date.getMonth() + 1
      const d = date.getDate()
      const hh = date.getHours()
      const mm = date.getMinutes()
      const ss = date.getSeconds()
      return `${y}-${m}-${d} ${hh}:${mm}:${ss}`
    }
  },
  //計算totaltime的換算
  computed: {
    minutes: function () {
      const minutes = Math.floor(this.totalTime / 60);   //因為totaltime是秒數 所以整除60得到幾分鐘
      return this.padTime(minutes);
    },
    seconds: function () {
      const seconds = this.totalTime - (this.minutes * 60);  //原本的秒數-換算成分鐘的時間 = 秒數 (剩下的尾數)
      return this.padTime(seconds);
    }
  },
  methods: {
    setUpSocket(socket) {
      socket.onopen = () => {
        this.$message({
          type: 'success',
          message: '聊天室連結成功'
        })
      };
      socket.onclose = () => {
        this.$message({
          type: 'warning',
          message: '連結斷開'
        })
        this.socket = null
      }
      socket.onmessage = event => {
        console.log('@socket onmessage')
        console.log(event.data)
        let j = JSON.parse(event.data)


        this.receiveMsg(JSON.parse(event.data))
      }
    },
    onExit() {
      window.location.href = '/'
    },
    clearInput() {
      this.msg = ''
      this.$refs.input.focus()
    },
    //傳送訊息
    sendMessage() {
      if (!this.msg) {
        this.$message({
          type: 'warning',
          message: '不能發送空消息'
        })
        this.$refs.input.focus()
        return
      }
      cusMsg = JSON.stringify({
        type: "Msg",
        msg: this.msg
      })

      const req = JSON.stringify({
        msg: cusMsg
      })

      this.socket && (this.socket.send(req), this.msg = '', this.$refs.input.focus())
    },
    //商品上架
    addItem() {
      // var item = new Item(this.formitems.itemName, this.formitems.itemQuantity, this.formitems.aSetQuantity, this.formitems.itemPrice, this.formitems.addPrice1, this.formitems.addPrice2, this.formitems.addPrice3, this.formitems.itemText)

      // ++this.itemID
      // item.id = this.itemID
      var temp
      var GetBackID
      var postObject = {
        itemName: this.formitems.itemName,
        itemQuantity: parseInt(this.formitems.itemQuantity),
        aSetQuantity: parseInt(this.formitems.aSetQuantity),
        itemPrice: parseInt(this.formitems.itemPrice),
        addPrice1: parseInt(this.formitems.addPrice1),
        addPrice2: parseInt(this.formitems.addPrice2),
        addPrice3: parseInt(this.formitems.addPrice3),
        itemText: this.formitems.itemText
      }

      console.log(postObject)
      service({
        url: "https://api.itlive.nctu.me/login/Item",
        method: "POST",
        data: postObject,
        withCredentials: true
      })
        .then((res) => {
          console.log(res)
          if (res.data.data != "err") {
            GetBackID = res.data.data.itemID // this string

            let item = new Item(this.formitems.itemName, this.formitems.itemQuantity, this.formitems.aSetQuantity, this.formitems.itemPrice, this.formitems.addPrice1, this.formitems.addPrice2, this.formitems.addPrice3, this.formitems.itemText, GetBackID)

            temp = item

            const deucedMsg = JSON.stringify({ type: "addItem", "item": temp })

            const req = JSON.stringify({
              msg: deucedMsg
            })
            this.socket && (this.socket.send(req), this.clearItemInput(), this.dialogFormVisible = false)

          }

        }).catch((err) => {
          console.log(err)
        })
      // const deucedMsg = JSON.stringify({ type: "addItem", "item": item })

      // const req = JSON.stringify({
      //   msg: deucedMsg
      // })

      // this.socket && (this.socket.send(req), this.clearItemInput(), this.dialogFormVisible = false)
    },
    //接收所有訊息
    receiveMsg(data) {
      if (data.type == 0) {
        var inner = JSON.parse(data.text)

        if (inner.type == 'Msg') {                          //接收聊天訊息
          console.log('用法正確')
          data.text = inner.msg;
          this.chatData.push(data);
        } else if (inner.type == 'addItem') {               //接收登入商品訊息
          data.text = inner.item
          console.log(data.text.addPrice2)
          this.items.push(data.text)
          console.log(this.items)
        } else if (inner.type == 'startbid') {              //接收開始競標的訊息
          var nowTime = Math.floor(Date.now() / 1000)
          data.text = inner.biditemform
          this.expTime = inner.expTime
          this.heightprice = inner.biditemform.itemPrice

          this.totalTime = this.expTime - nowTime
          this.record = this.totalTime
          this.itembid.push(data.text)
          this.startTimer()
        } else if (inner.type == 'renewHeightPrice') {      //接收加價的訊息
          data.bidlist = inner.bidlist
          this.bidtable.unshift(data.bidlist)
          this.reduecebidtable()
          var temptest = this.bidtable[0]
          this.tempheight.unshift(temptest)
          this.reduecetempheght()

          data.text = inner.addPrice
          this.chatData.push(data)
        } else if (inner.type == 'timeReset') {
          this.endBid()
        } else if (inner.type == 'bidSuccess') {
          this.bidSusDialogVisible = true
        } else if (inner.type == 'Sustempconfirm') {
          this.confirmBidSus()
        } else if (inner.type == 'bidFail') {
          this.bidFailDialogVisible = true
        } else if (inner.type == 'Failtempconfirm') {
          this.confirmBidFail()
        } else if (inner.type == 'ConsoleMsg') {
          console.log(inner.type.msg)

        }
      }
      else {
        if (data.type != 1) {
          this.chatData.push(data);
        }
        console.log('@func receiveMs')
        console.log(data)
      }
    },
    //清除商品上架表單的資料
    clearItemInput() {
      this.formitems.itemName = ''
      this.formitems.itemQuantity = ''
      this.formitems.aSetQuantity = ''
      this.formitems.itemPrice = ''
      this.formitems.addPrice1 = ''
      this.formitems.addPrice2 = ''
      this.formitems.addPrice3 = ''
      this.formitems.itemText = ''
    },
    //選擇欲競標商品的資訊
    putbidItem() {
      this.biditemForm = this.biditemSelect.itemName
      this.biditemForm.totalTime = this.biditemSelect.totalTime
      var nowTime = Math.floor(Date.now() / 1000)
      this.expTime = nowTime + this.biditemForm.totalTime

      var biditemForm = new BiditemForm(this.biditemForm.itemName, this.biditemForm.itemQuantity, this.biditemForm.aSetQuantity, this.biditemForm.itemPrice, this.biditemForm.addPrice1,
        this.biditemForm.addPrice2, this.biditemForm.addPrice3, this.biditemForm.itemText, this.biditemForm.id)

      const deucedMsg = JSON.stringify({ "type": "startbid", "biditemform": biditemForm, "expTime": this.expTime })

      const req = JSON.stringify({
        msg: deucedMsg
      })

      this.socket && (this.socket.send(req))
    },
    //倒數計時開始
    startTimer() {
      this.timer = setInterval(() => this.countdown(), 1000); //1000ms = 1 second
      this.showselect = false
      this.showBidBut = false
      this.showConBut = false
      this.showPauBut = true
      this.showResBut = true
    },
    //暫停後繼續倒數
    continueTimer() {
      const deucedMsg = JSON.stringify({ "type": "timeContinue" })
      const req = JSON.stringify({ msg: deucedMsg })
      this.socket && (this.socket.send(req))

      this.timer = setInterval(() => this.countdown(), 1000); //1000ms = 1 second
      this.showselect = false
      this.showBidBut = false
      this.showConBut = false
      this.showPauBut = true
      this.showResBut = true
    },
    //暫停倒數
    stopTimer: function () {
      const deucedMsg = JSON.stringify({ "type": "timeStop", "totaltime": this.totalTime })
      const req = JSON.stringify({ msg: deucedMsg })
      this.socket && (this.socket.send(req))

      clearInterval(this.timer);
      this.timer = null;
      this.showConBut = true
      this.showPauBut = false
    },
    //重置倒數(重新選擇欲競標商品)
    resetTimer: function () {
      const deucedMsg = JSON.stringify({ "type": "timeReset" })
      const req = JSON.stringify({ msg: deucedMsg })
      this.socket && (this.socket.send(req))

      clearInterval(this.timer);
      this.timer = null;
      // this.clearItemBid()
      this.totalTime = this.record;
      this.showselect = true
      this.showBidBut = true
      this.showConBut = false
      this.showResBut = false
      this.showPauBut = false

    },
    //minutes,seconds格式轉換
    padTime: function (time) {
      return (time < 10 ? '0' : '') + time;
    },
    countdown: function () {
      if (this.totalTime > 0) {
        this.totalTime--;  //如果要倒數的時間>0 倒數的時間減1秒
      }
      else {
        this.timeendDialogVisible = true
        // this.endBid()
        // alert('競標結束!');
      }
    },
    timeendHandelClose() {
      this.$confirm('確認關閉')
        .then(_ => {
          alert("確認關閉後視為成交!")
          this.bidSuccess()
          done();
        })
        .catch(_ => { });
    },
    //倒數計時結束或重置倒數清除正在競標商品中的欄位
    endBid() {
      clearInterval(this.timer);
      this.timer = null
      this.itembid.splice(0, this.itembid.length)
      this.bidtable.splice(0, this.bidtable.length)
      this.tempheight = []
      this.totalTime = 0
      this.showselect = true
      this.showBidBut = true
      this.showConBut = false
      this.showResBut = false
      this.showPauBut = false
    },
    bidSuccess() {
      var biditemForm = new BiditemForm(this.biditemForm.itemName, this.biditemForm.itemQuantity, this.biditemForm.aSetQuantity, this.biditemForm.itemPrice, this.biditemForm.addPrice1,
        this.biditemForm.addPrice2, this.biditemForm.addPrice3, this.biditemForm.itemText, this.biditemForm.id)
      console.log(biditemForm)
      var bidtableitem = new Bidtableitem(this.bidtable)
      console.log(bidtableitem)      
      clearInterval(this.timer)
      this.timer = null
      var winbid = this.tempheight

      service({
        url: "https://api.itlive.nctu.me/login/ConfirmHigh",
        method: "POST",
        data: {
          item_id:this.biditemForm.id,
          highest_price:parseInt(winbid[0].heightprice),
        },
        withCredentials: true
      }).then(res=>{
        const deucedMsg = JSON.stringify({ "type": "bidSuccess", "winbid": winbid })
        const req = JSON.stringify({ msg: deucedMsg })
        this.socket && (this.socket.send(req))
        console.log(winbid[0].heightprice)
  
        this.itembid.splice(0, this.itembid.length)
      })


    },
    confirmBidSus() {
      this.bidSusDialogVisible = false
      this.endBid()
      this.bidtable.splice(0, this.bidtable.length)
      this.tempheight = []
      this.timeendDialogVisible = false
    },
    SusHandleClose(done) {
      this.$confirm('確認關閉')
        .then(_ => {
          this.Sustempconfirm()
          done();
        })
        .catch(_ => { });
    },
    Sustempconfirm() {
      const deucedMsg = JSON.stringify({ "type": "Sustempconfirm" })
      const req = JSON.stringify({ msg: deucedMsg })
      this.socket && (this.socket.send(req))
    },
    bidFail() {
      var biditemForm = new BiditemForm(this.biditemForm.itemName, this.biditemForm.itemQuantity, this.biditemForm.aSetQuantity, this.biditemForm.itemPrice, this.biditemForm.addPrice1,
        this.biditemForm.addPrice2, this.biditemForm.addPrice3, this.biditemForm.itemText)

      console.log(this.bidtable)
      var bidtableitem = new Bidtableitem(this.bidtable)
      console.log(bidtableitem)

      axios.post('http://localhost:8888/failorders', {
        biditemForm: biditemForm,
        bidtableitem: bidtableitem
      })
        .then((res) => {
          console.log(res)
        }).catch((err) => {
          console.log(err)
        })

      const deucedMsg = JSON.stringify({ "type": "bidFail" })
      const req = JSON.stringify({ msg: deucedMsg })
      this.socket && (this.socket.send(req))

      clearInterval(this.timer)
      this.timer = null
      this.itembid.splice(0, this.itembid.length)
    },
    confirmBidFail() {
      this.bidFailDialogVisible = false
      this.endBid()
      this.bidtable.splice(0, this.bidtable.length)
      this.tempheight = []
      this.timeendDialogVisible = false
    },
    FailHandleClose(done) {
      this.$confirm('確認關閉')
        .then(_ => {
          this.Failtempconfirm()
          done();
        })
        .catch(_ => { });
    },
    Failtempconfirm() {
      const deucedMsg = JSON.stringify({ "type": "Failtempconfirm" })
      const req = JSON.stringify({ msg: deucedMsg })
      this.socket && (this.socket.send(req))
    },
    reduecebidtable() {
      var pos = 5, n = 1
      this.bidtable.splice(pos, n)
    },
    reduecetempheght() {
      var pos = 1, n = 1
      this.tempheight.splice(pos, n)
    },
    tableRowStyle({ row, rowIndex }) {
      return 'background-color:#FFFFFF'
    },
  }
})