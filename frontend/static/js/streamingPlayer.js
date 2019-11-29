function replaceStreamSrc(){


    service({
        url: 'https://api.itlive.nctu.me/login/GetStreamSrc/',
        method: 'GET',
    })
    document.getElementById("streamingPlayer").src=
}