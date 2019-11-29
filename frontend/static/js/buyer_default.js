//const name = Cookies.get('username') ? Cookies.get('username') : (window.location.href = '/')
const channelUrl = location.pathname.substr(9)
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