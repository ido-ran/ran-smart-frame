var accessToken = localStorage.getItem("accessToken");

var STATES = {
  INIT: 'init',
  LOADING: 'loading',
  LOADED: 'loaded',
  ERROR: 'error'
};

var CONSTS = {
  SHOW_PHOTO_DURATION_IN_SECONDS: 20,
  REFRESH_MEDIA_INTERVAL_IN_MINUTES: 5
};

var dataSource;
var mediaToShow;
var currMediaIndex = -1;
var imgElement;
var ribbonElement;
var lastUpdateAt;
var state = STATES.INIT;
var errorCount = 0;

var imgContainerElement;
var loadingContainerElement;

window.onload = function() {
  imgContainerElement = document.getElementById("imgcontainer");
  imgElement = document.getElementById('mainimg');
  loadingContainerElement = document.getElementById('loadingcontainer');
  ribbonElement = document.getElementById('ribbon');
  loadMedia();
}

function setState(newState) {
  loadingContainerElement.textContent = newState;
  state = newState;
}

function setErrorState() {
  setState(STATES.ERROR);
  // Try to reload media in 2 min
  errorCount++;

  // Wait at most 10 minutes for retry
  var retryDelay = Math.min(10, errorCount) * 60 * 1000;
  setTimeout(loadMedia, retryDelay);
}

function loadMedia() {

  if (state === STATES.LOADING) return;

  var isAlreadyLoaded = (state === STATES.LOADED);
  setState(STATES.LOADING);

  var xhr = new XMLHttpRequest();
  xhr.timeout = 8000; //60 * 1000; // timeout of 1 minute

  var url = '/photos';
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
      location = 'https://ran-smart-frame.appspot.com/authorize';
      return;
    } else if (this.status !== 200) {
      setErrorState();
      return;
    };

    errorCount = 0;
    lastUpdateAt = new Date();
    dataSource = JSON.parse(this.response);
    mediaToShow = dataSource.Media.slice(); // make a copy

    // TODO: we should not suffle if the dataSource is the same as the last one we got to ensure we show all the pictures before suffle again.
    shuffle(mediaToShow);

    setState(STATES.LOADED);

    if (!isAlreadyLoaded) {
      //debug - loadingContainerElement.style.display = 'none';
      imgContainerElement.style.display = 'block';

      showNextPhoto();
    }
  };

  xhr.send();
}

function showNextPhoto() {
  checkMediaFreshness();

  currMediaIndex++;
  if (currMediaIndex > mediaToShow.length - 1) {
    currMediaIndex = 0;
  }

  imgElement.src = mediaToShow[currMediaIndex].URL;

  var timePicTaken = new Date(parseInt(mediaToShow[currMediaIndex].Timestamp));
  var pointInTime = new Date("Jan 27 2015 3:00");
  var duration = moment.duration(timePicTaken - pointInTime);

  var durationString = '';
  if (duration.years() > 0) {
    durationString += duration.years() + 'y ';
  }
  if (duration.months() > 0) {
     durationString += duration.months() + 'm ';
  }
  if (duration.days() > 0) {
    durationString += duration.days() + 'd';
  }
  ribbonElement.innerText = durationString;

  setTimeout(showNextPhoto, CONSTS.SHOW_PHOTO_DURATION_IN_SECONDS * 1000);
}

function checkMediaFreshness() {
  var now = new Date();
  var diffMs = (now - lastUpdateAt);
  var diffMin = diffMs / 1000 / 60;

  if (diffMin > CONSTS.REFRESH_MEDIA_INTERVAL_IN_MINUTES) {
    loadMedia();
    return true;
  }

  return false;
}

function shuffle(array) {
  var currentIndex = array.length, temporaryValue, randomIndex;

  // While there remain elements to shuffle...
  while (0 !== currentIndex) {

    // Pick a remaining element...
    randomIndex = Math.floor(Math.random() * currentIndex);
    currentIndex -= 1;

    // And swap it with the current element.
    temporaryValue = array[currentIndex];
    array[currentIndex] = array[randomIndex];
    array[randomIndex] = temporaryValue;
  }

  return array;
}
