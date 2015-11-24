#!/bin/env python3


import time
import json

from TwitterAPI import TwitterAPI, TwitterError
from threading import Thread
from collections import deque
from flask import (Flask, redirect, url_for, send_file, render_template,
                   Response, request)


version = 3

consumer_key = '#'
consumer_secret = '#'
access_token_key = '#'
access_token_secret = '#'

searchstring = '#KoMa77'

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


def handle_twitter():

    api = TwitterAPI(
        consumer_key,
        consumer_secret,
        access_token_key,
        access_token_secret)

    while True:
        try:

            response = api.request('statuses/filter', {'track': searchstring})
            for item in response.get_iterator():
                if 'text' in item:
                    ring_buffer.append((item, time.time()))

        except TwitterError.TwitterError:
            print('Stream error!')
            time.sleep(10)
            print('Restart stream.')


if __name__ == "__main__":
    t = Thread(target=handle_twitter)
    t.start()
    app.run(host='0.0.0.0', port=5001)
