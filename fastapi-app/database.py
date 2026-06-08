import redis.asyncio as redis

redis_client = None 

def get_redis():
    return redis_client