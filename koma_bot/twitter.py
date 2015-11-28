import time

from TwitterAPI import TwitterAPI, TwitterError
from threading import Thread

class TwitterStream:
    def __init__(self, consumer_key, consumer_secret, access_token_key,
                 access_token_secret, track=None, follow=None):

        self.__api = TwitterAPI(consumer_key,
                                consumer_secret,
                                access_token_key,
                                access_token_secret)
        self.__track = track
        self.__follow = follow

        self.__data_hooks = []
        self.__error_hooks = []

        self.__live = False
        
        self.__thread = Thread(target=self.__worker)

    def __worker(self):
        while self.__live:  # Restart stream when an error occures.
            try:
                
                response = self.__api.request('statuses/filter',
                                              {'track': self.__track,
                                               'follow': self.__follow})
                for item in response.get_iterator():
                    if 'text' in item:
                        recv_time = time.time()
                        if not self.__live:
                            return
                        for hook in self.__data_hooks:
                            try:
                                hook(item, recv_time)
                            except Exception as e:
                                print(e.with_traceback(None))
                        if not self.__live:
                            return
                                
            except TwitterError.TwitterError:
                for hook in self.__error_hooks:
                    try:
                        hook(item, recv_time)
                    except e:
                        print(e)
                # An error occured, restart to stay alive
                print('Stream error!')
                time.sleep(10)
                print('Restart stream.')

    def add_data_hook(self, hook):
        if callable(hook):
            self.__data_hooks.append(hook)

    def add_error_hook(self, hook):
        if callable(hook):
            self.__error_hooks.append(hook)

    def start(self):
        self.__live = True
        if not self.__thread.is_alive():
            self.__thread.start()

    def stop(self):
        self.__live = False
