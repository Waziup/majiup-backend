from random import randint
import time
import requests

print("starting...")


def get_token():
    url = "http://wazigate.local/auth/token"

    username: str = "admin"
    password: str = "loragateway"

    response = requests.post(
        url=url,
        json={
            "username":username,
            "password":password
        },
        headers={
            "Content-Type":"application/json"
        }
    )

    token: str = response.text
    return token.replace('"','')

def post(val: int, token: str, iteration: int):
    # response = requests.post("http://localhost/devices/667885d0d6980465dbd5ec9f/sensors/667885e1d6980465dbd5eca0/value", headers={
    #     "Authorization": f"Bearer {token}", "Content-Type":"text/plain"
    # }, data=str(val))
    
    # print(token)
    response = requests.post("http://localhost/devices/667885d0d6980465dbd5ec9f/sensors/66950daad6980457ec07463f/value", data=str(val))

    if (response.status_code==200):
        print(f"[{iteration}]Success  --->   {val}")
        # res = response.json()
        # print(res)

    # elif response.status_code == 401:
    #     token = get_token()
    #     post()

    else:
        print("Failed to make request")
        print("Status Code: ", response.status_code)

values: int = 2900

for i in range(0, values):
    val: int = randint(50,700)
    post(val,token='', iteration=i)
    time.sleep(0.05)

def post_meta(token: str, data):
    response = requests.post("http://wazigate.local/devices/66afee0d68f3190a01e8329a/sensors/temperatureSensor_0/meta", headers={
        "Authorization": f"Bearer {token}", "Content-Type":"aplication/json"
    }, json=data)
    
    if (response.status_code==200):
        print(f"[ Success ]")
        # res = response.json()
        # print(res)

    # elif response.status_code == 401:
    #     token = get_token()
    #     post()

    else:
        print("Failed to make request")
        print("Status Code: ", response.status_code)

def get_meta(token: str):
    url = "http://wazigate.local/device/meta"
    response = requests.get(url=url, headers={
        "Authorization": f"Bearer {token}", "Content-Type":"application/json"
    })
    
    return response.json()

# token = get_token()

# meta = get_meta(token=token)

# print(meta)
# meta = {'createdBy': 'wazigate-lora', 'critical_max': 89, 'critical_min': 40, 'kind': 'WaterLevel', 'quantity': 'WaterLevel', 'unit': 'Millimetre', 'units': ''}
# post_meta(token=token, data=meta)

   
# for i in range(0,values):
#     val:int = randint(300,500)
#     post(val, token=get_token(), iteration=i)
#     # post(val, iteration=i, token="")
#     time.sleep(5)