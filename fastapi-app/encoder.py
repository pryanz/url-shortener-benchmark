import string

BASE62_ALPHABET = string.ascii_lowercase + string.ascii_uppercase + string.digits
BASE = len(BASE62_ALPHABET)

def encode(num : int) -> str :
    if num == 0:
        return BASE62_ALPHABET[0]
    arr =[]
    while num:
        num , rem = divmod(num,BASE)
        arr.append(BASE62_ALPHABET[rem]) 
    arr.reverse()
    return "".join(arr)

def decode(short_str :str) -> int:
    num = 0
    for char in short_str:
        num += num * BASE + BASE62_ALPHABET.index(char)
    
    return num
