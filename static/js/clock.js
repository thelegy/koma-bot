function update_clock() {
    time = new Date();

    hour = ('0' + time.getHours()).slice(-2);
    minute = ('0' + time.getMinutes()).slice(-2);
    second = ('0' + time.getSeconds()).slice(-2);

    if(document.querySelector('.clock .hour').textContent != hour)
        document.querySelector('.clock .hour').textContent = hour;
    if(document.querySelector('.clock .minute').textContent != minute)
        document.querySelector('.clock .minute').textContent = minute;
    if(document.querySelector('.clock .second').textContent != second)
        document.querySelector('.clock .second').textContent = second;

    setTimeout(update_clock, 1001 - time.getMilliseconds());
}

update_clock();
