return redis.call('zrevrange', KEYS[1], 0, -1, 'WITHSCORES')