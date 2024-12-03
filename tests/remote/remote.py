from random import randint
import time
import requests

# manually set cookie, retrieved from remote...
# retrieved_cookie =  "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJMbGE5X3BNeFVwc1QxTXpsM2dsdkRFTFZIa256OWZDb191Q3JheDh1dkJ3In0.eyJqdGkiOiIyYzNhN2EyYy0zYjUxLTQwODQtOWE5Mi0zYzM3YTBjMTI4YTciLCJleHAiOjE3MjA0MzM3MjcsIm5iZiI6MCwiaWF0IjoxNzIwNDMzMTI3LCJpc3MiOiJodHRwczovL2tleWNsb2FrLndheml1cC5pby9hdXRoL3JlYWxtcy93YXppdXAiLCJhdWQiOiJhcGktc2VydmVyIiwic3ViIjoiYjgzODM1YTQtMzM1ZS00ZWQzLWJkYWEtOTdmYzk0Nzk0OGVlIiwidHlwIjoiQmVhcmVyIiwiYXpwIjoiYXBpLXNlcnZlciIsImF1dGhfdGltZSI6MTcyMDQzMzEyNywic2Vzc2lvbl9zdGF0ZSI6IjNhYWQyYzc2LWFmNmMtNDc4Mi1iZjMxLTBkZmUxZTk1MTQxOSIsImFjciI6IjEiLCJhbGxvd2VkLW9yaWdpbnMiOlsiKiJdLCJyZWFsbV9hY2Nlc3MiOnsicm9sZXMiOlsib2ZmbGluZV9hY2Nlc3MiLCJyZWdpc3RlcmVkX3VzZXIiLCJ1bWFfYXV0aG9yaXphdGlvbiJdfSwicmVzb3VyY2VfYWNjZXNzIjp7InJlYWxtLW1hbmFnZW1lbnQiOnsicm9sZXMiOlsibWFuYWdlLXVzZXJzIiwidmlldy11c2VycyIsInF1ZXJ5LWdyb3VwcyIsInF1ZXJ5LXVzZXJzIl19LCJhY2NvdW50Ijp7InJvbGVzIjpbIm1hbmFnZS1hY2NvdW50IiwibWFuYWdlLWFjY291bnQtbGlua3MiLCJ2aWV3LXByb2ZpbGUiXX19LCJzY29wZSI6IiIsInNtc19jcmVkaXQiOiIxMDAiLCJuYW1lIjoiSm9zZXBoIE11c3lhIiwicHJlZmVycmVkX3VzZXJuYW1lIjoibXVjaWFqb2VAZ21haWwuY29tIiwiZ2l2ZW5fbmFtZSI6Ikpvc2VwaCIsImZhbWlseV9uYW1lIjoiTXVzeWEiLCJlbWFpbCI6Im11Y2lham9lQGdtYWlsLmNvbSJ9.ZYujW5hPOEp3b3_-E8r0pr5Xp30PZFI4T7hD1n4pPWxEZ_htayGE_zYXS4Ogyn-mAOF_XOKGwESLBWYPm15joPgoRR7vj7EKbuq1KG2Ljn8TzWkLRS_kBGaaKbc3207o00BBRW28zmLtvL8nugUUmT8ycAM7PoRMaL1hQq0cXBE4qpF8nCAhJ3oqlXO8mTeEFwKQod6qi-dHSr01ntiejQNUTdfBhSnhaCi6y-KhoCC9wt1hNpoZf_cJIq7sb_sj-VLNNIdRZku-j7BQyGP0bl7mRpjqNff0_zLFUC7sLHPvZFNMS3C0JZuvp_7-vWWttL9BDvUATkt7VESvYZzEFA"

# remote_url: str = "http://localhost/b827ebdfc68ad595/devices"
remote_url: str = "https://remote.waziup.io/b827ebdfc68ad595/devices"

def get_cookies(username:str, password: str) -> str:
  url = "https://api.waziup.io/api/v2/auth/token"
  # url = "http://127.0.0.1/auth/token"

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

  cookie: str = response.text
  # cookie: str = response.cookies
  
  return cookie

def get_devices(cookie:str="") -> str:
  response = requests.get(url=remote_url, cookies=cookie
    )
  
  if (response.status_code==200):
    print(f"[SUCCESS] -> {response.status_code}")
    res = response.text
    return res

  elif response.status_code == 401:
    print(f"[UNAUTHORIZED] -> {response.status_code}")
    # get cookie and refetch 
    # exit after 3 trials
    return "Error occured"
  else:
    print(f'[ERROR] -> {response.status_code}')
    return  "Error  occured"

print("Starting...")

retrieved_cookie = get_cookies(username="muciajoe@gmail.com", password="0757405701Jm")
print("COOKIE: -> ",retrieved_cookie)

tunnel_cookies = {'WaziupTunnelToken': retrieved_cookie}

devices = get_devices(cookie=tunnel_cookies)
print(devices)
  