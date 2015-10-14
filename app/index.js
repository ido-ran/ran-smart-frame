var accessToken = localStorage.getItem("accessToken");

var dataSource;
var currMediaIndex = -1;
var imgElement;
var lastUpdateAt;
var imgContainerElement;
var loadingContainerElement;

window.onload = function() {
  imgContainerElement = document.getElementById("imgcontainer");
  imgElement = document.getElementById('mainimg');
  loadingContainerElement = document.getElementById('loadingcontainer');
  loadMedia();
}

function loadMedia() {
  var xhr = new XMLHttpRequest();

  var url = '/photos';
  xhr.open('GET', url);
  xhr.setRequestHeader('Authorization', 'Bearer ' + accessToken);
  xhr.onload = function(e) {
    if (this.status === 401) {
      // unauthorized
      location = 'https://ran-smart-frame.appspot.com/authorize';
      return;
    }

    lastUpdateAt = new Date();
    dataSource = JSON.parse(this.response);

    loadingContainerElement.style.display = 'none';
    imgContainerElement.style.display = 'block';

    showNextPhoto();
  }

  xhr.send();
}

function showNextPhoto() {
  if (checkMediaFreshness()) return;

  currMediaIndex++;
  if (currMediaIndex > dataSource.Media.length - 1) {
    currMediaIndex = 0;
  }

  imgElement.src = dataSource.Media[currMediaIndex].URL;
  setTimeout(showNextPhoto, 5000);
}

function checkMediaFreshness() {
  var now = new Date();
  var diffMs = (now - lastUpdateAt);
  var diffMin = diffMs / 1000 / 60;

  if (diffMin > 1) {
    loadMedia();
    return true;
  }

  return false;
}
