var script_version = '{{ version }}';

function create_audio_element () {
    o = document.createElement('audio');
    o.addEventListener('ended', function(){
        play_next();
    });
    o.addEventListener('error', function(){
        play_next();
    });
    return o;
}

function play_next () {
    if (audioElement.error != null) {
        audioElement = create_audio_element();
    }
    if (audioElement.paused) {
        if (to_play.length > 0) {
            audioElement.src = '/sound_board/' + to_play.shift() + '.wav';
            audioElement.load();
            audioElement.play();
        }
    }
}

function update_sound (json) {
    var data = $.parseJSON(json)
    last_request = data['timestamp'];
    for (i=0; i<data['action'].length; i++) {
        to_play.push(data['action'][i]);
        play_next();
    }
    for (i = 0; i < data['tweets'].length; i++) {
        p = document.createElement('p');
        p.textContent = data['tweets'][i]['text'];
        $('body')[0].appendChild(p);
    }
}

var audioElement = create_audio_element();

var to_play = []

var eventSource = new EventSource("/stream/");
eventSource.onmessage = function(event) {
    update_sound(event.data)
};
