from koma_bot import actions_for
from app import WithApp


class TestTriggers(WithApp):
    def test_action(self):
        self.assertEquals(actions_for('roman'), ['roman'])
        self.assertEquals(actions_for('oh, man'), [])
