from gevent import (spawn, queue, sleep)
from time import (time)
from json import (dumps)


class SSE(object):

    KEEP_ALIVE_PERIOD = 60

    class Message(object):

        def __init__(self, data):
            self.__data = data

        def __str__(self):
            if self.__data is None:
                # Keep-Alive signal
                return ':\n\n'
            if type(self.__data) is str:
                # String
                lines = self.__data.split('\n')
                prepended_lines = ['data: '+l for l in lines]
                return '\n'.join(prepended_lines) + '\n\n'
            # Everything else
            return 'data: ' + dumps(self.__data) + '\n\n'

    def __init__(self, iterator):
        self.__iterator = iterator
        self.__queue = queue.Queue()
        self.__last_message = time()
        self.__alive = True

    def __iterate(self):
        for item in self.__iterator:
            self.__queue.put(self.Message(item))
        self.__queue.put(StopIteration)

    def __keep_alive_signal(self):
        while self.__alive:
            seconds_passed = time() - self.__last_message
            if seconds_passed >= self.KEEP_ALIVE_PERIOD:
                self.__queue.put(self.Message(None))
            sleep(max(0, self.KEEP_ALIVE_PERIOD - seconds_passed))

    def output(self):
        spawn(self.__iterate)
        spawn(self.__keep_alive_signal)

        for message in self.__queue:
            self.__last_message = time()
            yield str(message)
