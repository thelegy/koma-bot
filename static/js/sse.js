var eventSource = new EventSource("/api/v1/stream.json");

function create_audio_element() {
    o = document.createElement('audio');
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
