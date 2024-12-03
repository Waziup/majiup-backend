import unittest
import requests
import xmlrunner

majiup_url = "http://localhost:8081/api/tanks"
wazigate_url = "http://localhost/devices"

headers = {
    'Content-Type' : 'application/json',
}

class TestMajiupTanks(unittest.TestCase):
    
    # Set the base url for majiup and the tank ID to be used in this test class

    # The DELETE endpoint is not tested since the device is a real tank  (TO DO-Create a post API for Majiup)
    def setUp(self) -> None:
        tank_setup_data = {
            "name":"Testing Tank",
            "meta":{
                "settings": {
                    "height": 1650,
                    "capacity": 2000,
                    "maxalert": 1900,
                    "minalert": 100
                }        
            },
            "sensors":[
                {
                    "name" :"Temperature Sensor",
                    "value": 28,
                    "meta": {
                        "kind":"WaterThermometer"
                    }
                },
                {
                    "meta":{
                        "kind": "WaterLevel"
                    },
                    "name":"WaterLevel",
                    "value": 140
                },
                {
                    "name": "Water Quality Sensor",
                    "time": "2023-07-07T09:23:15.76Z",
                    "meta": {
                    "kind": "WaterPollutantSensor"
                    },
                    "value": 611
                }
            ],
            "actuators":[
                {
                    "name":"Pump",
                    "meta":{
                        "kind":"Motor"
                    },
                    "value":0
                }
            ]
        }
        self.base_url = majiup_url     
        res = requests.post(f"{wazigate_url}", json=tank_setup_data, headers=headers)
        self.tank_id = res.json()
        response2 = requests.get(f"{self.base_url}")
        tanks = response2.json()        
        self.tank_name = tanks[0]['name']
        # print("Running tests for {} with id {}".format(self.tank_name, self.tank_id))
    
    # Get all the tanks under majiup
    def test_get_tanks(self):
        response = requests.get(f"{self.base_url}")
        self.assertEqual(response.status_code, 200)
    
    # Get a specific tank by ID
    def test_get_tank_id(self):
        response = requests.get(f"{self.base_url}/{self.tank_id}")
        self.assertEqual(response.status_code, 200)
    
    # Get the sensors connected to majiup tank
    def test_get_sensors_in_tank(self):
        response = requests.get(f"{self.base_url}/{self.tank_id}/tank-sensors")
        self.assertEqual(response.status_code, 200)
    
    # Get the pump connected to the tank
    def test_get_pumps_in_tank(self):
        response = requests.get(f"{self.base_url}/{self.tank_id}/pumps")
        self.assertEqual(response.status_code, 200)
    
    # Return sensor historical data
    def test_get_sensor_history(self):
        response = requests.get(f"{self.base_url}/{self.tank_id}/tank-info")
        self.assertEqual(response.status_code, 200)
    
    # Change a tank name
    def test_change_tank_name(self):
        response = requests.post(f"{self.base_url}/{self.tank_id}/name", data=self.tank_name, headers=headers)
        self.assertEqual(response.status_code, 200)

    # Delete tank by ID (Tear down)
    # def test_delete_tank(self):
    #     response = requests.delete(f"{self.base_url}/{self.tank_id}")
    #     self.assertEqual(response.status_code, 200)

class TestMajiupTankMetaData(unittest.TestCase):
    def setUp(self) -> None:
        response = requests.get(f"{majiup_url}")
        tanks = response.json()
        self.tank_id = tanks[0]['id']
        self.tank_name = tanks[0]['name']
        self.base_url = majiup_url + str("/"+self.tank_id) + "/meta"
    
    def test_get_tank_meta_field(self):
        response = requests.get(f"{self.base_url}")
        self.assertEqual(response.status_code, 200)
    
    def test_change_tank_meta_field(self):
        metadata = {
            "settings": {
                "capacity": 3400,
                "height": 1650,
                "maxalert": 1900,
                "minalert": 100
            },
            "notifications":{		
                "Messages":[
                    {
                        "message":"Notification Message",
                        "read_status":False
                    }
                ]
            }
        }
        response = requests.post(f"{self.base_url}", json=metadata, headers=headers)
        msg = response.json()['message']
        self.assertEqual(msg, 'Meta field updated successfully')

class TestMajiupSensors(unittest.TestCase):
    def setUp(self) -> None:
        response = requests.get(f"{majiup_url}")
        tanks = response.json()
        self.tank_id = tanks[0]['id']
        self.tank_name = tanks[0]['name']
        self.base_url = majiup_url + str("/"+self.tank_id) + '/tank-sensors'

    # Get water sensor information on water level
    def test_get_water_level_data(self):
        response = requests.get(f"{self.base_url}/waterlevel")
        self.assertEqual(response.status_code, 200)

    # Get water current water level value
    def test_get_water_level_value(self):
        response = requests.get(f"{self.base_url}/waterlevel/value")
        self.assertEqual(response.status_code, 200)

    # Get water current water level values stored
    def test_get_water_level_values(self):
        response = requests.get(f"{self.base_url}/waterlevel/values")
        self.assertEqual(response.status_code, 200)
    
    # Get temperature sensor information
    def test_get_water_temp_data(self):
        response = requests.get(f"{self.base_url}/water-temperature")
        self.assertEqual(response.status_code, 200)

    # Get temperature value
    def test_get_water_temp_value(self):
        response = requests.get(f"{self.base_url}/water-temperature/value")
        self.assertEqual(response.status_code, 200)

    # Get temperature values stored
    def test_get_water_temp_values(self):
        response = requests.get(f"{self.base_url}/water-temperature/values")
        self.assertEqual(response.status_code, 200)

    # Get water quality sensor information
    def test_get_water_quality_data(self):
        response = requests.get(f"{self.base_url}/water-quality")
        self.assertEqual(response.status_code, 200)

    # Get water quality current value
    def test_get_water_quality_value(self):
        response = requests.get(f"{self.base_url}/water-quality/value")
        self.assertEqual(response.status_code, 200)

    # Get water quality values stored
    def test_get_water_quality_values(self):
        response = requests.get(f"{self.base_url}/water-quality/values")
        self.assertEqual(response.status_code, 200)
    
class TestMajiupPump(unittest.TestCase):
    def setUp(self) -> None:
        response = requests.get(f"{majiup_url}")
        tanks = response.json()
        self.tank_id = tanks[0]['id']
        self.tank_name = tanks[0]['name']
        self.base_url = majiup_url + str("/"+self.tank_id) + '/pumps'

    # Get the current pump status
    def test_get_pump_state(self):
        response = requests.get(f"{self.base_url}/state")
        self.assertEqual(response.status_code, 200)
    
    # Get the current pump status history
    def test_get_pump_state_states(self):
        response = requests.get(f"{self.base_url}/states")
        self.assertEqual(response.status_code, 200)
    
    # Get the change pump status
    def test_change_pump_state(self):
        
        # Turn the pump on
        response1 = requests.post(f"{self.base_url}/state", json={"value": 1}, headers=headers)
        self.assertEqual(response1.status_code, 200)

        # Turn the pump off
        response2 = requests.post(f"{self.base_url}/state", json={"value": 0}, headers=headers)
        self.assertEqual(response2.status_code, 200)

if __name__ == "__main__":
    with open("test_results.xml", "wb") as output:
        runner = xmlrunner.XMLTestRunner(output=output)
        unittest.main(testRunner=runner)
