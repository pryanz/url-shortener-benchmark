from fastapi import FastAPI
from contextlib import asynccontextmanager
import database
import redis.asyncio as redis

@asynccontextmanager
async def lifespan(app: FastAPI):
    database.redis_client = redis.Redis(host = 'localhost' , port = 6379 , decode_responses = True)
    yield
    await database.redis_client.aclose()

app = FastAPI(lifespan=lifespan)