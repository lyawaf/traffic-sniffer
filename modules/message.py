#!/usr/bin/env python3


class tcpMessage:
    """
    Basic message representation
    """

    direction = None
    data = None
    time = None

    def __init__(self, direction, data, time):
        self.direction = direction
        self.data = data
        self.time = time

    def __repr__(self):
        return "direction: {dir}, data: {data}, time: {time}".format(
            dir=self.direction, data=self.data, time=self.time
        )
