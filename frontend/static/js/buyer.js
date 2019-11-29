const app = new Vue({
    delimiters: ["&{", "}"],
    el: "#app",
    data() {
        return {
            bidTableVisible: false,
            itemTableVisible: false,
            bidSusDialogVisible: false,
            bidFailDialogVisible: false,
            shopbagTableVisible: false,
            pickupInfoFormVisible: false,
            checkInfoTableVisible: false,
            paymentFrameVisble: false,
            paymentFrameOuter: "",
            socket: null,
            name: name,
            count: 20,
            msg: "",
            chatData: [],
            itemID: 0,
            formitems: {
                itemName: "",
                itemQuantity: "",
                aSetQuantity: "",
                itemPrice: "",
                addPrice1: "",
                addPrice2: "",
                addPrice3: "",
                itemText: ""
            },
            items: [], //登入的商品資訊
            formLabelWidth: "100px",
            itembid: [], //目前競標的商品
            totalTime: "", //目前的競標時間
            record: "", //紀錄選擇的競標時間
            timer: null,
            times: [
                {
                    id: "1",
                    totalTime: 30,
                    title: "00:30"
                },
                {
                    id: "2",
                    totalTime: 1 * 60,
                    title: "01:00"
                },
                {
                    id: "3",
                    totalTime: 1 * 90,
                    title: "01:30"
                },
                {
                    id: "4",
                    totalTime: 2 * 60,
                    title: "02:00"
                },
                {
                    id: "5",
                    totalTime: 1 * 150,
                    title: "延長時間"
                }
            ],
            heightprice: "", //目前最高價
            addPrice1: "", //第一加價選項
            addPrice2: "", //第二加價選項
            addPrice3: "", //第三加價選項
            showaddPrice1: false,
            showaddPrice2: false,
            showaddPrice3: false,
            apitable: [],
            bidtable: [], //競標者名字及最高價
            tempheight: [],
            expTime: "",
            Sustempname: name + "Sustempconfrim",
            Failtempname: name + "Failtempconfrim",
            shopbagitem: [],
            tempitembid: [],
            pickupInfo: {
                name: '',
                phonenumber: '',
                country: '',
                address: ''
            },
            checktable: [],
            fouls: []
        };
    },
    mounted() {
        const socket = initSocket(name);
        this.setUpSocket(socket);
        this.socket = socket;
        this.$nextTick(replaceStreamSrc());
        //   axios.get('http://localhost:8888/users',{
        //       params:{
        //           name:name
        //       }
        //   }).then((res) => {
        //         this.fouls = res.data
        //     }).catch((err) => {
        //         console.log(err)
        //     })
        // axios.get('http://localhost:8888/users', {
        //     params: {
        //         name: name
        //     }
        // }).then((res) => {
        //     this.fouls = res.data
        // }).catch((err) => {
        //     console.log(err)
        // })
    },
    watch: {
        socket(val) {
            if (!val) {
                this.socket = initSocket(Cookies.get("username"));
                this.setUpSocket(this.socket);
            }
        },
        chatData() {
            // 滚动到最底部
            this.$nextTick(() => {
                this.$refs.chattable.bodyWrapper.scrollTop = this.$refs.chattable.bodyWrapper.scrollHeight;
            });
        },
        heightprice(val, oldval) {
            if (val == oldval) {
                this.addEqual = true;
            } else {
                this.addEqual = false;
            }
        }
    },
    filters: {
        formatDate(val) {
            const date = new Date(val);
            const y = date.getFullYear();
            const m = date.getMonth() + 1;
            const d = date.getDate();
            const hh = date.getHours();
            const mm = date.getMinutes();
            const ss = date.getSeconds();
            return `${y}-${m}-${d} ${hh}:${mm}:${ss}`;
        }
    },
    //計算totaltime的換算
    computed: {
        minutes: function () {
            const minutes = Math.floor(this.totalTime / 60);
            return this.padTime(minutes);
        },
        seconds: function () {
            const seconds = this.totalTime - (this.minutes * 60);
            return this.padTime(seconds);
        }
    },
    methods: {
        setUpSocket(socket) {
            socket.onopen = () => {
                this.$message({
                    type: "success",
                    message: "聊天室連接成功"
                });
            };
            socket.onclose = () => {
                this.$message({
                    type: "warning",
                    message: "連結斷開"
                });
                this.socket = null;
            };
            socket.onmessage = event => {
                // console.log("@socket onmessage");
                // console.log(event.data);
                var j = JSON.parse(event.data);

                this.receiveMsg(JSON.parse(event.data));
            };
        },
        onExit() {
            window.location.href = "/";
        },
        clearInput() {
            this.msg = "";
            this.$refs.input.focus();
        },
        //傳送訊息
        sendMessage() {
            if (!this.msg) {
                this.$message({
                    type: "warning",
                    message: "不能發送空消息"
                });
                this.$refs.input.focus();
                return;
            }
            cusMsg = JSON.stringify({
                type: "Msg",
                msg: this.msg
            });

            const req = JSON.stringify({
                msg: cusMsg
            });

            this.socket &&
                (this.socket.send(req), (this.msg = ""), this.$refs.input.focus());
        },
        //接收所有訊息
        receiveMsg(data) {
            if (data.type == 0) {
                var inner = JSON.parse(data.text);

                if (inner.type == "Msg") {
                    //接收聊天訊息
                    // console.log("用法正確");
                    data.text = inner.msg;
                    this.chatData.push(data);
                } else if (inner.type == "addItem") {
                    //接收登入商品訊息
                    data.text = inner.item;
                    this.items.push(data.text);
                } else if (inner.type == "startbid") {
                    //接收開始競標的訊息
                    var nowTime = Math.floor(Date.now() / 1000);

                    data.text = inner.biditemform;
                    this.expTime = inner.expTime;
                    this.heightprice = parseInt(inner.biditemform.itemPrice);

                    this.totalTime = this.expTime - nowTime;
                    this.record = this.totalTime;
                    this.itembid.push(data.text);
                    this.tempitembid = data.text;
                    // console.log(this.tempitembid, typeof this.tempitembid);
                    this.addPrice1 = inner.biditemform.addPrice1;
                    this.showaddPrice1 = true;
                    if (data.text.addPrice2 !== "") {
                        this.showaddPrice2 = true;
                        this.addPrice2 = inner.biditemform.addPrice2;
                    }
                    if (data.text.addPrice3 !== "") {
                        this.showaddPrice3 = true;
                        this.addPrice3 = inner.biditemform.addPrice3;
                    }
                    this.startTimer();
                } else if (inner.type == "timeContinue") {
                    //接收繼續計時
                    this.continueTimer();
                } else if (inner.type == "timeStop") {
                    //接收暫停計時
                    this.totalTime = inner.totaltime;
                    // console.log(this.totalTime);
                    this.stopTimer();
                    this.showaddPrice1 = false;
                    this.showaddPrice2 = false;
                    this.showaddPrice3 = false;
                } else if (inner.type == "timeReset") {
                    //接收重置計時的訊息
                    this.endBid();
                    this.showaddPrice1 = false;
                    this.showaddPrice2 = false;
                    this.showaddPrice3 = false;
                } else if (inner.type == "renewHeightPrice") {
                    //接收加價的訊息
                    this.showaddEqual = true;

                    data.bidlist = inner.bidlist;
                    this.bidtable.unshift(data.bidlist);
                    this.reduecebidtable();
                    var temptest = this.bidtable[0];
                    this.tempheight.unshift(temptest);
                    this.reduecetempheght();

                    this.heightprice = parseInt(data.bidlist.heightprice);
                    // console.log(typeof this.heightprice);
                    this.addEqualPrice = this.heightprice;
                    data.text = inner.addPrice;
                    this.chatData.push(data);
                } else if (inner.type == "bidSuccess") {
                    this.bidSuccess();
                    this.showaddPrice1 = false;
                    this.showaddPrice2 = false;
                    this.showaddPrice3 = false;

                    data.text = inner.winbid;
                    this.tempitembid.heightprice = data.text[0].heightprice;
                    // console.log(this.tempitembid, typeof this.tempitembid);
                    // console.log(typeof data.text[0]);
                    // console.log(data.text[0].heightprice);
                    // console.log(typeof data.text);
                    if (data.text[0].name == this.name) {
                        this.shopbagitem.push(this.tempitembid);
                        // console.log(this.shopbagitem);
                    }
                } else if (inner.type == this.Sustempname) {
                    this.endBid();
                    this.confirmBidSus();
                } else if (inner.type == "bidFail") {
                    this.bidFail();
                    this.showaddPrice1 = false;
                    this.showaddPrice2 = false;
                    this.showaddPrice3 = false;
                } else if (inner.type == this.Failtempname) {
                    this.endBid();
                    this.confirmBidFail();
                } else if (inner.type == 'ConsoleMsg') {
                    console.log(inner.type.msg)

                }
            }
            else {
                if (data.type != 1) {
                    this.chatData.push(data);
                }
                console.log("@func receiveMs");
                console.log(data);
            }
        },
        //清除商品上架表單的資料
        clearItemInput() {
            this.formitems.itemName = "";
            this.formitems.itemQuantity = "";
            this.formitems.aSetQuantity = "";
            this.formitems.itemPrice = "";
            this.formitems.addPrice1 = "";
            this.formitems.addPrice2 = "";
            this.formitems.addPrice3 = "";
            this.formitems.itemText = "";
        },
        //倒數計時開始
        startTimer() {
            this.timer = setInterval(() => this.countdown(), 1000); //1000ms = 1 second
        },
        //暫停後繼續倒數
        continueTimer() {
            this.timer = setInterval(() => this.countdown(), 1000); //1000ms = 1 second
        },
        //暫停倒數
        stopTimer: function () {
            clearInterval(this.timer);
            this.timer = null;
        },
        //重置倒數(重新選擇欲競標商品)
        resetTimer: function () {
            this.totalTime = this.record;
            clearInterval(this.timer);
            this.timer = null;
            this.clearItemBid();
        },
        //倒數計時結束或重置倒數清除正在競標商品中的欄位
        endBid() {
            clearInterval(this.timer);
            this.timer = null;
            this.itembid.splice(0, this.itembid.length);
            this.totalTime = 0;
        },
        //minutes,seconds格式轉換
        padTime: function (time) {
            return (time < 10 ? "0" : "") + time;
        },
        countdown: function () {
            if (this.totalTime > 0) {
                this.totalTime--; //如果要倒數的時間>0 倒數的時間減1秒
            } else {
                this.endBid();
                alert('競標結束等待賣家決議!');
            }
        },
        //第一加價選項
        addFucPrice1() {
            // if (this.fouls.length == 0) {
            //     alert('您尚未登入無法競標!');
            //     return
            // } else if (this.fouls[0].foul == 3) {
            //     alert('您棄標次數已達上限無法競標!');
            //     return
            // }

            this.heightprice += parseInt(this.addPrice1);
            var renewheightPrice = this.heightprice.toString();
            var addPrice = "+" + this.addPrice1.toString();

            var bidlist = new bidList(this.name, renewheightPrice);
            const deucedMsg = JSON.stringify({
                type: "renewHeightPrice",
                bidlist: bidlist,
                addPrice: addPrice
            });
            const req = JSON.stringify({
                msg: deucedMsg
            });

            this.socket && this.socket.send(req);
        },
        //第二加價選項
        addFucPrice2() {
            // if (this.fouls[0].foul == 3) {
            //     alert("您棄標次數已達上限無法競標!");
            //     return;
            // }
            this.heightprice += parseInt(this.addPrice2);
            var renewheightPrice = this.heightprice.toString();
            var addPrice = "+" + this.addPrice2.toString();

            var bidlist = new bidList(this.name, renewheightPrice);
            const deucedMsg = JSON.stringify({
                type: "renewHeightPrice",
                bidlist: bidlist,
                addPrice: addPrice
            });
            const req = JSON.stringify({
                msg: deucedMsg
            });

            this.socket && this.socket.send(req);
        },
        //第三加價選項
        addFucPrice3() {
            // if (this.fouls[0].foul == 3) {
            //     alert("您棄標次數已達上限無法競標!");
            //     return;
            // }
            this.heightprice += parseInt(this.addPrice3);
            var renewheightPrice = this.heightprice.toString();
            var addPrice = "+" + this.addPrice3.toString();

            var bidlist = new bidList(this.name, renewheightPrice);
            const deucedMsg = JSON.stringify({
                type: "renewHeightPrice",
                bidlist: bidlist,
                addPrice: addPrice
            });
            const req = JSON.stringify({
                msg: deucedMsg
            });

            this.socket && this.socket.send(req);
        },
        //bidsuccess
        bidSuccess() {
            this.bidSusDialogVisible = true;
            clearInterval(this.timer);
            this.timer = null;
        },
        confirmBidSus() {
            this.bidSusDialogVisible = false;
            this.itembid.splice(0, this.itembid.length);
            this.totalTime = 0;
            this.bidtable.splice(0, this.bidtable.length);
            this.tempheight = [];
        },
        SushandleClose(done) {
            this.$confirm("確認關閉")
                .then(_ => {
                    this.Sustempconfirm();
                    done();
                })
                .catch(_ => { });
        },
        Sustempconfirm() {
            const deucedMsg = JSON.stringify({ type: this.Sustempname });
            const req = JSON.stringify({ msg: deucedMsg });
            this.socket && this.socket.send(req);
        },
        //bidfail
        bidFail() {
            this.bidFailDialogVisible = true;
            clearInterval(this.timer);
            this.timer = null;
        },
        confirmBidFail() {
            this.bidFailDialogVisible = false;
            this.itembid.splice(0, this.itembid.length);
            this.totalTime = 0;
            this.bidtable.splice(0, this.bidtable.length);
            this.tempheight = [];
        },
        FailhandleClose(done) {
            this.$confirm("確認關閉")
                .then(_ => {
                    this.Failtempconfirm();
                    done();
                })
                .catch(_ => { });
        },
        Failtempconfirm() {
            const deucedMsg = JSON.stringify({ type: this.Failtempname });
            const req = JSON.stringify({ msg: deucedMsg });
            this.socket && this.socket.send(req);
        },
        apitest() {
            service({
                url: "/api/itemList",
                method: "get"
            }).then(res => {
                var inner = res.data;
                this.apitable = inner.items;
                console.log(inner);
                console.log(typeof inner);
            });
        },
        reduecebidtable() {
            var pos = 5,
                n = 1;
            var removeedItems = this.bidtable.splice(pos, n);
        },
        reduecetempheght() {
            var pos = 1,
                n = 1;
            var removeedItems = this.tempheight.splice(pos, n);
        },
        sortHighest() {
            var sortHigh = parseInt(this.bidtable.heightprice);
            this.bidtable.sort((a, b) => (a.sortHigh < b.sortHigh ? 1 : -1));
        },
        //聊天訊息背景顏色
        tableRowStyle({ row, rowIndex }) {
            return "background-color:#FFFFFF";
        },
        shopbag() {
            this.shopbagTableVisible = true;
        },
        checkout() { },
        cancelCheckout() {
            this.shopbagTableVisible = false;
        },
        writeInfo() {
            this.shopbagTableVisible = false;
            this.pickupInfoFormVisible = true;
        },
        backtoWriteInfo() {
            this.pickupInfoFormVisible = false;
            this.shopbagTableVisible = true;
        },
        checkInfo() {
            console.log(this.tempitembid)
            console.log(this.shopbagitem)
            // this.pickupInfo.photo = this.tempitembid.photo;
            // this.pickupInfo.heightprice = this.tempitembid.heightprice;
            // this.pickupInfo.itemName = this.tempitembid.itemName;
            // this.pickupInfo.itemQuantity = this.tempitembid.itemQuantity;
            this.checktable.push(this.pickupInfo)
            this.shopbagitem.forEach(element => {
                this.checktable.push(element);
            });

            this.pickupInfoFormVisible = false;
            this.checkInfoTableVisible = true;
        },
        backtoPickupInfo() {
            this.pickupInfoFormVisible = true;
            this.checkInfoTableVisible = false;
        },
        paymentStep() {
            // 接綠界Function
            console.log(this.checktable)
            var checktable = new checkTable(this.checktable)


            var postBody = {
                shipping_info: {
                    country: checktable.country,
                    address: checktable.address,
                    name: checktable.name,
                    phone_number: checktable.phonenumber,
                },
                items: []
            }
            this.shopbagitem.forEach(ele => {
                postBody.items.push(ele)
            })


            console.log(postBody)
            // axios.post('http://localhost:8888/shopcartorder', checktable)
            //     .then((res) => {
            //         console.log(res)
            //     }).catch((err) => {
            //         console.log(err)
            //     })


            service({
                url: "https://api.itlive.nctu.me/login/ConfirmOrder",
                method: "POST",
                data: postBody,
                withCredentials: true
            }).then(resp => {
                console.log(resp.data)
                if (resp.data.state === "ok") {
                    console.log(resp.data.spToken);
                    this.checkInfoTableVisible = false;
                    this.paymentFrameVisble = true;
                    var srcUrl = 'https://payment-stage.ecpay.com.tw/SP/SPCheckOut?MerchantID=2000132&SPToken=' + resp.data.data.spToken + '&PaymentType=CREDIT'
                    this.paymentFrameOuter = `<iframe src="${srcUrl}" id="paymentFrame" height="600px" width="100%" style="border:0;"></iframe>`
                    // document.getElementById("paymentFrame").src=srcUrl

                }
            })



            // service({
            //     url: 'https://api.itlive.nctu.me/login/FakePayment',
            //     method: 'GET'
            // }).then(resp=>{
            //     if(resp.data.state==="ok"){
            //         console.log(resp.data.spToken);
            //         this.checkInfoTableVisible=false;
            //         this.paymentFrameVisble=true;
            //         var srcUrl='https://payment-stage.ecpay.com.tw/SP/SPCheckOut?MerchantID=2000132&SPToken='+resp.data.spToken+'&PaymentType=CREDIT'
            //         this.paymentFrameOuter=`<iframe src="${srcUrl}" id="paymentFrame" height="600px" width="100%" style="border:0;"></iframe>`
            //         // document.getElementById("paymentFrame").src=srcUrl

            //     }
            // })

        }
    }
});
