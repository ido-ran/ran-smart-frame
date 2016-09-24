var accessToken = localStorage.getItem("accessToken");

var xhr = new XMLHttpRequest();
xhr.timeout = 8000; //60 * 1000; // timeout of 1 minute

var url = '/me';
xhr.open('GET', url);
xhr.setRequestHeader('Authorization', 'Bearer ' + accessToken);


xhr.ontimeout = function () {
  setErrorState();
};

xhr.onerror = function() {
  setErrorState();
}

xhr.onload = function(e) {
  if (this.status === 401) {
    // unauthorized
    console.log('response is 401');
    return;
  } else if (this.status !== 200) {
    console.log('response is error', this.status);
    return;
  };

  console.log('got response', this.response);
};

xhr.send();
