var accessToken = localStorage.getItem("accessToken");

var STATES = {
  INIT: 'init',
  LOADING: 'loading',
  LOADED: 'loaded',
  ERROR: 'error'
};

var dataSource;
var currMediaIndex = -1;
var imgElement;
var lastUpdateAt;
var state = STATES.INIT;
var errorCount = 0;

var imgContainerElement;
var loadingContainerElement;

window.onload = function() {
  imgContainerElement = document.getElementById("imgcontainer");
  imgElement = document.getElementById('mainimg');
  loadingContainerElement = document.getElementById('loadingcontainer');
  loadMedia();
}

function setState(newState) {
  loadingContainerElement.textContent = newState;
  state = newState;
}

function loadMedia() {
  if (state === STATES.LOADING) return;

  var isAlreadyLoaded = (state === STATES.LOADED);
  setState(STATES.LOADING);

  var xhr = new XMLHttpRequest();

  var url = '/photos';
  xhr.open('GET', url);
  xhr.setRequestHeader('Authorization', 'Bearer ' + accessToken);
  xhr.onload = function(e) {
    if (this.status === 401) {
      // unauthorized
      location = 'https://ran-smart-frame.appspot.com/authorize';
      return;
    } else if (this.status !== 200) {
      setState(STATES.ERROR);
      // Try to reload media in 2 min
      errorCount++;

      // Wait at most 10 minutes for retry
      var retryDelay = Math.min(10, errorCount) * 60 * 1000;
      setTimeout(loadMedia, retryDelay);
      return;
    }

    errorCount = 0;
    lastUpdateAt = new Date();
    dataSource = JSON.parse(this.response);

    setState(STATES.LOADED);

    if (!isAlreadyLoaded) {
      //debug - loadingContainerElement.style.display = 'none';
      imgContainerElement.style.display = 'block';

      showNextPhoto();
    }
  }

  xhr.send();
}

function showNextPhoto() {
  checkMediaFreshness();

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

  if (diffMin > 0.5) {
    loadMedia();
    return true;
  }

  return false;
}
