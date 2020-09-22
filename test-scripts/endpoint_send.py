from datetime import datetime, timedelta
import requests
import random
import string
import json
import hmac
import hashlib
import array
import time

from Crypto.Cipher import AES
from Crypto.Random import get_random_bytes

URL = "http://node1.puyuma.org:6666/mam/send"
SEED = "RVAFYAELBSCMJGAJLKCEXUMJJMHYTHZSXJTELSUSGTSIXXIMSEKCJCMNVIBAVOTFPLDOCSRZNKXZRTNSS"
PRIVATE_KEY = "n2r4u7x!A%D*G-KaPdSgVkYp3s6v8y/B"
DEVICE_ID = "1234567890ABCDEF"
IV = [43, 14, 145, 54, 106, 123, 234, 233, 175, 66, 106, 177, 224, 90, 248, 73]

HASH_ALGO = hashlib.sha256
SIG_SIZE = HASH_ALGO().digest_size


def generate_obd2_rand_data():
    vin_num = ''.join(random.choices(
        string.ascii_uppercase+string.digits, k=17))

    obd2_data = {
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

    return obd2_data


def aes_encrypt(data):
    cipher = AES.new(PRIVATE_KEY, AES.MODE_CBC, bytes(IV))
    pad = AES.block_size - len(data) % AES.block_size
    data = data + pad * chr(pad)
    ciphertext = cipher.encrypt(data)
    hamc_input = bytes(IV) + ciphertext
    sig = hmac.new(PRIVATE_KEY.encode(), hamc_input, HASH_ALGO).digest()
    return (ciphertext, sig)


def main():
    message = {"timestamp": str(
        datetime.now()), "device_id": DEVICE_ID, "obd2_data": generate_obd2_rand_data()}
    ciphertext, hmac_sig = aes_encrypt(json.dumps(message))

    bytes_arr = array.array('B', IV)  # IV(16 bytes)
    # timestamp(20 bytes)
    #bytes_arr += array.array('B', '{:020d}'.format(int(time.time())).encode())
    #bytes_arr += array.array('B', hmac_sig)  # hmac(32 bytes)
    # ciphertext length(10 bytes)
    #bytes_arr += array.array('B', '{:010d}'.format(len(ciphertext)).encode())
    # bytes_arr += array.array('B', ciphertext)

    print("message array:")
    print(bytes_arr)

    http_post_data = {
        "x-api-key": "3bff60110c9f0e40628031b16dfefe2eb45a6b0c2237b49ef5109f9aa8eb19a8", 
        "data": {
            "seed": SEED, 
            "message": "".join(map(chr, bytes_arr.tobytes()))
            }, 
            "protocol": "MAM_V1"
        }
    print(http_post_data)
    r = requests.post(URL, json=http_post_data)

    print(r.text)
    print("Finish")


if __name__ == "__main__":
    main()
