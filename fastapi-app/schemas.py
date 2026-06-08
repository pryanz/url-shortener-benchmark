from pydantic import BaseModel , HttpUrl

class URLShortenRequest(BaseModel):
    url : HttpUrl
