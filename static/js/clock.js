function updateClock() {
    var time = new Date();

    var hour = ('0' + time.getHours()).slice(-2);
    var minute = ('0' + time.getMinutes()).slice(-2);
    var second = ('0' + time.getSeconds()).slice(-2);
    var date = formatDate(time)

    if(document.querySelector('.clock .hour').textContent != hour)
        document.querySelector('.clock .hour').textContent = hour;
    if(document.querySelector('.clock .minute').textContent != minute)
        document.querySelector('.clock .minute').textContent = minute;
    if(document.querySelector('.clock .second').textContent != second)
        document.querySelector('.clock .second').textContent = second;
    if(document.querySelector('.clock .date').textContent != date)
        document.querySelector('.clock .date').textContent = date;

    setTimeout(updateClock, 1001 - time.getMilliseconds());

    updateTime();
    updateTimetable(time);
}

function formatDate(time) {
    var date = "";
    var dateFormat = document.querySelector(".clock-format");
    switch(time.getDay()) {
    case 0:
        date = dateFormat.querySelector(".weekday .sunday").textContent;
        break;
    case 1:
        date = dateFormat.querySelector(".weekday .monday").textContent;
        break;
    case 2:
        date = dateFormat.querySelector(".weekday .tuesday").textContent;
        break;
    case 3:
        date = dateFormat.querySelector(".weekday .wednesday").textContent;
        break;
    case 4:
        date = dateFormat.querySelector(".weekday .thursday").textContent;
        break;
    case 5:
        date = dateFormat.querySelector(".weekday .friday").textContent;
        break;
    case 6:
        date = dateFormat.querySelector(".weekday .saturday").textContent;
        break;
    }
    date += ", " + dateFormat.querySelector(".the").textContent + " "
    date +=(time.getYear() - 100) + '/' + (time.getMonth()+1) + '/' + time.getDate();
    return date;
}

function formatTimeDiff(time, format, date) {

        // Seconds
        if(time <= 1) {
            return format.second.singular;
        }
        if(time < 60) {
            return format.second.before + time + format.second.after;
        }

        time = Math.floor(time / 60);

        // Minutes
        if(time <= 1) {
            return format.minute.singular;
        }
        if(time < 60) {
            return format.minute.before + time + format.minute.after;
        }

        time = Math.floor(time / 60);

        // Hours
        if(time <= 1) {
            return format.hour.singular;
        }
        if(time < 24) {
            return format.hour.before + time + format.hour.after;
        }

        time = Math.floor(time / 24);

        // Days
        if(time <= 1) {
            return format.day.singular;
        }
        if(time < 7) {
            return format.day.before + time + format.day.after;
        }

    return (date.getYear() - 100) + '/' + date.getMonth() + '/' + date.getDate();

}

function updateTimetable(time) {
    var rows = document.querySelectorAll(".timetable tr")

    for(var i=0; i<rows.length; i++) {
        var start_time = new Date(rows[i].getAttribute("data-start"));
        if(start_time == "Invalid Date") {
            continue;
        }
        if(start_time-time > 3600000) {
            rows[i].classList.remove("active");
            continue;
        }
        var end_time = new Date(rows[i].getAttribute("data-end"));
        if(end_time == "Invalid Date") {
            continue;
        }
        if(end_time-time <= 0) {
            rows[i].classList.remove("active");
            continue;
        }
        rows[i].classList.add("active");
    }
}

function updateTime() {
    var tweets = document.querySelectorAll(".tweets .tweet");
    var time = new Date().getTime();

    var time_format = document.querySelector(".tweet-template .time-format");

    var format = {};

    format.second = {};
    format.second.singular = time_format.querySelector(".second .singular").textContent;
    format.second.before = time_format.querySelector(".second .before").textContent;
    format.second.after = time_format.querySelector(".second .after").textContent;

    format.minute = {};
    format.minute.singular = time_format.querySelector(".minute .singular").textContent;
    format.minute.before = time_format.querySelector(".minute .before").textContent;
    format.minute.after = time_format.querySelector(".minute .after").textContent;

    format.hour = {};
    format.hour.singular = time_format.querySelector(".hour .singular").textContent;
    format.hour.before = time_format.querySelector(".hour .before").textContent;
    format.hour.after = time_format.querySelector(".hour .after").textContent;

    format.day = {};
    format.day.singular = time_format.querySelector(".day .singular").textContent;
    format.day.before = time_format.querySelector(".day .before").textContent;
    format.day.after = time_format.querySelector(".day .after").textContent;

    for(var i=0; i<tweets.length; i++) {
        var tweetTime = Date.parse(tweets[i].getAttribute("data-tweetDate"));
        var tweetTimeDiff = Math.floor((time - tweetTime) / 1000);

        var tweetTimeString = formatTimeDiff(tweetTimeDiff, format, new Date(tweetTime));

        if(tweets[i].querySelector('.time').textContent !== tweetTimeString) {
            tweets[i].querySelector('.time').textContent = tweetTimeString;
        }

    }
}

updateClock();
