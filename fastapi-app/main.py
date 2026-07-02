from fastapi import FastAPI , Depends , HTTPException
from contextlib import asynccontextmanager
import database
import redis.asyncio as redis
import schemas
import encoder
from fastapi.responses import RedirectResponse

@asynccontextmanager
async def lifespan(app: FastAPI):
    database.redis_client = redis.Redis(host = 'redis-service' , port = 6379 , decode_responses = True)
    yield
    await database.redis_client.aclose()

app = FastAPI(lifespan=lifespan)

@app.post("/shorten", status_code = 201)
async def shorten(payload: schemas.URLShortenRequest, redis_client = Depends(database.get_redis)):
    
    url_id = await redis_client.incr("global:next_id")
    short_code = encoder.encode(url_id)
    await redis_client.set(f"url:{short_code}",str(payload.url))
    return {"short_code" : short_code}


@app.get("/{short_code}")
async def redirect(short_code:str , redis_client = Depends(database.get_redis)):
    long_url = await redis_client.get(f"url:{short_code}")

    if not long_url:
        raise HTTPException(status_code = 404 , detail="URL not found")

    else:
        return RedirectResponse(long_url , status_code=302)