import requests

resp = requests.post("[CLOUD_RUN_ENDPOINT]", files={'file': open('eight.png', 'rb')})

print(resp.json())
