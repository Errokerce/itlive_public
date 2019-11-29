 // 登入商品
 class Item {
  constructor(itemName, itemQuantity, aSetQuantity, itemPrice, addPrice1, addPrice2, addPrice3, itemText) {
      this.id = 0
      this.itemName = itemName
      this.itemQuantity = itemQuantity
      this.aSetQuantity = aSetQuantity
      this.itemPrice = itemPrice
      this.addPrice1 = addPrice1
      this.addPrice2 = addPrice2
      this.addPrice2 = addPrice3
      this.itemText = itemText
  }
}
// 選擇的競標商品
class SelectedItem {
  constructor(itemName, itemQuantity, aSetQuantity, itemPrice, addPrice1, addPrice2, addPrice3, itemText) {
      this.id = 0
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
//競標名單
class bidList {
  constructor(name, heightprice, addPrice) {
      this.name = name
      this.heightprice = heightprice
      this.addPrice = addPrice
  }
}


class checkTable {
  constructor(checktable) {
    this.checktable = checktable
  }
}