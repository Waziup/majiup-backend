import os
import openai
import requests
from dotenv import load_dotenv

# Now you can access the environment variables like this
import os

MAJIUP_URL = "http://localhost:8081/tanks"

KEY = os.getenv("SECRET_KEY")

openai.api_key = KEY

# Load environment variables from the .env file
load_dotenv()

def get_tank_data(url):
    resp = requests.get(url)
    data = resp.json()
    return data

def ask_copilot(query):
    
    max_tokens = 256

    query = query
    response = response = openai.Completion.create(
        model="text-davinci-003",
        prompt=query,
        temperature=1,
        max_tokens=max_tokens,
        top_p=1,
        frequency_penalty=0,
        presence_penalty=0,
        stop=None,
    )

    reply = response['choices'][0]['text']
    return reply

while True:
    tanks = get_tank_data(url=MAJIUP_URL)
    query = str(input("Ask Majiup Copilot: "))
    reply = ask_copilot(query + "\nTank data is {}".format(tanks))
    print("Answer: ", reply)