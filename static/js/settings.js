var volume = 0;

function toggle_settings() {
    document.querySelector(".settings-button").classList.toggle("active");
    document.querySelector(".settings-box").classList.toggle("active");
}

function update_volume() {
    volume = document.querySelector(".volume-control input").value;
    document.cookie = "volume=" + volume;
    audioElement.volume = volume / 100.0;
}

function load_volume() {
    volume = document.cookie.replace(/(?:(?:^|.*;\s*)volume\s*\=\s*([^;]*).*$)|^.*$/, "$1")
    if(volume === "") {
        volume = 100;
    }
    volume = volume * 1
    document.querySelector(".volume-control input").value = volume;
}

document.addEventListener("click", function(event) {
    if(!event.target.closest(".settings-button") && !event.target.closest(".settings-box")) {
        document.querySelector(".settings-button").classList.remove("active");
        document.querySelector(".settings-box").classList.remove("active");
    }
});
