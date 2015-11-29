import time

from flask.ext.testing import TestCase

import koma_bot


class WithApp(TestCase):
    @classmethod
    def create_app(cls):
        return koma_bot.create_app(testing=True)

    def inject_tweet(text, timestamp, retweet=False):
        pass
