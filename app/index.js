var accessToken = localStorage.getItem("accessToken");

var xhr = new XMLHttpRequest();

var url = '/photos';
xhr.open('GET', url);
xhr.setRequestHeader('Authorization', 'Bearer ' + accessToken);
xhr.onload = function(e) {
  console.log('result', this.response);
}

console.log('requesting data');
xhr.send();
