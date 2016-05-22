var eventSource = new EventSource("/api/v1/stream.json");

function isElementInViewport(el) {
    var rect = el.getBoundingClientRect();

    return (
        rect.top >= 0 &&
            rect.left >= 0 &&
            rect.bottom <= (window.innerHeight || document.documentElement.clientHeight) &&
            rect.right <= (window.innerWidth || document.documentElement.clientWidth)
    );
}

function updateViewport() {
    document.querySelector(".foot-anchor").scrollIntoView()
}

function deleteOldTweets(container) {
    var tweets = getAllTweets(container);
    var length = tweets.length;
    if(length > 310) {
        var deleteCount = length - 300;
        for(i=0; i<deleteCount; i++) {
            container.removeChild(tweets[i]);
        }
    }
}

function getAllTweets(container) {
    return container.querySelectorAll(".tweet");
}

function getNewestTweet(container) {
    var tweets = getAllTweets(container);
    if(tweets.length == 0) {
        return false;
    }
    return tweets[tweets.length-1];
}

function insertTweet(tweet, tweetId, container) {
    var tweets = getAllTweets(container);
    if(tweets.length < 2) {
        container.appendChild(tweet);
        return
    }

    for(var i=tweets.length-1; i>=tweets.length-10; i--) {
        if(i < 0) {
            tweets[0].parentNode.insertBefore(tweet, tweets[0]);
            return;
        }
        if(tweets[i].getAttribute("data-tweetId") == tweetId) {
            return;
        }
        if(tweets[i].getAttribute("data-tweetId") < tweetId) {
            if(tweets[i].nextSibling) {
                tweets[i].parentNode.insertBefore(tweet, tweets[i].nextSibling);
            } else {
                tweets[i].parentNode.appendChild(tweet)
            }
            return;
        }
    }

}

function createTweet(data) {
    var tweetTemplate = document.querySelector(".tweet-template .tweet");
    var tweet = tweetTemplate.cloneNode(true);
    var photo = null;

    for(var i in data.entities.Media) {
        var media = data.entities.Media[i];
        if(media.Type == "photo") {
            photo = media;
            break;
        }
    }

    tweet.setAttribute("data-tweetId", data.id);
    tweet.setAttribute("data-tweetDate", data.created_at);

    tweet.querySelector(".message").textContent = data.text;

    var user = tweet.querySelector(".user")
    user.querySelector("a").href = "https://twitter.com/" + data.user.screen_name;
    user.querySelector(".name").textContent = data.user.name;
    user.querySelector(".screenname").textContent = data.user.screen_name;
    user.querySelector("img").src = data.user.profile_image_url_https;

    if(photo != null) {
        tweet.querySelector(".photo").style = "display: none;";
        tweet.querySelector(".photo").src = photo.Media_url_https;
    }

    return tweet
}

function tweetContainer1Handler(event) {
    var container = document.querySelector(".tweetContainer1");
    var isScrolledDown = isElementInViewport(document.querySelector(".foot-anchor"));

    var data = JSON.parse(event.data);
    var tweet = createTweet(data);

    insertTweet(tweet, data.id, container);

    if(isScrolledDown) {
        updateViewport();
        deleteOldTweets(container);
    }
}

function tweetContainer2Handler(event) {
    var container = document.querySelector(".tweetContainer2");

    var data = JSON.parse(event.data);
    var tweet = createTweet(data);

    insertTweet(tweet, data.id, container);
}

function photo_onload(e) {
    var isScrolledDown = isElementInViewport(document.querySelector(".foot-anchor"));
    e.removeAttribute("style")
    if(isScrolledDown) {
        updateViewport();
    }
}

function create_audio_element() {
    var o = document.createElement('audio');
    o.addEventListener('ended', function(){
        play_next();
    });
    o.addEventListener('error', function(){
        play_next();
    });
    return o;
}

function play_next() {
    if (audioElement.error != null) {
        audioElement = create_audio_element();
    }
    if (audioElement.paused) {
        if (to_play.length > 0) {
            audioElement.src = '/sounds/' + to_play.shift() + '.wav';
            audioElement.load();
            audioElement.volume = volume / 100.0;
            audioElement.play();
        }
    }
}

function soundHandler(event) {
    to_play.push(event.data);
    play_next();
}

var audioElement = create_audio_element();

var to_play = [];

eventSource.addEventListener("sound", soundHandler, false);
eventSource.addEventListener("tweet1", tweetContainer1Handler, false);
eventSource.addEventListener("tweet2", tweetContainer2Handler, false);
