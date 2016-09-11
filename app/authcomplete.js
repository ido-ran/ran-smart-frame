(function(){
    var cookies;

    function readCookie(name,c,C,i){
        if(cookies){ return cookies[name]; }

        c = document.cookie.split('; ');
        cookies = {};

        for(i=c.length-1; i>=0; i--){
           C = c[i].split('=');
           cookies[C[0]] = C[1];
        }

        return cookies[name];
    }

    window.readCookie = readCookie; // or expose it however you want
})();

var accessToken;
var parts = location.search.substring(1).split('&');

    for (var i = 0; i < parts.length; i++) {
        var nv = parts[i].split('=');
        if (!nv[0]) continue;
        if (nv[0] == 'accesstoken') {
          accessToken = nv[1];
          break;
        }
    }

if (!accessToken) {
  location = '/app/auth_fail.html'
} else {
  localStorage.setItem("accessToken", accessToken);

  var initiator = window.readCookie("auth_initiate") || 'console';
  if (initiator === 'console') {
    location = '/app/console.html'
  } else {
    location = '/app/frame.html'
  }
}
