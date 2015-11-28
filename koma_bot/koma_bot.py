#!/bin/env python3


import time
import json

from twitter import TwitterStream

from configparser import ConfigParser

from collections import deque
from flask import (Flask, redirect, url_for, send_file, render_template,
                   Response, request)


VERSION = 3

app = Flask(__name__, template_folder='')

ring_buffer = deque(maxlen=20)


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
        version=VERSION,
        time=time.time())


@app.route('/stream/')
def stream():
    def gen():
        for i in list(ring_buffer):
            if i[1] >= last_time:
                yield i[0]
    json_o = {}
    json_o['timestamp'] = time.time()
    try:
        last_time = float(request.args.get('last_request'))
    except:
        last_time = 0
    json_o['tweets'] = list(gen())
    json_o['action'] = []

    for i in json_o['tweets']:
        if i['retweeted']:
            pass
        json_o['action'].extend(actions_for(i['text']))

    return Response(
        json.dumps(json_o),
        mimetype='application/json')


def handle_twitter(item, the_time):
    ring_buffer.append((item, the_time))
    print(item)


def actions_for(text):
    actions = []
    for pos in range(len(text)):
        sub = text.lower()[pos:]
        if sub.startswith('roh, man'):
            actions.append('roman')
        if sub.startswith('roman'):
            actions.append('roman')
        if sub.startswith('game'):
            actions.append('zonk')
        if sub.startswith('spiel'):
            actions.append('zonk')
        if sub.startswith('lost'):
            actions.append('zonk')
        if sub.startswith('ananas'):
            actions.append('ananas')
        if sub.startswith('orga'):
            actions.append('orga')
        if sub.startswith('ponny'):
            actions.append('jonny1')
        if sub.startswith('jonny'):
            actions.append('jonny2')
    return actions


def create_app(testing=False):
    config = ConfigParser()
    config.read('config.ini')

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

        app.run(host=config.get('Web', 'bind_ip', fallback='0.0.0.0'),
                port=config.get('Web', 'port', fallback=5001))
    return app


if __name__ == "__main__":
    create_app()
