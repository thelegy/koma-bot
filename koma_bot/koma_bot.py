import time
import json
import regex

from . twitter import TwitterStream
from . sse import SSE

from configparser import ConfigParser

from collections import deque
from flask import (Flask, redirect, url_for, send_file, render_template,
                   Response, request)

from gevent.wsgi import WSGIServer
from gevent.event import Event
from gevent import idle, Timeout, spawn, sleep


VERSION = 3
SOUND_TO_TRIGGERS = {'roman': ['roman', 'roh, man', ],
                     'zonk': ['game', 'spiel', 'lost', ],
                     'ananas:': ['ananas', ],
                     'orga': ['orga', ],
                     'ponny': ['jonny1', ],
                     'jonny': ['jonny2', ],
                     }

app = Flask(__name__, template_folder='')

ring_buffer = deque(maxlen=20)
triggers_to_sounds = {}


@app.route('/<path:path>')
def default(path):
    return redirect(url_for('soundboard'))


@app.route('/sound_board/<filename>.wav')
def sound_file(filename):
    return send_file(
        'sound_board/{}.wav'.format(filename))


@app.route('/')
def soundboard():
    return send_file('koma_bot.html')


@app.route('/script.js')
def script():
    return render_template(
        'script.js',
        version=VERSION)


@app.route('/stream/')
def stream():
    try:
        last_time = float(request.headers.get('Last-Event-ID'))
    except:
        last_time = time.time()

    def gen(last_time):
        while True:
            tweets = []
            for i in list(ring_buffer):
                if i[1] >= last_time:
                    tweets.append(i[0])

            if len(tweets) > 0:

                last_time = time.time()

                json_o = {}
                json_o['timestamp'] = last_time
                json_o['tweets'] = tweets
                json_o['action'] = []

                for i in tweets:
                    if i['retweeted']:
                        pass
                    json_o['action'].extend(actions_for(i['text']))

                yield json_o

            event.wait()

    sse = SSE(gen(last_time), 'timestamp')

    return Response(
        sse.output(),
        mimetype='text/event-stream')


def handle_twitter(item, the_time):
    ring_buffer.append((item, the_time))
    event.set()
    event.clear()


def actions_for(text):
    actions = []
    for match in app.trigger_regex.finditer(text, overlapped=True):
        actions.append(triggers_to_sounds[match.group(1).lower()])
    return actions


def create_app(testing=False):
    config = ConfigParser()
    config.read('config.ini')

    all_triggers = []
    for (sound, triggers) in SOUND_TO_TRIGGERS.items():
        all_triggers.extend([regex.escape(trigger) for trigger in triggers])
        triggers_to_sounds.update(dict([(trigger, sound)
                                        for trigger in triggers]))

    re = '(' + '|'.join(all_triggers) + ')'
    app.trigger_regex = regex.compile(re, regex.IGNORECASE)

    if testing:
        app.config['TESTING'] = True
    else:
        twitter = config['Twitter']
        twitter_stream = TwitterStream(twitter.get('consumer_key'),
                                       twitter.get('consumer_secret'),
                                       twitter.get('access_token_key'),
                                       twitter.get('access_token_secret'),
                                       track=twitter.get('track', '#KoMa77'),
                                       follow=twitter.get('follow'))
        twitter_stream.add_data_hook(handle_twitter)
        twitter_stream.start()

    return app


# This is TOO UGLY. Actually it is needed, because gevent has some kind of main
# loop and schedules things beforehand. It then waits for some external events
# and schedules the next "wakeup" in quite some time in the futute. When we
# receive an event from TwitterAPI, wich is running in a different thread I do
# not see a way no notify gevent to reschedule its things, so the event can
# actually be processed in the greenlets.
# This sleep loop keeps gevent from scheduling the next wakeup more than a
# quater second in the future yielding a drastically decreased response time.
# If you know a way, how to do this different in any way (that would be almost
# cerainly a cleaner one) please, PLEASE fix it!
def gevent_fix():
    while True:
        sleep(.25)


def create_server(testing=False):
    config = ConfigParser()
    config.read('config.ini')

    # not scalable; needs a better aproach
    global event
    event = Event()

    app = create_app(testing)

    http_server = WSGIServer((
        config.get('Web', 'bind_ip', fallback='0.0.0.0'),
        config.get('Web', 'port', fallback=5001)),
                             app)
    if not testing:
        spawn(gevent_fix)

        # this consumes the thread
        http_server.serve_forever()
