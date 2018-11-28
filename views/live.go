package views

var Str_live =`var server=location.href.split('/')[2]
var ws=new WebSocket('ws://'+server+'/live/ws')
ws.onmessage=function(e){
    console.log(e.data)
    location.href=location.href
}`
