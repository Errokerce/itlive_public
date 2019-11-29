var express = require("express");
var cors = require("cors");
var fs = require("fs");
var app = express();

const corsOptions = {
  origin: [
    "https://websocket.itlive.nctu.me",
    "https://wh.itlive.nctu.me",
    "https://api.itlive.nctu.me",
    "https://googleeads.g.doubleclick.net"
  ]
};

app.use(cors(corsOptions));

// app.use("static", express.static("static"));
app.use("/static/css", express.static("static/css"));
app.use("/static/js", express.static("static/js"));
app.use("/static/img", express.static("static/img"));

app.get("/", function (req, res) {
    res.send(render("Front/page/indexV2.html", null));
});
app.get("/channel/:chid", function (req, res) {
  let user = req.headers["user-agent"].toLocaleLowerCase()
  
  console.log(user)
  if(user.includes('macintosh')&& (!user.includes('chrome')&&!user.includes('firefox')))
    res.send(render("Front/page/spam.html",null))
  else if(user.includes('iphone'))
    res.send(render("Front/page/spam.html",null))
  else
  res.send(render("Front/page/buyer.html", null));
});
app.get("/seller/:chid", function (req, res) {
  let user = req.headers["user-agent"].toLocaleLowerCase()
  console.log(user)
  if(user.includes('safari'))
    res.send(render("Front/page/spam.html",null))
  else if(user.includes('iphone'))
    res.send(render("Front/page/spam.html",null))
  else
  res.send(render("Front/page/seller.html", null));
});
app.get("/checkout", function (req, res) {
  res.send(render("Front/page/checkout.html", null));
});
app.get("/OrderManagement", function (req, res) {
  res.send(render("Front/page/OrderManagement.html", null));
});

app.get("/PhoneVerify", function (req, res) {
  res.send(render("Front/page/phoneVerify-vue.html", null));
});
app.get("/regIsSeller", function (req, res) {
  res.send(render("Front/page/regIsSeller.html", null));
});
app.get("/regProfile", function (req, res) {
  res.send(render("Front/page/regProfile.html", null));
});
app.get("/setting-product", function (req, res) {
  res.send(render("Front/page/setting-product.html", null));
});
app.get("/PrivacyPolicy", function (req, res) {
  res.send(render("Front/page/pp.html", null));
});



app.listen(3000, function () {
  console.log("frontend listening on port 3000!");
});

function render(filename, params) {
  var data = fs.readFileSync(filename, "utf8");
  for (var key in params) {
    data = data.replace("{" + key + "}", params[key]);
  }
  return data;
}
