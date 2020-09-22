from datetime import datetime
import requests
import random
import string
import json

URL = "http://0.0.0.0:8000/mam/send"
SEED = ''.join(random.choices(string.ascii_uppercase+string.digits, k=81))
print("Seed:" + SEED)

for i in range(0, 30):
    vin_num = ''.join(random.choices(string.ascii_uppercase+string.digits, k=17))

    obd2_data = {
        "timestamp": str(datetime.now()),
        "vin": vin_num,
        "engine_load": random.randint(1, 100),
        "engine_coolant_temperature": random.randint(1, 100),
        "fuel_pressure": random.randint(1, 100),
        "engine_speed": random.randint(1, 100),
        "vehicle_speed": random.randint(1, 100),
        "intake_air_temperature": random.randint(1, 100),
        "mass_air_flow": random.randint(1, 100),
        "fuel_tank_level_input": random.randint(1, 100),
        "absolute_barometric_pressure": random.randint(1, 100),
        "control_module_voltage": random.randint(1, 100),
        "throttle_position": random.randint(1, 100),
        "ambient_air_temperature": random.randint(1, 100),
        "relative_accelerator_pedal_position": random.randint(1, 100),
        "engine_oil_temperature": random.randint(1, 100),
        "engine_fuel_rate": random.randint(1, 100),
        "service_distance": random.randint(1, 100),
        "anti_lock_barking_active": random.randint(1, 100),
        "steering_wheel_angle": random.randint(1, 100),
        "position_of_doors": random.randint(1, 100),
        "right_left_turn_signal_light": random.randint(1, 100),
        "alternate_beam_head_light": random.randint(1, 100),
        "high_beam_head_light": random.randint(1, 100),
    }

    obd2_json = json.dumps(obd2_data)
    data = {"seed": SEED, "message": obd2_json}
    r = requests.post(URL, json={
        "x-api-key": "3bff60110c9f0e40628031b16dfefe2eb45a6b0c2237b49ef5109f9aa8eb19a8", "data": data, "protocol": "MAM_V1"})
    print(str(i)+":"+r.text)
print("Finish")
