#!/bin/env python3


import time
import json

from twitter import TwitterStream

from configparser import ConfigParser

from collections import deque
from flask import (Flask, redirect, url_for, send_file, render_template,
                   Response, request)


version = 3

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
        version=version,
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
        for pos in range(len(i['text'])):
            sub = i['text'].lower()[pos:]
            if sub.startswith('roh, man'):
                json_o['action'].append('roman')
            if sub.startswith('roman'):
                json_o['action'].append('roman')
            if sub.startswith('game'):
                json_o['action'].append('zonk')
            if sub.startswith('spiel'):
                json_o['action'].append('zonk')
            if sub.startswith('lost'):
                json_o['action'].append('zonk')
            if sub.startswith('ananas'):
                json_o['action'].append('ananas')
            if sub.startswith('orga'):
                json_o['action'].append('orga')
            if sub.startswith('ponny'):
                json_o['action'].append('jonny1')
            if sub.startswith('jonny'):
                json_o['action'].append('jonny2')

    return Response(
        json.dumps(json_o),
        mimetype='application/json')


def handle_twitter(item, the_time):
    ring_buffer.append((item, the_time))
    print(item)


if __name__ == "__main__":
    config = ConfigParser()
    config.read('config.ini')
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
