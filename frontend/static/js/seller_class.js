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