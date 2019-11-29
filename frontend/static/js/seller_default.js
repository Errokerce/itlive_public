//const name = Cookies.get('username') ? Cookies.get('username') : (window.location.href = '/')
const channelUrl = location.pathname.substr(8)
function initSocket(username) {
  var url = 'wss://wh.itlive.nctu.me/s/' + channelUrl

  const socket = new WebSocket(url)
  return socket
}

function replaceStreamSrc() {
  service({
    url: `https://api.itlive.nctu.me/login/GetStreamSrc/${channelUrl}`,
    method: 'GET',
  }).then(resp => {

    document.getElementById("streamingPlayer").src = `https://www.youtube.com/embed/${resp.data}?playsinline=1`
  })
}
const service = axios.create()

      // 登入商品
      class Item {
        constructor(itemName, itemQuantity, aSetQuantity, itemPrice, addPrice1, addPrice2, addPrice3, itemText,id) {
            this.id = id
            this.itemName = itemName
            this.itemQuantity = itemQuantity
            this.aSetQuantity = aSetQuantity
            this.itemPrice = itemPrice
            this.addPrice1 = addPrice1
            this.addPrice2 = addPrice2
            this.addPrice3 = addPrice3
            this.itemText = itemText
        }
    }
    // 選擇的競標商品
    class BiditemForm {
        constructor(itemName, itemQuantity, aSetQuantity, itemPrice, addPrice1, addPrice2, addPrice3, itemText,id) {
            this.id = id
            this.itemName = itemName
            this.itemQuantity = itemQuantity
            this.aSetQuantity = aSetQuantity
            this.itemPrice = itemPrice
            this.addPrice1 = addPrice1
            this.addPrice2 = addPrice2
            this.addPrice3 = addPrice3
            this.itemText = itemText
        }
    }

    class Bidtableitem{
        constructor(bidtable){
            this.bidtable = bidtable
        }
    }